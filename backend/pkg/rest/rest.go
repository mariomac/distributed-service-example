package rest

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/mariomac/distributed-service-example/worker/pkg/gprc"
	"google.golang.org/grpc"
)

const (
	FactorialPath = "/factorial/"
)

var one = big.NewInt(1)

func FactorialService(workerAddr string, workers int, timeout time.Duration) http.HandlerFunc {
	conn, err := grpc.Dial(workerAddr)
	if err != nil {
		log.Fatalf("can't connect to worker: %s", err)
	}
	bigWorkers := big.NewInt(int64(workers))
	client := gprc.NewMultiplierClient(conn)
	return func(rw http.ResponseWriter, req *http.Request) {
		inputStr := req.URL.Path[len(FactorialPath):]
		input := &big.Int{}
		input, ok := input.SetString(inputStr, 0)
		if !ok {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(fmt.Sprintf("wrong input: %s", inputStr)))
			return
		}
		ctx, cancel := context.WithTimeout(req.Context(), timeout)
		defer cancel()
		resCh := make(chan *gprc.LoopResponse, workers)
		errsCh := make(chan error, workers)
		start := big.NewInt(1)
		sliceLen := (&big.Int{}).Div(input, bigWorkers)
		end := (&big.Int{}).Set(sliceLen)
		for i := 0; i < workers-1; i++ {
			lreq := &gprc.LoopRequest{From: start.Bytes(), To: end.Bytes()}
			go invokeWorker(ctx, client, lreq, errsCh, resCh)
			start.Set(end)
			end.Add(end, sliceLen)
		}
		go invokeWorker(ctx, client,
			&gprc.LoopRequest{From: end.Bytes(), To: input.Add(input, one).Bytes()},
			errsCh, resCh)

		result := big.NewInt(1)
		for i := 0; i < workers-1; i++ {
			select {
			case res := <-resCh:
				result.Mul(result, (&big.Int{}).SetBytes(res.Result))
			case err := <-errsCh:
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write([]byte(fmt.Sprintf("error calculating numbers: %s", err)))
				return
			case <-ctx.Done():
				rw.WriteHeader(http.StatusGatewayTimeout)
				return
			}
		}
		rw.Write([]byte(result.String()))
	}
}

func invokeWorker(ctx context.Context, client gprc.MultiplierClient, lreq *gprc.LoopRequest, errsCh chan error, resCh chan *gprc.LoopResponse) {
	lr, err := client.Loop(ctx, lreq)
	if err != nil {
		errsCh <- err
	}
	resCh <- lr
}
