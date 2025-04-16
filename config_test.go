package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    interface{}
		wantErr bool
	}{
		{
			name: "expected behaviour",
			input: []byte(`
			app:
				name: "Test"
				version: "9.3"
				mode: "dev"
			`),
			want: struct {
				Name    string `mapstructure:"name"`
				Version string `mapstructure:"version"`
				Mode    string `mapstructure:"mode"`
			}{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadConfig()
			if tt.wantErr != (err != nil) {
				t.Fatalf("%s failed. wantErr: %v, err: %v", tt.name, tt.wantErr, err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("%s mismatch. (-want, +got):\n%s", tt.name, diff)
			}
		})
	}
}
