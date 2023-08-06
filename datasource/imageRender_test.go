package datasource

import (
	"reflect"
	"testing"

	"smark.freecoop.net/grafana-email/config"
)

func TestPanelImage(t *testing.T) {
	config.Init()
	type args struct {
		orgID   string
		dID     string
		panelID string
		vars    map[string]string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "Test 2",
			args: args{orgID: "1", dID: "Gfgpou3Vk", panelID: "4"},
			want: []byte{1},
		},
		{
			name: "Test3",
			args: args{orgID: "1", dID: "Gfgpou3Vk", panelID: "2"},
			want: []byte{1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PanelImage(tt.args.orgID, tt.args.dID, tt.args.panelID, tt.args.vars); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PanelImage() = %v, want %v", got, tt.want)
			}
		})
	}
}
