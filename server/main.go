package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	port := flag.Int("port", -1, "The port number to use for this node")
	flag.Parse()

	if *port < 0 {
		fmt.Printf("Please provide a port number to bind to that is in the cluster file.\n")
		os.Exit(1)
	}

	_, err := parseClusterConfig(*port, "cluster.conf")
	if err != nil {
		fmt.Printf("Error reading cluster config: %v", err)
		os.Exit(1)
	}
}

func parseClusterConfig(myPort int, configPath string) (peers []int, err error) {
	bytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	for line := range strings.SplitSeq(string(bytes), "\n") {
		if line == "" {
			continue
		}
		peer, err := strconv.Atoi(line)
		if err != nil {
			fmt.Printf("Ignoring entry %s: could not parse to integer\n", line)
			continue
		}
		if peer != myPort {
			peers = append(peers, peer)
		}
	}
	return peers, nil
}
