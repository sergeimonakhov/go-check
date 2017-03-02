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
                c       string
                t       string
                lastSt  error = errors.New("zero step")
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
                if err != lastSt {
                        if err == nil {         // normal
                                c = "good"
                                t = " reachable"
                                conn.Close()
                        } else {                // not normal
                                c = "danger"
                                t = " unreachable"
                        }
                        lastSt = err            // key of success
                        am := alertMessage{
                                color:  c,
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
