package client

import (
	"context"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	originGRPC "google.golang.org/grpc"
)

const connectTimeout = 10 * time.Second

type LazyConnection struct {
	connection *originGRPC.ClientConn

	address string
	options []originGRPC.DialOption
}

func NewLazyConnection(address string, additionalOptions ...originGRPC.DialOption) *LazyConnection {
	options := DefaultDialOptions()
	options = append(
		options,
		originGRPC.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	options = append(options, additionalOptions...)

	return &LazyConnection{
		address: address,
		options: options,
	}
}

func (d *LazyConnection) connect(ctx context.Context) error {
	ctx, finish := context.WithTimeout(ctx, connectTimeout)
	defer finish()

	conn, err := originGRPC.DialContext(ctx, d.address, d.options...)
	if err != nil {
		return err
	}
	d.connection = conn
	return nil
}

func (d *LazyConnection) Invoke(ctx context.Context, method string, args, reply interface{},
	opts ...originGRPC.CallOption) error {
	if d.connection == nil {
		err := d.connect(ctx)
		if err != nil {
			return err
		}
	}
	return d.connection.Invoke(ctx, method, args, reply, opts...)
}

func (d *LazyConnection) NewStream(ctx context.Context, desc *originGRPC.StreamDesc, method string,
	opts ...originGRPC.CallOption) (originGRPC.ClientStream, error) {
	if d.connection == nil {
		err := d.connect(ctx)
		if err != nil {
			return nil, err
		}
	}
	return d.connection.NewStream(ctx, desc, method, opts...)
}

func (d *LazyConnection) Close() {
	if d.connection != nil {
		_ = d.connection.Close()
	}
}
