package main

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func downloadPresentation(w http.ResponseWriter, r *http.Request) {
	filePath := "/Users/admin/Documents/work/word-deliv/presentation.pdf"
	fileName := filepath.Base(filePath)

	log.Printf("Попытка открыть файл: %s", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Ошибка при открытии файла: %v", err)
		http.Error(w, "Файл не найден.", http.StatusNotFound)
		return
	}
	defer file.Close()

	log.Printf("Файл найден, отправка файла: %s", fileName)

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	w.Header().Set("Content-Type", "application/pdf")

	http.ServeFile(w, r, filePath)
}

func sendEmail(to, subject, body string) error {
	log.Println("Инициализация отправки письма.")

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}

	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		log.Printf("Ошибка при отправке письма: %v", err)
	} else {
		log.Println("Письмо успешно отправлено.")
	}
	return err
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodPost {
		log.Println("Получен POST-запрос на отправку формы.")

		name := r.FormValue("name")
		phone := r.FormValue("phone")
		email := r.FormValue("email")
		message := r.FormValue("message")

		if name == "" || phone == "" || email == "" || message == "" {
			log.Println("Ошибка: не все поля формы заполнены.")
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

		log.Println("Форма успешно обработана.")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status": "success", "message": "Форма успешно отправлена"}`)
	} else {
		log.Println("Ошибка: неправильный метод запроса.")
		http.Error(w, "Неправильный метод запроса.", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/submit_form", formHandler)
	http.HandleFunc("/download", downloadPresentation)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Запуск сервера на порту %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
