package main

import (
	transaction_processing "github.com/ebitezion/backend-framework/internal/transactions"
	"github.com/gin-gonic/gin"
)

var tp *transaction_processing.TransactionProcessor

func run_test() {
	tp = transaction_processing.NewTransactionProcessor()

	r := gin.Default()

	// Add accounts
	r.POST("/add-account", addAccountHandler)

	// Deposit endpoint
	r.POST("/deposit", depositHandler)

	// Withdraw endpoint
	r.POST("/withdraw", withdrawHandler)

	// Transfer endpoint
	r.POST("/transfer", transferHandler)

	r.Run(":8080")
}

func addAccountHandler(c *gin.Context) {
	var request struct {
		AccountID      string  `json:"account_id"`
		InitialBalance float64 `json:"initial_balance"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	tp.AddAccount(request.AccountID, request.InitialBalance)
	c.JSON(200, gin.H{"message": "Account added successfully"})
}

func depositHandler(c *gin.Context) {
	var request struct {
		AccountID string  `json:"account_id"`
		Amount    float64 `json:"amount"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	err := tp.Deposit(request.AccountID, request.Amount)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Deposit successful"})
}

func withdrawHandler(c *gin.Context) {
	var request struct {
		AccountID string  `json:"account_id"`
		Amount    float64 `json:"amount"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	err := tp.Withdraw(request.AccountID, request.Amount)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Withdrawal successful"})
}

func transferHandler(c *gin.Context) {
	var request struct {
		FromAccountID string  `json:"from_account_id"`
		ToAccountID   string  `json:"to_account_id"`
		Amount        float64 `json:"amount"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	err := tp.Transfer(request.FromAccountID, request.ToAccountID, request.Amount)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Transfer successful"})
}
