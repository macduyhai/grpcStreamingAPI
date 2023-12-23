package main

import (
	"context"
	"io"
	"log"
	"math"
	"net"
	"strconv"
	"sync"

	"google.golang.org/grpc"

	streampb "github.com/macduyhai/grpcStreamingServer/service/ewallet"
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
	log.Println("Check Numbers Prime")
	numsCheck := req.Numbers
	var wg sync.WaitGroup
	for i, n := range numsCheck {
		wg.Add(1)
		go func(num int32, i int) {
			defer wg.Done()
			log.Printf("SEND: %d - Check number:%v", i, num)
			msg := NumberIsPrime(num)
			stream.Send(&streampb.Response{
				Result: msg,
			})

		}(int32(n), i)

	}
	wg.Wait()
	return nil

}

// grpc Client streaming API

func (s *server) Average(stream streampb.ApiProto_AverageServer) error {
	log.Println("GetAverage")

	var sum int64
	var count int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			resp := &streampb.ResponseAverange{
				Result: float32(sum) / float32(count),
			}
			log.Printf("Average=%v", resp.Result)
			return stream.SendAndClose(resp)

		} else if err != nil {
			log.Fatalf("Error while rev Average:%v", err)
		}

		// log.Println(req.GetNumber())
		sum += req.GetNumber()
		count++
	}
}

// grpc Bi Directional API
func (s *server) FindMax(stream streampb.ApiProto_FindMaxServer) error {
	log.Println("FindMax")
	var max int64
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Printf("Max=%v", max)
			log.Println("Client finish streaming")

			return nil
		}
		if err != nil {
			log.Fatalf("Error while recv FindMax ")
			return err
		}
		num := req.GetNumber()
		if num > max {
			max = num
		}

		err = stream.Send(&streampb.ResponseFindmax{
			Max: max,
		})
		if err != nil {
			log.Fatalf("Send Max error: %v", err)
			return err
		}
	}

}

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
