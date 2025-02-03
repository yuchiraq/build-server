package user_api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"build-app/models"
)

// RegisterUser обрабатывает регистрацию нового пользовател€
func RegisterUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		// ƒл€ GET-запроса (временное решение дл€ тестировани€)
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
			// ƒл€ POST-запроса
			if err := c.ShouldBindJSON(&user); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
				return
			}
		}

		// ѕроверка, что логин не зан€т
		available, err := models.IsLoginAvailable(db, user.Login)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check login availability"})
			return
		}
		if !available {
			c.JSON(http.StatusConflict, gin.H{"error": "Login already exists"})
			return
		}

		// —оздание пользовател€
		if err := models.CreateUser(db, user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
	}
}

// LoginUser обрабатывает авторизацию пользовател€
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

		// ѕоиск пользовател€ по логину
		user, err := models.GetUserByLogin(db, input.Login)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login or password"})
			return
		}

		// ѕроверка парол€
		if err := models.CheckPassword(user.Password, input.Password); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login or password"})
			return
		}

		// ”бираем пароль из ответа
		user.Password = ""

		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": user})
	}
}

// CheckLoginAvailability провер€ет, свободен ли логин
func CheckLoginAvailability(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		login := c.Query("login")
		if login == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Login parameter is required"})
			return
		}

		available, err := models.IsLoginAvailable(db, login)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check login availability"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"available": available})
	}
}