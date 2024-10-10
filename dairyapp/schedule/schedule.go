package schedule

import (
	"dairyapp/database"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

// Grade представляет оценку по предмету.
type Grade struct {
	Subject string
	Score   string
}

// PageData представляет данные для отображения на странице.
type PageData struct {
	SelectedDate string
	Grades       []Grade
}

var store *sessions.CookieStore

func InitStore(s *sessions.CookieStore) {
	store = s
}

// Хранилище для оценок.
var gradesData = make(map[string][]Grade)

func ScheduleHand(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/dairy.html")
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		tmpl.Execute(c.Response(), nil)
		return c.Render(http.StatusOK, "diary.html", nil)
	}
	if c.Request().Method == http.MethodPost {
		date := c.FormValue("date")
		session, err := store.Get(c.Request(), "session-name")
		id := session.Values["userID"].(string)
		if err != nil {
			return http.ErrNoCookie
		}
		grades, err := database.GetSchedules(database.Databasename, date, id)
		if err != nil {
			return fmt.Errorf("Ошибка GetSchedules %v", err)
		}
		return c.JSON(http.StatusOK, grades)
	}
	return nil
}
