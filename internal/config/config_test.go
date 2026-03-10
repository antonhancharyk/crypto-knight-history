package config

import (
	"os"
	"testing"
)

func TestLoad_Valid(t *testing.T) {
	t.Setenv("DB_USER", "u")
	t.Setenv("DB_NAME", "n")
	t.Setenv("DB_HOST", "h")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("GRPC_HOST", "grpc:9090")
	t.Setenv("APP_SERVER_PORT", "9000")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() err = %v", err)
	}
	if cfg.DB.User != "u" || cfg.DB.Name != "n" {
		t.Errorf("DB: got user=%q name=%q", cfg.DB.User, cfg.DB.Name)
	}
	if cfg.DB.Host != "h" || cfg.DB.Port != "5433" {
		t.Errorf("DB: got host=%q port=%q", cfg.DB.Host, cfg.DB.Port)
	}
	if cfg.GRPC.Host != "grpc:9090" {
		t.Errorf("GRPC.Host = %q", cfg.GRPC.Host)
	}
	if cfg.Server.Port != "9000" {
		t.Errorf("Server.Port = %q", cfg.Server.Port)
	}
}

func TestLoad_Defaults(t *testing.T) {
	t.Setenv("DB_USER", "u")
	t.Setenv("DB_NAME", "n")
	t.Setenv("GRPC_HOST", "gh")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("APP_SERVER_PORT")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() err = %v", err)
	}
	if cfg.DB.Host != "localhost" || cfg.DB.Port != "5432" {
		t.Errorf("DB defaults: host=%q port=%q", cfg.DB.Host, cfg.DB.Port)
	}
	if cfg.Server.Port != "8080" {
		t.Errorf("Server.Port default = %q", cfg.Server.Port)
	}
}

func TestLoad_MissingDBUser(t *testing.T) {
	t.Setenv("DB_USER", "")
	t.Setenv("DB_NAME", "n")
	t.Setenv("GRPC_HOST", "gh")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() expected error when DB_USER empty")
	}
}

func TestLoad_MissingDBName(t *testing.T) {
	t.Setenv("DB_USER", "u")
	t.Setenv("DB_NAME", "")
	t.Setenv("GRPC_HOST", "gh")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() expected error when DB_NAME empty")
	}
}

func TestLoad_MissingGRPCHost(t *testing.T) {
	t.Setenv("DB_USER", "u")
	t.Setenv("DB_NAME", "n")
	t.Setenv("GRPC_HOST", "")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() expected error when GRPC_HOST empty")
	}
}
