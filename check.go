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
	port     string
	address  string
}

func checkError(err error) (int) {
	if err == nil {
		return 0
	} else {
		log.Printf("error: %s", err)
		return 1
	}

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
		lastE	int = -1
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
		e := checkError(err)
                if e != lastE {
                        if e == 0 {		// normal
                                k = "good"
                                t = " reachable"
                                conn.Close()
                        } else {		// not normal
                                k = "danger"
                                t = " unreachable"
                        }
                        lastE = e		// key of success
                        am := alertMessage{
                                color:  k,
                                text:   "Destination host " + *host + ":" + *port + t,
                                url:    *url,
                        }
                        am.sentSlack()
                        fmt.Printf("%s\n", "\"the message is send\"")
                } else {
                fmt.Printf("%s\n", "\"to do - nothing\"") // to do - nothin
                }
		time.Sleep(time.Duration(*interval) * time.Second)
	}
}
