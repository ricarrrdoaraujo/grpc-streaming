package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/ricarrrdoaraujo/grpc-go/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect the rpc server: %v", err)
	}

	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	//AddUser(client)
	// AddUserVerbose(client)
	// AddUsers(client)
	AddUserStreamBi(client)
}

func AddUser(client pb.UserServiceClient) {

	req := &pb.User{
		Id:    "0",
		Name:  "Joao",
		Email: "J@J.com",
	}

	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not do gRPC request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Joao",
		Email: "J@J.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not do gRPC request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not receive the message: %v", err)
		}
		fmt.Println("Status:", stream.Status, " - ", stream.GetUser())
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id:    "r1",
			Name:  "Ricardo",
			Email: "ric@ric.com",
		},
		&pb.User{
			Id:    "r2",
			Name:  "Ricardo2",
			Email: "ric2@ric.com",
		},
		&pb.User{
			Id:    "r3",
			Name:  "Ricardo3",
			Email: "ric3@ric.com",
		},
		&pb.User{
			Id:    "r4",
			Name:  "Ricardo4",
			Email: "ric4@ric.com",
		},
		&pb.User{
			Id:    "r5",
			Name:  "Ricardo5",
			Email: "ric5@ric.com",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBi(client pb.UserServiceClient) {

	stream, err := client.AddUserStreamBi(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	reqs := []*pb.User{
		&pb.User{
			Id:    "r1",
			Name:  "Ricardo",
			Email: "ric@ric.com",
		},
		&pb.User{
			Id:    "r2",
			Name:  "Ricardo2",
			Email: "ric2@ric.com",
		},
		&pb.User{
			Id:    "r3",
			Name:  "Ricardo3",
			Email: "ric3@ric.com",
		},
		&pb.User{
			Id:    "r4",
			Name:  "Ricardo4",
			Email: "ric4@ric.com",
		},
		&pb.User{
			Id:    "r5",
			Name:  "Ricardo5",
			Email: "ric5@ric.com",
		},
	}

	wait := make(chan int)

	//Sending data
	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	//Receiving data
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
			}
			fmt.Printf("Receiving user %v with status: %v\n", res.GetUser().GetName(), res.GetStatus())
		}
		close(wait)
	}()

	//Hold the running process
	<-wait

}
