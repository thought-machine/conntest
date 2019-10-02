package main

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"

	flags "github.com/jessevdk/go-flags"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"github.com/thought-machine/conntest/src/srvendpoints"
	"github.com/thought-machine/conntest/src/tcpconn"
)

var log = logrus.New()

var opts struct {
	HostPort         string  `long:"host_port" default:"8080" description:"Port to host on"`
	DestHost         string  `long:"dst_hst" default:"127.0.0.1:8080" description:"Destination host to target for tests"`
	TimeBetTests     float64 `long:"wait_time" default:"5" description:"Minimum time between individual tests"`
	RandTimeTest     float64 `long:"rand_secs" default:"5.0" description:"Maximum random time to be added to TimeBetTests"`
	ShortTestBytes   int     `long:"short_test_bytes" default:"10" description:"Bytes to use for short tests"`
	LongTestBytes    int     `long:"long_test_bytes" default:"10000" description:"Bytes to use for long tests"`
	TimesToSend      int     `long:"times_to_send" default:"0" description:"Number of times to send bytes"`
	DNSRetryInterval float64 `long:"DNS_retry_interval" default:"5.0" description:"Time between attempts to re-discover SRV records"`
	MaxDNSRetries    int     `long:"max_DNS_retries" default:"-1" description:"Maximum number of retries when attmpting to re-discover SRV records, use -1 for infinite retries"`
}

func init() {
	prometheus.MustRegister(tcpconn.ConnsHandledTotal)
	prometheus.MustRegister(tcpconn.RetransmitsCounterVec)
	prometheus.MustRegister(tcpconn.SndMssGaugeVec)
	prometheus.MustRegister(tcpconn.RcvMssGaugeVec)
	prometheus.MustRegister(tcpconn.LostPacketsCounterVec)
	prometheus.MustRegister(tcpconn.RetransCounterVec)
	prometheus.MustRegister(tcpconn.PmtuGaugeVec)
	prometheus.MustRegister(tcpconn.RttGaugeVec)
	prometheus.MustRegister(tcpconn.RttHistVec)
	prometheus.MustRegister(tcpconn.RttVarGaugeVec)
	prometheus.MustRegister(tcpconn.TotalRetransGaugeVec)
	prometheus.MustRegister(srvendpoints.TotalFailedSRVCounter)
}

func main() {

	_, err := flags.ParseArgs(&opts, os.Args)

	if err != nil {
		fmt.Printf("\n%s\n", err)
		os.Exit(1)
	}

	// Binding to all interfaces on IPv4 :'(
	addr := "0.0.0.0:" + opts.HostPort
	s, err := net.Listen("tcp4", addr)

	if err != nil {
		log.Fatal(err)
	}

	defer s.Close()

	// Means that we can stack up multiple servers/clients
	go tcpconn.DealWithTCPConnections(s)

	// Serves Prometheus metrics from default port :9990
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe("0.0.0.0:9990", nil)

	// Repeatedly send messages of specified size to the server
	// with time intervals plus a random amount up to that specified by opts.RandTimeTest
	for {
		err = srvendpoints.SendConcTCPConnections("tcp", "tcp", "conntest", opts.ShortTestBytes, opts.DNSRetryInterval, opts.MaxDNSRetries)
		if err != nil {
			log.Error(err)
		}
		// Only understands nanoseconds
		ti := int64(1e9 * (opts.TimeBetTests + (opts.RandTimeTest * rand.Float64())))
		log.Debug("Total time between tests: ", float64(ti)/float64(1e9), " seconds")
		time.Sleep(time.Duration(ti))
	}
}
