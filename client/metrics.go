package main

import (
	"fmt"
	"os"
	"time"

	"github.com/adithyabhatkajake/libchatter/crypto"
)

var (
	commitTimeMetric = make(map[crypto.Hash]time.Duration)
	metricCount      = uint64(1)
)

func printMetrics() {
	printDuration, err := time.ParseDuration("60s")
	if err != nil {
		panic(err)
	}
	var count int64 = 0
	for i := uint64(0); i < metricCount; i++ {
		<-time.After(printDuration)
		condLock.RLock()
		num := 0
		for _, dur := range commitTimeMetric {
			num++
			count += dur.Milliseconds()
		}
		fmt.Println("Metric")
		fmt.Printf("%d cmds in %d milliseconds\n", num, count)
		fmt.Printf("Throughput: %f\n", float64(num)/60.0)
		fmt.Printf("Latency: %f\n", float64(count)/float64(num))
		condLock.RUnlock()
	}
	os.Exit(0)
}
