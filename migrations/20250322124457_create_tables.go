package migrations

import (
	"gofr.dev/pkg/gofr/migration"
)

const (
	createUsers = `CREATE TABLE users (
  id INT PRIMARY KEY AUTO_INCREMENT,
  first_name VARCHAR(255),
  last_name VARCHAR(255),
  email VARCHAR(255),
  status ENUM('ACTIVE', 'INACTIVE') DEFAULT 'ACTIVE',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP DEFAULT null
);`

	createAccounts = `CREATE TABLE accounts (
  id INT PRIMARY KEY AUTO_INCREMENT,
  user_id INT NOT NULL,
  name VARCHAR(255),
  type ENUM('BANK', 'CASH', 'WALLET', 'CREDIT CARD') NOT NULL,
  balance FLOAT DEFAULT 0,
  status ENUM('ACTIVE', 'INACTIVE') DEFAULT 'ACTIVE',
  expense_categories TEXT NOT NULL,
  saving_categories TEXT NOT NULL ,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP DEFAULT null,
  FOREIGN KEY (user_id) REFERENCES users(id)
);`

	createTransactions = `CREATE TABLE transactions (
  id INT PRIMARY KEY AUTO_INCREMENT,
  user_id INT NOT NULL,
  account_id INT NOT NULL,
  amount FLOAT NOT NULL,
  type ENUM('INCOME', 'EXPENSE', 'SAVINGS') NOT NULL,
  category VARCHAR(255),
  description TEXT,
  transaction_date TIMESTAMP NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP DEFAULT null,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (account_id) REFERENCES accounts(id)
);`

	createSavings = `CREATE TABLE savings (
  id INT PRIMARY KEY AUTO_INCREMENT,
  user_id INT NOT NULL,
  transaction_id INT NOT NULL,
  type ENUM('FD', 'Mutual Funds', 'Stocks', 'Gold ETFs','Other') NOT NULL,
  category VARCHAR(255),
  amount FLOAT NOT NULL,
  current_value FLOAT,
  start_date DATE,
  maturity_date DATE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP DEFAULT null,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (transaction_id) REFERENCES transactions(id)
);`
)

func create_tables() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			_, err := d.SQL.Exec(createUsers)
			if err != nil {
				return err
			}

			_, err = d.SQL.Exec(createAccounts)
			if err != nil {
				return err
			}

			_, err = d.SQL.Exec(createTransactions)
			if err != nil {
				return err
			}

			_, err = d.SQL.Exec(createSavings)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
