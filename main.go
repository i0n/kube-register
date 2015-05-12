package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	apiEndpoint   string
	fleetEndpoint string
	metadata      string
	reverseLookup bool
	syncInterval  int
	healthzPort   string
	printVersion  bool
)

func init() {
	log.SetFlags(0)
	flag.StringVar(&apiEndpoint, "api-endpoint", "", "kubernetes API endpoint")
	flag.StringVar(&fleetEndpoint, "fleet-endpoint", "", "fleet endpoint")
	flag.StringVar(&metadata, "metadata", "k8s=kubelet", "comma-delimited key/value pairs")
	flag.BoolVar(&reverseLookup, "reverse-lookup", false, "execute reverse lookup for registering hostnames instead of hosts' public IPs")
	flag.StringVar(&healthzPort, "healthz-port", "10255", "the kubelet healthz port")
	flag.IntVar(&syncInterval, "sync-interval", 30, "sync interval")
	flag.BoolVar(&printVersion, "version", false, "print version and exit")
}

func main() {
	flag.Parse()
	if printVersion {
		fmt.Printf("kube-register %s\n", Version)
		os.Exit(0)
	}

	m, err := parseMetadata(metadata)
	if err != nil {
		log.Println(err)
	}
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		machines, err := getMachines(fleetEndpoint, healthzPort, m, reverseLookup)
		if err != nil {
			log.Println(err)
		}
		for _, machine := range machines {
			if err := register(apiEndpoint, machine); err != nil {
				log.Println(err)
			}
		}
		select {
		case c := <-signalChan:
			log.Println(fmt.Sprintf("captured %v exiting...", c))
			os.Exit(0)
		case <-time.After(time.Duration(syncInterval) * time.Second):
			// Continue syncing machines.
		}
	}
}
