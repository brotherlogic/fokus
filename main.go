package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	pb "github.com/brotherlogic/fokus/proto"
	"github.com/brotherlogic/fokus/server"

	auth_client "github.com/brotherlogic/auth/client"
)

var (
	port        = flag.Int("port", 8080, "Server port for grpc traffic")
	metricsPort = flag.Int("metrics_port", 8081, "Metrics port")
)

func main() {
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	authModule, err := auth_client.NewAuthInterceptor(ctx)
	cancel()
	if err != nil {
		log.Fatalf("Unable to get auth client: %v", err)
	}

	s := server.NewServer()

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%v", *metricsPort), nil)
		log.Fatalf("fokus is unable to serve metrics: %v", err)
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("fokus is unable to listen on the min grpc port %v: %v", *port, err)
	}
	gs := grpc.NewServer(grpc.UnaryInterceptor(authModule.AuthIntercept))
	pb.RegisterFokusServiceServer(gs, s)

	log.Printf("Serving on port :%d", *port)
	if err := gs.Serve(lis); err != nil {
		log.Fatalf("fokus is unable to serve grpc for some reason: %v", err)
	}
	log.Fatalf("fokus has closed the grpc port for some reason")
}
