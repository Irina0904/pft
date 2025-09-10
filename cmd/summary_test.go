package cmd

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

func Test_SummaryCmd_Success(t *testing.T) {
	tmpdir := t.TempDir()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	query := regexp.QuoteMeta(`SELECT category, SUM(amount) AS total 
	FROM transactions WHERE LTRIM(STRFTIME('%m', DATE(date)), '0') = ? AND category != 'salary'
	GROUP BY category;`)
	rows := sqlmock.NewRows([]string{"category", "total"}).
		AddRow("rent", 10000)

	mock.ExpectQuery(query).WithArgs("8").WillReturnRows(rows)

	filePath := filepath.Join(tmpdir, "summary_august.html")
	cmd := NewSummaryCmd(db, tmpdir)
	cmd.SetArgs([]string{"--month", "August"})
	cmd.Execute()

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		t.Errorf("file was not created: %s", filePath)
	} else if err != nil {
		t.Fatalf("error checking file status: %v", err)
	}
}

func Test_SummaryCmd_WrongMonth(t *testing.T) {
	tmpdir := t.TempDir()

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	filePath := filepath.Join(tmpdir, "summary_hello.html")
	cmd := NewSummaryCmd(db, tmpdir)
	cmd.SetArgs([]string{"--month", "Hello"})

	err = cmd.Execute()
	assert.Error(t, err)

	_, err = os.Stat(filePath)
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("file was created with wrong month input: %s", filePath)
	}
}

func Test_SummaryCmd_NoDataForMonth(t *testing.T) {
	tmpdir := t.TempDir()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	query := regexp.QuoteMeta(`SELECT category, SUM(amount) AS total 
	FROM transactions WHERE LTRIM(STRFTIME('%m', DATE(date)), '0') = ? AND category != 'salary'
	GROUP BY category;`)
	rows := sqlmock.NewRows([]string{})

	mock.ExpectQuery(query).WithArgs("4").WillReturnRows(rows)

	filePath := filepath.Join(tmpdir, "summary_april.html")
	cmd := NewSummaryCmd(db, tmpdir)
	cmd.SetArgs([]string{"--month", "april"})
	cmd.Execute()

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// We make sure that no summary is created when there is no data for the month 
	_, err = os.Stat(filePath)
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("file was created even though there was no data for the given month: %s", filePath)
	}
}
