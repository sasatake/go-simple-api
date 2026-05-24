package handler

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/mongodb"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := mongodb.Run(ctx, "mongo:6.0.2")
	if err != nil {
		log.Fatalf("failed to start mongo container: %v", err)
	}

	connStr, err := container.ConnectionString(ctx)
	if err != nil {
		_ = container.Terminate(ctx)
		log.Fatalf("failed to get connection string: %v", err)
	}
	uri = connStr

	code := m.Run()

	if err := container.Terminate(ctx); err != nil {
		log.Printf("warning: failed to terminate container: %v", err)
	}
	os.Exit(code)
}
