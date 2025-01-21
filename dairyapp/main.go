package main

import (
	"dairyapp/average"
	dairyfinalgrades "dairyapp/dairy-final-grades"
	"dairyapp/database"
	"dairyapp/grades"
	"dairyapp/homepage"
	"dairyapp/key"
	"dairyapp/login"
	"dairyapp/register"
	"dairyapp/schedule"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

var (
	store *sessions.CookieStore
)

func main() {
	key, err := key.GeneratePassword(32)
	if err != nil {
		fmt.Println("Ошибка образования ключа", err)
		return
	}
	store = sessions.NewCookieStore([]byte(key))

	// Инициализация хранилища для различных обработчиков
	login.InitStore(store)
	homepage.InitStore(store)
	grades.InitStore(store)
	// Настройка опций для cookie
	store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
	}

	e := echo.New()

	// Указываем директорию с статическими файлами
	e.Static("/subpath", "./templates/images")

	// Маршруты регистрации
	e.GET("/register", register.RegisterHandler)
	e.POST("/register", register.RegisterHandler)
	// Маршруты авторизации и деавторизации
	e.GET("/login", login.LoginHandler)
	e.POST("/login", login.LoginHandler)
	e.GET("/logout", login.LogoutHandler)
	// Маршруты основной страницы
	e.GET("/send-error", homepage.SendError, RoleMiddleware(store, "admin", "teacher", "student"))
	e.GET("/student/home", homepage.StudentHomepage, RoleMiddleware(store, "admin", "student"))
	e.GET("/teacher/home", homepage.TeacherHomepage, RoleMiddleware(store, "admin", "teacher"))
	e.GET("/diaryapp", homepage.DiaryPage, RoleMiddleware(store, "admin", "teacher", "student"))
	e.GET("/active-users", homepage.ActiveUsersHandler, RoleMiddleware(store, "admin", "teacher", "student"))
	e.GET("/average", average.AverageHandler, RoleMiddleware(store, "admin", "teacher"))
	e.POST("/average", average.AverageHandler)
	//Маршруты личного аккаунта полльзователя
	e.GET("/my-acc", homepage.AccountHomepage, RoleMiddleware(store, "admin", "teacher", "student"))
	//Маршруты внесения оценок
	e.POST("/students", grades.GetStudentsList)
	e.GET("/insert-grades", grades.InsertGradesHandler, RoleMiddleware(store, "admin", "teacher"))
	e.POST("/insert-grades", grades.InsertGradesHandler)
	//Маршруты оценок-записей для просмотра учителя
	e.GET("/grades", grades.AllGradesHTML)
	e.GET("/all-grades", grades.AllGrades, RoleMiddleware(store, "admin", "teacher"))
	//Маршрут удаление записи
	e.POST("/delete-grades", grades.DeleteGrades)
	//Маршрут для расписания
	e.GET("/diary-1", schedule.ScheduleHand1)
	e.GET("/diary-2", schedule.ScheduleHand2)
	e.GET("/diary-3", schedule.ScheduleHand3)
	e.GET("/diary-4", schedule.ScheduleHand4)
	e.GET("/diary-5", schedule.ScheduleHand5)
	e.GET("/diary-6", schedule.ScheduleHand6)
	e.GET("/diary-7", schedule.ScheduleHand7)
	e.GET("/diary-8", schedule.ScheduleHand8)
	// Итоговые отметки
	e.GET("/final-grades", schedule.StudentFinalGrades, RoleMiddleware(store, "admin", "student"))
	e.GET("/insert-final-grades", dairyfinalgrades.FinalGradesHandler, RoleMiddleware(store, "admin", "teacher"))
	e.POST("/insert-final-grades", dairyfinalgrades.FinalGradesHandler)
	// Дополнительные функции
	go homepage.CleanupInactiveUsers()
	// Открытие базы данных
	database.OpenDataBase(database.Databasename)

	// Запуск сервера
	if err := e.Start(":80"); err != nil {
		log.Printf("Сервер не запущен: %v", err)
	}
}

func RoleMiddleware(store *sessions.CookieStore, requiredRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session, err := store.Get(c.Request(), "session-name")
			if err != nil {
				log.Printf("Ошибка получения сессии: %v", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Ошибка сервера"})
			}

			role, ok := session.Values["userRole"].(string)
			if !ok {
				log.Println("Пользователь не имеет роли")
				return c.JSON(http.StatusForbidden, map[string]string{"message": "Доступ запрещен"})
			}

			// Проверка на наличие роли в массиве разрешенных ролей
			for _, r := range requiredRoles {
				if role == r {
					return next(c) // если хотя бы одна роль совпадает, продолжаем
				}
			}

			// Если ни одна роль не совпала
			log.Printf("Проверка роли: %s, Требуемые роли: %v", role, requiredRoles)
			return c.JSON(http.StatusForbidden, map[string]string{"message": "Доступ запрещен"})
		}
	}
}
