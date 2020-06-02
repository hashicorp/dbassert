package dbassert

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
)

// MockTesting provides a testing.T mock.
type MockTesting struct {
	err bool
	msg string
}

// HasError returns if the MockTesting has an error for a previous test.
func (m *MockTesting) HasError() bool {
	return m.err
}

// ErrorMsg returns if the MockTesting has an error msg for a previous
// test.
func (m *MockTesting) ErrorMsg() string {
	return m.msg
}

// Errorf provides a mock Errorf function.
func (m *MockTesting) Errorf(format string, args ...interface{}) {
	m.msg = fmt.Sprintf(format, args...)
	m.err = true
}

// FailNow provides a mock FailNow function.
func (m *MockTesting) FailNow() {
	m.err = true
}

// Reset will reset the MockTesting for the next test.
func (m *MockTesting) Reset() {
	m.err = false
	m.msg = ""
}

// AssertNoError asserts that the MockTesting has no current error.
func (m *MockTesting) AssertNoError(t *testing.T) {
	t.Helper()
	if m.HasError() {
		t.Error("Check should not fail")
	}
}

// AssertError asserts that the MockTesting has a current error.
func (m *MockTesting) AssertError(t *testing.T) {
	t.Helper()
	if !m.HasError() {
		t.Error("Check should fail")
	}
}

// TestSetup sets up the testing env, including starting a docker container
// running the db dialect and initializing the test database schema.
func TestSetup(t *testing.T, dialect string) (func() error, *sql.DB, string) {
	cleanup, url, _, err := StartDbInDocker(dialect)
	if err != nil {
		t.Fatal(err)
	}
	db, err := sql.Open(dialect, url)
	if err != nil {
		t.Fatal(err)
	}
	if err := initStore(t, db); err != nil {
		t.Fatal(err)
	}
	return cleanup, db, url
}

// StartDbInDocker starts up the dialect db in the local docker.
func StartDbInDocker(dialect string) (cleanup func() error, retURL, container string, err error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return func() error { return nil }, "", "", fmt.Errorf("could not connect to docker: %w", err)
	}

	var resource *dockertest.Resource
	var url string
	switch dialect {
	case "postgres":
		resource, err = pool.Run("postgres", "latest", []string{"POSTGRES_PASSWORD=secret", "POSTGRES_DB=watchtower"})
		url = "postgres://postgres:secret@localhost:%s?sslmode=disable"
	default:
		panic(fmt.Sprintf("unknown dialect %q", dialect))
	}
	if err != nil {
		return func() error { return nil }, "", "", fmt.Errorf("could not start resource: %w", err)
	}

	cleanup = func() error {
		return cleanupDockerResource(pool, resource)
	}

	url = fmt.Sprintf(url, resource.GetPort("5432/tcp"))

	if err := pool.Retry(func() error {
		db, err := sql.Open(dialect, url)
		if err != nil {
			return fmt.Errorf("error opening %s dev container: %w", dialect, err)
		}

		if err := db.Ping(); err != nil {
			return err
		}
		defer db.Close()
		return nil
	}); err != nil {
		return func() error { return nil }, "", "", fmt.Errorf("could not connect to docker: %w", err)
	}

	return cleanup, url, resource.Container.Name, nil
}

func initStore(t *testing.T, db *sql.DB) error {
	const (
		createDomainType = `
create domain dbasserts_public_id as text
check(
  length(trim(value)) > 10
);
comment on domain dbasserts_public_id is
'dbasserts test domain type';
`
		createTable = `
create table if not exists test_table_dbasserts (
  id bigint generated always as identity primary key,
  public_id dbasserts_public_id not null,
  nullable text,
  type_int int
);
comment on table test_table_dbasserts is
'dbasserts test table'
`
	)
	if _, err := db.Exec(createDomainType); err != nil {
		return err
	}
	if _, err := db.Exec(createTable); err != nil {
		return err
	}
	return nil
}

// cleanupDockerResource will clean up the dockertest resources
func cleanupDockerResource(pool *dockertest.Pool, resource *dockertest.Resource) error {
	var err error
	for i := 0; i < 10; i++ {
		err = pool.Purge(resource)
		if err == nil {
			return nil
		}
	}
	if strings.Contains(err.Error(), "No such container") {
		return nil
	}
	return fmt.Errorf("Failed to cleanup local container: %s", err)
}
