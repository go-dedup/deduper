package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/mauidude/deduper/minhash"
)

type config struct {
	path     string
	host     string
	port     int
	leader   string
	debug    bool
	bands    int
	rows     int
	shingles int
}

var cfg *config

func init() {
	cfg = &config{}

	flag.StringVar(&cfg.leader, "leader", "", "The HTTP host and port of the leader")
	flag.BoolVar(&cfg.debug, "debug", false, "Enable debug logging")
	flag.IntVar(&cfg.bands, "bands", 100, "Number of bands")
	flag.IntVar(&cfg.rows, "hashes", 2, "Number of hashes to use")
	flag.IntVar(&cfg.shingles, "shingles", 2, "Number of shingles")
}

func main() {
	log.SetFlags(0)
	flag.Parse()

	if cfg.debug {
	}


	rand.Seed(time.Now().UnixNano())

	// Set the data directory.
	if flag.NArg() == 0 {
		flag.Usage()
		log.Fatal("Data path argument required")
	}

	path := flag.Arg(0)
	if err := os.MkdirAll(path, 0744); err != nil {
		log.Fatalf("Unable to create path: %v", err)
	}

	log.SetFlags(log.LstdFlags)

	mh := minhash.New(cfg.bands, cfg.rows, cfg.shingles)

	log.Fatal(s.ListenAndServe(cfg.leader))
}
