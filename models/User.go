package models

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         string `json:"id"`
	Login      string `json:"login"`
	Password   string `json:"-"`
	FirstName  string `json:"firstName"`
	SecondName string `json:"secondName"`
	LastName   string `json:"lastName"`
	CompanyID  string `json:"companyID"`
}

// HashPassword �������� ������
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}
	return string(hashedPassword), nil
}

// CheckPassword ���������, ��������� �� ������ � �����
func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// CreateUser ��������� ������ ������������ � ���� ������
func CreateUser(db *sql.DB, user User) error {
	// �������� ������
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	query := `
		INSERT INTO users (id, login, password, first_name, second_name, last_name, company_id)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err = db.Exec(query, user.ID, user.Login, hashedPassword, user.FirstName, user.SecondName, user.LastName, user.CompanyID)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	return nil
}

// GetUserByLogin ���������� ������������ �� ������
func GetUserByLogin(db *sql.DB, login string) (User, error) {
	var user User
	query := `SELECT id, login, password, first_name, second_name, last_name, company_id FROM users WHERE login = ?`
	err := db.QueryRow(query, login).Scan(&user.ID, &user.Login, &user.Password, &user.FirstName, &user.SecondName, &user.LastName, &user.CompanyID)
	if err != nil {
		return User{}, fmt.Errorf("failed to get user: %v", err)
	}
	return user, nil
}

// IsLoginAvailable ���������, �������� �� �����
func IsLoginAvailable(db *sql.DB, login string) (bool, error) {
	_, err := GetUserByLogin(db, login)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil // ����� ��������
		}
		return false, fmt.Errorf("failed to check login: %v", err)
	}
	return false, nil // ����� �����
}