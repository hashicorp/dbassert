package gorm

import (
	"database/sql"

	dbassert "github.com/hashicorp/dbassert"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

// GormAsserts provides db assertions using Gorm as a db abstraction.
type GormAsserts struct {
	dbassert *dbassert.DbAsserts
	gormDb   *gorm.DB
}

// New will create a new GormAsserts.
func New(t dbassert.TestingT, db *sql.DB, dialect string) *GormAsserts {
	assert.NotNil(t, db, "db is nill")
	assert.NotEmpty(t, dialect, "dialect is not set")
	gormDb, err := gorm.Open(dialect, db)
	assert.NoError(t, err)

	return &GormAsserts{
		dbassert: &dbassert.DbAsserts{
			T:       t,
			Db:      db,
			Dialect: dialect,
		},
		gormDb: gormDb,
	}
}
func (a *GormAsserts) DbLog(enable bool) {
	a.gormDb.LogMode(enable)
}
