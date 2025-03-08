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
	"math"
	"time"

	grpcRetry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	originGRPC "google.golang.org/grpc"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	grpcKeepalive "google.golang.org/grpc/keepalive"
)

const (
	DefaultClientMaxReceiveMessageSize = 1024 * 1024 * 24
	DefaultClientMaxSendMessageSize    = math.MaxInt32

	DefaultClientKeepAliveTime    = 2 * time.Minute
	DefaultClientKeepAliveTimeout = 1*time.Minute + 9*time.Second

	DefaultClientConnectionMaxRetry = 3
	DefaultClientConnectionBackOff  = 450 * time.Millisecond
)

func DefaultKeepaliveClientOptions() grpcKeepalive.ClientParameters {
	return grpcKeepalive.ClientParameters{
		Time:                DefaultClientKeepAliveTime,
		Timeout:             DefaultClientKeepAliveTimeout,
		PermitWithoutStream: true,
	}
}

func DefaultRetryOptions() []grpcRetry.CallOption {
	return []grpcRetry.CallOption{
		grpcRetry.WithMax(DefaultClientConnectionMaxRetry),
		grpcRetry.WithBackoff(grpcRetry.BackoffLinear(DefaultClientConnectionBackOff)),
		grpcRetry.WithCodes(grpcCodes.Aborted, grpcCodes.Unavailable),
	}
}

func DefaultInterceptorsOptions() []originGRPC.UnaryClientInterceptor {
	return []originGRPC.UnaryClientInterceptor{
		grpcRetry.UnaryClientInterceptor(DefaultRetryOptions()...),
	}
}

func DefaultDialOptions() []originGRPC.DialOption {
	return []originGRPC.DialOption{
		originGRPC.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithContextDialer(Dialer), // use it if u need load balancing via dns
		originGRPC.WithBlock(),
		originGRPC.WithKeepaliveParams(DefaultKeepaliveClientOptions()),
		originGRPC.WithChainUnaryInterceptor(DefaultInterceptorsOptions()...),
	}
}
