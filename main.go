package main

import (
	"gofr.dev/pkg/gofr"
	"moneyManagement/middlewares"
	"moneyManagement/migrations"
	"moneyManagement/services/auth"
	"moneyManagement/stores/accounts"
	"moneyManagement/stores/recurringTransactions"
	"moneyManagement/stores/savings"
	"moneyManagement/stores/transactions"
	"moneyManagement/stores/users"

	validatorSvc "moneyManagement/services/Validator"
	accountService "moneyManagement/services/accounts"
	dashboardService "moneyManagement/services/dashboard"
	recurringTransactionService "moneyManagement/services/recurringTransactions"
	savingsService "moneyManagement/services/savings"
	transactionService "moneyManagement/services/transactions"
	usersService "moneyManagement/services/users"

	accountsHandler "moneyManagement/handler/accounts"
	authHandlers "moneyManagement/handler/auth"
	dashboardHandlers "moneyManagement/handler/dashboard"
	recurringTransactionsHandler "moneyManagement/handler/recurringTransactions"
	savingsHandler "moneyManagement/handler/savings"
	transactionsHandler "moneyManagement/handler/transactions"
	usersHandler "moneyManagement/handler/users"
)

func main() {
	app := gofr.New()

	app.Migrate(migrations.All())

	userStore := users.New()
	accountStore := accounts.New()
	transactionStore := transactions.New()
	savingStore := savings.New()
	recurringTransactionStore := recurringTransactions.New()

	userSvc := usersService.New(userStore)
	accountSvc := accountService.New(accountStore, userSvc)
	savingsSvc := savingsService.New(savingStore, transactionStore)
	transactionSvc := transactionService.New(transactionStore, accountSvc, savingsSvc, userSvc)
	dashboardSvc := dashboardService.New(accountSvc, transactionSvc, userSvc, savingsSvc)
	recurringTransactionSvc := recurringTransactionService.New(recurringTransactionStore, userSvc)
	authSvc := auth.New(app.Config.Get("REFRESH_SECRET"), app.Config.Get("ACCESS_SECRET"), app.Config.Get("GOOGLE_CLIENT_ID"),
		app.Config.Get("GOOGLE_CLIENT_SECRET"), app.Config.Get("REDIRECT_URL"))
	validator := validatorSvc.New(app.Config.Get("ACCESS_SECRET"))

	userHandler := usersHandler.New(userSvc)
	accountHandler := accountsHandler.New(accountSvc)
	savingHandler := savingsHandler.New(savingsSvc)
	transactionHandler := transactionsHandler.New(transactionSvc)
	dashboardHandler := dashboardHandlers.New(dashboardSvc)
	authHandler := authHandlers.New(authSvc, userSvc)
	recurringTransactionHandler := recurringTransactionsHandler.New(recurringTransactionSvc)

	app.UseMiddleware(middlewares.Authorization([]middlewares.ExemptPath{
		{Path: "^/google-token$", Method: "POST"},
		{Path: "^/login$", Method: "POST"},
		{Path: "^/refresh$", Method: "POST"},
	}, validator, userSvc))

	app.GET("/dashboard", dashboardHandler.Get)

	app.POST("/user", userHandler.Create)
	app.GET("/user", userHandler.GetAll)
	app.GET("/user/{id}", userHandler.GetByID)
	app.PUT("/user/{id}", userHandler.Update)
	app.DELETE("/user/{id}", userHandler.Delete)

	app.POST("/account", accountHandler.Create)
	app.GET("/account", accountHandler.GetAll)
	app.GET("/account/{id}", accountHandler.GetByID)
	app.PUT("/account/{id}", accountHandler.Update)
	app.DELETE("/account/{id}", accountHandler.Delete)

	app.POST("/savings", savingHandler.Create)
	app.GET("/savings", savingHandler.GetAll)
	app.GET("/savings/{id}", savingHandler.GetByID)
	app.PUT("/savings/{id}", savingHandler.Update)
	app.DELETE("/savings/{id}", savingHandler.Delete)

	app.POST("/transaction", transactionHandler.Create)
	app.GET("/transaction", transactionHandler.GetAll)
	app.GET("/transaction/{id}", transactionHandler.GetByID)
	app.PUT("/transaction/{id}", transactionHandler.Update)
	app.DELETE("/transaction/{id}", transactionHandler.Delete)

	app.POST("/recurring-transaction", recurringTransactionHandler.Create)
	app.GET("/recurring-transaction", recurringTransactionHandler.GetAll)
	app.GET("/recurring-transaction/{id}", recurringTransactionHandler.GetByID)
	app.PUT("/recurring-transaction/{id}", recurringTransactionHandler.Update)
	app.POST("/recurring-transaction/{id}/skip", recurringTransactionHandler.SkipNextRun)
	app.DELETE("/recurring-transaction/{id}", recurringTransactionHandler.Delete)

	app.POST("/google-token", authHandler.CreateToken)
	app.POST("/login", authHandler.Login)
	app.POST("/refresh", authHandler.Refresh)

	app.Run()
}
