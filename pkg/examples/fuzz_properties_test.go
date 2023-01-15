package examples

import (
	"encoding/json"
	"math"
	"sort"
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
)

type SomeStruct struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	SomePredicate bool   `json:"predicate"`
}

func TestSortProperty(t *testing.T) {
	for i := 0; i < 1000; i++ {
		arr := fuzzIntArray()
		sort.Ints(arr)

		last := math.MinInt
		assert.GreaterOrEqual(t, len(arr), 0)
		for _, v := range arr {
			assert.GreaterOrEqual(t, v, last)
			last = v
		}
	}
}

func TestJsonIsSymmetric(t *testing.T) {

	structs := fuzzStruct()
	for _, str := range structs {
		out, err := json.Marshal(str)
		assert.NoError(t, err)
		var result SomeStruct
		json.Unmarshal(out, &result)
		assert.Equal(t, str, result)
	}
	assert.GreaterOrEqual(t, len(structs), 500)
	assert.LessOrEqual(t, len(structs), 3000)

}

func fuzzIntArray() []int {
	f := fuzz.New()
	slice := make([]int, 0)
	f.Fuzz(&slice)
	return slice
}

func fuzzStruct() []SomeStruct {
	f := fuzz.New().NilChance(0).NumElements(500, 3000)
	slice := make([]SomeStruct, 0)
	f.Fuzz(&slice)
	return slice
}
