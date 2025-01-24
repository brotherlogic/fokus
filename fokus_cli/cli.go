package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/brotherlogic/fokus/proto"
)

func main() {
	conn, err := grpc.Dial("fokus.brotherlogic-backend.com:80", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Dial fail: %v", err)
	}

	fclient := pb.NewFokusServiceClient(conn)
	fok, err := fclient.GetFokus(context.Background(), &pb.GetFokusRequest{})

	if err != nil {
		log.Fatalf("Unable to get fokus: %v", err)
	}

	fmt.Printf("%v [%v]\n", fok, time.Unix(fok.GetGivenTime(), 0))
}
