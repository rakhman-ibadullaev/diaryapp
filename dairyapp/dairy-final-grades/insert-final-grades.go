package dairyfinalgrades

import (
	"dairyapp/database"
	"html/template"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type data1 struct {
	Student string `json:"student"`
	Subject string `json:"subject"`
	First   string `json:"first"`
	Second  string `json:"second"`
	Third   string `json:"third"`
	Fourth  string `json:"fourth"`
	Year    string `json:"year"`
	Exam    string `json:"exam"`
	Final   string `json:"final"`
}

func convert(data1 data1) (database.Data, error) {
	var data database.Data

	// Преобразуем строку Student в int64

	data.Student = data1.Student
	data.Subject = data1.Subject
	data.First = data1.First
	data.Second = data1.Second
	data.Third = data1.Third
	data.Fourth = data1.Fourth
	data.Year = data1.Year
	data.Exam = data1.Exam
	data.Final = data1.Final

	return data, nil
}
func FinalGradesHandler(c echo.Context) error {
	if c.Request().Method == http.MethodPost {
		var data1 data1
		// Декодирование JSON из тела запроса
		if err := c.Bind(&data1); err != nil {
			// Ошибка декодирования данных
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "Неверный формат данных",
			})
		}
		data, err := convert(data1)
		if err != nil {
			log.Printf("convert err is %v", err)
		}
		log.Printf("data is %v", data)

		err = database.InsertFinalGrades(database.Databasename, data)
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
		tmpl, err := template.ParseFiles("templates/insert-final-grades.html")
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		tmpl.Execute(c.Response(), nil)
		return c.Render(http.StatusOK, "insert-final-grades.html", nil)
	}
	return echo.ErrMethodNotAllowed
}
