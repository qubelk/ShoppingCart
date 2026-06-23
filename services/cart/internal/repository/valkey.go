package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/valkey-io/valkey-go"
)

func NewValkeyClient() (valkey.Client, error) {
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{"localhost:6379"},
	})

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Do(ctx, client.B().Ping().Build()).Error(); err != nil {
		return nil, fmt.Errorf("failed to connection to valkey: %w", err)
	}

	return client, nil
}
