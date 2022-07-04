package main

import (
	"context"
	"log"
	"math"
	"net"
	"strconv"
	"sync"

	streampb "github.com/macduyhai/grpcStreamingServer/streamproto"
	"google.golang.org/grpc"
)

type server struct {
	streampb.UnimplementedApiProtoServer
}

// Grpc type Unary  API
func (s *server) Greeter(ctx context.Context, req *streampb.HelloRequest) (*streampb.HelloResponse, error) {
	log.Println("Greeter")
	return &streampb.HelloResponse{
		Message: "Hello " + req.Name + ". I am Simon",
	}, nil
}

// Grpc type Server streaming API
func (s *server) CheckPrimeNumber(req *streampb.Request, stream streampb.ApiProto_CheckPrimeNumberServer) error {
	// log.Println("Check Numbers Prime")
	numsCheck := req.Numbers
	// log.Printf("List number check:%v", numsCheck)
	var wg sync.WaitGroup
	for i, n := range numsCheck {
		log.Printf("STT:%v", i)
		wg.Add(1)
		go func(num int32) {
			defer wg.Done()
			log.Printf("Check number:%v", num)
			msg := NumberIsPrime(num)
			stream.Send(&streampb.Response{
				Result: msg,
			})

		}(int32(n))
	}
	wg.Wait()
	return nil

}

// grpc Client streaming API

func main() {
	log.Println("Server Starting ...")

	lis, err := net.Listen("tcp", "0.0.0.0:8989")
	if err != nil {
		log.Fatalf("Error while create listen %v", err)
	}
	s := grpc.NewServer()
	streampb.RegisterApiProtoServer(s, &server{})

	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("Error while server %v", err)
	}

}

func NumberIsPrime(num int32) string {
	if num < 2 {
		return strconv.Itoa(int(num)) + ":Isn't prime number"
	}
	for i := 2; i <= int(math.Sqrt(float64(num))); i++ {
		if num%2 == 0 {
			return strconv.Itoa(int(num)) + ":Isn't prime number"
		}
	}
	return strconv.Itoa(int(num)) + ":Is prime number"
}
