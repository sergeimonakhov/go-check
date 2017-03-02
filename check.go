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

func (am *alertMessage) sentSlack() {
	err := slackhookgo.Send(
		am.url,
		slackhookgo.NewSlackMessage(
			"username",
			"backup",
		).AddAttachment(
			slackhookgo.MessageAttachment{
				Color: am.color,
				Text:  am.text,
				Title: "<!channel>",
			},
		),
	)
	checkIfError(err)
}

func checkIfError(err error) {
	if err == nil {
		return
	}
	log.Printf("error: %s", err)
}

func main() {
	var (
		sentUp   = 0
		sentDown = 0
		protocol = flag.String("protocol", "tcp", "protocol tcp/udp")
		host     = flag.String("host", "ya.ru", "destination host")
		port     = flag.String("port", "80", "destination port")
		interval = flag.Uint("interval", 5, "interval check seconds")
		url      = flag.String("url", "", "hook url")
	)

	flag.Parse()

	address := fmt.Sprintf("%s:%s", *host, *port)

	for {
		conn, err := net.DialTimeout(*protocol, address, 3*time.Second)
		checkIfError(err)
		if err == nil {
			conn.Close()
			if sentUp == 0 {
				am := alertMessage{
					color: "good",
					text:  "Destination host " + *host + ":" + *port + " reachable",
					url:   *url,
				}
				am.sentSlack()
				sentUp = 1
				sentDown = 0
			}
		} else {
			if sentDown == 0 {
				am := alertMessage{
					color: "danger",
					text:  "Destination host " + *host + ":" + *port + " unreachable",
					url:   *url,
				}
				am.sentSlack()
				sentUp = 0
				sentDown = 1
			}
		}
		time.Sleep(time.Duration(*interval) * time.Second)
	}
}
