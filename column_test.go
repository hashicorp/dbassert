package dbassert

import (
	"testing"
)

func TestDbAsserts_Column(t *testing.T) {
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
	cases := []struct {
		name   string
		column ColumnInfo
		want   bool
	}{
		{
			name: "nullable",
			column: ColumnInfo{
				TableName:  "test_table_dbasserts",
				Name:       "public_id",
				Default:    "",
				Type:       "text",
				DomainName: "dbasserts_public_id",
				IsNullable: false,
			},
			want: true,
		},
		{
			name: "bad type",
			column: ColumnInfo{
				TableName:  "test_table_dbasserts",
				Name:       "public_id",
				Default:    "",
				Type:       "bad type",
				DomainName: "dbasserts_public_id",
				IsNullable: false,
			},
			want: false,
		},
		{
			name: "bad column name",
			column: ColumnInfo{
				TableName:  "test_table_dbasserts",
				Name:       "bad_column_name",
				Default:    "",
				Type:       "text",
				DomainName: "dbasserts_public_id",
				IsNullable: false,
			},
			want: false,
		},
		{
			name: "bad table name",
			column: ColumnInfo{
				TableName:  "bad_name",
				Name:       "public_id",
				Default:    "",
				Type:       "text",
				DomainName: "dbasserts_public_id",
				IsNullable: false,
			},
			want: false,
		},
		{
			name: "bad domain",
			column: ColumnInfo{
				TableName:  "test_table_dbasserts",
				Name:       "public_id",
				Default:    "",
				Type:       "text",
				DomainName: "bad_domain",
				IsNullable: false,
			},
			want: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			mockery := new(MockTesting)
			a := New(mockery, conn, "postgres")

			if got := a.Column(tt.column); got != tt.want {
				t.Errorf("Column() = %v, want %v", got, tt.want)
			}
			switch {
			case tt.want:
				mockery.AssertNoError(t)
			default:
				mockery.AssertError(t)
			}
		})
	}
}
func TestDbAsserts_Domain(t *testing.T) {
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
	cases := []struct {
		name      string
		tableName string
		colName   string
		domain    string
		want      bool
	}{
		{
			name:      "nullable",
			tableName: "test_table_dbasserts",
			colName:   "nullable",
			domain:    "dbasserts_public_id",
			want:      false,
		},
		{
			name:      "public_id",
			tableName: "test_table_dbasserts",
			colName:   "public_id",
			domain:    "dbasserts_public_id",
			want:      true,
		},
		{
			name:      "bad_column",
			tableName: "test_table_dbasserts",
			colName:   "bad_column",
			domain:    "dbasserts_public_id",
			want:      false,
		},
		{
			name:      "bad_table",
			tableName: "bad_table",
			colName:   "public_id",
			domain:    "dbasserts_public_id",
			want:      false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			mockery := new(MockTesting)
			a := New(mockery, conn, "postgres")

			if got := a.Domain(tt.tableName, tt.colName, tt.domain); got != tt.want {
				t.Errorf("Domain() = %v, want %v", got, tt.want)
			}
			switch {
			case tt.want:
				mockery.AssertNoError(t)
			default:
				mockery.AssertError(t)
			}
		})
	}
}

func TestDbAsserts_Nullable(t *testing.T) {
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
	cases := []struct {
		name      string
		tableName string
		colName   string
		want      bool
	}{
		{
			name:      "nullable",
			tableName: "test_table_dbasserts",
			colName:   "nullable",
			want:      true,
		},
		{
			name:      "typeInt-bad-colname",
			tableName: "test_table_dbasserts",
			colName:   "typeInt",
			want:      false,
		},
		{
			name:      "type_int",
			tableName: "test_table_dbasserts",
			colName:   "type_int",
			want:      true,
		},
		{
			name:      "public_id",
			tableName: "test_table_dbasserts",
			colName:   "public_id",
			want:      false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			mockery := new(MockTesting)
			a := New(mockery, conn, "postgres")

			if got := a.Nullable(tt.tableName, tt.colName); got != tt.want {
				t.Errorf("Nullable() = %v, want %v", got, tt.want)
			}
			switch {
			case tt.want:
				mockery.AssertNoError(t)
			default:
				mockery.AssertError(t)
			}
		})
	}
}
