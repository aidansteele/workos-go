package auditlog

import (
	"strconv"
	"testing"
)

var limits = []struct {
	in  int
	out int
}{
	{1, 1},
	{-1, 10},
	{0, 10},
	{1001, 1000},
	{100, 100},
}

func TestEventsRequestParamsLimit(t *testing.T) {
	for _, tt := range limits {
		t.Run(strconv.Itoa(tt.in), func(t *testing.T) {
			params := ListRequestParams{
				Limit: tt.in,
			}
			if params.GetLimit() != tt.out {
				t.Errorf("got %q, wanted %q", params.GetLimit(), tt.out)
			}
		})
	}
}
