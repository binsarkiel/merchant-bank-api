package repository

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"merchant-bank-api/models"
)

// Data Path Constanta
const (
	usersDataPath        = "data/users.json"
	sessionsDataPath     = "data/sessions.json"
	transactionsDataPath = "data/transactions.json"
)

// JSON handling using mutex
var (
	usersMutex        sync.Mutex
	sessionsMutex     sync.Mutex
	transactionsMutex sync.Mutex
)

// LoadUsers loads users from the JSON file.
func LoadUsers() ([]models.User, error) {
	usersMutex.Lock()
	defer usersMutex.Unlock()

	data, err := os.ReadFile(usersDataPath)
	if err != nil {
		return nil, err
	}

	var users []models.User
	err = json.Unmarshal(data, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// SaveUsers saves users to the JSON file.
func SaveUsers(users []models.User) error {
	usersMutex.Lock()
	defer usersMutex.Unlock()

	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(usersDataPath, data, 0644)
}

// LoadSessions loads sessions from the JSON file.
func LoadSessions() ([]models.Session, error) {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	data, err := os.ReadFile(sessionsDataPath)
	if err != nil {
		return nil, err
	}

	var sessions []models.Session
	err = json.Unmarshal(data, &sessions)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

// SaveSessions saves sessions to the JSON file.
func SaveSessions(sessions []models.Session) error {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	data, err := json.MarshalIndent(sessions, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(sessionsDataPath, data, 0644)
}

// LoadTransactions loads transactions from the JSON file.
func LoadTransactions() ([]models.Transaction, error) {
	transactionsMutex.Lock()
	defer transactionsMutex.Unlock()

	data, err := os.ReadFile(transactionsDataPath)
	if err != nil {
		return nil, err
	}

	var transactions []models.Transaction
	err = json.Unmarshal(data, &transactions)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

// SaveTransactions saves transactions to the JSON file.
func SaveTransactions(transactions []models.Transaction) error {
	transactionsMutex.Lock()
	defer transactionsMutex.Unlock()

	data, err := json.MarshalIndent(transactions, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(transactionsDataPath, data, 0644)
}

// FindUserByUsername finds a user by their username.
func FindUserByUsername(username string) (*models.User, error) {
	users, err := LoadUsers()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.Username == username {
			return &user, nil
		}
	}

	return nil, errors.New("user not found")
}

// UpdateUserBalance updates a user's balance by username.
func UpdateUserBalance(username string, newBalance float64) error {
	users, err := LoadUsers()
	if err != nil {
		return err
	}

	for i, user := range users {
		if user.Username == username {
			users[i].AccountBalance = newBalance
			return SaveUsers(users)
		}
	}

	return errors.New("user not found")
}

// AddSession adds a new session to the sessions file.
func AddSession(session models.Session) error {
	sessions, err := LoadSessions()
	if err != nil {
		return err
	}

	sessions = append(sessions, session)
	return SaveSessions(sessions)
}

// AddTransaction adds a new transaction to the transactions file.
func AddTransaction(transaction models.Transaction) error {
	transactions, err := LoadTransactions()
	if err != nil {
		return err
	}

	transactions = append(transactions, transaction)
	return SaveTransactions(transactions)
}

// GetUserBalance retrieves the current balance of a user.
func GetUserBalance(username string) (float64, error) {
	user, err := FindUserByUsername(username)
	if err != nil {
		return 0, err
	}
	return user.AccountBalance, nil
}
