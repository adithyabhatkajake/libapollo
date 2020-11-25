package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"

	"github.com/adithyabhatkajake/libchatter/log"
)

var (
	cpuprofile  = flag.String("cpuprofile", "", "write cpu profile to `file`")
	memprofile  = flag.String("memprofile", "", "write memory profile to `file`")
	logLevelPtr = flag.Uint64("loglevel", uint64(log.InfoLevel),
		"Loglevels are one of \n0 - PanicLevel\n1 - FatalLevel\n2 - ErrorLevel\n3 - WarnLevel\n4 - InfoLevel\n5 - DebugLevel\n6 - TraceLevel")
	configFileStrPtr = flag.String("conf", "", "Path to config file")
	numCPU           = flag.Int("cpu", 2, "number of cpu cores to use for the node")
)

func ParseOptions() {
	// Parse flags
	flag.Parse()

	// Set number of cores
	runtime.GOMAXPROCS(*numCPU)

	logLevel := log.InfoLevel

	switch uint32(*logLevelPtr) {
	case 0:
		logLevel = log.PanicLevel
	case 1:
		logLevel = log.FatalLevel
	case 2:
		logLevel = log.ErrorLevel
	case 3:
		logLevel = log.WarnLevel
	case 4:
		logLevel = log.InfoLevel
	case 5:
		logLevel = log.DebugLevel
	case 6:
		logLevel = log.TraceLevel
	}

	// Log Settings
	log.SetLevel(logLevel)

	if *cpuprofile != "" {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		go func() {
			_ = <-sigs
			pprof.StopCPUProfile()
			f.Close() // error handling omitted for example
			os.Exit(0)
		}()
	}
}
