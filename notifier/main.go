package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	gomail "gopkg.in/mail.v2"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Environment variables from docker-compose
	dbHost := os.Getenv("MYSQL_HOST")
	dbName := os.Getenv("MYSQL_DATABASE")
	dbUser := os.Getenv("MYSQL_USER")
	dbPass := os.Getenv("MYSQL_PASSWORD")
	smtpServer := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	emailTo := os.Getenv("EMAIL_TO")
	checkInterval := os.Getenv("CHECK_INTERVAL")
	intervalDuration, err := time.ParseDuration(checkInterval)
	if err != nil {
		log.Fatalf("incorrect interval: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer db.Close()

	for {
		sendSummary(db, smtpServer, smtpPort, smtpUser, smtpPass, emailTo)
		log.Printf("Summary email sent. Sleeping %s ...\n", checkInterval)

		time.Sleep(intervalDuration)
	}
}

func sendSummary(db *sql.DB, smtpServer, smtpPort, smtpUser, smtpPass, emailTo string) {
	rows, err := db.Query("SELECT title FROM todos_todo WHERE isCompleted = 0")
	if err != nil {
		log.Printf("query error: %v", err)
		return
	}
	defer rows.Close()

	var todos []string
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			continue
		}
		todos = append(todos, fmt.Sprintf("- %s", title))
	}

	if len(todos) == 0 {
		log.Println("No pending todos — skipping email.")
		return
	}

	body := "Here are your pending tasks:<br><br>" + strings.Join(todos, "<br>")

	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", smtpUser)
	message.SetHeader("To", emailTo)
	message.SetHeader("Subject", "Unfinished TODO List"+time.Now().String())

	// Set email body
	message.SetBody("text/html", `
        <html>
            <body>
               `+body+`
            </body>
        </html>
    `)
	// Set up the SMTP dialer
	dialer := gomail.NewDialer(smtpServer, 587, smtpUser, smtpPass)

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		log.Printf("Failed to send message: %v", err)
		return
	}
}
