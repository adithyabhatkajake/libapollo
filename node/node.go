package main

import (
	"os"

	config "github.com/adithyabhatkajake/libapollo/config"
	"github.com/adithyabhatkajake/libapollo/consensus"
	"github.com/adithyabhatkajake/libchatter/io"
	"github.com/adithyabhatkajake/libchatter/log"
	"github.com/adithyabhatkajake/libchatter/net"
)

func main() {
	ParseOptions()

	log.Info("I am the replica.")

	Config := &config.NodeConfig{}
	io.ReadFromFile(Config, *configFileStrPtr)

	log.Debug("Finished reading the config file", os.Args[1])

	// Setup connections
	netw := net.Setup(Config, Config, Config)

	net.RetryWaitDuration = "10ms"
	net.RetryLimit = 300
	// Connect and send a test message
	netw.Connect()
	log.Debug("Finished connection to all the nodes")

	// Configure E2C protocol
	apl := &consensus.Apollo{}
	apl.Init(Config)
	apl.Setup(netw)

	// Start E2C
	apl.Start()

	// Disconnect
	netw.ShutDown()
}
