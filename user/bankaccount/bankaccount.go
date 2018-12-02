package bankaccount

import "database/sql"

type BankAccount struct {
	ID            int     `json:"id"`
	UserID        int     `json:"user_id"`
	AccountNumber string  `json:"account_no"`
	Name          string  `json:"name" binding:"required"`
	Balance       float64 `json:"balance"`
}

type BankService struct {
	DB *sql.DB
}

func (s *BankService) GetBankAccounts(id int) ([]BankAccount, error) {
	stmt := "SELECT * FROM bankaccounts WHERE user_id = $1 ORDER BY id DESC"
	rows, err := s.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}
	var accs []BankAccount
	for rows.Next() {
		var acc BankAccount
		err := rows.Scan(&acc.ID, &acc.UserID, &acc.AccountNumber, &acc.Name, &acc.Balance)
		if err != nil {
			return nil, err
		}
		accs = append(accs, acc)
	}
	return accs, nil
}

func (s *BankService) GetBankAccount(id int) (*BankAccount, error) {
	stmt := "SELECT * FROM bankaccounts WHERE id = $1"
	rows, err := s.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}
	var accs []*BankAccount
	for rows.Next() {
		var acc *BankAccount
		err := rows.Scan(&acc.ID, &acc.UserID, &acc.AccountNumber, &acc.Name, &acc.Balance)
		if err != nil {
			return nil, err
		}
		accs = append(accs, acc)
	}
	return accs[0], nil
}

func (s *BankService) Withdraw(amount float64, id int) error {
	acc, err := s.GetBankAccount(id)
	if err != nil {
		return err
	}
	stmt := "UPDATE bankaccounts SET balance = $1 WHERE id = $2"
	_, err = s.DB.Exec(stmt, acc.Balance-amount, acc.ID)
	return err
}

func (s *BankService) Deposit(amount float64, id int) error {
	acc, err := s.GetBankAccount(id)
	if err != nil {
		return err
	}
	stmt := "UPDATE bankaccounts SET balance = $1 WHERE id = $2"
	_, err = s.DB.Exec(stmt, acc.Balance+amount, acc.ID)
	return err
}

func (s *BankService) DeleteBankAccount(id int) error {
	stmt := "DELETE FROM bankaccounts WHERE id = $1"
	_, err := s.DB.Exec(stmt, id)
	return err
}

func (s *BankService) Transfer(accIDFrom int, accIDTo int, amount float64) error {
	err := s.Withdraw(amount, accIDFrom)
	if err != nil {
		return err
	}
	err = s.Deposit(amount, accIDTo)
	return err
}
