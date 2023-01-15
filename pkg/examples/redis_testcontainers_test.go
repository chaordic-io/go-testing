package examples

import (
	"context"
	"testing"

	"github.com/chaordic-io/go-testing/pkg/testhelpers"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

func TestRedis(t *testing.T) {

	ctx := context.Background()
	container, err := testhelpers.SetupRedis(ctx)
	assert.NoError(t, err)

	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	options, err := redis.ParseURL(container.URI)
	if err != nil {
		t.Fatal(err)
	}
	client := redis.NewClient(options)
	defer flushRedis(ctx, *client)

	pong, err := client.Ping().Result()
	assert.NoError(t, err)
	assert.Equal(t, pong, "PONG")

}

func flushRedis(ctx context.Context, client redis.Client) error {
	return client.FlushAll().Err()
}
