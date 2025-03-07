/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2025 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

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
	e          errorFormatterService

	address string
	options []originGRPC.DialOption
}

func NewLazyConnection(errFmtSvc errorFormatterService,
	address string,
	additionalOptions ...originGRPC.DialOption,
) *LazyConnection {
	options := DefaultDialOptions()
	options = append(
		options,
		originGRPC.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	options = append(options, additionalOptions...)

	return &LazyConnection{
		e:          errFmtSvc,
		address:    address,
		options:    options,
		connection: nil,
	}
}

func (d *LazyConnection) connect(ctx context.Context) (*originGRPC.ClientConn, error) {
	ctx, finish := context.WithTimeout(ctx, connectTimeout)
	defer finish()

	conn, err := originGRPC.DialContext(ctx, d.address, d.options...)
	if err != nil {
		return nil, d.e.ErrorOnly(err)
	}

	return conn, nil
}

func (d *LazyConnection) Invoke(ctx context.Context,
	method string,
	args, reply interface{},
	opts ...originGRPC.CallOption,
) error {
	if d.connection != nil {
		return d.connection.Invoke(ctx, method, args, reply, opts...)
	}

	conn, err := d.connect(ctx)
	if err != nil {
		return d.e.ErrorOnly(err)
	}

	d.connection = conn

	return d.connection.Invoke(ctx, method, args, reply, opts...)
}

func (d *LazyConnection) NewStream(ctx context.Context,
	desc *originGRPC.StreamDesc,
	method string,
	opts ...originGRPC.CallOption,
) (originGRPC.ClientStream, error) {
	if d.connection != nil {
		return d.connection.NewStream(ctx, desc, method, opts...)
	}

	conn, err := d.connect(ctx)
	if err != nil {
		return nil, d.e.ErrorOnly(err)
	}

	d.connection = conn

	return d.connection.NewStream(ctx, desc, method, opts...)
}

func (d *LazyConnection) Close() {
	if d.connection != nil {
		_ = d.connection.Close()
	}
}
