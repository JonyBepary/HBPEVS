package main

import (
	"flag"

	"github.com/sohelahmedjoni/pqhbevs_hac/internal/network"
)

// Flag ------------------------------------------------------------------------------------------

func parseFlags() *network.Config {
	c := &network.Config{}

	flag.StringVar(&c.RendezvousString, "rendezvous", "meetme", "Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	flag.StringVar(&c.ListenHost, "host", "0.0.0.0", "The bootstrap node host listen address\n")
	flag.StringVar(&c.ProtocolID, "pid", "/block/blockreq/0.0.1", "Sets a protocol id for stream headers")
	flag.IntVar(&c.ListenPort, "port", 4001, "node listen port")
	flag.IntVar(&c.ApiPort, "api", 5669, "node listen port")

	flag.Parse()
	return c
}
