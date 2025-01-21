package schedule

import (
	"dairyapp/database"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

type PageData struct {
	UserName string
}

var store *sessions.CookieStore

func InitStore(s *sessions.CookieStore) {
	store = s
}
func ScheduleHand1(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/onepage.html")
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		tmpl.Execute(c.Response(), nil)
		return c.Render(http.StatusOK, "onepage.html", nil)
	}
	return nil
}
func ScheduleHand2(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/onepage2.html")
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		tmpl.Execute(c.Response(), nil)
		return c.Render(http.StatusOK, "onepage2.html", nil)
	}
	return nil
}
func ScheduleHand3(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/onepage3.html")
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		tmpl.Execute(c.Response(), nil)
		return c.Render(http.StatusOK, "onepage3.html", nil)
	}
	return nil
}
func ScheduleHand4(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/onepage4.html")
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		tmpl.Execute(c.Response(), nil)
		return c.Render(http.StatusOK, "onepage4.html", nil)
	}
	return nil
}
func ScheduleHand5(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/onepage5.html")
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		tmpl.Execute(c.Response(), nil)
		return c.Render(http.StatusOK, "onepage5.html", nil)
	}
	return nil
}
func ScheduleHand6(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/onepage6.html")
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		tmpl.Execute(c.Response(), nil)
		return c.Render(http.StatusOK, "onepage6.html", nil)
	}
	return nil
}
func ScheduleHand7(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/onepage7.html")
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		tmpl.Execute(c.Response(), nil)
		return c.Render(http.StatusOK, "onepage7.html", nil)
	}
	return nil
}
func ScheduleHand8(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/onepage8.html")
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		tmpl.Execute(c.Response(), nil)
		return c.Render(http.StatusOK, "onepage8.html", nil)
	}
	return nil
}
func StudentFinalGrades(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		session, err := store.Get(c.Request(), "session-name")
		if err != nil {
			log.Printf("Ошибка получения сессии: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Ошибка получения сессии"})
		}

		userID, ok := session.Values["userID"].(string)
		if !ok {
			log.Println("Не удалось извлечь userID из сессии")
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Не удалось извлечь userID из сессии"})
		}

		username, err := database.GetUserByID(userID, database.Databasename)
		if err != nil {
			log.Printf("Ошибка получения имени пользователя для userID %s: %v", userID, err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Ошибка получения имени пользователя"})
		}

		data := PageData{
			UserName: username,
		}

		tmpl, err := template.ParseFiles("templates/final.html")
		if err != nil {
			log.Printf("Ошибка отображения шаблона: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Ошибка отображения страницы"})
		}

		// Render the template
		if err = tmpl.Execute(c.Response().Writer, data); err != nil {
			log.Printf("Ошибка отображения шаблона: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Ошибка отображения страницы"})
		}
	}
	return nil
}
