package main

import (
	"log"

	"github.com/grpc-example/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := pb.NewChatServiceClient(conn)

	ctx := metadata.NewOutgoingContext(context.TODO(), metadata.New(map[string]string{
		"Authorization": "Bearer cm9tYW46cHdk",
	}))

	response, err := c.SayHello(ctx, &pb.Message{
		Id:           1,
		Body:         "Hello From Client!",
		PhoneNumbers: []string{"111", "222"},
		PersonInfo: &pb.Person{
			Name:     "Roman",
			LastName: "Kosyi",
		},
	})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}

	log.Printf("Response from server: %s", response.Body)
}
