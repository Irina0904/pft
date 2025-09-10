package cmd

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
	"errors"

	"github.com/pft/internal/transaction"
	"github.com/spf13/cobra"
)

func NewAddCmd(db *sql.DB) *cobra.Command {
	addCmd := &cobra.Command {
		Use:   "add",
		Short: "Add a new transaction",
		Long: `
	Add a new transaction by specifying the transaction details`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Do Stuff Here
			newTransaction := transaction.Transaction{Date: time.Now()}
			date, _ := cmd.Flags().GetString("date")
			amount, _ := cmd.Flags().GetFloat64("amount")
			category, _ := cmd.Flags().GetString("category")
			description, _ := cmd.Flags().GetString("description")
			paymentMethod, _ := cmd.Flags().GetString("paymentmethod")


			if date != "" {
				transactionDate, err := parseDate(date)
				if err != nil {
					return fmt.Errorf("incorrect date value, %s ", err)
				}
				newTransaction.Date = transactionDate
			}
			if amount != 0 {
				newTransaction.Amount = amount
			}
			if category != "" {
				switch category {
				case "food":
					newTransaction.Category = transaction.Food
				case "salary":
					newTransaction.Category = transaction.Salary
				case "public transport":
					newTransaction.Category = transaction.PublicTransport
				case "rent":
					newTransaction.Category = transaction.Rent
				case "utilities":
					newTransaction.Category = transaction.Utilities
				default:
					categories := strings.Join(transaction.GetAvailableCategories(), ", ")
					return fmt.Errorf("wrong category selected, the available categories are: %s", categories)	
				}
			}
			if description != "" {
				newTransaction.Description = description
			}
			if paymentMethod != "" {
				switch paymentMethod {
				case "cash":
					newTransaction.PaymentMethod = transaction.Cash
				case "credit":
					newTransaction.PaymentMethod = transaction.Credit
				case "bank":
					newTransaction.PaymentMethod = transaction.Bank
				default:
					paymentMethods := strings.Join(transaction.GetAvailablePaymentMethods(), ", ")
					return fmt.Errorf(`wrong payment method selected, the available payment methods are: %s`, paymentMethods)	
				}
			}
			if db != nil {
				newTransaction.Add(db)
			}
			return nil
		},
	}

	addCmd.Flags().String("date", "", "Specify the date of the transaction in the format dd.mm.yyyy")
	addCmd.Flags().Float64("amount", 0, "Specify the transaction amount")
	addCmd.Flags().String("category", "", "Specify the transaction category")
	addCmd.Flags().String("description", "", "Specify the transaction description")
	addCmd.Flags().String("paymentmethod", "", "Specify the transaction payment method")

	return addCmd
}

func parseDate(date string) (time.Time, error) {
	values := strings.Split(date, ".")
	if len(values) != 3 {
		return time.Time{}, errors.New("wrong format. The format should be dd.mm.yyyy")
	}
	year, err := strconv.Atoi(values[2])
	if err != nil {
		return time.Time{}, err
	}
	month, err := strconv.Atoi(values[1])
	if err != nil {
		return time.Time{}, err
	} else if month <= 0 || month > 12 {
		return time.Time{}, errors.New("month is out of range. It should be an integer between 01 and 12")
	}
	day, err := strconv.Atoi(values[0])
	daysInMonth := getDaysInMonth(year, month)
	if err != nil {
		return time.Time{}, err
	} else if day <= 0 || day > daysInMonth {
		return time.Time{}, errors.New("day is out of range. It should be an integer between 01 and max days for this month")
	}

	dateTime := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	return dateTime, nil
}

func getDaysInMonth(year int, month int) int {
	t := time.Date(year, time.Month(month), 32, 0, 0, 0, 0, time.UTC)
	daysInMonth := 32 - t.Day()
	return daysInMonth
}
