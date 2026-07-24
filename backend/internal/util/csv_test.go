package util_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/util"
)

func TestTransactionsToCSV(t *testing.T) {
	txs := []model.Transaction{
		{
			ID:           1,
			Date:         "2024-01-15",
			Type:         "expense",
			CategoryName: "Makanan & Minuman",
			Description:  "Makan siang \"spesial\"",
			Amount:       55000,
		},
		{
			ID:           2,
			Date:         "2024-01-16",
			Type:         "income",
			CategoryName: "Gaji",
			Description:  "Gaji Bulanan",
			Amount:       10000000,
		},
	}

	out := util.TransactionsToCSV(txs)

	// Check BOM
	if !bytes.HasPrefix(out, []byte("\xEF\xBB\xBF")) {
		t.Error("CSV output missing UTF-8 BOM prefix")
	}

	content := string(out)
	if !strings.Contains(content, "Makan siang \"\"spesial\"\"") {
		t.Errorf("CSV escaping failed: %s", content)
	}

	if !strings.Contains(content, "Pengeluaran") || !strings.Contains(content, "Pemasukan") {
		t.Errorf("CSV missing transaction type labels: %s", content)
	}
}

func TestTransactionsToCSVSummary(t *testing.T) {
	metrics := []model.MonthlyMetric{
		{MonthName: "Januari 2024", Year: 2024, Income: 10000, Expense: 4000},
	}

	out := util.TransactionsToCSVSummary(metrics)
	content := string(out)

	if !strings.Contains(content, "Januari 2024,2024,10000,4000,6000") {
		t.Errorf("CSV summary content invalid: %s", content)
	}
}
