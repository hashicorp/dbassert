package dbassert

import (
	"database/sql"
	"strings"

	"github.com/stretchr/testify/assert"
)

// ColumnInfo defines a set of information about a column
type ColumnInfo struct {
	// TableName for the column
	TableName string

	// Name of the column
	Name string

	// Default value for the column
	Default string

	// Type of the column
	Type string

	// DomainName for the column
	DomainName string

	// IsNullable defines if the column can be null
	IsNullable bool
}

// ColumnNullable asserts colName in tableName is nullable
func (a *DbAsserts) ColumnNullable(tableName, colName string) bool {
	if h, ok := a.T.(THelper); ok {
		h.Helper()
	}
	dbColumn := a.getSchemaInfo(tableName, colName)
	if dbColumn.IsNullable {
		return true
	}
	assert.Fail(a.T, "column is not nullable", "%s: %s is not nullable", tableName, colName)
	return false
}

// ColumnDomain asserts colName in tableName is domainName
func (a *DbAsserts) ColumnDomain(tableName, colName, domainName string) bool {
	if h, ok := a.T.(THelper); ok {
		h.Helper()
	}
	dbColumn := a.getSchemaInfo(tableName, colName)
	if strings.EqualFold(domainName, dbColumn.DomainName) {
		return true
	}
	assert.Fail(a.T, "domain is not valid", "%s: %s is not %s", tableName, colName, domainName)
	return false
}

// Column asserts c ColumnInfo is valid
func (a *DbAsserts) Column(c ColumnInfo) bool {
	if h, ok := a.T.(THelper); ok {
		h.Helper()
	}
	dbColumn := a.getSchemaInfo(c.TableName, c.Name)
	if c != dbColumn {
		assert.Fail(a.T, "invalid column", "%s: %+v column is not valid in the db columns %+v", c.TableName, c, dbColumn)
		return false
	}
	return true
}

func (a *DbAsserts) getSchemaInfo(tableName, columnName string) ColumnInfo {
	if h, ok := a.T.(THelper); ok {
		h.Helper()
	}
	const query = `
select 
	table_name, 
	column_name, 
	column_default, 
	data_type, 
	domain_name,
	is_nullable
from information_schema.columns
where table_name = $1 and column_name = $2`
	row := a.Db.QueryRow(query, tableName, columnName)

	var table, colName, colType, colIsNullable string
	var colDefault, colDomainName sql.NullString
	err := row.Scan(&table, &colName, &colDefault, &colType, &colDomainName, &colIsNullable)
	assert.NoError(a.T, err)

	var nullable bool
	if colIsNullable == "YES" {
		nullable = true
	}
	return ColumnInfo{
		TableName:  tableName,
		Name:       colName,
		Default:    NullableString(colDefault),
		Type:       colType,
		DomainName: NullableString(colDomainName),
		IsNullable: nullable,
	}
}

// NullableString is a type alias for nullable database columns for strings.
func NullableString(value sql.NullString) string {
	if !value.Valid {
		return ""
	}
	return value.String
}
