// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package listeners

import (
	"fmt"
	"net"
	"strings"

	"github.com/DataDog/datadog-agent/comp/dogstatsd/packets"
	"github.com/DataDog/datadog-agent/comp/dogstatsd/replay"
	"github.com/DataDog/datadog-agent/pkg/config"
	"github.com/DataDog/datadog-agent/pkg/util/log"
)

// UDSStreamListener implements the StatsdListener interface for Unix Domain (streams)
type UDSStreamListener struct {
	UDSListener

	conn *net.UnixListener
}

// NewUDSStreamListener returns an idle UDS datagram Statsd listener
func NewUDSStreamListener(packetOut chan packets.Packets, sharedPacketPoolManager *packets.PoolManager, sharedOobPacketPoolManager *packets.PoolManager, cfg config.Reader, capture replay.Component) (*UDSStreamListener, error) {
	socketPath := cfg.GetString("dogstatsd_stream_socket")
	transport := "unix"

	address, err := setupSocketBeforeListen(socketPath, transport)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUnix(transport, address)
	if err != nil {
		return nil, fmt.Errorf("can't listen: %s", err)
	}

	err = setSocketWriteOnly(socketPath)
	if err != nil {
		return nil, err
	}

	l, err := NewUDSListener(packetOut, sharedPacketPoolManager, sharedOobPacketPoolManager, cfg, capture, transport)
	if err != nil {
		return nil, err
	}

	listener := &UDSStreamListener{
		UDSListener: *l,
		conn:        conn,
	}

	log.Infof("dogstatsd-uds-stream: %s successfully initialized", conn.Addr())
	return listener, nil
}

// Listen runs the intake loop. Should be called in its own goroutine
func (l *UDSStreamListener) Listen() {
	log.Infof("dogstatsd-uds-stream: starting to listen on %s", l.conn.Addr())
	for {
		conn, err := l.conn.AcceptUnix()
		if err != nil {
			if !strings.HasSuffix(err.Error(), " use of closed network connection") {
				log.Errorf("dogstatsd-uds: error accepting connection: %v", err)
			}
			break
		}
		go func() {
			_ = l.handleConnection(conn)
			if err != nil {
				log.Errorf("dogstatsd-uds-stream: error handling connection: %v", err)
			}
		}()
	}
}

// Stop closes the UDS connection and stops listening
func (l *UDSStreamListener) Stop() {
	_ = l.conn.Close()
	l.UDSListener.Stop()
}