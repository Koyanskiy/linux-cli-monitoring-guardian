package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type LogEntry struct {
	ID        int
	CreatedAt time.Time
	Message   string
}

func ConnectDB() {
	connStr := "host=127.0.0.1 port=5432 user=ivory_user password=secure_ivory_pass dbname=ivbase sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("\033[31m❌ Ошибка настройки базы данных!\033[0m")
		waitForKey()
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("\033[31m❌ База данных ivbase недоступна или пароль не подошел!\033[0m")
		waitForKey()
		return
	}

	fmt.Println("\033[32m✅ Успешное подключение к PostgreSQL базе ivbase!\033[0m")
	waitForKey()
}

func CreateLogTable() {
	connStr := "host=127.0.0.1 port=5432 user=ivory_user password=secure_ivory_pass dbname=ivbase sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return
	}
	defer db.Close()

	query := `
	CREATE TABLE IF NOT EXISTS system_logs (
		id SERIAL PRIMARY KEY,
		created_at TIMESTAMP DEFAULT NOW(),
		message TEXT
	);`

	_, err = db.Exec(query)
	if err != nil {
		fmt.Println("\033[31m❌ Не удалось создать таблицу логов в базе данных!\033[0m")
		fmt.Println("Ошибка:", err)
		waitForKey()
		return
	}
}

func InsertLogIntoDB(message string) {
	connStr := "host=127.0.0.1 port=5432 user=ivory_user password=secure_ivory_pass dbname=ivbase sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("⚠️ Ошибка подключения к базе при записи лога:", err)
		return
	}
	defer db.Close()

	query := "INSERT INTO system_logs (message) VALUES ($1);"
	_, err = db.Exec(query, message)
	if err != nil {
		fmt.Println("⚠️ Не удалось сохранить лог в PostgreSQL:", err)
		return
	}
	fmt.Println("\033[32m💾 Сбой сети успешно залогирован в базу данных PostgreSQL\033[0m")
}

func ClearLogTable() {
	connStr := "host=127.0.0.1 port=5432 user=ivory_user password=secure_ivory_pass dbname=ivbase sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("\033[31m❌ Ошибка подключения к базе данных!\033[0m")
		waitForKey()
		return
	}
	defer db.Close()

	_, err = db.Exec("TRUNCATE TABLE system_logs RESTART IDENTITY;")
	if err != nil {
		fmt.Println("\033[31m❌ Не удалось очистить таблицу в базе данных!\033[0m")
		fmt.Println("Ошибка:", err)
		waitForKey()
		return
	}

	fmt.Println("\033[32m🧹 Таблица system_logs в PostgreSQL успешно очищена\033[0m")
	waitForKey()

}

func PrintLogsFromDB() {
	connStr := "host=127.0.0.1 port=5432 user=ivory_user password=secure_ivory_pass dbname=ivbase sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("⚠️ Ошибка подключения к базе при чтении логов:", err)
		waitForKey()
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, created_at, message FROM system_logs ORDER BY id;")
	if err != nil {
		fmt.Println("⚠️ Не удалось получить логи из базы данных:", err)
		waitForKey()
		return
	}
	defer rows.Close()

	fmt.Println("\n\033[1;36m--- ИСТОРИЯ СБОЕВ ИЗ БАЗЫ ДАННЫХ PostgreSQL ---\033[0m")

	hasLogs := false

	for rows.Next() {
		hasLogs = true
		var id int
		var createdAt string
		var message string

		err := rows.Scan(&id, &createdAt, &message)
		if err != nil {
			fmt.Println("⚠️ Ошибка чтения строки лога:", err)
			continue
		}

		cleanTime := strings.Replace(createdAt, "T", " ", 1)
		cleanTime = strings.Replace(cleanTime, "Z", "", 1)
		if len(cleanTime) > 19 {
			cleanTime = cleanTime[:19]
		}

		fmt.Printf("🆔 [%d] 📅 [%s] 🚨 %s\n", id, cleanTime, message)
	}

	if !hasLogs {
		fmt.Println("📭 База данных пуста. Критических сбоев не обнаружено!")
	}

	fmt.Println()
	waitForKey()

}
