// Package dbassert provides a set of assertions for testing Go database
// applications.
//
// Example Usage:
//
// import (
// 	"testing"
//
// 	dbassert "github.com/hashicorp/dbassert/gorm"
// )
//
// func TestSomeDatabase(t *testing.T) {
// 	conn, err := sql.Open("postgres", "postgres://postgres:secret@localhost:%s?sslmode=disable")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer conn.Close()
//   	db, err := gorm.Open("postgres", conn)
//   	m := testModel{}
//   	if err = db.Create(&m).Error; err != nil {
//     		t.Fatal(err)
//   	}
// 	dbassert := dbassert.New(t, conn, "postgres")
// 	dbassert.FieldIsNull(&someModel, "someField")
// }

package gorm
