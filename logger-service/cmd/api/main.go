package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/nhathuych/go-microservices-sandbox/logger-service/data"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	serverPort = "80"
	rpcPort    = "5001"
	gRpcPort   = "50051"
	mongoURL   = "mongodb://mongo:27017"
)

var client *mongo.Client

type Config struct {
	Model data.Model
}

func main() {
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Model: data.New(client),
	}

	// Register the RPC Server
	// Go dùng Reflection để tự động ánh xạ các method public của struct này
	// thành các "hàm từ xa" có thể gọi được (ví dụ: "RPCServer.LogInfo").
	err = rpc.Register(new(RPCServer))
	// Chạy RPC server trên một Goroutine riêng (non-blocking) để không làm treo luồng chính.
	// Điều này cho phép ứng dụng lắng nghe đồng thời cả RPC (cổng 5001) và HTTP (cổng 80) bên dưới.
	go app.rpcListen()

	// Start the gRPC server
	go app.gRPCListen()

	app.serve()
}

func (app *Config) rpcListen() error {
	log.Println("Starting RPC server on port", rpcPort)
	// Thiết lập lắng nghe kết nối TCP tại port 5001.
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		// Chấp nhận kết nối mới từ Client (như Broker-service).
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}

func (app *Config) serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", serverPort),
		Handler: app.route(),
	}

	fmt.Println("Starting logging web service on port", serverPort)
	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}

func connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		AuthSource: "admin",
		Username:   "admin",
		Password:   "password",
	})

	c, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Println("Error connecting: ", err)
		return nil, err
	}

	if err := c.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}

	log.Println("Connected to mongo.")
	return c, nil
}
