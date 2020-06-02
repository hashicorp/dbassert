package gorm

import (
	"testing"

	dbassert "github.com/hashicorp/dbassert"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestDbAsserts_Domain(t *testing.T) {
	t.Parallel()
	cleanup, conn, _ := dbassert.TestSetup(t, "postgres")
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
		model     interface{}
		fieldName string
		domain    string
		want      bool
	}{
		{
			name:      "nullable",
			model:     &TestModel{},
			fieldName: "nullable",
			domain:    "dbasserts_public_id",
			want:      false,
		},
		{
			name:      "public_id",
			model:     &TestModel{},
			fieldName: "PublicId",
			domain:    "dbasserts_public_id",
			want:      true,
		},
		{
			name:      "bad_column",
			model:     &TestModel{},
			fieldName: "BadField",
			domain:    "dbasserts_public_id",
			want:      false,
		},
		{
			name:      "nil_model",
			model:     nil,
			fieldName: "public_id",
			domain:    "dbasserts_public_id",
			want:      false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			mockery := new(dbassert.MockTesting)
			a := New(mockery, conn, "postgres")

			if got := a.Domain(tt.model, tt.fieldName, tt.domain); got != tt.want {
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
	cleanup, conn, _ := dbassert.TestSetup(t, "postgres")
	defer func() {
		if err := cleanup(); err != nil {
			t.Error(err)
		}
		if err := conn.Close(); err != nil {
			t.Error(err)
		}
	}()
	cases := []struct {
		name  string
		model interface{}
		fName string
		want  bool
	}{
		{
			name:  "nullable",
			model: &TestModel{},
			fName: "nullable",
			want:  true,
		},
		{
			name:  "badField",
			model: &TestModel{},
			fName: "badField",
			want:  false,
		},
		{
			name:  "TypeInt",
			model: &TestModel{},
			fName: "TypeInt",
			want:  true,
		},
		{
			name:  "PublicId",
			model: &TestModel{},
			fName: "PublicId",
			want:  false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			mockery := new(dbassert.MockTesting)
			a := New(mockery, conn, "postgres")

			if got := a.Nullable(tt.model, tt.fName); got != tt.want {
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

func TestDbAsserts_IsNull(t *testing.T) {
	t.Parallel()
	cleanup, conn, _ := dbassert.TestSetup(t, "postgres")
	defer func() {
		if err := cleanup(); err != nil {
			t.Error(err)
		}
		if err := conn.Close(); err != nil {
			t.Error(err)
		}
	}()

	db, err := gorm.Open("postgres", conn)
	assert.NoError(t, err)

	v := 1
	m := CreateTestModel(t, db, nil, &v)

	cases := []struct {
		name  string
		model interface{}
		fName string
		want  bool
	}{
		{
			name:  "nullable",
			model: m,
			fName: "nullable",
			want:  true,
		},
		{
			name:  "badField",
			model: m,
			fName: "badField",
			want:  false,
		},
		{
			name:  "TypeInt",
			model: m,
			fName: "TypeInt",
			want:  false,
		},
		{
			name:  "nil model",
			model: nil,
			fName: "PublicId",
			want:  false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			mockery := new(dbassert.MockTesting)
			a := New(mockery, conn, "postgres")

			if got := a.IsNull(tt.model, tt.fName); got != tt.want {
				t.Errorf("IsNull() = %v, want %v", got, tt.want)
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

func TestDbAsserts_NotNull(t *testing.T) {
	t.Parallel()
	cleanup, conn, _ := dbassert.TestSetup(t, "postgres")
	defer func() {
		if err := cleanup(); err != nil {
			t.Error(err)
		}
		if err := conn.Close(); err != nil {
			t.Error(err)
		}
	}()

	db, err := gorm.Open("postgres", conn)
	assert.NoError(t, err)

	v := 1
	m := CreateTestModel(t, db, nil, &v)

	cases := []struct {
		name  string
		model interface{}
		fName string
		want  bool
	}{
		{
			name:  "nullable",
			model: m,
			fName: "nullable",
			want:  false,
		},
		{
			name:  "badField",
			model: m,
			fName: "badField",
			want:  false,
		},
		{
			name:  "TypeInt",
			model: m,
			fName: "TypeInt",
			want:  true,
		},
		{
			name:  "nil model",
			model: nil,
			fName: "PublicId",
			want:  false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			mockery := new(dbassert.MockTesting)
			a := New(mockery, conn, "postgres")

			if got := a.NotNull(tt.model, tt.fName); got != tt.want {
				t.Errorf("NotNull() = %v, want %v", got, tt.want)
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
