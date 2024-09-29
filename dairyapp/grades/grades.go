package grades

import (
	"dairyapp/database"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

var store *sessions.CookieStore

func InitStore(s *sessions.CookieStore) {
	store = s
}

func InsertGradesHandler(c echo.Context) error {
	if c.Request().Method == http.MethodPost {
		var data struct {
			Grade        string `json:"grade"`
			Subject      string `json:"subject"`
			Name_Student string `json:"student"`
			Work_type    string `json:"workType"`
			Date         string `json:"date"`
		}
		// Декодирование JSON из тела запроса
		if err := c.Bind(&data); err != nil {
			// Ошибка декодирования данных
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "Неверный формат данных",
			})
		}
		// Внесение оценок в базу данных
		session, err := store.Get(c.Request(), "session-name")
		id_teacher := session.Values["userID"].(string)

		log.Printf("data is %v,%v,%v,%v,%v,%v", data.Grade, data.Subject, data.Work_type, data.Date, data.Name_Student, id_teacher)

		err = database.InsertGrades(database.Databasename, data.Grade, data.Subject, data.Work_type, data.Date, data.Name_Student, id_teacher)
		if err != nil {
			log.Printf("Ошибка внесения оценки %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"message": "Не удалось внести запись. Попробуйте позже",
			})
		}

		// Успешный ответ
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Успешно",
		})
	}

	if c.Request().Method == http.MethodGet {
		// Отображение HTML-шаблона
		tmpl, err := template.ParseFiles("templates/insert-grades.html")
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		tmpl.Execute(c.Response(), nil)
		return c.Render(http.StatusOK, "insert-grades.html", nil)
	}

	return echo.ErrMethodNotAllowed // или можно вернуть 405
}
func GetStudentsList(c echo.Context) error {
	// Запрос к базе данных
	if c.Request().Method == http.MethodPost {
		var data struct {
			Class string `json:"class"`
		}
		// Декодирование JSON из тела запроса
		if err := c.Bind(&data); err != nil {
			// Ошибка декодирования данных
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "Неверный формат данных",
			})
		}
		log.Printf("Class is %v", data.Class)
		users, err := database.GetUsers(database.Databasename, data.Class)
		if err != nil {
			log.Printf("Ошибка GetStudents %v", err)
			return err
		}

		// Формирование списка учеников
		students := []string{}
		for _, user := range users {
			students = append(students, user.Name)
		}
		// Формирование ответа сервера
		log.Printf("Список учеников %v", students)
		return c.JSON(http.StatusOK, students)
	}
	return nil
}

func DeleteGrades(c echo.Context) error {
	var data struct {
		ID string `json:"id"`
	}
	// Декодирование JSON из тела запроса
	if err := c.Bind(&data); err != nil {
		// Ошибка декодирования данных
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Неверный формат данных",
		})
	}
	log.Printf("ID is %v", data.ID)
	err := database.DeleteGrades(database.Databasename, data.ID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Успешно удалена запись",
	})
}

func AllGradesHTML(c echo.Context) error {
	tmpl, err := template.ParseFiles("templates/grades.html")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	tmpl.Execute(c.Response(), nil)
	return c.Render(http.StatusOK, "grades.html", nil)
}

func AllGrades(c echo.Context) error {
	session, err := store.Get(c.Request(), "session-name")
	id_teacher := session.Values["userID"].(string)

	grades, err := database.AllGrades(database.Databasename, id_teacher)
	log.Printf("Class равен %v", grades)
	if err != nil {
		log.Printf("Ошибка AllGrades %v", err)
	}
	return c.JSON(http.StatusOK, grades)
}
