package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	configContent := `
service:
  name: test-service
  port: 9090
  log_level: debug
database:
  driver: postgres
  host: testdb
  port: 5433
  user: testuser
  password: testpass
  dbname: testdb
  sslmode: disable
jwt:
  secret: test-secret
  duration_hours: 48
consul:
  enabled: true
  address: consul:8501
  service_id: test-service-1
etcd:
  enabled: true
  endpoints:
    - etcd:2379
tracing:
  jaeger_enabled: true
  jaeger_endpoint: http://jaeger:14268/api/traces
  otlp_enabled: false
  otlp_endpoint: otlp:4317
metrics:
  port: 9091
`
	tmpFile, err := os.CreateTemp("", "config*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.Write([]byte(configContent)); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()

	tests := []struct {
		name       string
		configPath string
		wantCfg    *Config
		wantErr    bool
	}{
		{
			name:       "Valid config file",
			configPath: tmpFile.Name(),
			wantCfg: &Config{
				Service: ServiceConfig{
					Name: "test-service",
					Port: 9090,
					LogLevel: "debug",
				},
				Database: DatabaseConfig{
					Driver:   "postgres",
					Host:     "testdb",
					Port:     5433,
					User:     "testuser",
					Password: "testpass",
					DBName:   "testdb",
					SSLMode:  "disable",
				},
				JWT: JWTConfig{
					Secret:        "test-secret",
					DurationHours: 48,
				},
				Consul: ConsulConfig{
					Enabled:   true,
					Address:   "consul:8501",
					ServiceID: "test-service-1",
				},
				Etcd: EtcdConfig{
					Enabled:   true,
					Endpoints: []string{"etcd:2379"},
				},
				Tracing: TracingConfig{
					JaegerEnabled:  true,
					JaegerEndpoint: "http://jaeger:14268/api/traces",
					OTLPEnabled:    false,
					OTLPEndpoint:   "otlp:4317",
				},
				Metrics: MetricsConfig{
					Port: "9091",
				},
			},
			wantErr: false,
		},
		{
			name:       "Invalid config file",
			configPath: "nonexistent.yaml",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := LoadConfig(tt.configPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if cfg.Service.Name != tt.wantCfg.Service.Name || cfg.Service.Port != tt.wantCfg.Service.Port ||
					cfg.Service.LogLevel != tt.wantCfg.Service.LogLevel {
					t.Errorf("Service config = %+v, want %+v", cfg.Service, tt.wantCfg.Service)
				}
				if cfg.Database.Driver != tt.wantCfg.Database.Driver || cfg.Database.Host != tt.wantCfg.Database.Host ||
					cfg.Database.Port != tt.wantCfg.Database.Port || cfg.Database.User != tt.wantCfg.Database.User ||
					cfg.Database.Password != tt.wantCfg.Database.Password || cfg.Database.DBName != tt.wantCfg.Database.DBName ||
					cfg.Database.SSLMode != tt.wantCfg.Database.SSLMode {
					t.Errorf("Database config = %+v, want %+v", cfg.Database, tt.wantCfg.Database)
				}
				if cfg.JWT.Secret != tt.wantCfg.JWT.Secret || cfg.JWT.DurationHours != tt.wantCfg.JWT.DurationHours {
					t.Errorf("JWT config = %+v, want %+v", cfg.JWT, tt.wantCfg.JWT)
				}
				if cfg.Consul.Enabled != tt.wantCfg.Consul.Enabled || cfg.Consul.Address != tt.wantCfg.Consul.Address ||
					cfg.Consul.ServiceID != tt.wantCfg.Consul.ServiceID {
					t.Errorf("Consul config = %+v, want %+v", cfg.Consul, tt.wantCfg.Consul)
				}
				if cfg.Etcd.Enabled != tt.wantCfg.Etcd.Enabled || len(cfg.Etcd.Endpoints) != len(tt.wantCfg.Etcd.Endpoints) ||
					cfg.Etcd.Endpoints[0] != tt.wantCfg.Etcd.Endpoints[0] {
					t.Errorf("Etcd config = %+v, want %+v", cfg.Etcd, tt.wantCfg.Etcd)
				}
				if cfg.Tracing.JaegerEnabled != tt.wantCfg.Tracing.JaegerEnabled ||
					cfg.Tracing.JaegerEndpoint != tt.wantCfg.Tracing.JaegerEndpoint ||
					cfg.Tracing.OTLPEnabled != tt.wantCfg.Tracing.OTLPEnabled ||
					cfg.Tracing.OTLPEndpoint != tt.wantCfg.Tracing.OTLPEndpoint {
					t.Errorf("Tracing config = %+v, want %+v", cfg.Tracing, tt.wantCfg.Tracing)
				}
				if cfg.Metrics.Port != tt.wantCfg.Metrics.Port {
					t.Errorf("Metrics config = %+v, want %+v", cfg.Metrics, tt.wantCfg.Metrics)
				}
			}
		})
	}
}

func TestDatabaseConfig_GetDSN(t *testing.T) {
	cfg := &DatabaseConfig{
		Host:     "testdb",
		Port:     5433,
		User:     "testuser",
		Password: "testpass",
		DBName:   "testdb",
		SSLMode:  "disable",
	}
	expectedDSN := "host=testdb port=5433 user=testuser password=testpass dbname=testdb sslmode=disable"
	if dsn := cfg.GetDSN(); dsn != expectedDSN {
		t.Errorf("GetDSN() = %v, want %v", dsn, expectedDSN)
	}
}