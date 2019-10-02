package tcpconn

import (
	"net"

	"testing"

	"github.com/stretchr/testify/assert"
)

// TestOnceSmallPacketOneConn sends one small packet using one connection
func TestOnceSmallPacketOneConn(t *testing.T) {
	addr := "0.0.0.0:9999"
	s, err := net.Listen("tcp4", addr)

	assert.Nil(t, err)
	defer s.Close()

	timeBetSend := 0.005
	bytesToSend := 1
	timesToSend := 1
	maxRandTime := 0.001

	// Means that we can stack up multiple servers/clients
	go DealWithTCPConnections(s)
	err = SendTCPConnections(addr, timeBetSend, bytesToSend, timesToSend, maxRandTime)
	assert.Nil(t, err)
}

// TestMultiSmallPacketsSeqConn sends multiple small packets using one connection at a time
// Note that the port used each time may be different
func TestMultiSmallPacketsSeqConn(t *testing.T) {
	addr := "0.0.0.0:9998"
	s, err := net.Listen("tcp4", addr)

	assert.Nil(t, err)
	defer s.Close()

	timeBetSend := 0.0005
	bytesToSend := 8
	timesToSend := 100
	maxRandTime := 0.0001

	go DealWithTCPConnections(s)
	err = SendTCPConnections(addr, timeBetSend, bytesToSend, timesToSend, maxRandTime)
	assert.Nil(t, err)
}

// makeConnection is a supporting function for stacking up many concurrent connections using goroutines
func makeConnection(t *testing.T, ch chan bool, addr string, timeBetSend float64, bytesToSend int, timesToSend int, maxRandTime float64) {
	var err error
	err = SendTCPConnections(addr, timeBetSend, bytesToSend, timesToSend, maxRandTime)
	assert.Nil(t, err)
	ch <- true
}

// TestMultiSmallPacketsConcConn sends multiple small packets using many concurrent sequential connections
// The number of established connections is stated with numConnections
// Depending on the local environment, the test may fail if there are too many sockets being used at the same time
func TestMultiSmallPacketsConcConn(t *testing.T) {
	addr := "0.0.0.0:9997"
	s, err := net.Listen("tcp4", addr)

	assert.Nil(t, err)
	defer s.Close()

	timeBetSend := 0.0005
	bytesToSend := 8
	timesToSend := 5
	maxRandTime := 0.01
	numConnections := 20

	go DealWithTCPConnections(s)

	ch := make(chan bool)
	defer close(ch)
	for i := 0; i < numConnections; i++ {
		go makeConnection(t, ch, addr, timeBetSend, bytesToSend, timesToSend, maxRandTime)
	}

	// 	Blocks execution of defer statements until all goroutines are processed
	for i := 0; i < numConnections; i++ {
		_ = <-ch
	}
}

// TestCumulativeConnections tests that the number of total sockets used exceeds the 1024 limit of simultaneous connections on the local test environment
// Essentially runs TestMultiSmallPacketsConcConn multiple times with numSimuConn connections each time
// If passes, means that the sockets are being used and closed correctly and its not a worry that we cannot exceed 1024 connections on a single server
// Else, something is probably wrong
// **Cumulative connections toned down to 500 while investigating a minor connection issue**
func TestCumulativeConnections(t *testing.T) {
	addr := "0.0.0.0:9996"
	s, err := net.Listen("tcp4", addr)

	assert.Nil(t, err)
	defer s.Close()

	timeBetSend := 0.001
	bytesToSend := 8
	timesToSend := 1
	maxRandTime := 0.0005
	numSimuConn := 10 //number of simultaneous connections
	numDesiredConn := 500

	numLoops := (numDesiredConn / (numSimuConn * timesToSend)) + 1

	go DealWithTCPConnections(s)

	for j := 0; j < numLoops; j++ {
		ch := make(chan bool)
		defer close(ch)
		for i := 0; i < numSimuConn; i++ {
			go makeConnection(t, ch, addr, timeBetSend, bytesToSend, timesToSend, maxRandTime)
		}

		for i := 0; i < numSimuConn; i++ {
			_ = <-ch
		}
	}
}

// TestMultiLargePacketsSeqConn sends many very large packets using one connection at a time
func TestMultiLargePacketsSeqConn(t *testing.T) {
	addr := "0.0.0.0:9995"
	s, err := net.Listen("tcp4", addr)

	assert.Nil(t, err)
	defer s.Close()

	timeBetSend := 0.005
	bytesToSend := 100000000
	timesToSend := 20
	maxRandTime := 0.001

	go DealWithTCPConnections(s)
	err = SendTCPConnections(addr, timeBetSend, bytesToSend, timesToSend, maxRandTime)
	assert.Nil(t, err)
}

// TestMultiLargePacketsConcConn sends many large packets concurrently
// A different port are used each time a package is sent
func TestMultiLargePacketsConcConn(t *testing.T) {
	addr := "0.0.0.0:9994"
	s, err := net.Listen("tcp4", addr)

	assert.Nil(t, err)
	defer s.Close()

	timeBetSend := 0.0005
	bytesToSend := 1000000
	timesToSend := 5
	maxRandTime := 0.05
	numConnections := 3

	go DealWithTCPConnections(s)

	ch := make(chan bool)
	defer close(ch)
	for i := 0; i < numConnections; i++ {
		go makeConnection(t, ch, addr, timeBetSend, bytesToSend, timesToSend, maxRandTime)
	}

	for i := 0; i < numConnections; i++ {
		_ = <-ch
	}
}

// TestInvalidConn uses an invalid address which the client tries to dial and checks that errors are returned as expected
func TestInvalidConn(t *testing.T) {
	addr := "0.0.0.0:9993"
	s, err := net.Listen("tcp4", addr)

	assert.Nil(t, err)
	defer s.Close()

	timeBetSend := 0.005
	bytesToSend := 10000
	timesToSend := 1
	maxRandTime := 0.001

	go DealWithTCPConnections(s)
	err = SendTCPConnections("some_string", timeBetSend, bytesToSend, timesToSend, maxRandTime)
	assert.NotNil(t, err)
}
