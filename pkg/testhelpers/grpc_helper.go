package testhelpers

import (
	"fmt"
	"net"
	"testing"
	"time"

	"cloud.google.com/go/rpcreplay"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func newRecorder(name string) (*rpcreplay.Recorder, error) {
	return rpcreplay.NewRecorder(fmt.Sprintf("testdata/%s.replay", name), nil)
}

func newReplayer(name string) (*rpcreplay.Replayer, error) {
	return rpcreplay.NewReplayer(fmt.Sprintf("testdata/%s.replay", name))
}

// ConnectOrRecord connects to GRPC if a local port is open and records, otherwise replays a previous connection
func ConnectOrRecord(t *testing.T, name string, port int) (*grpc.ClientConn, func()) {
	var conn *grpc.ClientConn
	timeout := time.Second
	url := fmt.Sprintf("localhost:%d", port)
	c, err := net.DialTimeout("tcp", url, timeout)
	if err != nil {
		rep, e := newReplayer(name)
		assert.NoError(t, e)

		conn, err = rep.Connection()
		assert.NoError(t, err)
		return conn, func() {
			err = rep.Close()
			assert.NoError(t, err)
		}
	}

	fmt.Printf("testing against %s\n", url) //nolint

	err = c.Close()
	assert.NoError(t, err)
	rec, err := newRecorder(name)
	assert.NoError(t, err)
	opts := append(rec.DialOptions(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err = grpc.Dial(url, opts...)
	assert.NoError(t, err)
	return conn, func() {
		err = rec.Close()
		assert.NoError(t, err)
	}
}
