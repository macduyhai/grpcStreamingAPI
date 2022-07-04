package main

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"time"

	streampb "github.com/macduyhai/grpcStreamingServer/streamproto"
	"google.golang.org/grpc"
)

const (
	name     = "Anna"
	pathTest = "numberstest.json"
)

type ListNums struct {
	nums []int32 `json:"nums"`
}

func main() {
	log.Println("Client Starting")
	conn, err := grpc.Dial("0.0.0.0:8989", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Connect to host error:%v", err)
	} else {
		log.Println("Connect done")
	}
	defer conn.Close()

	c := streampb.NewApiProtoClient(conn)

	SayHello(c)
	// Test streaming server API
	ln := LoadDataTest(pathTest)
	log.Println(ln.nums)
	CheckNumberPrime(c, ln.nums)

}
func SayHello(p streampb.ApiProtoClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	log.Println(name + " :Send --> Hi")
	resp, err := p.Greeter(ctx, &streampb.HelloRequest{
		Name: name,
	})
	if err != nil {
		log.Fatalf("Send request error: %v", err)
	}

	log.Printf("Message: %s", resp.Message)
}
func CheckNumberPrime(s streampb.ApiProtoClient, arr []int32) {
	log.Println("Check Numbers Prime")
	log.Println(arr)
	stream, err := s.CheckPrimeNumber(context.Background(), &streampb.Request{
		// Numbers: []int32{1, 2, 3, 4},
		Numbers: arr,
	})

	if err != nil {
		log.Fatalf("Request error:%v", err)
	} else {
		log.Println("Request done")
	}
	for {
		resp, recvErr := stream.Recv()
		if recvErr == io.EOF {
			log.Println("Server finish streaming")
			return
		}
		log.Println(resp.Result)

	}

}
func LoadDataTest(path string) (result ListNums) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Read file Test error : %v", err)
	}
	// log.Println(string(file))
	// log.Println(result)
	err = json.Unmarshal(file, &result)
	if err != nil {
		log.Fatalf("Parser error : %v", err)
	}
	log.Println(result)
	return result
}
