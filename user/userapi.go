package user

import (
	"bankapi/user/bankaccount"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	FindByID(id int) (*User, error)
	All() ([]User, error)
	Insert(u *User) error
	Update(u *User) error
	Delete(u *User) error
}

type Handler struct {
	userService UserService
}

func (h *Handler) allUser(c *gin.Context) {
	users, err := h.userService.All()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *Handler) getUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	user, err := h.userService.FindByID(id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) createUser(c *gin.Context) {
	var u User
	err := c.ShouldBindJSON(&u)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = h.userService.Insert(&u)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, u)
}

func (h *Handler) updateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	u, err := h.userService.FindByID(id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var update struct {
		FirstName *string `json:"first_name" binding:"required"`
		LastName  *string `json:"last_name" binding:"required"`
	}

	err = c.ShouldBindJSON(&update)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	if update.FirstName != nil {
		u.FirstName = *update.FirstName
	}
	if update.LastName != nil {
		u.LastName = *update.LastName
	}

	err = h.userService.Update(u)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (h *Handler) deleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	err = h.userService.Delete(&User{
		ID: id,
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (h *Handler) addBankAccount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var acc bankaccount.BankAccount
	err = c.ShouldBindJSON(&acc)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	user, err := h.userService.FindByID(id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	acc.UserID = user.ID
	acc.Balance = 0
	acc.Name = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	acc.AccountNumber = "125635" //To be implemented: generate unique acc no.
	err = h.bankingService.AddBankAccount(&acc)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, acc)
}

func StartServer(addr string, db *sql.DB) error {
	r := gin.Default()
	h := &Handler{
		userService: &Service{
			DB: db,
		},
	}
	bh := &bankaccount.BankingHandler{
		bankingService: &bankaccount.BankService{
			DB: db,
		},
	}

	r.GET("/users", h.allUser)
	r.GET("/users/:id", h.getUser)
	r.POST("/users", h.createUser)
	r.PUT("/users/:id", h.updateUser)
	r.DELETE("/users/:id", h.deleteUser)
	r.POST("/users/:id/bankAccounts", h.addBankAccount)
	r.GET("/users/:id/bankAccounts", bh.getBankAccounts)
	r.DELETE("/bankAccounts/:id", bh.deleteBankAccount)
	r.PUT("/bankAccounts/:id/withdraw", bh.withdraw)
	r.PUT("/bankAccounts/:id/deposit", bh.deposit)
	r.POST("/transfers", bh.transfers)

	return r.Run(addr)
}
