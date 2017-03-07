package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

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
	port     int
	address  string
}

func checkError(err error) int {
	if err == nil {
		return 0
	}
	log.Printf("error: %s", err)
	return 1
}

func (c *connection) conn() (net.Conn, int) {
	conn, err := net.DialTimeout(c.protocol, c.address, 3*time.Second)
	errInt := checkError(err)
	return conn, errInt
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
		k        string
		t        string
		lastE    int
		protocol = flag.String("protocol", "tcp", "protocol tcp/udp")
		host     = flag.String("host", "ya.ru", "destination host")
		port     = flag.Int("port", 80, "destination port")
		interval = flag.Uint("interval", 5, "interval check seconds")
		url      = flag.String("url", "", "hook url")
	)

	flag.Parse()

	for {
		c := connection{
			protocol: *protocol,
			address:  fmt.Sprintf("%s:%v", *host, *port),
		}
		conn, err := c.conn()
		if err != lastE {
			if err == 0 { // normal
				k = "good"
				t = "reachable"
				conn.Close()
			} else { // not normal
				k = "danger"
				t = "unreachable"
			}
			lastE = err // key of success
			am := alertMessage{
				color: k,
				text:  fmt.Sprintf("Destination host %s:%v %s\n", *host, *port, t),
				url:   *url,
			}
			am.sentSlack()
		}
		time.Sleep(time.Duration(*interval) * time.Second)
	}
}
