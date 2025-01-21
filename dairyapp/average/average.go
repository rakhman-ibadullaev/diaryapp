package average

import (
	"dairyapp/database"
	"log"
	"net/http"
	"text/template"

	"github.com/labstack/echo/v4"
)

type Marks struct {
	Math     float32
	Rus      float32
	Ing      float32
	Chem     float32
	Bio      float32
	Ob       float32
	Info     float32
	Sport    float32
	Phy      float32
	Obizr    float32
	Tech     float32
	Creature float32
	Hist     float32
	Geo      float32
	Lit      float32
}

func AverageHandler(c echo.Context) error {
	if c.Request().Method == http.MethodPost {
		var data struct {
			Name_Student string `json:"student"`
		}
		// Декодирование JSON из тела запроса
		if err := c.Bind(&data); err != nil {
			// Ошибка декодирования данных
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "Неверный формат данных",
			})
		}

		var marks Marks
		var err error
		marks.Bio, err = database.AverageMark(database.Databasename, data.Name_Student, "Биология")
		if err != nil {
			log.Printf("Error in averagehand %v", err)
		}
		marks.Math, err = database.AverageMark(database.Databasename, data.Name_Student, "Математика")
		if err != nil {
			log.Printf("Error in averagehand %v", err)
		}
		marks.Rus, err = database.AverageMark(database.Databasename, data.Name_Student, "Русский язык")
		if err != nil {
			log.Printf("Error in averagehand %v", err)
		}
		marks.Tech, err = database.AverageMark(database.Databasename, data.Name_Student, "Технология")
		if err != nil {
			log.Printf("Error in averagehand %v", err)
		}
		marks.Phy, err = database.AverageMark(database.Databasename, data.Name_Student, "Физика")
		if err != nil {
			log.Printf("Error in averagehand %v", err)
		}
		marks.Chem, err = database.AverageMark(database.Databasename, data.Name_Student, "Химия")
		if err != nil {
			log.Printf("Error in averagehand %v", err)
		}
		marks.Hist, err = database.AverageMark(database.Databasename, data.Name_Student, "История")
		if err != nil {
			log.Printf("Error in averagehand %v", err)
		}
		marks.Ob, err = database.AverageMark(database.Databasename, data.Name_Student, "Обществознание")
		if err != nil {
			log.Printf("Error in averagehand %v", err)
		}
		marks.Obizr, err = database.AverageMark(database.Databasename, data.Name_Student, "ОБиЗР")
		if err != nil {
			log.Printf("Error in averagehand %v", err)
		}
		marks.Ing, err = database.AverageMark(database.Databasename, data.Name_Student, "Английский язык")
		if err != nil {
			log.Printf("Error in averagehand %v", err)
		}
		marks.Info, err = database.AverageMark(database.Databasename, data.Name_Student, "Информатика")
		if err != nil {
			log.Printf("Error in averagehand %v", err)
		}
		marks.Sport, err = database.AverageMark(database.Databasename, data.Name_Student, "Физкультура")
		if err != nil {
			log.Printf("Error in averagehand %v", err)
		}
		marks.Creature, err = database.AverageMark(database.Databasename, data.Name_Student, "Изобразительное искусство")
		if err != nil {
			log.Printf("Error in averagehand %v", err)
		}
		marks.Geo, err = database.AverageMark(database.Databasename, data.Name_Student, "География")
		if err != nil {
			log.Printf("Error in averagehand %v", err)
		}
		marks.Lit, err = database.AverageMark(database.Databasename, data.Name_Student, "Литература")
		if err != nil {
			log.Printf("Error in averagehand %v", err)
		}
		return c.JSON(http.StatusOK, marks)
	}
	if c.Request().Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/average.html")
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		tmpl.Execute(c.Response(), nil)
		return c.Render(http.StatusOK, "average.html", nil)
	}
	return echo.ErrMethodNotAllowed
}
