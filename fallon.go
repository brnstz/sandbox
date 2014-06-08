package main

// Check for new Jimmy Fallon tickets

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	twilio           = `https://api.twilio.com/2010-04-01/Accounts/%s/Calls`
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
	user,
	pw,
	mailserver,
	twilioSid,
	twilioAuth,
	twilioFrom,
	twilioTo string
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
		callTwilio()
	} else {
		log.Println("Tickets not available")
	}
}

func callTwilio() error {
	vals := url.Values{}
	// FIXME: make real URL for Twilio to hit?
	vals.Set("Url", "http://www.brnstz.com/")
	vals.Set("To", twilioTo)
	vals.Set("From", twilioFrom)
	valStr := vals.Encode()

	req, err := http.NewRequest("POST", fmt.Sprintf(twilio, twilioSid),
		strings.NewReader(valStr))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(valStr)))

	if err != nil {
		log.Printf("Couldn't create Twilio req: %v", err)
		return err
	}

	req.SetBasicAuth(twilioSid, twilioAuth)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Printf("Couldn't run Twilio req: %v", err)
		return err
	}

	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Couldn't read Twilio resp: %v", err)
		return err
	}
	log.Printf("Twilio response:\n%s\n", all)

	resp.Body.Close()

	return nil
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
	twilioSid = os.Getenv("TWILIO_SID")
	twilioAuth = os.Getenv("TWILIO_AUTH")
	twilioFrom = os.Getenv("TWILIO_FROM")
	twilioTo = os.Getenv("TWILIO_TO")

	for {
		runOne()
		time.Sleep(timeout)
	}
}
