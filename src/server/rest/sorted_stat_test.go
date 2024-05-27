package rest

import (
	"testing"

	"github.com/DervexDev/ghloc/src/service/loc_count"
	"github.com/stretchr/testify/assert"
)

func TestSortedStatMarshaling(t *testing.T) {
	stat := &loc_count.StatTree{
		LOC: 100,
		LOCByLangs: map[string]loc_count.LinesNumber{
			"Go":  10,
			"PHP": 50,
		},
		Children: map[string]*loc_count.StatTree{
			"file1": {
				LOC: 500,
			},
			"file2": {
				LOC: 100,
			},
		},
	}
	j, err := (*SortedStat)(stat).MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, `{"loc":100,"locByLangs":{"PHP":50,"Go":10},"children":{"file1":500,"file2":100}}`, string(j))
}
