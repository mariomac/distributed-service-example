package server

import (
	context "context"
	"math/big"

	"github.com/mariomac/distributed-service-example/worker/pkg/gprc"
)

var one = (&big.Int{}).SetInt64(1)

type MultiplyServer struct {
	gprc.MultiplierServer
}

func (m *MultiplyServer) Loop(_ context.Context, request *gprc.LoopRequest) (*gprc.LoopResponse, error) {
	start := &big.Int{}
	start.SetBytes(request.From)
	result := &*start
	end := &big.Int{}
	end.SetBytes(request.To)
	for start.Cmp(end) < 0 {
		start.Add(start, one)
		result.Mul(result, start)
	}
	return &gprc.LoopResponse{Result: result.Bytes()}, nil
}

func (m *MultiplyServer) mustEmbedUnimplementedMultiplierServer() {}
