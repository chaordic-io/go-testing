package examples

import (
	"context"
	"os"
	"testing"

	"github.com/chaordic-io/go-testing/pkg/gen/api/gotesting/api/v1"
	"github.com/chaordic-io/go-testing/pkg/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestAPISearch(t *testing.T) {
	rt, close, err := testhelpers.MakeRecorder("search", t)
	defer close()
	assert.NoError(t, err)

	assert.NoError(t, os.Setenv("EDAMAM_APP_ID", "DUMMY_APP_ID"))
	assert.NoError(t, os.Setenv("EDAMAM_APP_KEY", "DUMMY_APP_KEY"))

	ctx := context.Background()
	s := New()
	s.client.Transport = rt
	r, err := s.Search(ctx, &api.SearchRequest{Query: "lasagna"})
	assert.NoError(t, err)
	assert.Equal(t, 20, len(r.Recipes))

	for _, receipe := range r.Recipes {
		assert.NotEmpty(t, receipe.Uri)
		assert.NotEmpty(t, receipe.Label)
		assert.NotEmpty(t, receipe.Image)
	}
}
