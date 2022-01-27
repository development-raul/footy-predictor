package seasons

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/development-raul/footy-predictor/src/clients/mysql/footy_db"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSeasonDao_Create(t *testing.T) {
	testCases := []struct {
		title       string
		funcMock    func(sqlmock.Sqlmock)
		expectedErr error
	}{
		{
			title: "error Client.Exec",
			funcMock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("INSERT INTO seasons").
					WithArgs(1).
					WillReturnError(errors.New("test Exec"))
			},
			expectedErr: errors.New("test Exec"),
		},
		{
			title: "success",
			funcMock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("INSERT INTO seasons").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.title, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()
			footy_db.Client = sqlx.NewDb(db, "sqlmock")
			testCase.funcMock(mock)

			err = SeasonDao.Create(1)

			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}

func TestSeasonDao_Find(t *testing.T) {
	testCases := []struct {
		title       string
		funcMock    func(sqlmock.Sqlmock)
		expectedRes *Season
		expectedErr error
	}{
		{
			title: "error Client.Get",
			funcMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT (.+) FROM seasons").
					WithArgs(1).
					WillReturnError(errors.New("test Get"))
			},
			expectedErr: errors.New("test Get"),
		},
		{
			title: "success",
			funcMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT (.+) FROM seasons").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			expectedRes: &Season{
				ID: 1,
			},
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.title, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()
			footy_db.Client = sqlx.NewDb(db, "sqlmock")
			testCase.funcMock(mock)

			res, err := SeasonDao.Find(1)

			assert.Equal(t, testCase.expectedRes, res)
			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}

func TestSeasonDao_List(t *testing.T) {
	testCases := []struct {
		title       string
		funcMock    func(sqlmock.Sqlmock)
		expectedRes []Season
		expectedErr error
	}{
		{
			title: "error Client.Select",
			funcMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT (.+) FROM seasons").
					WithArgs(1).
					WillReturnError(errors.New("error Select"))
			},
			expectedErr: errors.New("error Select"),
		},
		{
			title: "success",
			funcMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT (.+) FROM seasons").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			expectedErr: nil,
			expectedRes: []Season{
				{
					ID: 1,
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.title, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()
			footy_db.Client = sqlx.NewDb(db, "sqlmock")
			testCase.funcMock(mock)

			res, err := SeasonDao.List(&ListSeasonInput{
				ID:    1,
				Order: "asc",
			})

			assert.Equal(t, testCase.expectedRes, res)
			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}

func TestSeasonDao_Delete(t *testing.T) {
	testCases := []struct {
		title       string
		funcMock    func(sqlmock.Sqlmock)
		expectedErr error
	}{
		{
			title: "error Client.Exec",
			funcMock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("DELETE FROM seasons").
					WithArgs(1).
					WillReturnError(errors.New("test Exec"))
			},
			expectedErr: errors.New("test Exec"),
		},
		{
			title: "success",
			funcMock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("DELETE FROM seasons").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.title, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()
			footy_db.Client = sqlx.NewDb(db, "sqlmock")
			testCase.funcMock(mock)

			err = SeasonDao.Delete(1)

			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}
