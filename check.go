package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/lowstz/slackhookgo"
)

func checkIfError(err error) {
	if err == nil {
		return
	}

	//fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	log.Printf("error: %s", err)
	//os.Exit(1)
}

func sendToSlack(color, text, url *string) {
	var attachment slackhookgo.MessageAttachment
	attachment.Color = *color
	attachment.Text = *text
	attachment.Title = "<!channel>"
	msg := slackhookgo.NewSlackMessage("slack-bot", "backup")
	msg.IconEmoji = ":exclamation:"
	msg.AddAttachment(attachment)
	err := slackhookgo.Send(*url, msg)
	checkIfError(err)
}

func checkFlags() {
	flag.Parse()
	/*	if *host == "" {
			flag.PrintDefaults()
			log.Fatal("host missing, exiting.")
		}

		if *port == "" {
			flag.PrintDefaults()
			log.Fatal("port missing, exiting.")
		}
	*/
}

func main() {
	var color string
	sentUp := 0
	sentDown := 0
	protocol := flag.String("protocol", "tcp", "protocol tcp/udp")
	host := flag.String("host", "ya.ru", "destination host")
	port := flag.String("port", "80", "destination port")
	interval := flag.Uint("interval", 5, "interval check seconds")
	url := flag.String("url", "", "hook url")

	flag.Parse()
	checkFlags()

	address := fmt.Sprintf("%s:%s", *host, *port)

	for {
		conn, err := net.DialTimeout(*protocol, address, 3*time.Second)
		checkIfError(err)
		if err == nil {
			conn.Close()
			if sentUp == 0 {
				color = "good"
				text := "Destination host " + *host + ":" + *port + " reachable"
				sendToSlack(&color, &text, url)
				sentUp = 1
				sentDown = 0
			}
		} else {
			if sentDown == 0 {
				color = "danger"
				text := "Destination host " + *host + ":" + *port + " unreachable"
				sendToSlack(&color, &text, url)
				sentUp = 0
				sentDown = 1
			}
		}
		time.Sleep(time.Duration(*interval) * time.Second)
	}
}
