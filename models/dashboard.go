package models

type ChartData struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Color string  `json:"color"`
}

type Dashboard struct {
	TotalIncome      float64     `json:"totalIncome"`
	TotalExpense     float64     `json:"totalExpense"`
	TotalSavings     float64     `json:"totalSavings"`
	RemainingBalance float64     `json:"remainingBalance"`
	ExpenseBreakdown []ChartData `json:"expenseBreakdown"`
	SavingsBreakdown []ChartData `json:"savingsBreakdown"`
	IncomeBreakdown  []ChartData `json:"incomeBreakdown"`
}
