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
	testCases := []struct {
		title       string
		funcMock    func(sqlmock.Sqlmock)
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

			err = CountryDao.Create(&Country{
				AsID:   1,
				Code:   "code",
				Name:   "name",
				Flag:   "flag",
				Active: true,
			})

			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}

func TestCountryDao_Update(t *testing.T) {
	testCases := []struct {
		title       string
		funcMock    func(sqlmock.Sqlmock)
		expectedErr error
	}{
		{
			title: "error Client.NamedExec",
			funcMock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("UPDATE countries SET").
					WithArgs(
						"code",
						"name",
						"flag",
						true,
						1).
					WillReturnError(errors.New("test NamedExec"))
			},
			expectedErr: errors.New("test NamedExec"),
		},
		{
			title: "success",
			funcMock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("UPDATE countries SET").
					WithArgs(
						"code",
						"name",
						"flag",
						true,
						1).
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

			err = CountryDao.Update(&UpdateCountryInput{
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

func TestCountryDao_FindByID(t *testing.T) {
	testCases := []struct {
		title       string
		funcMock    func(sqlmock.Sqlmock)
		expectedRes *CountryOutput
		expectedErr error
	}{
		{
			title: "error Client.Get",
			funcMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT (.+) FROM countries").
					WithArgs(1).
					WillReturnError(errors.New("test Get"))
			},
			expectedErr: errors.New("test Get"),
		},
		{
			title: "success",
			funcMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT (.+) FROM countries").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{
						"id",
						"as_id",
						"code",
						"name",
						"flag",
						"active",
					}).AddRow(
						1,
						1,
						"code",
						"name",
						"flag",
						1,
					))
			},
			expectedRes: &CountryOutput{
				ID:     1,
				AsID:   1,
				Code:   "code",
				Name:   "name",
				Flag:   "flag",
				Active: true,
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

			res, err := CountryDao.FindByID(1)

			assert.Equal(t, testCase.expectedRes, res)
			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}

func TestCountryDao_FindByAsID(t *testing.T) {
	testCases := []struct {
		title       string
		funcMock    func(sqlmock.Sqlmock)
		expectedRes *CountryOutput
		expectedErr error
	}{
		{
			title: "error Client.Get",
			funcMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT (.+) FROM countries").
					WithArgs(1).
					WillReturnError(errors.New("test Get"))
			},
			expectedErr: errors.New("test Get"),
		},
		{
			title: "success",
			funcMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT (.+) FROM countries").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{
						"id",
						"as_id",
						"code",
						"name",
						"flag",
						"active",
					}).AddRow(
						1,
						1,
						"code",
						"name",
						"flag",
						1,
					))
			},
			expectedRes: &CountryOutput{
				ID:     1,
				AsID:   1,
				Code:   "code",
				Name:   "name",
				Flag:   "flag",
				Active: true,
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

			res, err := CountryDao.FindByAsID(1)

			assert.Equal(t, testCase.expectedRes, res)
			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}

func TestCountryDao_List(t *testing.T) {
	testCases := []struct {
		title         string
		funcMock      func(sqlmock.Sqlmock)
		expectedRes   []CountryOutput
		expectedTotal int64
		expectedErr   error
	}{
		{
			title: "error Client.Select",
			funcMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT (.+) FROM countries").
					WithArgs("code", "%name%").
					WillReturnError(errors.New("error Select"))
			},
			expectedErr: errors.New("error Select"),
		},
		{
			title: "error GetTableTotalRowsArgs",
			funcMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT (.+) FROM countries").
					WithArgs("code", "%name%").
					WillReturnRows(sqlmock.NewRows([]string{
						"id",
						"as_id",
						"code",
						"name",
						"flag",
						"active",
					}).AddRow(
						1,
						1,
						"code",
						"name",
						"flag",
						1,
					))
				m.ExpectQuery("SELECT (.+) FROM countries").
					WithArgs("code", "%name%").
					WillReturnError(errors.New("error GetTableTotalRowsArgs"))
			},
			expectedErr: errors.New("error GetTableTotalRowsArgs"),
		},
		{
			title: "success",
			funcMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT (.+) FROM countries").
					WithArgs("code", "%name%").
					WillReturnRows(sqlmock.NewRows([]string{
						"id",
						"as_id",
						"code",
						"name",
						"flag",
						"active",
					}).AddRow(
						1,
						1,
						"code",
						"name",
						"flag",
						1,
					))
				m.ExpectQuery("SELECT (.+) FROM countries").
					WithArgs("code", "%name%").
					WillReturnRows(sqlmock.NewRows([]string{
						"total",
					}).AddRow(1))
			},
			expectedErr: nil,
			expectedRes: []CountryOutput{
				{
					ID:     1,
					AsID:   1,
					Code:   "code",
					Name:   "name",
					Flag:   "flag",
					Active: true,
				},
			},
			expectedTotal: 1,
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

			res, total, err := CountryDao.List(&ListCountryInput{
				Code:   "code",
				Name:   "name",
				Active: true,
			})

			assert.Equal(t, testCase.expectedRes, res)
			assert.Equal(t, testCase.expectedTotal, total)
			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}

func TestCountryDao_Delete(t *testing.T) {
	testCases := []struct {
		title       string
		funcMock    func(sqlmock.Sqlmock)
		expectedErr error
	}{
		{
			title: "error Client.Exec",
			funcMock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("DELETE FROM countries").
					WithArgs(1).
					WillReturnError(errors.New("test Exec"))
			},
			expectedErr: errors.New("test Exec"),
		},
		{
			title: "success",
			funcMock: func(m sqlmock.Sqlmock) {
				m.ExpectExec("DELETE FROM countries").
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

			err = CountryDao.Delete(1)

			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}
