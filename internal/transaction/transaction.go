package transaction

import (
	"database/sql"
	"log"
	"time"
)

type Category string

const (
	Food            Category = "food"
	Salary          Category = "salary"
	PublicTransport Category = "public transport"
	Rent            Category = "rent"
	Utilities       Category = "utilities"
)

type PaymentMethod string

const (
	Cash   PaymentMethod = "cash"
	Credit PaymentMethod = "credit"
	Bank   PaymentMethod = "bank"
)

type Transaction struct {
	Date          time.Time     `json:"date"`
	Amount        float64       `json:"amount"`
	Category      Category      `json:"category"`
	Description   string        `json:"description"`
	PaymentMethod PaymentMethod `json:"payment_method"`
}

func (t *Transaction) Add(db *sql.DB) {

	_, err := db.Exec(`INSERT INTO transactions 
	(date, 
	amount, 
	category, 
	description, 
	payment_method) VALUES (?, ?, ?, ?, ?)`,
		t.Date, t.Amount, t.Category, t.Description, t.PaymentMethod)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Transaction added!")
}

func GetAvailableCategories() []string {
	return []string{"food", "salary", "public transport", "rent", "utilities"}
}

func GetAvailablePaymentMethods() []string {
	return []string{"cash", "credit", "bank"}
}
