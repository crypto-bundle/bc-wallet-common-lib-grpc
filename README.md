# bc-wallet-common-lib-grpc

## Description

Library for manage grpc-server and client

Library contains:
* GRPC-**server** init options and helper functions
* GRPC-**client** init options and helper functions
* Examples of create GRPC-server and client

## Usage example

Examples of create connection and write database communication code

### Init and start GRPC-server

```go
package main

import (
	pbApi "gitlab.heronodes.io/bc-platform/bc-wallet-common-lib-grpc/pkg/grpc/grpc_handlers/proto"

	commonGRPCServer "gitlab.heronodes.io/bc-platform/bc-wallet-common-lib-grpc/pkg/server"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	grpcServer *grpc.Server
	handlers   pbApi.ApiServerHandlers

	listener net.Listener
}

func (s *Server) ListenAndServe(ctx context.Context) (err error) {
	options := commonGRPCServer.DefaultServeOptions()
	msgSizeOptions := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(commonGRPCServer.DefaultServerMaxReceiveMessageSize),
		grpc.MaxSendMsgSize(commonGRPCServer.DefaultServerMaxSendMessageSize),
	}
	options = append(options, msgSizeOptions...)
	options = append(options, grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer())))

	s.grpcServer = grpc.NewServer(options...)
	reflection.Register(s.grpcServer)
	pbApi.RegisterApiServerHandlers(s.grpcServer, s.handlers)

	return s.grpcServer.Serve(s.listener)
}

func main() {
	...

	listenConn, err := net.Listen("tcp", appCfg.GetBindPort())
	if err != nil {
		panic(err)
	}

	apiHandlers, err := grpcHandlers.New(ctx, dependecyService)
	if err != nil {
		panic(err)
	}

	srv := &Server{
		handlers: apiHandlers,
		listener: listenConn,
	}
	
	go func() {
		err = srv.ListenAndServe(ctx)
		if err != nil {
			panic(err)
		}
	}()
	
	...

}
```

## Licence

**bc-wallet-common-lib-grpc** has a proprietary license.

Switched to proprietary license from MIT - [CHANGELOG.MD - v0.0.3](./CHANGELOG.md)