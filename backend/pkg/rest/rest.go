package rest

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/mariomac/distributed-service-example/worker/pkg/gprc"
	"google.golang.org/grpc"
)

const (
	FactorialPath = "/factorial/"
)

var one = big.NewInt(1)

func FactorialService(workerAddr string, workers int, timeout time.Duration) http.HandlerFunc {
	conn, err := grpc.Dial(workerAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("can't connect to worker: %s", err)
	}
	bigWorkers := big.NewInt(int64(workers))
	client := gprc.NewMultiplierClient(conn)
	return func(rw http.ResponseWriter, req *http.Request) {
		actualWorkers := big.NewInt(int64(workers))
		inputStr := req.URL.Path[len(FactorialPath):]
		input := &big.Int{}
		input, ok := input.SetString(inputStr, 0)
		if !ok {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(fmt.Sprintf("wrong input: %s\n", inputStr)))
			return
		}
		ctx, cancel := context.WithTimeout(req.Context(), timeout)
		defer cancel()
		var sliceLen *big.Int
		start := big.NewInt(1)
		if input.Cmp(actualWorkers) < 0 {
			log.Printf("%s < %d", input, actualWorkers)
			actualWorkers.SetInt64(1)
			sliceLen = (&big.Int{}).Set(input)
		} else {
			sliceLen = (&big.Int{}).Div(input, bigWorkers)
		}
		end := (&big.Int{}).Set(sliceLen)
		awn := int(actualWorkers.Int64())
		resCh := make(chan *gprc.LoopResponse, awn)
		errsCh := make(chan error, awn)
		for i := 0; i < awn-1; i++ {
			sstart, send := (&big.Int{}).Set(start), (&big.Int{}).Set(end)
			log.Printf("%s -> %s", sstart, send)
			go invokeWorker(ctx, client, sstart, send, errsCh, resCh)
			start.Set(end).Add(start, one)
			end.Add(end, sliceLen)
		}
		go invokeWorker(ctx, client, start, input, errsCh, resCh)

		result := big.NewInt(1)
		for i := 0; i < awn; i++ {
			select {
			case res := <-resCh:
				ires := (&big.Int{}).SetBytes(res.Result)
				log.Printf("worker %d returned %s", i, ires.String())
				result.Mul(result, ires)
			case err := <-errsCh:
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write([]byte(fmt.Sprintf("error calculating numbers: %s\n", err)))
				return
			case <-ctx.Done():
				rw.WriteHeader(http.StatusGatewayTimeout)
				return
			}
		}
		rw.Write([]byte(result.String()))
		rw.Write([]byte{'\n'})
	}
}

func invokeWorker(ctx context.Context, client gprc.MultiplierClient, start, end *big.Int, errsCh chan error, resCh chan *gprc.LoopResponse) {
	log.Printf("sending to worker: (%s, %s)", start.String(), end.String())
	lr, err := client.Loop(ctx, &gprc.LoopRequest{From: start.Bytes(), To: end.Bytes()})
	if err != nil {
		errsCh <- err
	}
	resCh <- lr
}
