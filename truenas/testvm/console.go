package testvm

import (
	"bufio"
	"errors"
	"io"
	"net"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type console struct {
	sock     string
	listener net.Listener
	wg       sync.WaitGroup
}

func newConsole(t *testing.T) *console {
	t.Helper()
	sock := filepath.Join(t.TempDir(), "console.sock")
	l, err := net.Listen("unix", sock)
	require.NoError(t, err)
	c := &console{
		sock:     sock,
		listener: l,
	}

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			conn, err := c.listener.Accept()
			if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {
				return
			}
			if err != nil {
				t.Logf("failed to accept console socket connection: %v", err)
				continue
			}
			t.Logf("console socket connected: %s", conn.LocalAddr())

			c.wg.Add(1)
			go func() {
				defer c.wg.Done()

				reader := bufio.NewReader(conn)
				for {
					s, err := reader.ReadString('\n')
					if err != nil {
						if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {
							return
						}
						t.Logf("failed to read from console socket: %v", err)
						continue
					}
					s = strings.TrimSpace(stripANSI(s))
					if s != "" {
						t.Logf("[vm-console] %s", s)
					}
				}
			}()
		}
	}()
	return c
}

func (c *console) Addr() string {
	return "unix:" + c.sock
}

func (c *console) Close() {
	_ = c.listener.Close()
	c.wg.Wait()
}

var ansiRegex = regexp.MustCompile("[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))")

func stripANSI(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}
