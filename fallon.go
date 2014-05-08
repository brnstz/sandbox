package main

// Check for new Jimmy Fallon tickets

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"
)

const (
	startHTMLComment = `<!--`
	stopHTMLComment  = `-->`
	jimmyUrl         = `http://www.showclix.com/event/thetonightshowstarringjimmyfallon`
	notifyEmail      = `brnstz@gmail.com`
	timeout          = time.Duration(5 * time.Second)
	lastFile         = `/tmp/jimmy`
	port             = 587
	msg              = `From: brnstz@gmail
To: brnstz@gmail.com
Subject: Jimmy Fallon tickets

%s
`
)

var (
	user       string
	pw         string
	mailserver string
)

func removeComments(body string) string {
	for {
		posStart := strings.Index(body, startHTMLComment)
		posStop := strings.Index(body, stopHTMLComment)
		stopLen := len(stopHTMLComment)

		if posStart >= 0 && posStop >= 0 && posStop > posStart {

			// Both posStop and posStop found, and stop is after start
			body = strings.Join(
				[]string{
					body[0:posStart],
					body[posStop+stopLen : len(body)],
				},
				``,
			)

		} else {

			// Stop removing comments, found the end of comments or
			// something we don't understand
			break
		}
	}

	return body
}

func mailStuff(content string) {
	auth := smtp.PlainAuth("", user, pw, mailserver)

	err := smtp.SendMail(
		fmt.Sprintf("%v:%v", mailserver, port),
		auth,
		notifyEmail,
		[]string{notifyEmail},
		[]byte(fmt.Sprintf(msg, content)),
	)

	if err != nil {
		log.Println(err)
	}
}

func readOld() (string, bool) {
	fh, err := os.Open(lastFile)
	if err != nil {
		mailStuff(err.Error())
		return "", false
	}
	defer fh.Close()

	oldBytes, err := ioutil.ReadAll(fh)
	if err != nil {
		// Possible new file, so return true
		return "", true
	}

	return string(oldBytes), true
}

func bodyDiff(body string) bool {
	old, good := readOld()

	// Write new content to file
	ioutil.WriteFile(lastFile, []byte(body), 0666)

	// Couldn't figure it out, or body is same as old
	if !good || old == body {
		return false
	}

	return true
}

func runOne() {
	resp, err := http.Get(jimmyUrl)

	if err != nil {
		mailStuff(err.Error())
		return
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	body := removeComments(string(bytes))

	if bodyDiff(body) {
		mailStuff(fmt.Sprintf("Tickets may be available! %v", jimmyUrl))
	} else {
		log.Println("Tickets not available")
	}
}

func main() {
	user = os.Getenv("JF_USER")
	pw = os.Getenv("JF_PW")
	mailserver = os.Getenv("JF_SMTP")

	for {
		runOne()
		time.Sleep(timeout)
	}
}
