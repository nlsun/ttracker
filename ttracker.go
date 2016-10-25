package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/nlsun/ttracker/format"
	"github.com/nlsun/ttracker/output"
	"github.com/nlsun/ttracker/tracker"
)

// XXX have an endpoint that will trigger a timestamp to be written to file
// XXX have an endpoint that will cause this to flush all writes to disk

var spec string
var prefix string
var verbose bool
var pollRateStr string

func main() {
	flag.StringVar(&spec, "spec", "", "name1,pid1;name2,pid2;...")
	flag.StringVar(&prefix, "prefix", "", "The prefix of the log file name")
	flag.BoolVar(&verbose, "verbose", false, "Set verbose")
	flag.StringVar(&pollRateStr, "rate", "2s", "Polling rate")
	flag.Parse()

	if verbose {
		log.SetLevel(log.DebugLevel)
	}
	if spec == "" {
		log.Fatal("Must provide a spec")
	}
	pollRate, err := time.ParseDuration(pollRateStr)
	if err != nil {
		log.Fatal(err)
	}
	if prefix == "" {
		prefix = fmt.Sprintf("%d", unixtime())
	}
	targets, err := tracker.ParseSpec(spec)
	if err != nil {
		log.Fatal(err)
	}
	log.WithField("targets", targets).Debug()

	logger, err := output.InitLogger(prefix)
	if err != nil {
		log.Fatal(err)
	}
	logger.LogStats(format.Title(targets))
	ticker := time.NewTicker(pollRate)
	res := make(map[tracker.Target]string)
	for _ = range ticker.C {
		timestamp := strconv.Itoa(int(unixtime()))
		for _, tgt := range targets {
			res[tgt] = tracker.RunOnce(tgt)
		}
		err = logger.LogStats(format.Line(res, targets, timestamp))
		if err != nil {
			log.WithError(err).Debug()
		}
	}
}

func unixtime() int32 {
	return int32(time.Now().Unix())
}
