on: push

jobs:
  unit-tests:
    name: Unit tests
    runs-on: ubuntu-latest
    steps:
      - name: "Checkout"
        uses: actions/checkout@master
      - uses: actions/setup-go@master
        with:
          go-version: "1.13.3"
      - name: "Unit tests"
        run: make test

  integration-tests:
    name: Integration tests
    runs-on: ubuntu-latest
    services:
      # Label used to access the service container
      postgres:
        # Docker Hub image
        image: postgres
        # Provide the password for postgres
        env:
          POSTGRES_PASSWORD: postgres
        # Set health checks to wait until postgres has started
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          # Maps tcp port 5432 on service container to the host
          - 5432:5432
    steps:
      - name: "Checkout"
        uses: actions/checkout@master
      - uses: actions/setup-go@master
        with:
          go-version: "1.13.3"
      - name: Install sql-migrate
        run: go get -v github.com/rubenv/sql-migrate/...
      - name: Apply migrations
        run: make db-ci-migrate-up
      - name: "Integration tests"
        run: make test-integration
        env:
          TEST_DB_URI: postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
