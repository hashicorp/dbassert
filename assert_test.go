package dbassert

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
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
	mockery := new(MockTesting)
	type args struct {
		t       TestingT
		db      *sql.DB
		dialect string
	}
	tests := []struct {
		name    string
		args    args
		want    *DbAsserts
		wantErr bool
	}{
		{
			name: "postgres",
			args: args{
				t:       mockery,
				db:      conn,
				dialect: "postgres",
			},
			want: &DbAsserts{
				T:       mockery,
				Db:      conn,
				Dialect: "postgres",
			},
		},
		{
			name: "unsupported-dialect",
			args: args{
				t:       mockery,
				db:      conn,
				dialect: "mysql",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "nil-db",
			args: args{
				t:       mockery,
				db:      nil,
				dialect: "postgres",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "nil-t",
			args: args{
				t:       nil,
				db:      conn,
				dialect: "postgres",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer mockery.Reset()
			got := func() *DbAsserts {
				defer func() {
					if r := recover(); r != nil {
						if err, ok := r.(error); ok {
							if errors.Is(err, ErrNilTestingT) {
								mockery.err = true
							}
						}
					}
				}()
				return New(tt.args.t, tt.args.db, tt.args.dialect)
			}()
			if tt.wantErr {
				mockery.AssertError(t)
			} else {
				mockery.AssertNoError(t)
			}
			assert.Equal(t, got, tt.want)
		})
	}
}
