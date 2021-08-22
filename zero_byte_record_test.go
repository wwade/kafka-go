package kafka

import (
	"context"
	"net"
	"runtime/debug"
	"testing"
	"time"
)

// XXX tmp test file until we figure out where to put this test

type mockConn struct {
	net.Conn

	t *testing.T
}

func TestZeroByteRecord(t *testing.T) {
	r := NewReader(ReaderConfig{
		Brokers:   []string{"unused"},
		Topic:     "topic-A",
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
		Dialer: &Dialer{
			DialFunc: func(ctx context.Context, network string, address string) (net.Conn, error) {
				return &mockConn{t: t}, nil
			},
		},
	})
	defer r.Close()

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			t.Logf("readMessage failed: %s", err)
			break
		}
		t.Logf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}

	t.Fatal("XXX always fail")
}

func (c *mockConn) Write(b []byte) (n int, err error) {
	c.t.Logf("Write(%s)", b)
	//c.t.Logf("Write(%v)", b)
	return len(b), nil
}

func (c *mockConn) Read(b []byte) (n int, err error) {
	//c.t.Logf("Read(%d)", len(b))
	// XXX dump stack trace to get caller
	c.t.Logf("Read stack:\n%s", debug.Stack())

	output := []byte("hi")
	return copy(b, output), nil
}

func (c *mockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (c *mockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *mockConn) Close() error {
	return nil
}
