package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/rpc"
	"time"

	"github.com/nhathuych/go-microservices-sandbox/broker-service/event"
	"github.com/nhathuych/go-microservices-sandbox/broker-service/logs"
	"github.com/nhathuych/go-microservices-sandbox/broker-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Deep Health Check: kiểm tra coi App Go này có kết nối được với thế giới bên ngoài (RabbitMQ) không
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

// API Gateway: Nhận request tập trung và điều phối sang các service nội bộ (Auth, Log, Mail...) dựa vô field 'action'.
func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	if err := app.readJSON(w, r, &requestPayload); err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		app.logItemViaRPC(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	case "get-user":
		app.GetUserByIDViaGRPC(w, requestPayload.User)
	default:
		app.errorJSON(w, errors.New("unknow action"))
	}
}

// Proxy xác thực: Forward yêu cầu login sang Authentication Service và chuẩn hóa kết quả trả về cho Client.
func (app *Config) authenticate(w http.ResponseWriter, authPayload AuthPayload) {
	jsonData, _ := json.MarshalIndent(authPayload, "", "\t")

	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	var jsonFromService jsonResponse
	if err := json.NewDecoder(response.Body).Decode(&jsonFromService); err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Authenticated!",
		Data:    jsonFromService.Data,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	request, err := http.NewRequest("POST", "http://logger-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) {
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	request, err := http.NewRequest("POST", "http://mail-service/send", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling mail service"))
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Message sent to" + msg.To,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logEventViaRabbitMQ(w http.ResponseWriter, l LogPayload) {
	if err := app.pushToQueue(l.Name, l.Data); err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "logged via RabbitMQ",
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.RabbitMQ)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, _ := json.MarshalIndent(&payload, "", "\t")

	return emitter.Push(string(j), "log.INFO")
}

func (app *Config) logItemViaRPC(w http.ResponseWriter, l LogPayload) {
	// Thiết lập kết nối TCP tới Logger Service qua mạng nội bộ Docker.
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	rpcPayload := RPCPayload{
		Name: l.Name,
		Data: l.Data,
	}

	var result string
	// Gọi hàm LogInfo của RPCServer đã đăng ký bên phía Logger Service qua mạng nội bộ Docker.
	if err := client.Call("RPCServer.LogInfo", rpcPayload, &result); err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: result,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) LogViaGRPC(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	if err := app.readJSON(w, r, &requestPayload); err != nil {
		app.errorJSON(w, err)
		return
	}

	conn, err := grpc.NewClient(
		"logger-service:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer conn.Close()

	c := logs.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: requestPayload.Log.Name,
			Data: requestPayload.Log.Data,
		},
	})
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "logged",
	}
	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) GetUserByIDViaGRPC(w http.ResponseWriter, userPayload UserPayload) {
	// Thiết lập kết nối tới Authentication Service
	conn, err := grpc.NewClient(
		"authentication-service:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer conn.Close()

	// Tạo client từ service đã định nghĩa trong proto
	c := proto.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Gọi hàm gRPC
	response, err := c.GetUserByID(ctx, &proto.UserRequest{
		Id: int64(userPayload.ID),
	})
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Trả về kết quả cho Client (Postman/Browser)
	payload := jsonResponse{
		Error:   false,
		Message: "User retrieved via gRPC",
		Data:    response.User,
	}
	app.writeJSON(w, http.StatusAccepted, payload)
}
