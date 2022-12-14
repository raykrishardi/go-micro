package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/raykrishardi/log-service/data"
	"github.com/raykrishardi/log-service/logs"
	"google.golang.org/grpc"
)

// Need to have the first line for every service in GRPC which is used for backward compatibility
type LogServer struct {
	logs.UnimplementedLogServiceServer             // MUST HAVE for all grpc service
	Models                             data.Models // in order to have access to Insert to mongoDB function
}

// We are now using the generated source code from protoc
func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	// Using the generated source code
	input := req.GetLogEntry() // value of input will be input.name and input.data just like in the proto file

	// write the log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{
			Result: "failed",
		}
		return res, err
	}

	res := &logs.LogResponse{
		Result: "logged!",
	}
	return res, nil
}

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for grpc: %v", err)
	}

	s := grpc.NewServer()

	logs.RegisterLogServiceServer(s, &LogServer{Models: app.Models})

	log.Printf("gRPC server started on port %s", gRpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for grpc: %v", err)
	}
}
