package main

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"

	streampb "github.com/macduyhai/grpcStreamingServer/service/ewallet"
)

const (
	name     = "Simon"
	pathTest = "numberstest.json"
)

type ListNums struct {
	Nums []int32 `json:"nums"`
}

func main() {
	log.Println("Client Starting")

	//TODO: Init Connection
	conn, err := grpc.Dial("0.0.0.0:8989", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Connect to host error:%v", err)
	} else {
		log.Println("Connect done")
	}
	defer conn.Close()

	//TODO: Register Client
	c := streampb.NewApiProtoClient(conn)

	//TODO: Process logic

	// Unary API
	SayHello(c)

	//// Streaming server API
	//ln := LoadDataTest(pathTest)
	//CheckPrime(c, ln.Nums)
	//
	////Streaming Client API
	//GetAverage(c, ln.Nums)
	//
	// Bi Directional API
	//FindMax(c, ln.Nums)

}

// SayHello : Unary API
func SayHello(p streampb.ApiProtoClient) {
	startTime := time.Now()
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
	log.Printf("Execution time %s\n", time.Since(startTime))
}

// CheckNumberPrime : Server streaming  API
func CheckPrime(c streampb.ApiProtoClient, arr []int32) {
	startTime := time.Now()
	log.Println("Check Numbers Prime")
	// log.Println(arr)
	stream, err := c.CheckPrimeNumber(context.Background(), &streampb.Request{
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
			log.Printf("Execution time %s\n", time.Since(startTime))
			return
		}
		log.Println(resp.Result)

	}

}

// Streaming Client API

func GetAverage(c streampb.ApiProtoClient, arr []int32) {
	startTime := time.Now()
	log.Println("Get Average")
	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalf("get Average err:%v", err)
	}
	// arrtest := []int{1, 2, 3, 4, 1232131231231231212}
	// for _, num := range arrtest {
	for _, num := range arr {
		err := stream.Send(&streampb.RequestAverange{
			Number: int64(num),
		})
		if err != nil {
			log.Fatalf("Send average request error:%v", err)
		}

	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while rev from server:%v", err)
	}
	log.Printf("Avarage response -> %v", resp)
	log.Printf("Number request: %v - Execution time %s\n", len(arr), time.Since(startTime))

}

//FindMax : BI Directional API

func FindMax(s streampb.ApiProtoClient, arr []int32) {
	startTime := time.Now()
	log.Println("Find Max")
	stream, err := s.FindMax(context.Background())
	if err != nil {
		log.Fatalf("FinMax err init stream:%v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	// rountine send request
	go func() {

		log.Printf("Array input :%v", arr)
		for _, num := range arr {

			err := stream.Send(&streampb.RequestFindMax{
				Number: int64(num),
			})
			if err != nil {
				log.Fatalf("Error while send FindMax:%v", err)
			}
		}
		stream.CloseSend()
		wg.Done()
	}()

	// routine recv message
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				log.Println("Server stop streaming")
				break
			}
			if err != nil {
				log.Fatalf("Error while recv Max:%v", err)
				break
			}
			log.Printf("Max:%v", resp.Max)
		}
		wg.Done()
	}()

	wg.Wait()
	log.Printf("Execution time %s\n", time.Since(startTime))

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
	// log.Println(result)
	return result
}
