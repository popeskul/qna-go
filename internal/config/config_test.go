package config

import (
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		folder   string
		filename string
	}
	tests := []struct {
		name  string
		args  args
		want  *Config
		error bool
	}{
		{
			name: "ok",
			args: args{
				folder:   ".",
				filename: "config_test",
			},
			want: &Config{
				DB: Postgres{
					Host:     "localhost",
					Port:     5432,
					User:     "postgres",
					Password: "******",
					DBName:   "postgres",
					SSLMode:  "disable",
				},
				Server: struct {
					Port int `mapstructure:"port"`
				}{
					Port: 8080,
				},
			},
			error: false,
		},
		{
			name: "bad",
			args: args{
				folder:   ".",
				filename: "config1",
			},
			want:  nil,
			error: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.folder, tt.args.filename)
			if (err != nil) != tt.error {
				t.Errorf("New() error = %v, wantErr %v", err, tt.error)
			}

			if got != nil {
				got.DB.Password = "******"
			}

			if got != nil && tt.want != nil {
				if got.DB.Host != tt.want.DB.Host {
					t.Errorf("New() got = %v, want %v", got.DB.Host, tt.want.DB.Host)
				}
				if got.DB.Port != tt.want.DB.Port {
					t.Errorf("New() got = %v, want %v", got.DB.Port, tt.want.DB.Port)
				}
				if got.DB.User != tt.want.DB.User {
					t.Errorf("New() got = %v, want %v", got.DB.User, tt.want.DB.User)
				}
				if got.DB.Password != tt.want.DB.Password {
					t.Errorf("New() got = %v, want %v", got.DB.Password, tt.want.DB.Password)
				}
				if got.DB.DBName != tt.want.DB.DBName {
					t.Errorf("New() got = %v, want %v", got.DB.DBName, tt.want.DB.DBName)
				}
				if got.DB.SSLMode != tt.want.DB.SSLMode {
					t.Errorf("New() got = %v, want %v", got.DB.SSLMode, tt.want.DB.SSLMode)
				}
				if got.Server.Port != tt.want.Server.Port {
					t.Errorf("New() got = %v, want %v", got.Server.Port, tt.want.Server.Port)
				}
			}

		})
	}
}
