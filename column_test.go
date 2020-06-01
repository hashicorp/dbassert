package dbassert

import (
	"testing"
)

func Test_Column(t *testing.T) {
	t.Parallel()
	cleanup, conn, _ := TestSetup(t, "postgres")
	defer func() {
		if err := cleanup(); err != nil {
			t.Error(err)
		}
		if err := conn.Close(); err != nil {
			t.Error(err)
		}
	}()
	t.Run("columns", func(t *testing.T) {
		assert := New(t, conn, "postgres")
		assert.Column(
			ColumnInfo{
				TableName:  "test_table_dbasserts",
				Name:       "public_id",
				Default:    "",
				Type:       "text",
				DomainName: "dbasserts_public_id",
				IsNullable: false,
			},
		)
	})
}
func Test_ColumnNullable(t *testing.T) {
	t.Parallel()
	cleanup, conn, _ := TestSetup(t, "postgres")
	defer func() {
		if err := cleanup(); err != nil {
			t.Error(err)
		}
		if err := conn.Close(); err != nil {
			t.Error(err)
		}
	}()
	t.Run("nullable", func(t *testing.T) {
		mockery := new(MockTesting)
		assert := New(mockery, conn, "postgres")
		assert.Nullable("test_table_dbasserts", "nullable")
		mockery.AssertNoError(t)
	})
	t.Run("typeInt", func(t *testing.T) {
		mockery := new(MockTesting)
		assert := New(mockery, conn, "postgres")
		assert.Nullable("test_table_dbasserts", "typeInt")
		mockery.AssertError(t)
	})
}

func Test_ColumnDomain(t *testing.T) {
	t.Parallel()
	cleanup, conn, _ := TestSetup(t, "postgres")
	defer func() {
		if err := cleanup(); err != nil {
			t.Error(err)
		}
		if err := conn.Close(); err != nil {
			t.Error(err)
		}
	}()
	t.Run("public_id", func(t *testing.T) {
		mockery := new(MockTesting)
		assert := New(mockery, conn, "postgres")
		assert.Domain("test_table_dbasserts", "public_id", "dbasserts_public_id")
		mockery.AssertNoError(t)
	})
	t.Run("nullable", func(t *testing.T) {
		mockery := new(MockTesting)
		assert := New(mockery, conn, "postgres")
		assert.Domain("test_table_dbasserts", "nullable", "dbasserts_public_id")
		mockery.AssertError(t)
	})
}
