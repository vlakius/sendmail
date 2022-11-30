package main

import (
	"flag"
	"fmt"

	"crypto/tls"

	"gopkg.in/gomail.v2"
)

type Mail struct {
	From     string
	To       []string
	Subject  string
	Body     string
	Attach   string
	Smtp     string
	Port     int
	Username string
	Password string
}

func printUsage() {
	usage :=
		`
LBGoMail v0.0.1 - Semplice strumento CLI per inviare email
---------------------------------------------------------
Utilizzo:
lbgomail --from <mittente> --to <destinatario> --subject <oggetto> --body <corpo> --smtp <smtp> --port <porta>
Opzioni:
--from <mittente>		Indirizzo email mittente
--to <destinatario>		Indirizzo email destinatario
--subject <oggetto>		Oggetto email
--body <corpo>			Corpo email
--attach <allegato>		Allegato email, inserire la path del file da allegare (opzionale)
--smtp <smtp>			SMTP server
--port <porta>			Porta server (default 25)
--username <username>	Username (opzionale)
--password <password>	Password (opzionale)

	`
	fmt.Println(usage)
}

func sendMail(mail Mail) error {
	m := gomail.NewMessage()
	m.SetHeader("From", mail.From)
	m.SetHeader("To", mail.To...)
	m.SetHeader("Subject", mail.Subject)
	m.SetBody("text/html", mail.Body)
	if mail.Attach != "" {
		m.Attach(mail.Attach)
	}

	d := gomail.NewDialer(mail.Smtp, mail.Port, mail.Username, mail.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	d.TLSConfig.InsecureSkipVerify = true
	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("Errore invio email: %v\n", err)
		printUsage()
		return err
	}
	return nil
}

func parseArgs() Mail {
	// parse args from cli and return a Mail struct
	from := flag.String("from", "", "Mittente email")
	to := flag.String("to", "", "Destinatari/o email (valori multipli separati da virgola)")
	subject := flag.String("subject", "", "Oggetto email")
	body := flag.String("body", "", "Corpo email")
	attach := flag.String("attach", "", "Allegato email (inserire la path)")
	smtp := flag.String("smtp", "", "SMTP server")
	port := flag.Int("port", 25, "Porta server")
	username := flag.String("username", "", "username")
	password := flag.String("password", "", "password")
	flag.Parse()

	return Mail{
		From:     *from,
		To:       []string{*to},
		Subject:  *subject,
		Body:     *body,
		Attach:   *attach,
		Smtp:     *smtp,
		Port:     *port,
		Username: *username,
		Password: *password,
	}

}

func main() {
	mail := parseArgs()
	if mail.From == "" || mail.To[0] == "" || mail.Smtp == "" {
		printUsage()
		return
	}

	err := sendMail(mail)
	if err != nil {
		printUsage()
		return
	}
	for _, v := range mail.To {
		fmt.Printf("Email inviata con successo all' indirizzo %s\n", v)

	}

}
