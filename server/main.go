package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"unicode"

	pb "github.com/j-griffith/cinderclient-grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

// server is used to implement cinderattacher.AttacherServer.
type server struct{}

func removeSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

func parseAttachOutput(o string) (map[string]string, error) {
	m := make(map[string]string)
	lines := strings.Split(strings.TrimSpace(string(o)), "\n")
	for _, l := range lines {
		if !strings.Contains(l, "+---") && !strings.Contains(l, "Property | Value") {
			l = removeSpaces(l)
			l = strings.Trim(l, "|")
			entries := strings.Split(l, "|")
			m[string(entries[0])] = string(entries[1])
		}

	}
	return m, nil
}

// Attach implements cinderattacher.AttachVolume
func (s *server) Attach(ctx context.Context, in *pb.AttachRequest) (*pb.AttachResponse, error) {
	fmt.Printf("issuing Attach request for volume id: %s\n", in.Id)
	c := exec.Command("cinder", "local-attach", "--mountpoint", "/dev/x", in.Id)
	out, err := c.CombinedOutput()
	if err != nil {
		fmt.Printf("error response from attach command: %s\n", out)
		fmt.Printf("error message: %+v\n", err)
		return &pb.AttachResponse{}, err
	}

	m, err := parseAttachOutput(string(out))
	return &pb.AttachResponse{PublishInfo: m}, nil
}

// Detach implements cinderattacher.DetachVolume
func (s *server) Detach(ctx context.Context, in *pb.DetachRequest) (*pb.DetachResponse, error) {
	fmt.Printf("Our attach request is for id: %s\n", in.Id)
	return nil, nil
}
func parseSocket(arg string) ([]string, error) {
	r := strings.Split(arg, ",")
	if len(r) != 2 {
		return nil, fmt.Errorf("invalid socket argument: %s", arg)
	}
	fmt.Printf("Using provided args to create the socket: Type=%s, Name=%s\n", r[0], r[1])
	return r, nil
}

func main() {
	// Here be a BUG!  If using unix socket and you kill the server we're not cleaning up the socket file, need to figure that out
	requestedSocket := flag.String("socket", "unix,/tmp/cinderattach-socket", "Socket to use for grpc connection.  To use tcp use the form: `tcp,<ip-address>:<port>`")
	flag.Parse()
	socket, err := parseSocket(*requestedSocket)
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}
	lis, err := net.Listen(socket[0], socket[1])
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAttacherServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
