package main

import (
	"context"
	"log"
	"time"

	"github.com/nhathuych/go-microservices-sandbox/logger-service/data"
)

type RPCServer struct{}

type RPCPayload struct {
	Name string
	Data string
}

func (R *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error writing to mongo:", err)
		return err
	}

	// Gán kết quả vô *resp. Lưu ý: Đây chỉ là quy ước cú pháp của net/rpc để Go tự đóng gói (serialize)
	// và gởi dữ liệu qua mạng, các service không hề dùng chung vùng nhớ hay con trỏ thiệt
	*resp = "Processed payload via RPC:" + payload.Name

	return nil
}
