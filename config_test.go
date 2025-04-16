package main

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type AppConfig struct {
	Sql struct {
		Address string `mapstructure:"address" validate:"required"`
		Port    string `mapstructure:"port" validate:"required"`
		DbName  string `mapstructure:"dbName" validate:"required"`
	} `mapstructure:"sql" validate:"required"`
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		desiredType any
		want        any
		wantErr     bool
	}{
		{
			name: "expected behaviour",
			input: []byte(`
sql:
  address: "127.0.0.1"
  port: ":9"
  dbName: "test"
`),
			want: AppConfig{
				Sql: struct {
					Address string "mapstructure:\"address\" validate:\"required\""
					Port    string "mapstructure:\"port\" validate:\"required\""
					DbName  string "mapstructure:\"dbName\" validate:\"required\""
				}{
					Address: "127.0.0.1",
					Port:    ":9",
					DbName:  "test",
				},
			},
			wantErr: false,
		},
		{
			name: "error on empty field",
			input: []byte(`
sql:
  address: "127.0.0.1"
  port: ":9"
`),
			desiredType: t,
			want:        AppConfig{},
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadConfig[AppConfig](bytes.NewReader(tt.input), "yaml")
			if tt.wantErr != (err != nil) {
				t.Fatalf("%s failed. wantErr: %v, err: %v", tt.name, tt.wantErr, err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("%s mismatch. (-want, +got):\n%s", tt.name, diff)
			}
		})
	}
}
