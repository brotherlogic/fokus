package main

import (
	"context"
	"fmt"
	"log"
	"os"
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

	sig := ""
	if len(os.Args) > 2 {
		switch os.Args[1] {
		case "code":
			sig = "type-code"
		case "process":
			sig = "type-process"
		case "research":
			sig = "type-research"
		}
	}

	fclient := pb.NewFokusServiceClient(conn)
	fok, err := fclient.GetFokus(context.Background(), &pb.GetFokusRequest{Label: sig})

	if err != nil {
		log.Fatalf("Unable to get fokus: %v", err)
	}

	fmt.Printf("%v [%v]\n", fok, time.Unix(fok.GetGivenTime(), 0))
}
