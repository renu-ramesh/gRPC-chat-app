package main

import (
	"chat_app_grpc/chatpb"
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

type ChatServiceServer struct {
	chatpb.UnimplementedChatServiceServer
	channel map[string][]chan *chatpb.Message
}

func (s *ChatServiceServer) SendMessage(msgStream chatpb.ChatService_SendMessageServer) error {
	msg, err := msgStream.Recv()

	if err == io.EOF {
		return nil
	}

	if err != nil {
		return err
	}

	go func() {
		streams := s.channel[msg.Channel.Name]
		for _, msgChan := range streams {
			msgChan <- msg
		}
	}()

	return nil
}

func (s *ChatServiceServer) JoinChannel(ch *chatpb.Channel, msgStream chatpb.ChatService_JoinChannelServer) error {

	msgChannel := make(chan *chatpb.Message)
	s.channel[ch.Name] = append(s.channel[ch.Name], msgChannel)

	for {
		select {
		case <-msgStream.Context().Done():
			return nil
		case msg := <-msgChannel:
			fmt.Printf("Message History [ message: %v from %v ]\n", msg.Message, msg.Sender)
			msgStream.Send(msg)
		}
	}
}

func newServer() *ChatServiceServer {
	s := &ChatServiceServer{
		channel: make(map[string][]chan *chatpb.Message),
	}
	return s
}

func main() {
	fmt.Println("--- SERVER APP ---")
	lis, err := net.Listen("tcp", "localhost:5400")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	chatpb.RegisterChatServiceServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
