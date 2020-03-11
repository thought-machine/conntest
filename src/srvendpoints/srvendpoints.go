package srvendpoints

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"github.com/thought-machine/conntest/src/tcpconn"
)

var log = logrus.New()

// Counts number of failed SRV discoveries
var (
	TotalFailedSRVCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "conntest_failed_SRV_discoveries_counter",
		},
	)
)

// DiscoverEndpoints uses SRV records to discover available endpoints
func DiscoverEndpoints(service, protocol, name string, retryIntervalSecs float64, maxRetries int, failed int) ([]string, error) {
	var err error
	_, srv, serr := net.LookupSRV(service, protocol, name)
	for serr != nil {
		failed++
		log.Debug("Failed SRV discovery attempts: ", failed)
		TotalFailedSRVCounter.Add(1)
		// Use maxRetries = -1 to retry indefinitely
		if maxRetries == -1 {
		} else if failed > maxRetries {
			err = fmt.Errorf("Attempt to discover SRV record timed out after %v retries", (failed - 1))
			return make([]string, 0), err
		}
		err = fmt.Errorf("Cannot find SRV record, retrying in %v seconds...(attempt %v)", retryIntervalSecs, failed)
		log.Error(err)
		// Only understands nanoseconds
		ti := int64(1e9 * retryIntervalSecs)
		log.Debug("Time between retries: ", float64(ti)/float64(1e9), " seconds")
		time.Sleep(time.Duration(ti))
		_, srv, serr = net.LookupSRV(service, protocol, name)
	}
	endpoints := make([]string, len(srv))
	for i := 0; i < len(srv); i++ {
		log.Debug("Available endpoints: ", srv[i].Target, ":", srv[i].Port)
		endpoints[i] = srv[i].Target + ":" + strconv.Itoa(int(srv[i].Port))
	}
	log.Debug("Discovered endpoints: ", endpoints)
	return endpoints, serr
}

// makethConnection is a supporting function for stacking up many concurrent connections using goroutines
func makethConnection(ch chan bool, endpoint string, nodeName string, testBytes int) {
	err := tcpconn.SendTCPConnection(endpoint, testBytes, nodeName)
	if err != nil && strings.TrimSpace(err.Error()) != "EOF" {
		log.Error(err)
	}
	ch <- true
	return
}

// SendConcTCPConnections sends packets using concurrent sequential connections
func SendConcTCPConnections(service string, protocol string, name string, nodeName string, testBytes int, retryIntervalSecs float64, maxRetries int) error {
	ch := make(chan bool)
	defer close(ch)
	defer log.Debug("Channel closed")
	endpoints, err := DiscoverEndpoints(service, protocol, name, retryIntervalSecs, maxRetries, 0)
	if err != nil {
		return err
	}
	for i := 0; i < len(endpoints); i++ {
		go makethConnection(ch, endpoints[i], nodeName, testBytes)
	}
	// blocks further execution until connections to all endpoints are completed
	for i := 0; i < len(endpoints); i++ {
		_ = <-ch
	}
	return err
}
