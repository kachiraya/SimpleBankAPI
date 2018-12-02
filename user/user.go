package user

import (
	"database/sql"
	"simplebankapi-heroku/user/bankaccount"
)

type Service struct {
	DB *sql.DB
}

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

func (s *Service) Insert(u *User) error {
	stmt := `INSERT INTO users(first_name, last_name)
		 values ($1, $2) RETURNING id`
	row := s.DB.QueryRow(stmt, u.FirstName, u.LastName)
	err := row.Scan(&u.ID)

	return err
}

func (s *Service) All() ([]User, error) {
	stmt := "SELECT id, first_name, last_name FROM users ORDER BY id DESC"
	rows, err := s.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	var us []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.FirstName, &u.LastName)
		if err != nil {
			return nil, err
		}
		us = append(us, u)
	}
	return us, nil
}

func (s *Service) FindByID(id int) (*User, error) {
	stmt := "SELECT id, first_name, last_name FROM users WHERE id = $1"
	row := s.DB.QueryRow(stmt, id)
	var u User
	err := row.Scan(&u.ID, &u.FirstName, &u.LastName)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Service) FindByName(firstName, lastName string) (*User, error) {
	stmt := "SELECT id, first_name, last_name FROM users WHERE first_name = $1, last_name = $2"
	row := s.DB.QueryRow(stmt, firstName, lastName)
	var u User
	err := row.Scan(&u.ID, &u.FirstName, &u.LastName)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Service) Update(u *User) error {
	stmt := "UPDATE users SET first_name = $1, last_name = $2 WHERE id = $3"
	_, err := s.DB.Exec(stmt, u.FirstName, u.LastName, u.ID)
	return err
}

func (s *Service) Delete(u *User) error {
	stmt := "DELETE FROM users WHERE id = $1"
	_, err := s.DB.Exec(stmt, u.ID)
	return err
}

func (s *Service) AddBankAccount(acc *bankaccount.BankAccount) error {
	stmt := `INSERT INTO bankaccounts(user_id, account_no, name, balance)
		 values ($1, $2, $3, $4) RETURNING id`
	row := s.DB.QueryRow(stmt, acc.UserID, acc.AccountNumber, acc.Name, acc.Balance)
	err := row.Scan(&acc.ID)
	return err
}
