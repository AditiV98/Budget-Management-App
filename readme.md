# Budget Management System

- A modern and secure Budget Management System that helps users track and manage their finances effectively across multiple accounts.
- Built with React (frontend), Golang (backend), and MySQL (database), with authentication powered by OAuth 2.0 and JWT.

## ✨ Features
- 🔒 OAuth 2.0 Login — Secure authentication using OAuth providers (Google)

- 🔑 JWT-based Authentication — Stateless session management with JSON Web Tokens

- 🏦 Multiple Account Management — Manage different accounts (Savings, Cash, Credit Cards, etc.)

- 💸 Transactions Management — Add, edit, and delete transactions seamlessly

- 📄 Transaction CSV Upload — Import bulk transactions from a CSV file

- 🗂 Transaction Categorization — Classify transactions into categories (Food, Rent, Salary, etc.)

- 🔎 Advanced Filtering — Filter transactions by category, type (income/expense), and date range

- 🔁 Recurring Transactions — Automatically manage repeated transactions like monthly bills or salaries

# 🛤 API Endpoints

## 📊 Dashboard
| Method | Endpoint    | Description |
|:------:|:-----------:|:------------|
| GET    | `/dashboard` | Fetch user dashboard data (summary of accounts, transactions, savings) |

---

## 👤 User Management
| Method | Endpoint     | Description              |
|:------:|:------------:|:-------------------------|
| POST   | `/user`       | Create a new user         |
| GET    | `/user`       | Get all users             |
| GET    | `/user/{id}`  | Get a specific user by ID |
| PUT    | `/user/{id}`  | Update user by ID         |
| DELETE | `/user/{id}`  | Delete user by ID         |

---

## 🏦 Account Management
| Method | Endpoint       | Description           |
|:------:|:--------------:|:----------------------|
| POST   | `/account`      | Create a new account  |
| GET    | `/account`      | Get all accounts      |
| GET    | `/account/{id}` | Get account by ID     |
| PUT    | `/account/{id}` | Update account by ID  |
| DELETE | `/account/{id}` | Delete account by ID  |

---

## 💰 Savings Management
| Method | Endpoint        | Description           |
|:------:|:---------------:|:----------------------|
| POST   | `/savings`       | Create a new saving goal |
| GET    | `/savings`       | Get all saving goals  |
| GET    | `/savings/{id}`  | Get saving goal by ID |
| PUT    | `/savings/{id}`  | Update saving goal by ID |
| DELETE | `/savings/{id}`  | Delete saving goal by ID |

---

## 💳 Transaction Management
| Method | Endpoint            | Description           |
|:------:|:-------------------:|:----------------------|
| POST   | `/transaction`        | Create a new transaction |
| GET    | `/transaction`        | Get all transactions  |
| GET    | `/transaction/{id}`   | Get transaction by ID |
| PUT    | `/transaction/{id}`   | Update transaction by ID |
| DELETE | `/transaction/{id}`   | Delete transaction by ID |

---

## 🔁 Recurring Transaction Management
| Method | Endpoint                    | Description                     |
|:------:|:----------------------------:|:-------------------------------|
| POST   | `/recurring-transaction`      | Create a recurring transaction |
| GET    | `/recurring-transaction`      | Get all recurring transactions |
| GET    | `/recurring-transaction/{id}` | Get recurring transaction by ID |
| PUT    | `/recurring-transaction/{id}` | Update recurring transaction by ID |
| DELETE | `/recurring-transaction/{id}` | Delete recurring transaction by ID |

---

## 🔐 Authentication
| Method | Endpoint        | Description                          |
|:------:|:---------------:|:------------------------------------|
| POST   | `/google-token`  | Exchange Google OAuth token for access |
| POST   | `/login`         | Login with OAuth provider token     |
| POST   | `/refresh`       | Refresh access token using refresh token |

