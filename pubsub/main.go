package main

import (
	"context"
	"encoding/json"
	"log"
	"net"

	"google.golang.org/grpc"

	"overcompute.io/pubsub/pkg/config"
	ws "overcompute.io/pubsub/pkg/websocket"
	pb "overcompute.io/pubsub/pkg/websocket/publisher"
)

type IPublisherServiceServer struct {
	pb.UnimplementedPublisherServiceServer
}

// Initialize http server for websockets and gRPC server
func initServer() error {

	conf, err := config.ParseConfig()

	if err != nil {
		return err
	}

	creds, err := config.ParseTransportCreds()

	if err != nil {
		return err
	}

	wsconf := ws.NewServer(conf)
	// since we need two instances of servers we need to run one of them in separate go routine
	go wsconf.InitWebsocket()

	lis, err := net.Listen("tcp", conf.GRPC_PORT)
	if err != nil {
		return err
	}

	s := grpc.NewServer(grpc.Creds(*creds))

	pb.RegisterPublisherServiceServer(s, &IPublisherServiceServer{})

	log.Println("Listening and serving gRPC on", conf.GRPC_PORT)

	if err := s.Serve(lis); err != nil {
		return err
	}

	return nil

}

func main() {

	// Handle server initialization error
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
