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
			name:     "defaults",
			page:     0,
			pageSize: 0,
			want:     PageInput{Limit: DefaultPageSize, Offset: 0},
		},
		{
			name:     "calculates offset",
			page:     3,
			pageSize: 25,
			want:     PageInput{Limit: 25, Offset: 50},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ParseInput(tt.page, tt.pageSize))
		})
	}
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
			name:  "defaults",
			total: 41,
			want:  PageOutput{Page: 1, PageSize: DefaultPageSize, TotalPages: 3, TotalSize: 41},
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, CalcOutput(tt.page, tt.pageSize, tt.total))
		})
	}
}
