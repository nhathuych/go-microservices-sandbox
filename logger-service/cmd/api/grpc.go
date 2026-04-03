package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/nhathuych/go-microservices-sandbox/logger-service/data"
	"github.com/nhathuych/go-microservices-sandbox/logger-service/logs"
	"google.golang.org/grpc"
)

// LogServer: Cấu trúc đóng vai trò là "Implementation" của gRPC service đã định nghĩa trong file .proto.
type LogServer struct {
	// Nhúng (embed) để đảm bảo tính tương thích ngược.
	// Nếu file .proto thêm hàm mới mà code chưa kịp viết, server vẫn không bị crash.
	logs.UnimplementedLogServiceServer
	// Truy cập vào Database Models để thực hiện ghi log vào MongoDB.
	Models data.Model
}

// WriteLog: Triển khai logic thực tế cho hàm RPC đã khai báo trong file .proto.
// Nhận vào context (quản lý timeout/cancel) và con trỏ LogRequest (dữ liệu từ Client).
func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	// Mapping dữ liệu từ định dạng gRPC (Generated code) sang định dạng Model nội bộ của App.
	input := req.GetLogEntry() // Lấy dữ liệu LogEntry về từ Request
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	// Ghi vô DB thông qua logic đã viết ở package data.
	if err := l.Models.LogEntry.Insert(logEntry); err != nil {
		res := &logs.LogResponse{Result: "failed"}
		return res, err
	}

	// Trả về phản hồi thành công cho Client dưới dạng con trỏ LogResponse.
	res := &logs.LogResponse{Result: "logged!"}
	return res, nil
}

// gRPCListen: Khởi tạo và vận hành gRPC server.
func (app *Config) gRPCListen() {
	// Mở cổng lắng nghe TCP.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	// Khởi tạo một thực thể gRPC server mới.
	srv := grpc.NewServer()

	// Đăng ký LogServer (vừa định nghĩa ở trên) với gRPC server.
	logServer := &LogServer{Models: app.Model}

	// Kết nối logic code vào hạ tầng mạng gRPC.
	logs.RegisterLogServiceServer(srv, logServer)

	log.Printf("gRPC Server started on port %s", gRpcPort)

	// Bắt đầu tiếp nhận và phục vụ các yêu cầu gRPC thông qua listener đã mở.
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}
