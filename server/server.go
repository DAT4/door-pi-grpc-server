package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	door "door-pi-grpc-server/pictrl"
	pb "github.com/dat4/grpc-test/mygrpc"
	"google.golang.org/grpc"
)

type doorServer struct {
	pb.UnimplementedDoorServiceServer
	mu sync.Mutex
}

func (s *doorServer) Login(ctx context.Context, user *pb.User) (*pb.Token, error) {
	fmt.Println(user.Username, user.Password)
	return nil, nil
}
func (s *doorServer) OpenDoor(stream pb.DoorService_OpenDoorServer) (err error) {
	door.OpenPin(17, door.HIGH)
	for {
		x, err := stream.Recv()
		if err != nil {
			break
		}
		fmt.Println(x)
	}
	door.OpenPin(17, door.LOW)
	if err == io.EOF {
		stream.SendAndClose(&pb.DoorResponse{Ok: "ok"})
	}
	fmt.Println("done")
	return nil
}

func newServer() *doorServer {
	s := &doorServer{}
	return s
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterDoorServiceServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
