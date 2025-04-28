package dashboard

import (
	"fmt"
	"gofr.dev/pkg/gofr"
	"math/rand"
	"moneyManagement/filters"
	"moneyManagement/models"
	"moneyManagement/services"
	"time"
)

type dashboardService struct {
	accountSvc      services.Account
	transactionsSvc services.Transactions
	userSvc         services.User
}

func New(accountSvc services.Account, transactionsSvc services.Transactions, userSvc services.User) services.Dashboard {
	return &dashboardService{accountSvc: accountSvc, transactionsSvc: transactionsSvc, userSvc: userSvc}
}

func (s *dashboardService) Get(ctx *gofr.Context, f *filters.Transactions) (models.Dashboard, error) {
	var dashboard models.Dashboard

	userID, _ := ctx.Value("userID").(int)

	f.UserID = userID

	transactions, err := s.transactionsSvc.GetAll(ctx, f)
	if err != nil {
		return models.Dashboard{}, err
	}

	account, err := s.accountSvc.GetByID(ctx, f.AccountID)
	if err != nil {
		return models.Dashboard{}, err
	}

	expenseMap := make(map[string]float64)
	incomeMap := make(map[string]float64)
	savingsMap := make(map[string]float64)

	for _, txn := range transactions {
		switch txn.Type {
		case models.EXPENSE:
			dashboard.TotalExpense += txn.Amount
			expenseMap[txn.Category] += txn.Amount
		case models.INCOME:
			dashboard.TotalIncome += txn.Amount
			incomeMap[txn.Category] += txn.Amount
		case models.SAVINGS:
			dashboard.TotalSavings += txn.Amount
			savingsMap[txn.Category] += txn.Amount
		}
	}

	dashboard.RemainingBalance = account.Balance

	dashboard.ExpenseBreakdown = mapToChartData(expenseMap)
	dashboard.IncomeBreakdown = mapToChartData(incomeMap)
	dashboard.SavingsBreakdown = mapToChartData(savingsMap)

	return dashboard, nil
}

func mapToChartData(data map[string]float64) []models.ChartData {
	var chartData []models.ChartData

	var assignedColors = make(map[string]string)
	var usedColors = make(map[string]bool)

	for category, value := range data {
		color := getCategoryColor(category, assignedColors, usedColors)
		chartData = append(chartData, models.ChartData{
			Name:  category,
			Value: value,
			Color: color,
		})
	}

	return chartData
}

// Generates a unique pastel color in HEX
func getUniquePastelColor(usedColors map[string]bool) string {
	rand.Seed(time.Now().UnixNano())

	for {
		h := rand.Intn(360)
		s, v := 0.5, 0.9
		r, g, b := hsvToRgb(float64(h), s, v)
		hex := fmt.Sprintf("#%02X%02X%02X", r, g, b)

		if !usedColors[hex] {
			usedColors[hex] = true
			return hex
		}
	}
}

// Converts HSV to RGB
func hsvToRgb(h, s, v float64) (uint8, uint8, uint8) {
	c := v * s
	x := c * (1 - absMod(h/60.0, 2) - 1)
	m := v - c

	var r, g, b float64

	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	return uint8((r + m) * 255), uint8((g + m) * 255), uint8((b + m) * 255)
}

func absMod(a, b float64) float64 {
	return a - float64(int(a/b))*b
}

func getCategoryColor(category string, assignedColors map[string]string, usedColors map[string]bool) string {
	if color, exists := assignedColors[category]; exists {
		return color
	}
	color := getUniquePastelColor(usedColors)
	assignedColors[category] = color
	return color
}
