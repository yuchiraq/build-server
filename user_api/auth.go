package user_api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"build-app/models"
)

// RegisterUser обрабатывает регистрацию нового пользователя
func RegisterUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		// Для GET-запроса (временное решение для тестирования)
		if c.Request.Method == http.MethodGet {
			user = models.User{
				ID:         c.Query("id"),
				Login:      c.Query("login"),
				Password:   c.Query("password"),
				FirstName:  c.Query("firstName"),
				SecondName: c.Query("secondName"),
				LastName:   c.Query("lastName"),
				CompanyID:  c.Query("companyID"),
			}
		} else {
			// Для POST-запроса
			if err := c.ShouldBindJSON(&user); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
				return
			}
		}

		// Проверка, что логин не занят
		available, err := models.IsLoginAvailable(db, user.Login)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check login availability" + err})
			return
		}
		if !available {
			c.JSON(http.StatusConflict, gin.H{"error": "Login already exists"})
			return
		}

		// Создание пользователя
		if err := models.CreateUser(db, user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
	}
}

// LoginUser обрабатывает авторизацию пользователя
func LoginUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// Поиск пользователя по логину
		user, err := models.GetUserByLogin(db, input.Login)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login or password"})
			return
		}

		// Проверка пароля
		if err := models.CheckPassword(user.Password, input.Password); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login or password"})
			return
		}

		// Убираем пароль из ответа
		user.Password = ""

		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": user})
	}
}

// CheckLoginAvailability проверяет, свободен ли логин
func CheckLoginAvailability(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		login := c.Query("login")
		if login == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Login parameter is required"})
			return
		}

		available, err := models.IsLoginAvailable(db, login)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check login availability" + err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"available": available})
	}
}