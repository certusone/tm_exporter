package main

import (
	"bytes"
	"encoding/hex"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tendermint/tendermint/rpc/client"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	valAddrHex = os.Getenv("ADDRESS")
	rpcHost    = os.Getenv("RPC")
	lAddr      = os.Getenv("LADDR")
)

func main() {
	go pollRoutine()

	http.Handle("/metrics", promhttp.Handler())

	panic(http.ListenAndServe(lAddr, nil))
}

func pollRoutine() {
	c := client.NewHTTP(rpcHost, "/ws")
	missCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   "tm",
		Subsystem:   "mon",
		Name:        "misses_count",
		Help:        "Amount of missed blocks since exporter start",
		ConstLabels: nil,
	})
	prometheus.MustRegister(missCounter)

	heightGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   "tm",
		Subsystem:   "mon",
		Name:        "height",
		Help:        "Current block height that the monitor has processed",
		ConstLabels: nil,
	})
	prometheus.MustRegister(heightGauge)

	validatorAddress, err := hex.DecodeString(valAddrHex)
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(time.Second)
	var lastHeight int64
OUTER:
	for {
		select {
		case <-ticker.C:
			res, err := c.Block(nil)
			if err != nil {
				log.Printf("block could not be requested; err=%v\n", err)
				continue
			}

			if res.Block.Height > lastHeight {
				lastHeight = res.Block.Height
				heightGauge.Set(float64(lastHeight))

				for _, pre := range res.Block.LastCommit.Precommits {
					if pre != nil && bytes.Equal(validatorAddress, pre.ValidatorAddress) && !pre.BlockID.IsZero() {
						// We signed and the vote is not zero (nil votes also count as downtime)
						continue OUTER
					}
				}

				// Seems like we missed the block
				missCounter.Inc()
			}
		}
	}
}
