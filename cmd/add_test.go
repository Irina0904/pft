package cmd

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_AddCmd_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO transactions").
		WithArgs(time.Date(2025, 4, 5, 0, 0, 0, 0, time.Local), float64(4000), "salary", "test", "bank").
		WillReturnResult(sqlmock.NewResult(1, 1))

	cmd := NewAddCmd(db)
	cmd.SetArgs([]string{
		"--date", "05.04.2025",
		"--amount", "4000",
		"--category", "salary",
		"--description", "test",
		"--paymentmethod", "bank",
	})
	cmd.Execute()

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_AddCmd_WrongCategory(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cmd := NewAddCmd(db)
	cmd.SetArgs([]string{
		"--date", "05.04.2025",
		"--amount", "4000",
		"--category", "subscription",
		"--description", "test",
		"--paymentmethod", "bank",
	})

	err = cmd.Execute()
	assert.Error(t, err)
}

func Test_AddCmd_WrongPaymentMethod(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cmd := NewAddCmd(db)
	cmd.SetArgs([]string{
		"--date", "05.04.2025",
		"--amount", "4000",
		"--category", "rent",
		"--description", "test",
		"--paymentmethod", "swish",
	})

	err = cmd.Execute()
	assert.Error(t, err)
}

func Test_AddCmd_ParseDate(t *testing.T) {

	tests := []struct {
		name         string
		givenDate    string
		expectedDate time.Time
		expectError  bool
	}{
		{
			name:        "Correct date",
			givenDate:   "28.02.2025",
			expectedDate: time.Date(2025, 2, 28, 0, 0, 0, 0, time.Local),
			expectError: false,
		},
		{
			name:        "Wrong date format",
			givenDate:   "28/02/2025",
			expectError: true,
		},
		{
			name:        "Month out of range",
			givenDate:   "28.13.2025",
			expectError: true,
		},
		{
			name:        "Day out of range",
			givenDate:   "32.02.2025",
			expectError: true,
		},
		{
			name:        "Day out of range for February",
			givenDate:   "31.02.2025",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			cmd := NewAddCmd(db)
			cmd.SetArgs([]string{
				"--date", tc.givenDate,
				"--amount", "4000",
				"--category", "food",
				"--description", "test",
				"--paymentmethod", "bank",
			})

			if tc.expectError {
				err = cmd.Execute()
				assert.Error(t, err)
			} else {
				mock.ExpectExec("INSERT INTO transactions").
					WithArgs(
						tc.expectedDate,
						float64(4000),
						"food",
						"test",
						"bank").
					WillReturnResult(sqlmock.NewResult(1, 1))
				cmd.Execute()
				// we make sure that all expectations were met
				if err := mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			}
		})
	}
}
