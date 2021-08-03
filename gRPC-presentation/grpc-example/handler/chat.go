package handler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/grpc-example/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Chat struct {
	pb.UnimplementedChatServiceServer
}

func (s *Chat) SayHello(ctx context.Context, in *pb.Message) (*pb.Message, error) {
	log.Printf("Receive message from client: %v", in)

	user, _ := ctx.Value("user").(string)
	return &pb.Message{
		LastUpdated: timestamppb.New(time.Now().UTC()),
		Body: fmt.Sprintf("Hello From %s the Server:ChatHandler!", user),
		}, nil
}