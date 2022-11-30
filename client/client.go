package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"chat_app_grpc/chatpb"
	"chat_app_grpc/internal/common"
	"chat_app_grpc/internal/db"

	"google.golang.org/grpc"
)

var channelDetails = flag.String("channel", "", "Channel name for chatting")
var userName = flag.String("username", "", "Senders name")
var tcpServer = flag.String("server", ":5400", "Tcp server")

func joinChannel(ctx context.Context, client chatpb.ChatServiceClient) {

	if *channelDetails != "" {
		err := common.ValidateUserChannel(*userName, *channelDetails)
		if err != nil {
			log.Fatalf("Invalid group name: %v", err)
		}
		fmt.Printf("Joined channel: %v \n", *channelDetails)
	}
	channel := chatpb.Channel{Name: *channelDetails, SendersName: *userName}
	stream, err := client.JoinChannel(ctx, &channel)
	if err != nil {
		log.Fatalf("client.JoinChannel(ctx, &channel) throws: %v", err)
	}

	waitc := make(chan struct{})

	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive message from channel joining.")
			}

			if *userName != in.Sender {
				fmt.Printf("%v : %v \n", in.Sender, in.Message)
			}
		}
	}()
	<-waitc
}
func sendMessage(ctx context.Context, client chatpb.ChatServiceClient, message string) {
	stream, err := client.SendMessage(ctx)
	if err != nil {
		log.Printf("Cannot send message: error: %v", err)
	}
	msg := chatpb.Message{
		Channel: &chatpb.Channel{
			Name:        *channelDetails,
			SendersName: *userName},
		Message: message,
		Sender:  *userName,
	}
	stream.Send(&msg)

}

func main() {

	flag.Parse()
	db.ConnectDatabase()

	if *userName != "" {
		err := common.ValidateUser(*userName)
		if err != nil {
			log.Fatalf("Fail to authenticate user: %v", err)
		}

		fmt.Println("--- CLIENT APP ---")
		var opts []grpc.DialOption
		opts = append(opts, grpc.WithBlock(), grpc.WithInsecure())

		conn, err := grpc.Dial(*tcpServer, opts...)
		if err != nil {
			log.Fatalf("Fail to dail: %v", err)
		}

		defer conn.Close()
		ctx := context.Background()
		client := chatpb.NewChatServiceClient(conn)

		go joinChannel(ctx, client)

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			go sendMessage(ctx, client, scanner.Text())
		}

	} else {
		log.Fatalf("Fail to authenticate user")
	}

}
