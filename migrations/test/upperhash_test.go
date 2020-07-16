// +build integration
package upperhash_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

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

func UpperHash(t *testing.T) {
	row := conn.QueryRow("SELECT upper_md5($1::text)", "hello, world")
	var hash string
	if err := row.Scan(&hash); err != nil {
		require.NoError(t, err)
	}
	assert.Equal(t, "E4D7F1B4ED2E42D15898F4B27B019DA4", hash)
}
