package main

import (
	"log"
	"net"

	"github.com/ricarrrdoaraujo/grpc-go/pb"
	"github.com/ricarrrdoaraujo/grpc-go/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// 1 - Cria o servidor
// 2 - Cria os métodos que quer implementar
// 3 - Registra

func main() {

	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, services.NewUserService())
	reflection.Register(grpcServer) // modo reflection para ser lido pelo client

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Could not serve: %v", err)
	}
}

/*
	Comando do client - https://github.com/ktr0731/evans
	evans -r repl --host localhost --port 50051
*/
