package bankaccount

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BankingService interface {
	GetBankAccounts(id int) ([]*BankAccount, error)
	GetBankAccount(id int) (*BankAccount, error)
	DeleteBankAccount(id int) error
	Withdraw(amount float64, id int) error
	Deposit(amount float64, id int) error
	Transfer(accIDFrom int, accIDTo int, amount float64) error
}

type BankingHandler struct {
	bankingServices BankingService
}

func (bh *BankingHandler) getBankAccounts(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	accounts, err := bh.bankingServices.GetBankAccounts(id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func (bh *BankingHandler) withdraw(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var update struct {
		Amount float64 `json:"amount" binding:"required"`
	}
	err = c.ShouldBindJSON(&update)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	err = bh.bankingServices.Withdraw(update.Amount, id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "Withdraw success")
}

func (bh *BankingHandler) deleteBankAccount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	err = bh.bankingServices.DeleteBankAccount(id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, "Bank Account Deleted")
}

func (bh *BankingHandler) transfers(c *gin.Context) {

	var update struct {
		Amount    float64 `json:"amount" binding:"required"`
		AccIdFrom int     `json:"from" binding:"required"`
		AccIdTo   int     `json:"to" binding:"required"`
	}

	err := c.ShouldBindJSON(&update)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	err = bh.bankingServices.Transfer(update.AccIdFrom, update.AccIdTo, update.Amount)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "Transfers success")
}

func (bh *BankingHandler) deposit(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var update struct {
		Amount float64 `json:"amount" binding:"required"`
	}
	err = c.ShouldBindJSON(&update)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	err = bh.bankingServices.Deposit(update.Amount, id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "Deposit Success")
}
