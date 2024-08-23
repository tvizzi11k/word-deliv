package main

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
)

func sendEmail(to, subject, body string) error {
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	return err
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	// Настройка CORS заголовков
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Обработка OPTIONS-запроса для CORS
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Обработка запроса POST
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		phone := r.FormValue("phone")
		email := r.FormValue("email")
		message := r.FormValue("message")

		if name == "" || phone == "" || email == "" || message == "" {
			http.Error(w, "Все поля формы обязательны к заполнению.", http.StatusBadRequest)
			return
		}

		subject := "Новая заявка с сайта World Delivery"
		body := fmt.Sprintf("Имя: %s\nТелефон: %s\nПочта: %s\nСообщение:\n%s", name, phone, email, message)

		to := "zgalkin319@gmail.com"

		err := sendEmail(to, subject, body)
		if err != nil {
			log.Printf("Ошибка отправки письма: %v", err)
			http.Error(w, "Ошибка отправки заявки.", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status": "success", "message": "Форма успешно отправлена"}`)
	} else {
		http.Error(w, "Неправильный метод запроса.", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/submit_form", formHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Запуск сервера на порту %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
