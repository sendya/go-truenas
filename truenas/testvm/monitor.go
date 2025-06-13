package testvm

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type monitor struct {
	sock     string
	listener net.Listener
	writeCh  chan []byte
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

func newMonitor(t *testing.T) *monitor {
	t.Helper()
	sock := filepath.Join(t.TempDir(), "monitor.sock")
	l, err := net.Listen("unix", sock)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(t.Context())
	m := &monitor{
		sock:     sock,
		listener: l,
		writeCh:  make(chan []byte, 10),
		ctx:      ctx,
		cancel:   cancel,
	}

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		for {
			select {
			case <-m.ctx.Done():
				return
			default:
			}

			conn, err := m.listener.Accept()
			if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {
				return
			}
			if err != nil {
				t.Logf("failed to accept monitor socket connection: %v", err)
				continue
			}
			t.Logf("monitor socket connected: %s", conn.LocalAddr())

			// Send QMP capabilities negotiation immediately upon connection
			m.sendQMPCapabilities(conn)

			m.handleConnection(t, conn)
		}
	}()
	return m
}

func (m *monitor) handleConnection(t *testing.T, conn net.Conn) {
	connCtx, connCancel := context.WithCancel(m.ctx)

	// Reader
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		defer conn.Close()
		defer connCancel()

		reader := bufio.NewReader(conn)
		for {
			select {
			case <-connCtx.Done():
				return
			default:
			}

			line, err := reader.ReadBytes('\n')
			if err != nil {
				if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {
					return
				}
				t.Logf("failed to read from monitor socket: %v", err)
				return
			}

			// Parse QMP JSON response
			var qmpMsg map[string]any
			if err := json.Unmarshal(line, &qmpMsg); err != nil {
				t.Logf("failed to parse QMP message: %v", err)
				continue
			}

			// Log QMP messages (skip QMP capability negotiation noise)
			if _, hasReturn := qmpMsg["return"]; hasReturn {
				t.Logf("[qmp-response] %s", strings.TrimSpace(string(line)))
			} else if _, hasError := qmpMsg["error"]; hasError {
				t.Logf("[qmp-error] %s", strings.TrimSpace(string(line)))
			} else if _, hasEvent := qmpMsg["event"]; hasEvent {
				t.Logf("[qmp-event] %s", strings.TrimSpace(string(line)))
			}
		}
	}()

	// Writer
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		defer connCancel()

		for {
			select {
			case <-connCtx.Done():
				return
			case b, ok := <-m.writeCh:
				if !ok {
					return
				}
				t.Logf("monitor write: %s", string(b))
				if _, err := conn.Write(b); err != nil {
					t.Logf("write to monitor socket: %v", err)
					return
				}
			}
		}
	}()
}

func (m *monitor) Addr() string {
	return "unix:" + m.sock
}

func (m *monitor) Close() {
	m.cancel()             // Signal all goroutines to stop
	_ = m.listener.Close() // Close listener to unblock Accept()
	close(m.writeCh)       // Close write channel
	m.wg.Wait()            // Wait for all goroutines to finish
}

func (m *monitor) Write(s string) {
	m.writeCh <- []byte(s + "\n")
}

// WriteQMP sends a QMP command as JSON
func (m *monitor) WriteQMP(cmd map[string]any) {
	data, _ := json.Marshal(cmd)
	m.writeCh <- append(data, '\n')
}

// sendQMPCapabilities performs the initial QMP handshake on a connection
func (m *monitor) sendQMPCapabilities(conn net.Conn) {
	cmd := map[string]any{
		"execute": "qmp_capabilities",
	}
	data, _ := json.Marshal(cmd)
	data = append(data, '\n')
	_, _ = conn.Write(data)
}
