package process

import (
	"reflect"
	"testing"
)

func TestParseTasklistCSV(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   []int
	}{
		{
			name:   "single match",
			output: `"DeadByDaylight-Win64-Shipping.exe","12345","Console","1","123,456 K"` + "\r\n",
			want:   []int{12345},
		},
		{
			name: "multiple matches",
			output: `"DeadByDaylight-Win64-Shipping.exe","111","Console","1","100,000 K"` + "\r\n" +
				`"DeadByDaylight-Win64-Shipping.exe","222","Console","1","100,000 K"` + "\r\n",
			want: []int{111, 222},
		},
		{
			name:   "no matching process (English)",
			output: "INFO: No tasks are running which match the specified criteria.\r\n",
			want:   nil,
		},
		{
			name:   "no matching process (Russian)",
			output: "INFO: Не найдено задач с указанными критериями.\r\n",
			want:   nil,
		},
		{
			name:   "empty output",
			output: "",
			want:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTasklistCSV(tt.output)
			if err != nil {
				t.Fatalf("parseTasklistCSV() returned error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTasklistCSV() = %v, want %v", got, tt.want)
			}
		})
	}
}
