package tcpconn

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	"github.com/brucespang/go-tcpinfo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

// Set up simple server side metrics to be exported
var (
	// Total number of connections handled by the server
	ConnsHandledTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "conntest_connections_handled_total",
		},
	)
)

// Set up socket level statistics as metrics
var (
	RetransmitsCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "conntest_tcp_retransmits_counter",
		},
		[]string{
			// IP address of the target we are sending tests to
			"dst_ip",
			// IP address of the source the tests are being sent from
			"src_ip",
			// Name of current node
			"node_name",
		},
	)

	SndMssGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "conntest_tcp_send_message_gauge",
		},
		[]string{
			"dst_ip",
			"src_ip",
			"node_name",
		},
	)

	RcvMssGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "conntest_tcp_receive_message_gauge",
		},
		[]string{
			"dst_ip",
			"src_ip",
			"node_name",
		},
	)

	LostPacketsCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "conntest_tcp_lost_packets_counter",
		},
		[]string{
			"dst_ip",
			"src_ip",
			"node_name",
		},
	)

	RetransCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "conntest_tcp_retrans_counter",
		},
		[]string{
			"dst_ip",
			"src_ip",
			"node_name",
		},
	)

	PmtuGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "conntest_tcp_pmtu_gauge",
		},
		[]string{
			"dst_ip",
			"src_ip",
			"node_name",
		},
	)

	RttGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "conntest_tcp_round_trip_time_seconds_gauge",
		},
		[]string{
			"dst_ip",
			"src_ip",
			"node_name",
		},
	)

	RttHistVec = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "conntest_tcp_round_trip_time_seconds_hist",
			Buckets: prometheus.ExponentialBuckets(1e-9, 10, 10),
		},
		[]string{
			"dst_ip",
			"src_ip",
			"node_name",
		},
	)

	RttVarGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "conntest_tcp_round_trip_time_variance_gauge",
		},
		[]string{
			"dst_ip",
			"src_ip",
			"node_name",
		},
	)

	TotalRetransGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "conntest_tcp_total_retrans_gauge",
		},
		[]string{
			"dst_ip",
			"src_ip",
			"node_name",
		},
	)
)

// HandleTCPConnection deals with our TCP based protocol, closes the connection once it finishes serving the client
func HandleTCPConnection(c net.Conn) error {
	log.Debug("Serving ", c.RemoteAddr().String())
	var err error
	defer log.Debug("Finished serving ", c.RemoteAddr().String())
	for {
		data, err := ReceiveViaProtocol(c)
		if (err != nil) && (err != io.EOF) {
			log.Error(err)
			break
		}

		temp := strings.TrimSpace(string(data))
		if temp == "EOS" {
			break
		}
	}
	return err
}

// DealWithTCPConnections ensures we can deal with multiple clients without blocking
func DealWithTCPConnections(s net.Listener) error {
	for {
		c, err := s.Accept()
		if err != nil {
			log.Error("Error accepting connection: ", err)
			return err
		}
		go HandleTCPConnection(c)
	}
}

// SendTCPConnection sends bytesToSend bytes to destHost
func SendTCPConnection(destHost string, bytesToSend int) error {
	c, err := net.Dial("tcp", destHost)
	if err != nil {
		return err
	}
	log.Debug("Local addr: ", c.LocalAddr())
	defer log.Debug("Client finished sending to ", destHost)
	defer c.Close()

	// Query for socket info
	socketInfo, err := tcpinfo.GetsockoptTCPInfo(&c)
	if err != nil {
		errMsg := "Error while attempting to fetch TCP info: " + err.Error()
		err = errors.New(errMsg)
		return err
	}

	// Get local IPs to be registered as a Prometheus label
	localHostName, err := os.Hostname()
	if err != nil {
		errMsg := "Error while attempting to look up local host name: " + err.Error()
		err = errors.New(errMsg)
		return err
	}
	localIPs, err := net.LookupHost(localHostName)
	if err != nil {
		errMsg := "Error while attempting to look up local IP address: " + err.Error()
		err = errors.New(errMsg)
		return err
	}
	var localIPsBuilder strings.Builder
	for _, IP := range localIPs {
		fmt.Fprintf(&localIPsBuilder, "%v, ", IP)
	}
	localIPsStr := localIPsBuilder.String()
	log.Debug("Discovered IPs: ", localIPsStr)

	// Look up node name
	nodeName, foundBool := os.LookupEnv("NODE_NAME")
	if foundBool == false {
		nodeName = "UNKNOWN"
		log.Warning("No node name was discovered!")
	} else {
		log.Debug("Discovered node name: ", nodeName)
	}

	// Register relevant socket info
	RetransmitsCounterVec.WithLabelValues(destHost, localIPsStr, nodeName).Add(float64(socketInfo.Retransmits))
	SndMssGaugeVec.WithLabelValues(destHost, localIPsStr, nodeName).Set(float64(socketInfo.Snd_mss))
	RcvMssGaugeVec.WithLabelValues(destHost, localIPsStr, nodeName).Set(float64(socketInfo.Rcv_mss))
	LostPacketsCounterVec.WithLabelValues(destHost, localIPsStr, nodeName).Add(float64(socketInfo.Lost))
	RetransCounterVec.WithLabelValues(destHost, localIPsStr, nodeName).Add(float64(socketInfo.Retrans))
	PmtuGaugeVec.WithLabelValues(destHost, localIPsStr, nodeName).Set(float64(socketInfo.Pmtu))
	RttGaugeVec.WithLabelValues(destHost, localIPsStr, nodeName).Set(float64(socketInfo.Rtt) / 1e9)
	RttHistVec.WithLabelValues(destHost, localIPsStr, nodeName).Observe(float64(socketInfo.Rtt) / 1e9)
	RttVarGaugeVec.WithLabelValues(destHost, localIPsStr, nodeName).Set(float64(socketInfo.Rttvar))
	TotalRetransGaugeVec.WithLabelValues(destHost, localIPsStr, nodeName).Set(float64(socketInfo.Total_retrans))

	strToSend := strings.Repeat("a", bytesToSend)
	err = SendViaProtocol(c, []byte(strToSend))
	if err != nil {
		return err
	}
	err = SendViaProtocol(c, []byte("EOS"))
	if err != nil {
		return err
	}
	return err
}

// SendViaProtocol sends data over connection c using our custom protocol
func SendViaProtocol(c net.Conn, data []byte) error {

	dataWNL := append(data, '\n')
	log.Debug("Client sent: ", string(dataWNL))

	// Writes to and receives from server
	_, err := c.Write(dataWNL)
	if err != nil {
		return err
	}
	netData, err := bufio.NewReader(c).ReadString('\n')
	log.Debug(":", netData, ":")
	tempNetdata := strings.TrimSpace(string(netData))
	if err != nil {
		return err
	}
	if tempNetdata == "ACK" {
		err = nil
	}
	log.Debug("Client finished sending to ", c.RemoteAddr())
	return err
}

// ReceiveViaProtocol runs on server with HandleTCPConnection to receive messages
// from clients via our custom protocol
func ReceiveViaProtocol(c net.Conn) (string, error) {
	log.Debug("Server receiving from ", c.RemoteAddr().String(), "\n")

	netData, err := bufio.NewReader(c).ReadString('\n')
	log.Debug("Server received: ", netData)

	if (err != nil) && (err != io.EOF) {
		// Down to debug level as we don't care whether the client stops sending
		log.Debug(err)
		return "", err
	}

	result := "ACK\n"
	c.Write([]byte(string(result)))

	ConnsHandledTotal.Inc()
	log.Debug("Server finished receiving from ", c.RemoteAddr().String(), "\n")

	switch strings.TrimSpace(netData) {
	case "":
		// If receives nothing, keep the connection alive for a little bit before killing it
		// otherwise will cause a "use of closed network connection error"
		// further investigation about this issue in the future is preferable
		// Caused by the readiness and liveness probes
		go func() {
			time.Sleep(time.Duration(5e9))
			log.Debug("Waiting for 5 seconds before closing connection...")
			c.Close()
		}()
	case "EOS":
		// Packet has finished sending correctly though the desired process, safe to murder the connection
		c.Close()
	default:
		// If receives data that is not the above cases (i.e. the test packet), do nothing,
		// keep the connection alive because the EOS that follows this will terminate it
	}

	return netData, nil
}

// SendTCPConnections repeatedly send messages of size bytesToSend to the server
// with specified time intervals plus a random amount up to a second
// Time interval counted in seconds
// **Functionality duplicated and improved in conntools, no longer in used in main**
func SendTCPConnections(destHost string, timeBetweenTestsS float64, bytesToSend int, timesToSend int, randSecs float64) error {
	var terr error

	runForever := false
	if timesToSend == 0 {
		runForever = true
	}
	i := 0

	for {
		log.Debug("Sending TCP test to ", destHost, "\n")
		err := SendTCPConnection(destHost, bytesToSend)
		if (err != nil) && (err != io.EOF) {
			log.Error(err)
			// We return the last error encountered
			terr = err
		}

		if !runForever {
			i++
			if i >= timesToSend {
				log.Debug("Finished sending ", timesToSend, " times to server")
				break
			}
		}

		// Only understands nanoseconds.......
		ti := int64(1e9 * (timeBetweenTestsS + (randSecs * rand.Float64())))
		log.Debug("Total time between tests: ", float64(ti)/float64(1e9), " seconds")
		time.Sleep(time.Duration(ti))
	}
	return terr
}
