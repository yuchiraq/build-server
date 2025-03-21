package models

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         string `json:"id"`
	Login      string `json:"login"`
	Password   string `json:"password"`
	FirstName  string `json:"firstName"`
	SecondName string `json:"secondName"`
	LastName   string `json:"lastName"`
	CompanyID  string `json:"companyID"`
}

// HashPassword хэширует пароль
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("failed to hash password: %v", err)
		return "", fmt.Errorf("failed to hash password: %v", err)
	}
	return string(hashedPassword), nil
}

// CheckPassword проверяет, совпадает ли пароль с хэшем
func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// CreateUser добавляет нового пользователя в базу данных
func CreateUser(db *sql.DB, user User) error {
	// Хэшируем пароль
	fmt.Println("User: ", user.Login, "password: ", user.Password)
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		fmt.Println("failed to hash password: %v", err)
		return fmt.Errorf("failed to hash password: %v", err)
	}

	query := `
		INSERT INTO users (login, password, firstName, secondName, lastName, companyID)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err = db.Exec(query, user.Login, hashedPassword, user.FirstName, user.SecondName, user.LastName, user.CompanyID)
	if err != nil {
		fmt.Println("failed to create user: %v", err)
		return fmt.Errorf("failed to create user: %v", err)
	}
	return nil
}

// GetUserByLogin возвращает пользователя по логину
func GetUserByLogin(db *sql.DB, login string) (User, error) {
	var user User
	query := `SELECT id, login, password, firstName, secondName, lastName, companyID FROM users WHERE login = ?`
	err := db.QueryRow(query, login).Scan(&user.ID, &user.Login, &user.Password, &user.FirstName, &user.SecondName, &user.LastName, &user.CompanyID)
	if err != nil {
		fmt.Println("failed to get user: %v", err)
		return User{}, fmt.Errorf("failed to get user: %v", err)
	}
	return user, nil
}

// IsLoginAvailable проверяет, свободен ли логин
func IsLoginAvailable(db *sql.DB, login string) (bool, error) {
	var user User
	query := `SELECT id, login FROM users WHERE login = ?`
	err := db.QueryRow(query, login).Scan(&user.ID, &user.Login)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil // Логин свободен
		}
		fmt.Println("failed to check login: %v", err)
		return false, fmt.Errorf("failed to check login: %v", err)
	}
	return false, nil // Логин занят
}