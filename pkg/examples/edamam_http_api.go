package examples

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/chaordic-io/go-testing/pkg/gen/api/gotesting/api/v1"
)

// Edamam service struct implements the gRPC API service for Recipe search
type Edamam struct {
	api.UnimplementedRecipeServiceServer
	client *http.Client
}

// SearchResponse struct, for the Recipe API
type SearchResponse struct {
	From int `json:"from"`
	To   int `json:"to"`
	Hits []struct {
		Recipe Recipe `json:"recipe"`
	} `json:"hits"`
}

// Recipe struct, for the Recipe API
type Recipe struct {
	URI   string `json:"uri"`
	Label string `json:"label"`
	Image string `json:"image"`
}

func New() *Edamam {
	return &Edamam{client: &http.Client{Timeout: 10 * time.Second}}
}

func (e *Edamam) Search(ctx context.Context, in *api.SearchRequest) (*api.SearchResponse, error) {
	params := url.Values{
		"app_id":  []string{os.Getenv("EDAMAM_APP_ID")},
		"app_key": []string{os.Getenv("EDAMAM_APP_KEY")},
		"q":       []string{in.Query},
		"type":    []string{"public"},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://api.edamam.com/api/recipes/v2?%s", params.Encode()), nil)
	if err != nil {
		return nil, err
	}

	res, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = res.Body.Close()
		if err != nil {
			log.Printf("can't close response body, %v", err)
		}
	}()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search query responded status %d", res.StatusCode)
	}
	decoder := json.NewDecoder(res.Body)

	s := &SearchResponse{}
	err = decoder.Decode(s)
	if err != nil {
		return nil, err
	}

	r := &api.SearchResponse{Recipes: []*api.Recipe{}}
	for _, re := range s.Hits {
		r.Recipes = append(r.Recipes, &api.Recipe{Uri: re.Recipe.URI, Label: re.Recipe.Label, Image: re.Recipe.Label})
	}

	return r, nil
}
