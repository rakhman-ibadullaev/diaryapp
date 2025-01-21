package login

import (
	"dairyapp/database"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

var loginAttempts = make(map[string]int)
var store *sessions.CookieStore

func LoginHandler(c echo.Context) error {
	if c.Request().Method == http.MethodPost {
		email := c.FormValue("email")
		password := c.FormValue("password")
		log.Printf("Login attempt: email=%s", email)

		if loginAttempts[email] >= 5 {
			return c.JSON(http.StatusTooManyRequests, map[string]interface{}{
				"success": false,
				"message": "Too many login attempts",
			})
		}

		user := &database.User{}
		user, result, err := database.LoginUser(database.Databasename, email, password)
		if err != nil {
			log.Printf("Login error: %v", err)
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"message": "Ошибка входа: " + err.Error(),
			})
		}

		if result {
			// Успешный вход
			log.Println("Successful login")
			session, _ := store.Get(c.Request(), "session-name")
			session.Values["authenticated"] = true
			session.Values["userID"] = user.ID
			session.Values["userRole"] = user.Role
			session.Values["userEmail"] = email
			session.Values["userPassword"] = password

			log.Printf("Проверка роли: %s, Проверка ID: %v", user.Role, user.ID)

			err = session.Save(c.Request(), c.Response())
			if err != nil {
				log.Printf("Ошибка сохранения сессии при успешном входе: %v", err)
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"success":  false,
					"message":  "Ошибка сохранения сессии",
					"redirect": "/login",
				})
			}

			if user.Role == "teacher" {
				session.Options.MaxAge = 1800 // 30 minutes for teachers
				return c.JSON(http.StatusOK, map[string]interface{}{
					"success":  true,
					"message":  "Успешный вход",
					"redirect": "/teacher/home",
				})
			} else {
				session.Options.MaxAge = 600 // 10 minutes for students
				return c.JSON(http.StatusOK, map[string]interface{}{
					"success":  true,
					"message":  "Успешный вход",
					"redirect": "/student/home",
				})
			}
		} else {
			log.Println("Invalid login or password")
			loginAttempts[email]++
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"success":  false,
				"message":  "Неверный пароль или логин.",
				"redirect": "/login",
			})
		}
	}
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		log.Println("Ошибка при парсинге шаблона:", err)
		return c.String(http.StatusInternalServerError, "Ошибка сервера")
	}

	// Выполнение шаблона регистрации
	if err := tmpl.Execute(c.Response().Writer, nil); err != nil {
		log.Println("Ошибка отображения шаблона login.html:", err)
		return c.String(http.StatusInternalServerError, "Ошибка сервера")
	}
	return err
}
func LogoutHandler(c echo.Context) error {
	// Retrieve the session using the context
	session, err := store.Get(c.Request(), "session-name")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	// Set MaxAge to -1 to delete the session
	session.Options.MaxAge = -1
	err = session.Save(c.Request(), c.Response())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	// Redirect to the login page
	return c.Redirect(http.StatusSeeOther, "/login")
}

// InitStore initializes the session store
func InitStore(s *sessions.CookieStore) {
	store = s
}
