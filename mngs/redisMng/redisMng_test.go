package redisMng

import (
	"context"
	"net"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/wiidz/goutil/structs/configStruct"
)

func TestNewRedisMngAuthenticatesAndSelectsDatabase(t *testing.T) {
	t.Parallel()

	server := miniredis.RunT(t)
	server.RequireUserAuth("crm-user", "secret")
	host, port, err := net.SplitHostPort(server.Addr())
	if err != nil {
		t.Fatalf("split redis address: %v", err)
	}

	mng, err := NewRedisMng(context.Background(), &configStruct.RedisConfig{
		Host:      host,
		Port:      port,
		Username:  "crm-user",
		Password:  "secret",
		Database:  3,
		MaxActive: 2,
		MaxIdle:   1,
	})
	if err != nil {
		t.Fatalf("NewRedisMng() error = %v", err)
	}
	t.Cleanup(func() {
		_ = mng.Client.Close()
	})

	if err := mng.Set(context.Background(), "login:42", "token", 0); err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	got, err := server.DB(3).Get("login:42")
	if err != nil {
		t.Fatalf("read selected database: %v", err)
	}
	if got != "token" {
		t.Fatalf("stored value = %q, want %q", got, "token")
	}
	if server.DB(0).Exists("login:42") {
		t.Fatal("value was unexpectedly written to database 0")
	}
}
