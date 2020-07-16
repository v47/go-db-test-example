package storage

import (
	"context"
	"database/sql"
	"fmt"
)

// Client represents postgres db persistance layer client
type Client struct {
	conn *sql.DB
}

// NewClient creates a new postgres client
func NewClient(conn *sql.DB) *Client {
	return &Client{
		conn: conn,
	}
}

// SaveUser ...
func (c *Client) SaveUser(ctx context.Context, name string, docName string) (int, error) {
	row := c.conn.QueryRowContext(ctx, "INSERT INTO users (name, doc_name) VALUES ($1, $2) RETURNING id", name, docName)

	var id int

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

// FetchUser ...
func (c *Client) FetchUser(ctx context.Context, id int) (string, string, string, error) {
	var (
		name    string
		docName string
		docHash string
	)

	err := c.conn.QueryRowContext(ctx,
		`SELECT
			name,
			doc_name,
			doc_hash
		FROM users
		WHERE id = $1;`,
		id).Scan(
		&name,
		&docName,
		&docHash,
	)

	if err != nil {
		return "", "", "", fmt.Errorf("failed executing db query: %s", err)
	}

	return name, docName, docHash, nil
}
