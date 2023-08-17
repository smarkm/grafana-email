package datasource

import (
	"reflect"
	"testing"

	"smark.freecoop.net/grafana-email/config"
)

func TestDashboardPanels(t *testing.T) {
	config.Init("../config.json")

	type args struct {
		orgID string
		dID   string
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Test Case 1",
			args: args{
				orgID: "1",
				dID:   "Gfgpou3Vk",
			},
			want: []int{4, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DashboardPanels(tt.args.orgID, tt.args.dID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DashboardPanels() = %v, want %v", got, tt.want)
			}
		})
	}
}
