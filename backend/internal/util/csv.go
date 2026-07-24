package util

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/akrom/finance-backend/internal/model"
)

// csvEscape wraps a string in quotes and escapes internal double-quotes per RFC 4180.
func csvEscape(s string) string {
	s = strings.ReplaceAll(s, `"`, `""`)
	return `"` + s + `"`
}

// TransactionsToCSV converts a slice of transactions to a UTF-8 CSV byte slice.
// The first row is a header row. No external library is required.
func TransactionsToCSV(transactions []model.Transaction) []byte {
	var buf bytes.Buffer

	// BOM for Excel compatibility
	buf.WriteString("\xEF\xBB\xBF")

	// Header
	headers := []string{"ID", "Tanggal", "Jenis", "Kategori", "Deskripsi", "Jumlah (Rp)"}
	buf.WriteString(strings.Join(headers, ","))
	buf.WriteString("\r\n")

	// Rows
	for _, t := range transactions {
		typeLabel := "Pemasukan"
		if t.Type == "expense" {
			typeLabel = "Pengeluaran"
		}
		row := []string{
			fmt.Sprintf("%d", t.ID),
			csvEscape(t.Date),
			csvEscape(typeLabel),
			csvEscape(t.CategoryName),
			csvEscape(t.Description),
			fmt.Sprintf("%d", t.Amount),
		}
		buf.WriteString(strings.Join(row, ","))
		buf.WriteString("\r\n")
	}

	return buf.Bytes()
}

// TransactionsToCSVSummary returns a CSV byte slice that contains a monthly
// summary: month, total_income, total_expense, balance.
func TransactionsToCSVSummary(metrics []model.MonthlyMetric) []byte {
	var buf bytes.Buffer

	buf.WriteString("\xEF\xBB\xBF")
	buf.WriteString("Bulan,Tahun,Total Pemasukan (Rp),Total Pengeluaran (Rp),Saldo (Rp)\r\n")

	for _, m := range metrics {
		balance := m.Income - m.Expense
		row := fmt.Sprintf("%s,%d,%d,%d,%d\r\n",
			m.MonthName, m.Year, m.Income, m.Expense, balance)
		buf.WriteString(row)
	}

	return buf.Bytes()
}
