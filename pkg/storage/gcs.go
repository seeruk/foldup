package storage

import (
	"context"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/storage/v1"
)

const gcsScopeReadWrite = "https://www.googleapis.com/auth/devstorage.read_write"

func NewGCSClient() (*storage.Service, error) {
	client, err := google.DefaultClient(context.Background(), gcsScopeReadWrite)
	if err != nil {
		return nil, err
	}

	return storage.New(client)
}
