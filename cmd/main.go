package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	docs "overcompute.io/pubsub/docs"
	"overcompute.io/pubsub/pkg/utils"
	ws "overcompute.io/pubsub/pkg/websocket"
	pb "overcompute.io/pubsub/pkg/websocket/publisher"
)

type IPublisherServiceServer struct {
	pb.UnimplementedPublisherServiceServer
}

// Initialize http server for websockets and gRPC server
func initServer() error {

	// since we need two instances of servers we need to run one of them in separate go routine
	go ws.InitWebsocket()

	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		return err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	lis, err := net.Listen("tcp", utils.ParseConfig().GRPC_PORT)
	if err != nil {
		return err
	}

	s := grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)))

	pb.RegisterPublisherServiceServer(s, &IPublisherServiceServer{})

	log.Println("Listening and serving gRPC on", utils.ParseConfig().GRPC_PORT)

	if err := s.Serve(lis); err != nil {
		return err
	}

	return nil

}

func main() {

	docs.SwaggerInfo.Title = "Pubsub API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Description = "Pubsub service for allowing user to listen to the queued tasks results"
	docs.SwaggerInfo.Host = "localhost:6001"

	// Handler server initialization error
	if err := initServer(); err != nil {
		log.Fatalln("Error initalizing server : ", err)
	}

}

// Publish to client message handler
func (s *IPublisherServiceServer) PublishToClient(ctx context.Context, in *pb.PublishToClient) (*pb.Response, error) {

	// Get the message and convert to raw json
	var payloadMap map[string]interface{}
	json.Unmarshal([]byte(in.GetMessage()), &payloadMap)

	payload, err := json.Marshal(payloadMap)
	if err != nil {
		log.Printf("Error parsing payload : %v", err)
		return nil, err
	}

	ws.Pool.PublishToClient(in.GetUid(), in.GetTopic(), payload)

	return &pb.Response{Status: true}, nil
}
