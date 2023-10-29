package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/tessellated-io/mail-in-rebates/paymaster/bursar"
	proto "github.com/tessellated-io/mail-in-rebates/paymaster/server/proto"
	"github.com/tessellated-io/mail-in-rebates/paymaster/tracker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type paymasterServer struct {
	server         *grpc.Server
	bursar         *bursar.Bursar
	addressTracker *tracker.AddressTracker

	// TODO: this should live in the bursar.
	lock sync.Mutex
}

func StartPaymasterServer(
	bursar *bursar.Bursar,
	addressTracker *tracker.AddressTracker,
	port int,
) {
	server := grpc.NewServer()
	service := &paymasterServer{
		bursar:         bursar,
		server:         server,
		addressTracker: addressTracker,

		lock: sync.Mutex{},
	}
	proto.RegisterFundingServiceServer(server, service)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("Failed to listen on port %d: %v\n", port, err)
		return
	}

	log.Printf("gRPC server listening on port %d...", port)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
		return
	}
}

func (s *paymasterServer) Fund(ctx context.Context, req *proto.FundingRequest) (*proto.FundingResponse, error) {
	fmt.Printf("💰 Received request for funds (address=%s, prefix=%s)\n", req.Address, req.AddressPrefix)

	// Spawn a thread
	err := s.asyncFund(ctx, req)
	if err != nil {
		return nil, err
	}

	return &proto.FundingResponse{}, nil
}

// Use an async function that acquires a lock, then lets funds settle. This is slow, but fixes the question of nonce conflicts.
// TODO: handle return value?
func (s *paymasterServer) asyncFund(ctx context.Context, req *proto.FundingRequest) error {
	// Lock
	s.lock.Lock()
	defer s.lock.Unlock()
	fmt.Printf("💰 Acquired lock to send funds. Waiting for txs to settle. (address=%s, prefix=%s)\n", req.Address, req.AddressPrefix)
	time.Sleep(20 * time.Second) // like 3ish blocks
	fmt.Printf("💰 Proceeding to process funds with txs settled. (address=%s, prefix=%s)\n", req.Address, req.AddressPrefix)

	hash, err := s.bursar.SendFunds(req.Address, req.AddressPrefix)
	if err != nil {
		fmt.Printf("🛑 Failed to send funds: %s (address: %s)\n", err, req.Address)
		return err
	}
	fmt.Printf("💰 Sent funds in hash %s (address=%s)\n", hash, req.Address)

	// Get IP
	p, ok := peer.FromContext(ctx)
	clientIP := "unknown"
	if ok {
		clientIP = p.Addr.String()
	}
	err = s.addressTracker.AddAddress(fmt.Sprintf("%s %s", req.Address, clientIP))
	if err != nil {
		return err
	}
	return nil
}
