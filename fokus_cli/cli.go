package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	pb "github.com/brotherlogic/fokus/proto"
)

func main() {
	conn, err := grpc.Dial("fokus.brotherlogic-backend.com:80")
	if err != nil {
		log.Fatalf("Dial fail: %v", err)
	}

	fclient := pb.NewFokusServiceClient(conn)
	fok, err := fclient.GetFokus(context.Background(), &pb.GetFokusRequest{})

	if err != nil {
		log.Fatalf("Unable to get fokus: %v", err)
	}

	fmt.Printf("%v\n", fok)
}
