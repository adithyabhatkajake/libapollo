package main

import (
	"flag"

	"github.com/adithyabhatkajake/libchatter/log"
)

var (
	// ConfFile is the config file directory
	ConfFile = flag.String("conf", "", "Path to client config file")
	// Batch the number of commands in a block
	Batch = flag.Uint64("batch", 10, "Number of commands to wait for")
	// Payload is the number of blocks in a command
	Payload     = flag.Uint64("payload", 2, "Number of bytes to get as response")
	Count       = flag.Uint64("metric", metricCount, "Number of metrics to collect before exiting")
	LogLevelPtr = flag.Uint64("loglevel", uint64(log.InfoLevel),
		"Loglevels are one of \n0 - PanicLevel\n1 - FatalLevel\n2 - ErrorLevel\n3 - WarnLevel\n4 - InfoLevel\n5 - DebugLevel\n6 - TraceLevel")
)

// ParseOptions parses options from the command line
func ParseOptions() {
	flag.Parse()

	// Setup Logger
	switch uint32(*LogLevelPtr) {
	case 0:
		log.SetLevel(log.PanicLevel)
	case 1:
		log.SetLevel(log.FatalLevel)
	case 2:
		log.SetLevel(log.ErrorLevel)
	case 3:
		log.SetLevel(log.WarnLevel)
	case 4:
		log.SetLevel(log.InfoLevel)
	case 5:
		log.SetLevel(log.DebugLevel)
	case 6:
		log.SetLevel(log.TraceLevel)
	}

	BufferCommands = *Batch
	metricCount = *Count
}
