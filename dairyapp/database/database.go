package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// Databasename - имя базы данных.
var Databasename string = "DAIRYDATABASE.db"

type User struct {
	ID             string
	Name           string
	Email          string
	HashedPassword string
	Role           string
}
type Grades struct {
	ID          string
	Class       string
	Grade       string
	Subject     string
	Work_type   string
	Date        string
	ID_Student  string
	StudentName string
}

func LoginUser(dbname, email, password string) (*User, bool, error) {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return &User{}, false, err
	}
	defer db.Close()

	var user User
	query := `SELECT id, password, role FROM users WHERE email = ?`
	err = db.QueryRow(query, email).Scan(&user.ID, &user.HashedPassword, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return &user, false, nil
		}
		return &user, false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err == nil {
		log.Printf("Успешное выполнение авторизации: %v", err)
		return &user, true, nil
	}

	return &user, false, nil
}

// RegisterUser регистрирует нового пользователя.
func RegisterUser(dbname, username, email, password, role, class string) error {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return err
	}
	defer db.Close()
	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	stmt, err := db.Prepare("INSERT INTO users(username, email, password, role, class) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, email, hashedPassword, role, class)
	return err
}

// Проверка пользователя (уникален ли его username и email)
func VerificationUser(dbname, username, email string) (error, bool) {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return err, false
	}
	defer db.Close()

	var count int
	query := "SELECT COUNT(*) FROM users WHERE username = ? OR email = ?"
	err = db.QueryRow(query, username, email).Scan(&count)
	if err != nil {
		return err, false
	}

	if count != 0 {
		return err, false
	}
	log.Printf("Верификация true")
	return err, true
}

// OpenDataBase открывает соединение с базой данных и возвращает его.
func OpenDataBase(dbname string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return nil, err
	}

	// Создаем таблицу users
	sqlStmtUsers := `
	CREATE TABLE IF NOT EXISTS users (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	username TEXT NOT NULL UNIQUE,
    	email TEXT NOT NULL UNIQUE,
    	password TEXT NOT NULL,
		role TEXT NOT NULL,
		class TEXT
	);`
	_, err = db.Exec(sqlStmtUsers)
	if err != nil {
		db.Close()
		return nil, err
	}
	log.Println("Создана/обновлена таблица users")

	// Создаем таблицу grades
	sqlStmtGrades := `
	CREATE TABLE IF NOT EXISTS grades (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		grade TEXT NOT NULL,
		subject TEXT NOT NULL,
		work_type TEXT,
		date TEXT NOT NULL,
		id_student INTEGER NOT NULL,
		id_teacher INTEGER NOT NULL,
		FOREIGN KEY(id_teacher) REFERENCES users(id)
		FOREIGN KEY(id_student) REFERENCES users(id)
	);`
	_, err = db.Exec(sqlStmtGrades)
	if err != nil {
		log.Printf("Ошибки grades %v", err)
		db.Close()
		return nil, err
	}
	log.Println("Создана/обновлена таблица grades") // Исправлено

	return db, nil // Возвращаем открытое соединение и nil для ошибки
}
func GetUserByID(user_id string, dbname string) (string, error) {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		log.Printf("Ошибка открытия базы данных в GetUserByID: %v", err)
		return "", err
	}
	defer db.Close()

	var username string
	query := `SELECT username FROM users WHERE id = ?`
	err = db.QueryRow(query, user_id).Scan(&username)
	if err != nil {
		log.Printf("Ошибка выполнения GetUserByID: %v", err)
		return "", nil
	}
	log.Printf("Успешное выполнение GetUserByID: %v", err)
	return username, err
}
func GetClassByID(user_id string, dbname string) (string, error) {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		log.Printf("Ошибка открытия базы данных в GetClassByID: %v", err)
		return "", err
	}
	defer db.Close()

	var class string
	query := `SELECT class FROM users WHERE id = ?`
	err = db.QueryRow(query, user_id).Scan(&class)
	if err != nil {
		log.Printf("Ошибка выполнения GetClassByID: %v", err)
		return "", nil
	}
	log.Printf("Успешное выполнение GetClassByID: %v", err)
	return class, err
}
func GetUserByName(username, dbname string) (*User, error) {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		log.Printf("Ошибка открытия базы данных в GetUserByName: %v", err)
		return nil, err // Возвращаем nil вместо пустого объекта
	}
	defer db.Close()

	var user User
	query := `SELECT id, username FROM users WHERE username = ?` // Получаем id и username
	err = db.QueryRow(query, username).Scan(&user.ID, &user.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Пользователь не найден: %s", username)
			return nil, nil // Возвращаем nil, если пользователь не найден
		}
		log.Printf("Ошибка выполнения GetUserByName: %v", err)
		return nil, err
	}
	log.Printf("Успешное выполнение GetUserByName: найден пользователь %s, id: %v", user.Name, user.ID)
	return &user, nil // Возвращаем указатель на найденного пользователя
}
func InsertGrades(dbname, grade, subject, work_type, date, name_student, id_teacher string) error {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return err
	}
	defer db.Close()

	var user *User
	user = &User{}
	log.Printf("Name is %v", name_student)
	user, err = GetUserByName(name_student, dbname)
	if err != nil {
		return err
	}

	stmt, err := db.Prepare("INSERT INTO grades(grade, subject, work_type, date, id_student, id_teacher) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	log.Printf("User is %v, %v", user.ID, user.Name)

	_, err = stmt.Exec(grade, subject, work_type, date, user.ID, id_teacher)
	if err != nil {
		return err
	}

	return nil
}
func GetUsers(dbname, class string) ([]User, error) {
	// Открываем соединение с базой данных
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return nil, fmt.Errorf("ошибка при открытии базы данных: %w", err)
	}
	defer db.Close() // Закрываем соединение в конце функции

	// SQL-запрос для выборки пользователей
	query := `SELECT id, username FROM users WHERE role = 'student' AND class = ?`

	// Выполняем запрос к базе данных
	rows, err := db.Query(query, class)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer rows.Close() // Закрываем rows после использования

	// Инициализируем слайс пользователей
	users := []User{}

	// Проходим по строкам результата
	for rows.Next() {
		var user User
		// Сканируем строки в структуру
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строки: %w", err)
		}
		// Добавляем пользователя в слайс
		users = append(users, user)
	}

	// Проверка на ошибки после завершения обработки результата
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка во время итерации: %w", err)
	}

	return users, nil // Возвращаем список пользователей
}
func AllGrades(dbname, id string) ([]Grades, error) {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return nil, fmt.Errorf("ошибка при открытии базы данных: %w", err)
	}
	defer db.Close()
	query := `SELECT id, grade, subject, work_type, date, id_student FROM grades WHERE id_teacher = ? `
	// Выполняем запрос к базе данных
	rows, err := db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer rows.Close() // Закрываем rows после использования

	grades := []Grades{}
	for rows.Next() {
		var grade Grades
		// Сканируем строки в структуру
		if err := rows.Scan(&grade.ID, &grade.Grade, &grade.Subject, &grade.Work_type, &grade.Date, &grade.ID_Student); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строки: %w", err)
		}
		grade.StudentName, err = GetUserByID(grade.ID_Student, dbname)
		if err != nil {
			log.Printf("Error in AllGrades %v", err)
		}
		grade.Class, err = GetClassByID(grade.ID_Student, dbname)
		if err != nil {
			log.Printf("Error in AllGrades %v", err)
		}
		grades = append(grades, grade)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка во время итерации: %w", err)
	}
	log.Printf("Grades is %v", grades)
	return grades, nil
}
func DeleteGrades(dbname, id string) error {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("DELETE FROM grades WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
func GetSchedules(dbname, date, id string) ([]Grades, error) {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	query := `SELECT grade, subject, work_type, date FROM grades WHERE id_student = ? AND date = ? `
	// Выполняем запрос к базе данных
	rows, err := db.Query(query, id, date)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer rows.Close()
	grades := []Grades{}
	for rows.Next() {
		var grade Grades
		// Сканируем строки в структуру
		if err := rows.Scan(&grade.Grade, &grade.Subject, &grade.Work_type); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строки: %w", err)
		}
		log.Printf("%v", grade)
		grades = append(grades, grade)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка во время итерации: %w", err)
	}
	log.Printf("Grades is %v", grades)
	return grades, nil
}
