package main

import (
	"log"
	"time"

	pb "./identity_struct"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:8081", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewIdentifierClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.IssueToken(ctx, &pb.IssueTokenRequest{Uid: "some hex"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	req, err := c.ParseToken(context.Background(), &pb.ParseTokenRequest{AccessToken: r.AccessToken})
	if err != nil {
		log.Fatalf("Aaaa %s", err.Error())
	}
	log.Printf("Access token: %v", req.Uid)
}
