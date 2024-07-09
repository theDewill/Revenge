package Revenge

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

//Interface Implementations

type Message interface {
	Send()
	Schedule()
}

//Strictural Intefaces
type interface_scheduler struct {
	server net.TCPListener
}


type Channel struct {
	party_1,party_2 Runner
	msgBuff []Message //Slice of messages
}


type Channel_Manager struct {
	Channel_pool map([string]Channel)
}


func New () *interface_scheduler {
	return interface_scheduler {}
}

func (IS *interface_scheduler) start() {
	
	listner , err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Println("Error in starting the server")
	}

	IS.server = grpc.NewServer()
	if err := IS.server.Serve(listner); err != nil {
		fmt.Println("Grpc server error")
	}


}