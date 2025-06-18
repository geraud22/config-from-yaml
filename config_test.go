package cfy

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

type MockValidator struct {
	called bool
}

func (m *MockValidator) Struct(any) error {
	m.called = true
	return nil
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		desiredType any
		want        any
		wantErr     bool
		validator   *MockValidator
	}{
		{
			name: "correctly load all fields - case insensitive",
			input: []byte(`
SQL:
  address: "127.0.0.1"
  port: ":9"
  DBnAME: "test"
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
		{
			name:      "inject custom validator",
			wantErr:   false,
			want:      AppConfig{},
			validator: &MockValidator{},
		},
	}

	for _, tt := range tests {
		var v Validator
		if tt.validator != nil {
			v = tt.validator
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadConfig[AppConfig](bytes.NewReader(tt.input), "yaml", v)
			if tt.wantErr != (err != nil) {
				t.Fatalf("wantErr: %v, err: %v", tt.wantErr, err)
			}
			if tt.validator != nil {
				if !tt.validator.called {
					t.Fatalf("wanted to use mock validator, but it was not called")
				}
			}
			if err == nil {
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Fatalf("Mismatch. (-want, +got):\n%s", diff)
				}
			}
		})
	}
}
