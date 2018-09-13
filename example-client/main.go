package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/j-griffith/cinderclient-grpc/proto"
	"google.golang.org/grpc"
)

func main() {
	sock := flag.String("socket", "unix:///tmp/cinderattach-socket", "Socket to use for grpc connection.  To use tcp use the form: `<ip-address>:<port>`")
	volumeUUID := flag.String("volume", "", "The UUID of the Cinder Volume to perform the local attach on.")
	flag.Parse()

	conn, err := grpc.Dial(*sock, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewAttacherClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Attach(ctx, &pb.AttachRequest{Id: *volumeUUID})

	if err != nil {
		log.Fatalf("well that sucks: %v", err)
	}
	log.Printf("attach response: %s", r.PublishInfo["path"])
}
