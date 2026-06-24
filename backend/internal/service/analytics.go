package service

import (
	"fmt"
	"math"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
)

func GetAdvancedAnalytics(userID int) (*model.AdvancedAnalyticsResponse, error) {
	now := time.Now()
	currentMonth := int(now.Month())
	currentYear := now.Year()

	// 1. Fetch historical metrics (last 6 months)
	monthlyMetrics, err := repository.GetHistoricalMetrics(userID, 6)
	if err != nil {
		return nil, fmt.Errorf("service: get historical metrics: %w", err)
	}

	// 2. Fetch budget comparisons for current month
	budgetComps, err := repository.GetBudgetComparison(userID, currentMonth, currentYear)
	if err != nil {
		return nil, fmt.Errorf("service: get budget comparison: %w", err)
	}

	// 3. Fetch top categories for current month
	topExpenses, err := repository.GetTopExpenses(userID, currentMonth, currentYear)
	if err != nil {
		return nil, fmt.Errorf("service: get top expenses: %w", err)
	}

	// 4. Calculate Current Month Metrics
	var currentMetric model.MonthlyMetric
	for _, m := range monthlyMetrics {
		if m.Month == currentMonth && m.Year == currentYear {
			currentMetric = m
			break
		}
	}

	// Savings Rate
	var savingsRate float64
	if currentMetric.Income > 0 {
		savingsRate = float64(currentMetric.Income-currentMetric.Expense) / float64(currentMetric.Income) * 100
	}

	// Budget Adherence
	var budgetAdherence float64 = 100.0
	var budgetedCategoriesCount int
	var withinBudgetCount int
	var overBudgetCategories []string

	for _, bc := range budgetComps {
		budgetedCategoriesCount++
		if bc.ActualAmount <= bc.BudgetAmount {
			withinBudgetCount++
		} else {
			// Find category name if needed
			overBudgetCategories = append(overBudgetCategories, fmt.Sprintf("Kategori %d", bc.CategoryID))
		}
	}

	if budgetedCategoriesCount > 0 {
		budgetAdherence = (float64(withinBudgetCount) / float64(budgetedCategoriesCount)) * 100
	}

	// 5. Forecast next month's income and expense (Weighted Moving Average)
	var forecastIncome, forecastExpense int64
	n := len(monthlyMetrics)
	if n > 0 {
		var sumExpenseWeight, sumIncomeWeight float64
		var totalExpenseWeighted, totalIncomeWeighted float64

		// Weights for up to past 3 months: 0.5 (most recent), 0.3, 0.2
		weights := []float64{0.5, 0.3, 0.2}
		weightIdx := 0

		for i := n - 1; i >= 0 && weightIdx < len(weights); i-- {
			w := weights[weightIdx]
			totalExpenseWeighted += float64(monthlyMetrics[i].Expense) * w
			totalIncomeWeighted += float64(monthlyMetrics[i].Income) * w
			sumExpenseWeight += w
			sumIncomeWeight += w
			weightIdx++
		}

		if sumExpenseWeight > 0 {
			forecastExpense = int64(math.Round(totalExpenseWeighted / sumExpenseWeight))
		}
		if sumIncomeWeight > 0 {
			forecastIncome = int64(math.Round(totalIncomeWeighted / sumIncomeWeight))
		}
	}

	// 6. Calculate Financial Health Score (0-100)
	healthScore := 50 // Base score

	// Savings Rate contribution (max +30 points, min -20 points)
	if savingsRate >= 20 {
		healthScore += 30
	} else if savingsRate >= 10 {
		healthScore += 20
	} else if savingsRate >= 0 {
		healthScore += 10
	} else if savingsRate < 0 {
		healthScore -= 20
	}

	// Budget Adherence contribution (max +20 points)
	healthScore += int(math.Round(budgetAdherence * 0.20))

	// Historical stability (max +10 points)
	hasNegativeMonth := false
	for _, m := range monthlyMetrics {
		if m.Income < m.Expense {
			hasNegativeMonth = true
		}
	}
	if !hasNegativeMonth {
		healthScore += 10
	}

	// Ensure health score is within bounds
	if healthScore > 100 {
		healthScore = 100
	}
	if healthScore < 0 {
		healthScore = 0
	}

	// Health Rating text
	var healthRating string
	if healthScore >= 80 {
		healthRating = "Sangat Sehat"
	} else if healthScore >= 60 {
		healthRating = "Sehat"
	} else if healthScore >= 40 {
		healthRating = "Cukup Sehat"
	} else {
		healthRating = "Perlu Perhatian"
	}

	// 7. Generate Actionable Insights
	var insights []string

	if healthScore >= 80 {
		insights = append(insights, "Kondisi keuangan Anda luar biasa! Teruskan kedisiplinan pencatatan dan pengelolaan ini.")
	}

	if savingsRate < 10 {
		insights = append(insights, fmt.Sprintf("Rasio tabungan Anda (%.1f%%) berada di bawah target aman 10%%. Pertimbangkan untuk meminimalkan pengeluaran keinginan.", savingsRate))
	} else if savingsRate >= 20 {
		insights = append(insights, fmt.Sprintf("Hebat! Anda menabung %.1f%% dari pemasukan Anda bulan ini. Ini adalah langkah cepat menuju kebebasan finansial.", savingsRate))
	}

	if budgetAdherence < 80 {
		insights = append(insights, fmt.Sprintf("Anda telah melampaui batas anggaran pada %.0f%% kategori anggaran Anda. Kurangi pengeluaran non-esensial.", 100-budgetAdherence))
	}

	if forecastExpense > forecastIncome {
		insights = append(insights, fmt.Sprintf("Peringatan: Proyeksi pengeluaran bulan depan (Rp %d) melebihi proyeksi pemasukan (Rp %d). Silakan rencanakan pemotongan anggaran.", forecastExpense, forecastIncome))
	}

	if len(topExpenses) > 0 {
		insights = append(insights, fmt.Sprintf("Kategori pengeluaran terbesar Anda bulan ini adalah '%s' senilai Rp %d. Review apakah ini benar-benar kebutuhan utama.", topExpenses[0].CategoryName, topExpenses[0].TotalAmount))
	}

	return &model.AdvancedAnalyticsResponse{
		HealthScore:     healthScore,
		HealthRating:    healthRating,
		SavingsRate:     savingsRate,
		BudgetAdherence: budgetAdherence,
		ForecastExpense: forecastExpense,
		ForecastIncome:  forecastIncome,
		Insights:        insights,
		MonthlyMetrics:  monthlyMetrics,
	}, nil
}
