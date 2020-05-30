package dbassert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_testSetup(t *testing.T) {
	assert := assert.New(t)
	cleanup, db, url := TestSetup(t, "postgres")
	defer func() {
		if err := cleanup(); err != nil {
			t.Error(err)
		}
		if err := db.Close(); err != nil {
			t.Error(err)
		}
	}()
	assert.NotNil(db)
	assert.NotNil(cleanup)
	assert.NotEmpty(url)
}
