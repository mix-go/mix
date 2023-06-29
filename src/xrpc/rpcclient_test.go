package xrpc

import (
	"context"
	"fmt"
	pb "github.com/mix-go/xrpc/testdata"
	"log"
	"net/http"
	"strings"
	"testing"
)

func TestNewGrpcClient(t *testing.T) {
	conn, err := NewGrpcClient("127.0.0.1:50000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := pb.NewOrderClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), CallTimeout)
	resp, err := client.RequestForRelease(ctx, &pb.ReleaseRequest{
		OrderNumber: "123456789",
	})
	fmt.Println(resp)
}

func TestNewGatewayClient(t *testing.T) {
	tlsConf, err := LoadTLSClientConfig("/certificates/ca.pem", "/certificates/client.pem", "/certificates/client.key")
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConf,
		},
	}
	resp, err := client.Post("http://127.0.0.1:50001/v1/request_for_release", "application/json", strings.NewReader(`{"order_number":"123456789"}`))
	fmt.Println(resp.Body)
}
