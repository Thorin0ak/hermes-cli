package internal

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewOrchestrator(t *testing.T) {
	type args struct {
		config *Config
	}

	hConf := GetConfig()

	tests := []struct {
		name    string
		args    args
		want    *Orchestrator
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
		{
			name:    "should return a new orchestrator",
			args:    args{config: hConf},
			want:    nil,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewOrchestrator(tt.args.config)
			if !tt.wantErr(t, err, fmt.Sprintf("NewOrchestrator(%v)", tt.args.config)) {
				return
			}
			assert.Equalf(t, tt.want, got, "NewOrchestrator(%v)", tt.args.config)
		})
	}
}
