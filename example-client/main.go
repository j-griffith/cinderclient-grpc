package main

import (
	"context"
	"log"
	"time"

	pb "github.com/j-griffith/cinderclient-grpc/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("unix:///tmp/grpc-example-socket", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewAttacherClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Attach(ctx, &pb.AttachRequest{Id: "foobar-volume"})

	if err != nil {
		log.Fatalf("well that sucks: %v", err)
	}
	log.Printf("attach response: %s", r.PublishInfo["path"])
}
