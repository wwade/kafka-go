package kafka

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"
)

// XXX tmp test file until we figure out where to put this test

type mockConn struct {
	net.Conn
}

func TestZeroByteRecord(t *testing.T) {
	r := NewReader(ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     "topic-A",
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
		Dialer: &Dialer{
			DialFunc: func(ctx context.Context, network string, address string) (net.Conn, error) {
				return &mockConn{}, nil
			},
		},
	})
	defer r.Close()

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}

	t.Fatal("XXX always fail")
}

func (c *mockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (c *mockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *mockConn) Write(b []byte) (n int, err error) {
	return len(b), nil
}

func (c *mockConn) Read(b []byte) (n int, err error) {
	return 0, nil
}

func (c *mockConn) Close() error {
	return nil
}
