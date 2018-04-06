package main

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	pb "github.com/alextanhongpin/traefik-grpc/proto"
)

type echoServer struct{}

func (s *echoServer) Echo(ctx context.Context, msg *pb.EchoRequest) (*pb.EchoResponse, error) {
	return &pb.EchoResponse{
		Text: msg.Text,
	}, nil
}

func main() {
	sslCert := os.Getenv("SSL_CERT")
	sslKey := os.Getenv("SSL_KEY")
	port := os.Getenv("PORT")

	log.Println(sslCert, sslKey)

	//
	// CRED
	//
	// BackendCert, _ := ioutil.ReadFile("/runs/secrets/backend.cert")
	// BackendKey, _ := ioutil.ReadFile("/runs/secrets/backend.key")

	// Generate Certificate struct
	cert, err := tls.X509KeyPair([]byte(sslCert), []byte(sslKey))
	if err != nil {
		log.Fatalf("failed to parse certificate: %v", err)
	}

	// Create credentials
	creds := credentials.NewServerTLSFromCert(&cert)

	// Use Credentials in gRPC server options
	serverOption := grpc.Creds(creds)

	//
	// SERVER
	//
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %s", err.Error())
	}

	grpcServer := grpc.NewServer(serverOption)
	defer grpcServer.Stop()

	pb.RegisterEchoServiceServer(grpcServer, &echoServer{})
	reflection.Register(grpcServer)
	log.Printf("listening to server at port *:%v. press ctrl + c to cancel.\n", port)
	log.Fatal(grpcServer.Serve(lis))
}
