package main

import (
	"flag"
	"net"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/nlsun/ttracker/format"
	"github.com/nlsun/ttracker/hook"
	"github.com/nlsun/ttracker/output"
	"github.com/nlsun/ttracker/tracker"
	"github.com/nlsun/ttracker/util"
)

var spec string
var prefix string
var verbose bool
var pollRateStr string
var bindAddr string

func main() {
	flag.StringVar(&spec, "spec", "", "name1,pid1;name2,pid2;...")
	flag.StringVar(&prefix, "prefix", "", "The prefix of the log file name")
	flag.BoolVar(&verbose, "verbose", false, "Set verbose")
	flag.StringVar(&pollRateStr, "rate", "2s", "Polling rate")
	flag.StringVar(&bindAddr, "addr", "127.0.0.1:0", "Bind Address")
	flag.Parse()

	unixstamp := util.Unixtime

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
	targets, err := tracker.ParseSpec(spec)
	if err != nil {
		log.Fatal(err)
	}
	lsnr, err := net.Listen("tcp", bindAddr)
	if err != nil {
		log.Fatal(err)
	}
	hookServer := hook.NewServer(lsnr)
	log.WithField("targets", targets).Debug()
	if prefix == "" {
		prefix = unixstamp()
	}
	logger, err := output.InitLogger(prefix, unixstamp)
	if err != nil {
		log.Fatal(err)
	}
	logger.LogAddr(hookServer.Addr())
	logger.LogStats(format.Title(targets))
	log.Info("Tracker Initialized")

	go func(s hook.Server, l output.Logger) {
		for {
			log.WithError(s.Serve(l)).Info("hook server crashed")
		}
	}(hookServer, logger)

	ticker := time.NewTicker(pollRate)
	res := make(map[tracker.Target]string)
	for _ = range ticker.C {
		stamp := unixstamp()
		for _, tgt := range targets {
			res[tgt] = tracker.RunOnce(tgt)
		}
		err = logger.LogStats(format.Line(res, targets, stamp))
		if err != nil {
			log.WithError(err).Debug()
		}
	}
}
