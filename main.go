package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"gopkg.in/suntong/deduper.v1/command"
	"gopkg.in/suntong/deduper.v1/minhash"
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
	threshold float64
}

type testCase struct {
	id    string
	value string
}

var cfg *config

func init() {
	cfg = &config{}

	flag.StringVar(&cfg.leader, "leader", "", "The HTTP host and port of the leader")
	flag.BoolVar(&cfg.debug, "debug", false, "Enable debug logging")
	flag.IntVar(&cfg.bands, "bands", 100, "Number of bands")
	flag.IntVar(&cfg.rows, "hashes", 2, "Number of hashes to use")
	flag.IntVar(&cfg.shingles, "shingles", 2, "Number of shingles")
	flag.Float64Var(&cfg.threshold, "threshold", 0.5, "Threshold")
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
		log.Fatal("Data string argument required")
	}

	testStr := flag.Arg(0)
	log.SetFlags(log.LstdFlags)

	mh := minhash.New(cfg.bands, cfg.rows, cfg.shingles)

	tests := []testCase{
		{"p1", "hello world foo baz bar zomg"},
		{"p2", "goodbye world foo qux bar zomg"},
		{"p3", "entirely unrelated"},
	}
	for _, tt := range tests {
		command.NewWriteCommand(tt.id, tt.value).Apply(mh)
	}

	matches := mh.FindSimilar(
		strings.NewReader(testStr), cfg.threshold)

	json.NewEncoder(os.Stdout).Encode(matches)
}

/*

$ deduper
Usage of deduper:
  -bands int
        Number of bands (default 100)
  -debug
        Enable debug logging
  -hashes int
        Number of hashes to use (default 2)
  -leader string
        The HTTP host and port of the leader
  -shingles int
        Number of shingles (default 2)
  -threshold float
        Threshold (default 0.5)
Data string argument required

$ deduper "hello world foo baz"
[{"id":"p1","similarity":1}]

$ deduper "world foo baz"
[]

$ deduper "unrelated"
[]

$ deduper "entirely unrelated"
[{"id":"p3","similarity":1}]

$ deduper "entire unrelate"
[]

$ deduper "unrelated" -threshold 0.00001
[]

$ deduper "foo qux bar zomg"
[]

$ deduper "foo qux bar zomg" -threshold 0.0000000001
[]

$ deduper "foo qux bar zomg world"
[]

$ deduper "foo qux bar zomg world goodbye"
[]

$ deduper "goodbye world foo qux bar zomg"
[{"id":"p2","similarity":1}]

*/
