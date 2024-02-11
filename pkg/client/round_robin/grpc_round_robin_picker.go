package round_robin

import (
	"crypto/rand"
	"log"
	"math/big"
	"sync/atomic"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

// Name is the name of round_robin balancer.
const Name = "round_robin_crypto_bundle"

// newBuilder creates a new roundrobin balancer builder.
func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, &rrPickerBuilder{}, base.Config{HealthCheck: true})
}

// nolint:gochecknoinits // ok. It is just copy of origin func
func init() {
	balancer.Register(newBuilder())
}

type rrPickerBuilder struct {
}

func (*rrPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	scs := make([]balancer.SubConn, 0, len(info.ReadySCs))
	scsInfo := make([]base.SubConnInfo, 0, len(info.ReadySCs))
	for sc := range info.ReadySCs {
		scs = append(scs, sc)
		scsInfo = append(scsInfo, info.ReadySCs[sc])
	}

	randBigInt := big.NewInt(int64(len(scs)))
	next, err := rand.Int(rand.Reader, randBigInt)
	if err != nil {
		next = big.NewInt(0)
	}

	return &rrPicker{
		subConns:       scs,
		subConnsInfo:   scsInfo,
		subsConnsCount: uint64(len(scs)),
		// Start at a random index, as the same RR balancer rebuilds a new
		// picker when SubConn states change, and we don't want to apply excess
		// load to the first server in the list.
		next: next.Uint64(),
	}
}

type rrPicker struct {
	// subConns is the snapshot of the roundrobin balancer when this picker was
	// created. The slice is immutable. Each Get() will do a round robin
	// selection from it and return the selected SubConn.
	subConns       []balancer.SubConn
	subConnsInfo   []base.SubConnInfo
	subsConnsCount uint64

	next uint64

	logger *log.Logger
}

func (p *rrPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	next := atomic.AddUint64(&p.next, 1)

	return balancer.PickResult{
		SubConn: p.subConns[next%p.subsConnsCount],
	}, nil
}
