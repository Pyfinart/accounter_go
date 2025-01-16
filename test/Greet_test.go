package test

import (
	v1 "accounter_go/api/helloworld/v1"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"testing"
	"time"
)

func callOpenValidateHTTP(conn *http.Client) {
	client := v1.NewGreeterHTTPClient(conn)
	request := &v1.HelloRequest{
		Name: "kratos",
	}
	reply, err := client.SayHello(context.Background(), request)
	if err != nil {
		fmt.Println("2")
		fmt.Println(err.Error())
	} else {
		fmt.Println("3")
		fmt.Println(fmt.Sprintf("[http] OpenValidate %+v\n", reply))
		//log.Println(reply.Errcode)
		op := protojson.MarshalOptions{EmitUnpopulated: true}
		//log.Println(reply.List[1].IsCancelButton)
		fmt.Println(op.Format(reply))
		fmt.Println(protojson.Format(reply))
	}
}
func TestOpenValidate(t *testing.T) {
	connHTTP, err := http.NewClient(
		context.Background(),
		//http.WithEndpoint("discovery:///wallet.api"),
		http.WithEndpoint("127.0.0.1:8000"),
		http.WithBlock(),
		http.WithMiddleware(
			metadata.Client(),
		),
		http.WithTimeout(time.Second*60),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer connHTTP.Close()
	callOpenValidateHTTP(connHTTP)
}
