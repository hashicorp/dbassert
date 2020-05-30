package gorm

import (
	"errors"
	"fmt"
	"strings"

	dbassert "github.com/hashicorp/dbassert"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

// FieldIsNull asserts that the modelFieldName is null in the db.
func (a *GormAsserts) FieldIsNull(model interface{}, modelFieldName string) bool {
	if h, ok := a.dbassert.T.(dbassert.THelper); ok {
		h.Helper()
	}

	colName, err := findColumnName(a.gormDb, model, modelFieldName)
	assert.NoError(a.dbassert.T, err)

	where := fmt.Sprintf("%s is null", colName)
	if err := a.gormDb.Where(where).First(model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			assert.NoError(a.dbassert.T, errors.New("field is not null"))
		}
		assert.NoError(a.dbassert.T, err)
		return false
	}
	return true
}

// FieldNotNull asserts that the modelFieldName is not null in the db.
func (a *GormAsserts) FieldNotNull(model interface{}, modelFieldName string) bool {
	if h, ok := a.dbassert.T.(dbassert.THelper); ok {
		h.Helper()
	}
	colName, err := findColumnName(a.gormDb, model, modelFieldName)
	assert.NoError(a.dbassert.T, err)
	where := fmt.Sprintf("%s is not null", colName)
	if err := a.gormDb.Where(where).First(model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			assert.NoError(a.dbassert.T, errors.New("field is null"))
		}
		assert.NoError(a.dbassert.T, err)
		return false
	}
	return true
}

// FieldNullable asserts that the modelFieldName nullable in the database.
func (a *GormAsserts) FieldNullable(model interface{}, modelFieldName string) bool {
	if h, ok := a.dbassert.T.(dbassert.THelper); ok {
		h.Helper()
	}
	colName, err := findColumnName(a.gormDb, model, modelFieldName)
	assert.NoError(a.dbassert.T, err)
	return a.dbassert.ColumnNullable(tableName(a.gormDb, model), colName)
}

// FieldDomain asserts that the modelFieldName is the domainName in the database.
func (a *GormAsserts) FieldDomain(model interface{}, modelFieldName, domainName string) bool {
	if h, ok := a.dbassert.T.(dbassert.THelper); ok {
		h.Helper()
	}
	colName, err := findColumnName(a.gormDb, model, modelFieldName)
	assert.NoError(a.dbassert.T, err)
	return a.dbassert.ColumnDomain(tableName(a.gormDb, model), colName, domainName)
}

func tableName(db *gorm.DB, model interface{}) string {
	return db.NewScope(model).TableName()
}

// findColumnName will find the model's db column name using the fieldName parameter
func findColumnName(db *gorm.DB, model interface{}, fieldName string) (string, error) {
	for _, f := range db.NewScope(model).GetStructFields() {
		if strings.EqualFold(fieldName, f.Name) {
			return f.DBName, nil
		}
	}
	return "", errors.New("modelFieldName not found in model")
}
