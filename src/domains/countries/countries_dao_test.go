package countries

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/development-raul/footy-predictor/src/clients/mysql/footy_db"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestCountryDao_Create(t *testing.T) {
	testCases := []struct{
		title string
		funcMock func(sqlmock.Sqlmock)
		expectedErr error
	}{
		{
			title: "error Client.NamedExec",
			funcMock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("INSERT INTO countries").
					WithArgs(
						1,
						"code",
						"name",
						"flag",
						true).
					WillReturnError(errors.New("test NamedExec"))
			},
			expectedErr: errors.New("test NamedExec"),
		},
		{
			title: "error LastInsertId",
			funcMock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("INSERT INTO countries").
					WithArgs(
						1,
						"code",
						"name",
						"flag",
						true).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("test LastInsertId")))
			},
			expectedErr: errors.New("test LastInsertId"),
		},
		{
			title: "success",
			funcMock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("INSERT INTO countries").
					WithArgs(
						1,
						"code",
						"name",
						"flag",
						true).
					WillReturnResult(sqlmock.NewResult(1,1))
			},
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.title, func(t *testing.T){
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()
			footy_db.Client = sqlx.NewDb(db, "sqlmock")
			testCase.funcMock(mock)

			err = CountryDao.Create(&Country{
				ID:     1,
				Code:   "code",
				Name:   "name",
				Flag:   "flag",
				Active: true,
			})

			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}