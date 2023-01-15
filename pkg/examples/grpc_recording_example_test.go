package examples

import (
	"context"
	"testing"

	"github.com/chaordic-io/go-testing/pkg/gen/api/gotesting/api/v1"
	"github.com/chaordic-io/go-testing/pkg/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestGRPCRecording(t *testing.T) {
	conn, close := testhelpers.ConnectOrRecord(t, "recipe_server", 50001)
	defer close()

	ctx := context.Background()
	c := api.NewRecipeServiceClient(conn)
	_, err := c.Search(ctx, &api.SearchRequest{Query: "lasagne"})
	assert.NoError(t, err)
}
