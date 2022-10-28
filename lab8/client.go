package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/smtp"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"strings"
)

type Mail struct {
	Sender  string
	To      []string
	Subject string
	Body    string
}

func main() {
	db, err1 := sql.Open("mysql", "iu9networkslabs:Je2dTYr6@tcp(students.yss.su)/iu9networkslabs")
	if err1 != nil {
		fmt.Println("\033[38;5;125mGOT ERROR TRYING TO CONNECT TO THE DATABASE\u001B[0m")
	}

	user := "*email*"
	password := "*app password*"
	subject := "*subject*"

	sender := "*email*"
	rows, err := db.Query("select email, name from *table*)
	defer rows.Close()
	var text string
	var name string
	for rows.Next() {
		to := []string{}

		err := rows.Scan(&text, &name)
		fmt.Println(text)
		to = append(to, text)
		body := "<!DOCTYPE html>\n" +
			"<html>\n" +
			"<head>\n" +
			"<style>\n" +
			"table, th, td, h2, html {\n" +
			"background-color: rgb(150, 212, 212);\n" +
			"}\n" +
			"</style>\n" +
			"</head>\n" +
			"<body>\n\n" +
			"<table style=\"width:100%\">\n" +
			"<tr>\n" +
			"<td> " +
			"<b style=\"font-size:30px\"> " +
			"Здравствуйте, " +
			name +
			"</b>" +
			"</td>\n" +
			"</tr>\n" +
			"<tr>\n" +
			"<td> " +
			"<i style=\"font-size:20px\"> " +
			"А я содержательная часть" +
			"</i>" +
			"</td>\n" +
			"</tr>\n" +
			"</table>\n\n" +
			"</body>\n" +
			"</html>\n\n\n\n"

		request := Mail{
			Sender:  sender,
			To:      to[len(to)-1:],
			Subject: subject,
			Body:    body,
		}

		addr := "smtp.gmail.com:587"
		host := "smtp.gmail.com"
		if err != nil {
			fmt.Println(err)
		}
		msg := BuildMessage(request)
		auth := smtp.PlainAuth("", user, password, host)
		err2 := smtp.SendMail(addr, auth, sender, to, []byte(msg))

		if err != nil {
			log.Fatal(err2)
		}
		time.Sleep(30 * time.Second)
	}
	fmt.Println(err)

	//fmt.Println(msg)
}

func BuildMessage(mail Mail) string {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", mail.Sender)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)

	return msg
}
