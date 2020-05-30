# dbassert

The `dbassert` package provides some helpful functions to help you write better
tests when writing Go database applications.  The package supports both sql.DB
and Gorm assertions. 

### Example sql.DB asserts usage:

```go
package your_brilliant_pkg

import (
    "testing"
    dbassert "github.com/hashicorp/dbassert"
)

func TestSomeDb(t *testing.T) {
	conn, err := sql.Open("postgres", "postgres://postgres:secret@localhost:%s?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	
	dbassert := dbassert.New(t, conn, "postgres")
    
	// assert that the db column is nullable
	dbassert.ColumnNullable("some_table", "some_column")

	// assert that the db column is a particular domain type
	dbassert.ColumnDomain("test_table_dbasserts", "public_id", "dbasserts_public_id")

}
```
### Example Gorm asserts usage:

```go
package your_brilliant_pkg

import (
    "testing"
    dbassert "github.com/hashicorp/dbassert/gorm"
)

func TestSomeGormModel(t *testing.T) {
	conn, err := sql.Open("postgres", "postgres://postgres:secret@localhost:%s?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	db, err := gorm.Open("postgres", conn)
 	m := testModel{}
	if err = db.Create(&m).Error; err != nil {
    	t.Fatal(err)
	}
	dbassert := dbassert.New(t, conn, "postgres")
    
	// assert that the db field is null
	dbassert.FieldIsNull(&someModel, "SomeField")

	// assert that the db field is not null
	dbassert.FieldNotNull(&someModel, "SomeField")

	// assert that the db field nullable
	dbassert.FieldNullable(&someModel, "SomeField")

	// assert that the db field is a particular domain type
	dbassert.FieldDomain(&someModel, "SomeField", "some_domain_type")
}
```

