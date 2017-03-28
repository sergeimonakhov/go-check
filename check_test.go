package main

import (
	"testing"
	"errors"
	"os"
)

func TestCheckError(t *testing.T) {
	var (
		e error = nil
		r bool = false
	)

	if r != CheckError(e) { // first stage: nil = false
		t.Error("%d != %d", r, e)
	}

	e = errors.New("test error")
	r = true

	if r != CheckError(e) { // second stage: "test error" = true
		t.Error("%d != %d", r, e)
	}
}

func TestConn(t *testing.T) {
	c := connection {
		protocol: "tcp",
		address:  "127.0.0.1:0", // the value should not work
	}

	conn, e := c.Conn()
	if e != true { // get the error
		t.Error("%d != %d", e, true)
	}

	if conn != nil { // don't get the connection
                t.Error("%d != %d", conn, nil)
        }
}

func TestSentSlack(t *testing.T) {
	am := alertMessage {
		url:	 "https://127.0.0.1:0/nothere", // the value should not work
		text: 	 "mesg text",
		color: 	 "good",
		channel: "general",
	}

	e := am.SentSlack()
	if e != true { // get the error
                t.Error("%d != %d", e, true)
        }

}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}
