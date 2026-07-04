package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseInput(t *testing.T) {
	tests := []struct {
		name     string
		page     int32
		pageSize int32
		want     PageInput
	}{
		{
			name:     "both zero means no limit",
			page:     0,
			pageSize: 0,
			want:     PageInput{},
		},
		{
			name:     "page set defaults page size",
			page:     2,
			pageSize: 0,
			want:     PageInput{Limit: DefaultPageSize, Offset: DefaultPageSize},
		},
		{
			name:     "page size set defaults page",
			page:     0,
			pageSize: 25,
			want:     PageInput{Limit: 25, Offset: 0},
		},
		{
			name:     "calculates offset",
			page:     3,
			pageSize: 25,
			want:     PageInput{Limit: 25, Offset: 50},
		},
		{
			name:     "negative inputs treated as no limit",
			page:     -5,
			pageSize: -1,
			want:     PageInput{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ParseInput(tt.page, tt.pageSize))
		})
	}
}

func TestRequestNoLimit(t *testing.T) {
	req := Request{}

	assert.Equal(t, PageInput{}, req.Input())
	assert.Equal(t, PageOutput{Page: 1, PageSize: 41, TotalPages: 1, TotalSize: 41}, req.Output(41))
	assert.Equal(t, PageOutput{Page: 1, PageSize: 0, TotalPages: 0, TotalSize: 0}, req.Output(0))
}

func TestCalcOutput(t *testing.T) {
	tests := []struct {
		name     string
		page     int32
		pageSize int32
		total    int64
		want     PageOutput
	}{
		{
			name:  "both zero collapses to single page",
			total: 41,
			want:  PageOutput{Page: 1, PageSize: 41, TotalPages: 1, TotalSize: 41},
		},
		{
			name:     "calculates pages",
			page:     2,
			pageSize: 10,
			total:    20,
			want:     PageOutput{Page: 2, PageSize: 10, TotalPages: 2, TotalSize: 20},
		},
		{
			name:     "zero total",
			page:     1,
			pageSize: 10,
			total:    0,
			want:     PageOutput{Page: 1, PageSize: 10, TotalPages: 0, TotalSize: 0},
		},
		{
			name:     "clamps totals to max int32",
			page:     1,
			pageSize: 1,
			total:    int64(maxInt32) + 1,
			want:     PageOutput{Page: 1, PageSize: 1, TotalPages: maxInt32, TotalSize: maxInt32},
		},
		{
			name:     "exact multiple of page size",
			page:     1,
			pageSize: 10,
			total:    20,
			want:     PageOutput{Page: 1, PageSize: 10, TotalPages: 2, TotalSize: 20},
		},
		{
			name:     "negative total clamps to zero",
			page:     1,
			pageSize: 10,
			total:    -5,
			want:     PageOutput{Page: 1, PageSize: 10, TotalPages: 0, TotalSize: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, CalcOutput(tt.page, tt.pageSize, tt.total))
		})
	}
}
