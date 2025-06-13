package truenas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/puzpuzpuz/xsync/v3"
)

type Options struct {
	Username string
	Password string
	APIKey   string
	Debug    bool
}

type Client struct {
	// Type-safe API clients
	Auth         *AuthClient
	Pool         *PoolClient
	Dataset      *DatasetClient
	Service      *ServiceClient
	System       *SystemClient
	Network      *NetworkClient
	SMB          *SMBClient
	NFS          *NFSClient
	SSH          *SSHClient
	Smart        *SmartClient
	VM           *VMClient
	Job          *JobClient
	VMDevice     *VMDeviceClient
	User         *UserClient
	Group        *GroupClient
	Alert        *AlertClient
	AlertService *AlertServiceClient
	Boot         *BootClient
	Certificate  *CertificateClient
	Cronjob      *CronjobClient
	Disk         *DiskClient
	APIKey       *APIKeyClient
	Filesystem   *FilesystemClient
	Sharing      *SharingClient

	// Internal state
	url         string
	conn        *websocket.Conn
	opts        Options
	mu          sync.RWMutex
	writeMu     sync.Mutex
	msgID       atomic.Int64
	pending     *xsync.MapOf[string, chan Message]
	errCh       chan error
	reconnectCh chan struct{}
	closed      atomic.Bool
	wg          sync.WaitGroup
}

type Message struct {
	ID     string          `json:"id,omitempty"`
	Msg    string          `json:"msg,omitempty"`
	Method string          `json:"method,omitempty"`
	Params []any           `json:"params,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *ErrorMsg       `json:"error,omitempty"`
}

func (m *Message) Unmarshal(v any) error {
	if err := json.Unmarshal(m.Result, v); err != nil {
		return fmt.Errorf("unmarshal result: %s: %w", string(m.Result), err)
	}
	return nil
}

type ErrorMsg struct {
	Message string `json:"message,omitempty"`
	Code    int    `json:"error,omitempty"`
	Reason  string `json:"reason,omitempty"`
	Type    string `json:"errorType,omitempty"`
}

func (e *ErrorMsg) Error() string {
	var parts []string
	if e.Code > 0 {
		parts = append(parts, fmt.Sprintf("code: %d", e.Code))
	}
	if e.Message != "" {
		parts = append(parts, fmt.Sprintf("message: %s", e.Message))
	}
	if e.Reason != "" {
		parts = append(parts, fmt.Sprintf("reason: %s", e.Reason))
	}
	if e.Type != "" {
		parts = append(parts, fmt.Sprintf("type: %s", e.Type))
	}
	if len(parts) == 0 {
		return "TrueNAS API error"
	}
	return fmt.Sprintf("TrueNAS API error (%s)", strings.Join(parts, ", "))
}

// NewClient builds a new TrueNAS Client.
// Close() should be called to clean up resources when the client is no longer needed.
func NewClient(endpoint string, opts Options) (*Client, error) {
	c := &Client{
		url:         endpoint,
		opts:        opts,
		pending:     xsync.NewMapOf[string, chan Message](),
		errCh:       make(chan error, 1),
		reconnectCh: make(chan struct{}, 1),
	}

	// Initialize type-safe API clients
	c.Auth = NewAuthClient(c)
	c.Pool = NewPoolClient(c)
	c.Dataset = NewDatasetClient(c)
	c.Service = NewServiceClient(c)
	c.System = NewSystemClient(c)
	c.Network = NewNetworkClient(c)
	c.SMB = NewSMBClient(c)
	c.NFS = NewNFSClient(c)
	c.SSH = NewSSHClient(c)
	c.Smart = NewSmartClient(c)
	c.VM = NewVMClient(c)
	c.VMDevice = NewVMDeviceClient(c)
	c.User = NewUserClient(c)
	c.Group = NewGroupClient(c)
	c.Alert = NewAlertClient(c)
	c.Job = NewJobClient(c)
	c.AlertService = NewAlertServiceClient(c)
	c.Boot = NewBootClient(c)
	c.Certificate = NewCertificateClient(c)
	c.Cronjob = NewCronjobClient(c)
	c.Disk = NewDiskClient(c)
	c.APIKey = NewAPIKeyClient(c)
	c.Filesystem = NewFilesystemClient(c)
	c.Sharing = NewSharingClient(c)

	if err := c.connect(); err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	c.wg.Add(1)
	go c.connectionManager()

	if err := c.authenticate(); err != nil {
		_ = c.Close()
		return nil, fmt.Errorf("authentication: %w", err)
	}

	return c, nil
}

func (c *Client) connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	u, err := url.Parse(c.url)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	dialer := &websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}

	conn, _, err := dialer.Dial(u.String(), http.Header{})
	if err != nil {
		return fmt.Errorf("websocket dial: %s: %w", u.String(), err)
	}

	msg := map[string]any{
		"msg":     "connect",
		"version": "1",
		"support": []string{"1"},
	}
	if c.opts.Debug {
		fmt.Printf("send: %s\n", tryMarshal(msg))
	}
	c.writeMu.Lock()
	err = conn.WriteJSON(msg)
	c.writeMu.Unlock()
	if err != nil {
		conn.Close()
		return fmt.Errorf("send connect request: %w", err)
	}

	var resp struct {
		Msg     string `json:"msg"`
		Session string `json:"session"`
	}
	if err := conn.ReadJSON(&resp); err != nil {
		conn.Close()
		return fmt.Errorf("read connection response: %w", err)
	}
	if c.opts.Debug {
		fmt.Printf("recv: %s\n", tryMarshal(resp))
	}
	if !strings.EqualFold(resp.Msg, "connected") {
		conn.Close()
		return fmt.Errorf("connection failed: %s", resp.Msg)
	}
	if resp.Session == "" {
		conn.Close()
		return fmt.Errorf("connected but did not receive a session")
	}
	c.conn = conn
	c.closed.Store(false)
	return nil
}

func (c *Client) authenticate() error {
	// Skip authentication if no credentials provided
	if c.opts.APIKey == "" && c.opts.Username == "" && c.opts.Password == "" {
		return nil
	}

	var method string
	var params []any

	if c.opts.APIKey != "" {
		method = "auth.login_with_api_key"
		params = []any{c.opts.APIKey}
	} else {
		method = "auth.login"
		params = []any{c.opts.Username, c.opts.Password}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	var success bool
	if err := c.Call(ctx, method, params, &success); err != nil {
		return fmt.Errorf("call %s: %w", method, err)
	}
	if success {
		return nil
	}
	return fmt.Errorf("auth unsuccessful")
}

func (c *Client) Call(ctx context.Context, method string, params []any, v any) error {
	msgID := fmt.Sprintf("%d", c.msgID.Add(1))

	msg := &Message{
		ID:     msgID,
		Msg:    "method",
		Method: method,
		Params: params,
	}

	resultCh := make(chan Message, 1)

	c.pending.Store(msgID, resultCh)

	defer func() {
		c.pending.Delete(msgID)
	}()

	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()

	if conn == nil || c.closed.Load() {
		return fmt.Errorf("not connected")
	}

	if c.opts.Debug {
		fmt.Printf("send: %s\n", tryMarshal(msg))
	}
	c.writeMu.Lock()
	err := conn.WriteJSON(msg)
	c.writeMu.Unlock()
	if err != nil {
		return fmt.Errorf("send message: %s: %w", msgID, err)
	}

	select {
	case err := <-c.errCh:
		return err
	case result, ok := <-resultCh:
		if !ok {
			// Channel was closed, client is shutting down
			return fmt.Errorf("client closed")
		}
		if result.Error != nil {
			return result.Error
		}
		if v != nil {
			return result.Unmarshal(v)
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// CallJob calls a job endpoint and waits for completion, returning the job result
func (c *Client) CallJob(ctx context.Context, method string, params []any, v any) error {
	var jobID int
	if err := c.Call(ctx, method, params, &jobID); err != nil {
		return fmt.Errorf("call %s: %w", method, err)
	}

	job, err := c.Job.Wait(ctx, jobID)
	if err != nil {
		return fmt.Errorf("wait for job %d (%s): %w", jobID, method, err)
	}

	if v != nil && job.Result != nil {
		resultBytes, err := json.Marshal(job.Result)
		if err != nil {
			return fmt.Errorf("marshal job result: %w", err)
		}
		if err := json.Unmarshal(resultBytes, v); err != nil {
			return fmt.Errorf("unmarshal job result: %w", err)
		}
	}

	return nil
}

func (c *Client) connectionManager() {
	defer c.wg.Done()
	defer func() {
		if c.opts.Debug {
			fmt.Println("connectionManager exiting")
		}
	}()

	c.wg.Add(1)
	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()
	go c.readLoop(conn)

	for !c.closed.Load() {
		<-c.reconnectCh
		if c.closed.Load() {
			return
		}

		if c.opts.Debug {
			fmt.Println("attempting to reconnect...")
		}

		if err := c.reconnect(); err != nil {
			if c.opts.Debug {
				fmt.Printf("reconnection failed, retrying: %v\n", err)
			}
			if !c.closed.Load() {
				// Wait briefly before retrying, but check for close frequently
				<-time.After(100 * time.Millisecond)
				if !c.closed.Load() {
					select {
					case c.reconnectCh <- struct{}{}:
					default:
					}
				}
			}
			continue
		}
		if c.opts.Debug {
			fmt.Println("reconnected successfully")
		}

		c.wg.Add(1)
		c.mu.RLock()
		conn := c.conn
		c.mu.RUnlock()
		go c.readLoop(conn)
	}
}

func (c *Client) readLoop(conn *websocket.Conn) {
	defer c.wg.Done()
	defer func() {
		if c.opts.Debug {
			fmt.Println("readLoop exiting")
		}
	}()

	if conn == nil {
		return
	}

	for !c.closed.Load() {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			if c.closed.Load() {
				return
			}
			// Check for connection errors that should trigger reconnection
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) ||
				websocket.IsUnexpectedCloseError(err) ||
				strings.Contains(err.Error(), "connection reset") ||
				strings.Contains(err.Error(), "broken pipe") ||
				strings.Contains(err.Error(), "use of closed network connection") {
				if c.opts.Debug {
					fmt.Printf("connection lost: %v\n", err)
				}
				select {
				case c.reconnectCh <- struct{}{}:
				default:
				}
				return
			}
			if c.opts.Debug {
				fmt.Printf("recv err: %v\n", err)
			}
			select {
			case c.errCh <- fmt.Errorf("read message: %w", err):
			default:
			}
			continue
		}
		if c.opts.Debug {
			fmt.Printf("recv: %s\n", tryMarshal(msg))
		}

		if msg.ID != "" {
			if ch, exists := c.pending.Load(msg.ID); exists {
				ch <- msg
			}
		}
	}
}

func (c *Client) reconnect() error {
	if err := c.connect(); err != nil {
		return err
	}
	return c.authenticate()
}

func (c *Client) Close() error {
	c.closed.Store(true)

	// Cancel all pending requests by closing their channels
	c.pending.Range(func(id string, ch chan Message) bool {
		close(ch)
		c.pending.Delete(id)
		return true
	})

	c.mu.Lock()
	if c.conn != nil {
		// Set read deadline to immediately unblock any pending reads
		_ = c.conn.SetReadDeadline(time.Now())
		_ = c.conn.Close() // Ignore close errors - connection might already be closed
		c.conn = nil
	}
	c.mu.Unlock()

	close(c.reconnectCh)
	c.wg.Wait()
	return nil
}

func tryMarshal(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func Ptr[T any](v T) *T {
	return &v
}

func value[T any](v *T) T {
	var zero T
	if v == nil {
		return zero
	}
	return *v
}
