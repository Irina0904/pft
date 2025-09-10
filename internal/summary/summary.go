package summary

import (
	"database/sql"
	"errors"
	"math"
	"strconv"
	"time"
	"fmt"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type Entry struct {
	category string
	total    float64
}

func generatePieItems(entries []Entry) []opts.PieData {
	items := make([]opts.PieData, 0)
	for _, entry := range entries {
		items = append(items, opts.PieData{Name: entry.category, Value: math.Abs(entry.total)})
	}
	return items
}

func PieSummary(entries []Entry, month string) *charts.Pie {
	pie := charts.NewPie()
	month = strings.ToUpper(month[:1]) + strings.ToLower(month[1:])
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: fmt.Sprintf("Your Spending in %s", month)}),
	)

	pie.AddSeries("pie", generatePieItems(entries))
	return pie
}

func GetSummary(month string, db *sql.DB) ([]Entry, error) {
	givenMonth, errParse := time.Parse("January", month)
	if errParse != nil {
		return []Entry{}, errors.New("Failed to parse full month, " + errParse.Error())
	}
	m := strconv.Itoa(int(givenMonth.Month()))

	rows, err := db.Query(
		`SELECT category, SUM(amount) AS total 
	 FROM transactions 
	 WHERE LTRIM(STRFTIME('%m', DATE(date)), '0') = ? AND category != 'salary' 
	 GROUP BY category;`, m)
	if err != nil {
		return []Entry{}, errors.New("failed to retrieve summary: " + err.Error())
	}
	defer rows.Close()

	var entries []Entry

	for rows.Next() {
		var entry Entry
		if err := rows.Scan(&entry.category, &entry.total); err != nil {
			return []Entry{}, errors.New("failed to scan entries: " + err.Error())
		}
		entries = append(entries, entry)
	}

	return entries, nil
}
