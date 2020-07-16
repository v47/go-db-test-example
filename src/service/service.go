package service

import "context"

// Storage ...
type Storage interface {
	SaveUser(ctx context.Context, name string, docName string) (int, error)
	FetchUser(ctx context.Context, id int) (string, string, string, error)
}
