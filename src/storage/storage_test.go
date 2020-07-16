// +build integration

package storage_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/v47/go-db-test-example/src/storage"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var conn *sql.DB

func TestMain(m *testing.M) {
	dbURI := os.Getenv("TEST_DB_URI")
	if dbURI == "" {
		fmt.Println("environment variable TEST_DB_URI unset")
		os.Exit(1)
	}

	// connect to DB
	var err error
	conn, err = sql.Open("postgres", dbURI)
	if err != nil {
		fmt.Printf("error connect to DB at %q: %s\n", dbURI, err)
		os.Exit(1)
	}

	// check DB connection
	for i := 0; i < 10; i++ {
		err = conn.Ping()
		if err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if err != nil {
		fmt.Printf("error ping DB at %q: %s\n", dbURI, err)
		os.Exit(1)
	}

	// run tests
	os.Exit(m.Run())
}

func TestSaveUser(t *testing.T) {
	clearDB := func() {
		_, err := conn.Exec(`TRUNCATE users RESTART IDENTITY CASCADE;`)
		if err != nil {
			panic(err)
		}
	}
	defer clearDB()
	c := storage.NewClient(conn)
	id, err := c.SaveUser(context.Background(), "the_user", "hello, world")

	require.NoError(t, err)
	assert.Equal(t, 1, id)
}

func TestFetchUser(t *testing.T) {
	clearDB := func() {
		_, err := conn.Exec(`TRUNCATE users RESTART IDENTITY CASCADE;`)
		if err != nil {
			panic(err)
		}
	}
	defer clearDB()

	c := storage.NewClient(conn)

	id, err := c.SaveUser(context.Background(), "the_user", "hello, world")
	require.NoError(t, err)

	name, docName, docHash, err := c.FetchUser(context.Background(), id)
	require.NoError(t, err)

	assert.Equal(t, "the_user", name)
	assert.Equal(t, "hello, world", docName)
	assert.Equal(t, "E4D7F1B4ED2E42D15898F4B27B019DA4", docHash)
}
