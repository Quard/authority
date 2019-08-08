package internal_api

import (
	"log"
	"net"

	"github.com/Quard/authority/internal/storage"
	grpc "google.golang.org/grpc"
)

type Opts struct {
	bind string
}

type internalAPIServer struct {
	opts       Opts
	listener   net.Listener
	grpcServer *grpc.Server
	storage    storage.Storage
}

func NewInternalAPIServer(bind string, stor storage.Storage) internalAPIServer {
	srv := internalAPIServer{storage: stor, opts: Opts{bind: bind}}
	var err error
	srv.listener, err = net.Listen("tcp", srv.opts.bind)
	if err != nil {
		log.Fatal(err)
	}
	var opts []grpc.ServerOption
	srv.grpcServer = grpc.NewServer(opts...)
	RegisterInternalAPIServer(srv.grpcServer, &srv)

	return srv
}

func (srv internalAPIServer) Run() {
	log.Printf("INTERNAL gRPC server listen on: %s", srv.opts.bind)
	srv.grpcServer.Serve(srv.listener)
}
