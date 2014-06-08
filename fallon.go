package main

// Check for new Jimmy Fallon tickets

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	startHTMLComment = `<!--`
	stopHTMLComment  = `-->`
	jimmyUrl         = `http://www.showclix.com/event/thetonightshowstarringjimmyfallon`
	notifyEmail      = `brnstz@gmail.com`
	timeout          = time.Duration(30 * time.Minute)
	lastFile         = `/tmp/jimmy`
	port             = 587
	msg              = `From: brnstz@gmail.com
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

func bodyDiff(body string) (string, error) {
	old, good := readOld()

	if !good {
		return "", fmt.Errorf("Can't read old file")
	}

	// Write new content to file
	ioutil.WriteFile(lastFile, []byte(body), 0666)

	return diff(old, body)
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

	d, err := bodyDiff(body)

	// err non-nil indicates there is a diff
	if err != nil {
		mailStuff(fmt.Sprintf("Tickets may be available!\r\n%+v\r\n%v\r\n%v", err, d, jimmyUrl))
	} else {
		log.Println("Tickets not available")
	}
}

// Stolen/modified from https://github.com/bradfitz/camlistore/blob/master/pkg/test/diff.go
// If there is a difference between and a and b, return the diff and an
// error value. Otherwise, return blank string and nil error.
// We send -b option which ignores changes in the amount of whitespace.
func diff(a, b string) (string, error) {
	if a == b {
		return "", nil
	}
	ta, err := ioutil.TempFile("", "")
	if err != nil {
		return "", err
	}
	tb, err := ioutil.TempFile("", "")
	if err != nil {
		return "", err
	}
	defer os.Remove(ta.Name())
	defer os.Remove(tb.Name())
	ta.WriteString(a)
	tb.WriteString(b)
	ta.Close()
	tb.Close()
	out, err := exec.Command("diff", "-bu", ta.Name(), tb.Name()).CombinedOutput()
	if len(out) > 0 && err == nil {
		// There is a diff
		return string(out), fmt.Errorf("There is a diff")
	} else if err != nil {
		// There is an error
		return string(out), err
	} else {
		// No diff and no error
		return "", nil
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
