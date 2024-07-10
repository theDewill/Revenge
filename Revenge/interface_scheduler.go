package Revenge

import (
	"context"
	"fmt"
	"net"
	"ssego/Revenge/pb"

	"google.golang.org/grpc"
)

//Interface Implementations

type Message interface {
	Send() (string, error)
	Schedule() (string, error)
}

// Strictural Intefaces
type interface_scheduler struct {
	Server      *grpc.Server
	Connections Channel_Manager
}

type Channel struct {
	party_1, party_2 Runner
	MsgBuff          []Message //Slice of messages්‍
}

type Channel_Manager struct {
	Channel_pool map[string]Channel
}

type ChannelServiceServer struct {
	pb.UnimplementedChannelServiceServer
	channelManager *Channel_Manager
}

func (s *ChannelServiceServer) SendMessage(ctx context.Context, msg *pb.Message) (*pb.Message, error) {
	// Implement the logic to handle incoming messages
	// You might want to route it to another runner based on your logic
	return &pb.Message{Content: "Received: " + msg.Content}, nil
}

func NewIS() *interface_scheduler {
	return &interface_scheduler{}
}

func (IS *interface_scheduler) start() {

	listner, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Println("Error in starting the server")
	}

	IS.Server = grpc.NewServer()
	if err := IS.Server.Serve(listner); err != nil {
		fmt.Println("Grpc server error")
	}

}
