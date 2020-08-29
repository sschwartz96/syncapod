package services

import (
	"github.com/sschwartz96/syncapod/internal/config"
	"github.com/sschwartz96/syncapod/internal/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	// should we load in config??? or just mock everything?
	config := &config.Config{GRPCPort: 123}
	s := grpc.NewServer()
}

func TestAuthenticate() {

}
