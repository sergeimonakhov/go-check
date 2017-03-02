package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"
	"errors"

	"github.com/lowstz/slackhookgo"
)

type alertMessage struct {
	url   string
	text  string
	color string
}

type connection struct {
	protocol string
	host     string
	port     string
	address  string
}

func checkError(err error) {
	if err == nil {
		return
	}
	log.Printf("error: %s", err)
}

func (c *connection) conn() (net.Conn, error) {
	conn, err := net.DialTimeout(c.protocol, c.address, 3*time.Second)
	checkError(err)
	return conn, err
}

func (a *alertMessage) sentSlack() {
	err := slackhookgo.Send(
		a.url,
		slackhookgo.NewSlackMessage(
			"username",
			"backup",
		).AddAttachment(
			slackhookgo.MessageAttachment{
				Color: a.color,
				Text:  a.text,
				Title: "<!channel>",
			},
		),
	)
	checkError(err)
}

func main() {
	var (
		k 	string
		t	string
		lastSt  error = errors.New("zero step")
		protocol = flag.String("protocol", "tcp", "protocol tcp/udp")
		host     = flag.String("host", "ya.ru", "destination host")
		port     = flag.String("port", "80", "destination port")
		interval = flag.Uint("interval", 5, "interval check seconds")
		url      = flag.String("url", "", "hook url")
	)

	flag.Parse()

	for {
		c := connection{
			protocol: *protocol,
			address:  fmt.Sprintf("%s:%s", *host, *port),
		}
		conn, err := c.conn()
                if err != lastSt {
                        if err == nil {         // normal
                                k = "good"
                                t = " reachable"
                                conn.Close()
                        } else {                // not normal
                                k = "danger"
                                t = " unreachable"
                        }
                        lastSt = err            // key of success
                        am := alertMessage{
                                color:  k,
                                text:   "Destination host " + *host + ":" + *port + t,
                                url:    *url,
                        }
                        am.sentSlack()
                        fmt.Printf("the message is send\n")
                } else {
                fmt.Printf("to do - nothing\n") // to do - nothin
                }
		time.Sleep(time.Duration(*interval) * time.Second)
	}
}
