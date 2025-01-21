package homepage

import (
	"dairyapp/database"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

type PageData struct {
	UserName string
	Email    string
	Password string
	Date     string
}
type Session struct {
	ID         string
	Username   string
	LastActive time.Time
}

var (
	store       *sessions.CookieStore
	activeUsers = make(map[string]time.Time)
	mu          sync.Mutex
)

func InitStore(s *sessions.CookieStore) {
	store = s
}

func authenticate(c echo.Context) (*sessions.Session, error) {
	session, err := store.Get(c.Request(), "session-name")
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Check if the user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		return nil, echo.NewHTTPError(http.StatusForbidden, map[string]string{"message": "Нет доступа", "redirect": "/login"})
	}

	// Get the userID from the session
	userIDInterface := session.Values["userID"]
	userID, ok := userIDInterface.(string) // Type assert to string
	if !ok {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Неверный формат ID пользователя")
	}

	// Retrieve user by ID
	username, err := database.GetUserByID(userID, database.Databasename)
	if err != nil {
		log.Printf("Ошибка выполнения GetUserByID - %v", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Ошибка получения пользователя")
	}

	// Update the active users map with the current time
	activeUsers[username] = time.Now()

	// Save session
	if err := session.Save(c.Request(), c.Response()); err != nil {
		log.Printf("Ошибка сохранения сессии - %v", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Ошибка сохранения сессии")
	}

	return session, nil
}

// ActiveUsersHandler retrieves and returns the list of active users
func ActiveUsersHandler(c echo.Context) error {
	users := make([]string, 0, len(activeUsers))
	for username := range activeUsers {
		users = append(users, username)
	}
	return c.JSON(http.StatusOK, users)
}
func CleanupInactiveUsers() {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()
	for userID, lastActive := range activeUsers {
		if now.Sub(lastActive) > 20*time.Minute { // Удаляем неактивного пользователя
			delete(activeUsers, userID)
			log.Printf("Пользователь очищен - %v", userID)
		}
	}
	log.Printf("Cleanup закончен")
	log.Printf("Список пользователей - %v", activeUsers)
}
func renderPage(c echo.Context, userID string, templateFile string) error {
	username, err := database.GetUserByID(userID, database.Databasename)
	if err != nil {
		log.Printf("Ошибка получения username из базы данных: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Ошибка получения данных пользователя"})
	}

	data := PageData{
		UserName: username,
	}

	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		log.Printf("Ошибка отображения шаблона: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Ошибка отображения страницы"})
	}

	// Render the template
	if err = tmpl.Execute(c.Response().Writer, data); err != nil {
		log.Printf("Ошибка отображения шаблона: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Ошибка отображения страницы"})
	}

	return nil
}
func DiaryPage(c echo.Context) error {
	tmpl, err := template.ParseFiles("templates/diaryapp.html")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	tmpl.Execute(c.Response(), nil)
	return c.Render(http.StatusOK, "diaryapp.html", nil)
}
func SendError(c echo.Context) error {
	tmpl, err := template.ParseFiles("templates/send-error.html")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	tmpl.Execute(c.Response(), nil)
	return c.Render(http.StatusOK, "send-error.html", nil)
}

// StudentHomepage handles rendering the homepage for students
func StudentHomepage(c echo.Context) error {
	_, err := authenticate(c)
	if err != nil {
		return err
	}

	session, err := store.Get(c.Request(), "session-name")
	if err != nil {
		return err
	}

	userID, ok := session.Values["userID"].(string)
	if !ok {
		log.Println("не удалось извлечь userID из сессии")
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "не удалось извлечь userID из сессии"})
	}

	log.Printf("Homepage for student: userID=%s", userID)
	return renderPage(c, userID, "templates/student-home.html")
}

func TeacherHomepage(c echo.Context) error {
	_, err := authenticate(c)
	if err != nil {
		return err
	}
	session, err := store.Get(c.Request(), "session-name")
	if err != nil {
		return err
	}

	userID, ok := session.Values["userID"].(string)
	if !ok {
		log.Println("не удалось извлечь userID из сессии")
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "не удалось извлечь userID из сессии"})
	}
	log.Printf("Homepage for teacher: %v", err)
	return renderPage(c, userID, "templates/teacher-home.html")
}
func AccountHomepage(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		_, err := authenticate(c)
		if err != nil {
			return err
		}
		session, err := store.Get(c.Request(), "session-name")
		if err != nil {
			return err
		}

		userID, ok := session.Values["userID"].(string)
		if !ok {
			log.Println("не удалось извлечь userID из сессии")
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "не удалось извлечь userID из сессии"})
		}
		email, ok := session.Values["userEmail"].(string)
		if !ok {
			log.Println("не удалось извлечь userEmail из сессии")
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "не удалось извлечь userID из сессии"})
		}
		password, ok := session.Values["userPassword"].(string)
		if !ok {
			log.Println("не удалось извлечь userPassword из сессии")
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "не удалось извлечь userID из сессии"})
		}
		username, err := database.GetUserByID(userID, database.Databasename)
		if err != nil {
			log.Printf("Ошибка получения username из базы данных: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Ошибка получения данных пользователя"})
		}
		date, err := database.GetDate(database.Databasename, userID)
		if err != nil {
			return err
		}
		log.Printf("%v", date)
		data := PageData{
			UserName: username,
			Email:    email,
			Password: password,
			Date:     date,
		}

		tmpl, err := template.ParseFiles("templates/account.html")
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
