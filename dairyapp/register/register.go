package register

import (
	"dairyapp/database"
	"html/template"
	"log"
	"net/http"
	"unicode/utf8"

	"github.com/labstack/echo/v4"
)

func RegisterHandler(c echo.Context) error {
	// Проверка на метод POST
	if c.Request().Method == http.MethodPost {
		var data struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
			Role     string `json:"role"`
			Class    string `json:"class"`
		}

		// Декодирование JSON данных из тела запроса
		if err := c.Bind(&data); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "Неверный формат данных",
			})
		}

		// Проверка пользователя в базе данных
		err, result := database.VerificationUser(database.Databasename, data.Username, data.Email)
		if err != nil {
			log.Println("Ошибка верификации пользователя:", err)
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"message": "Ошибка верификации",
			})
		}

		if result {
			if utf8.RuneCountInString(data.Password) >= 6 {
				err = database.RegisterUser(database.Databasename, data.Username, data.Email, data.Password, data.Role, data.Class)
				if err != nil {
					return c.JSON(http.StatusInternalServerError, map[string]interface{}{
						"success": false,
						"message": "Ошибка регистрации пользователя",
					})
				}

				return c.JSON(http.StatusOK, map[string]interface{}{
					"success":  true,
					"message":  "Успешная регистрация",
					"redirect": "/login",
				})
			} else {
				return c.JSON(http.StatusOK, map[string]interface{}{
					"success": false,
					"message": "Пароль должен быть больше 6 символов!",
				})
			}
		}
		if !result {
			return c.JSON(http.StatusConflict, map[string]interface{}{
				"success": false,
				"message": "Пользователь с таким именем или почтой уже занят",
			})
		}
	}

	// Обработка метода GET для отображения HTML формы регистрации
	tmpl, err := template.ParseFiles("templates/register.html")
	if err != nil {
		log.Println("Ошибка при парсинге шаблона:", err)
		return c.String(http.StatusInternalServerError, "Ошибка сервера")
	}

	// Выполнение шаблона регистрации
	if err := tmpl.Execute(c.Response().Writer, nil); err != nil {
		log.Println("Ошибка отображения шаблона register.html:", err)
		return c.String(http.StatusInternalServerError, "Ошибка сервера")
	}

	return nil
}
