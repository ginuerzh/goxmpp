package main

import (
	//"bufio"
	"flag"
	"fmt"
	xmpp "github.com/ginuerzh/goxmpp"
	"github.com/ginuerzh/goxmpp/client"
	"github.com/ginuerzh/goxmpp/core"
	"github.com/ginuerzh/goxmpp/xep"
	"log"
	"os"
	//"strings"
	"crypto/tls"
	//"encoding/base64"
	//"io/ioutil"
)

var server = flag.String("server", "talk.google.com:443", "server")
var username = flag.String("username", "", "username")
var password = flag.String("password", "", "password")
var notls = flag.Bool("notls", false, "No TLS")
var debug = flag.Bool("debug", false, "debug output")

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: example [options]\n")
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()
	if *username == "" || *password == "" {
		flag.Usage()
	}

	talk := client.NewClient(*server, *username, *password,
		&client.Options{NoTLS: *notls, Debug: *debug, TlsConfig: &tls.Config{InsecureSkipVerify: true}})

	if err := talk.Init(); err != nil {
		log.Fatal(err)
	}

	exit := make(chan int, 1)

	go func() {
		talk.Run(func(st xmpp.Stan) {
			fmt.Println(st)
			if st.Name() == "message" {
				msg := st.(*xmpp.Stanza)
				body := ""
				subject := ""
				for _, e := range msg.Elements {
					if e.Name() == "body" {
						body = e.(*core.MsgBody).Body
						break
					}
					if e.Name() == "subject" {
						subject = e.(*core.MsgSubject).Subject
					}
				}
				if len(body) > 0 {
					talk.Send(xmpp.NewMessage("chat", msg.From, body, subject))
				}
			}
		})
		exit <- 1
	}()

	talk.Send(xmpp.NewStanza("presence"))

	_, err := talk.SendIQ(xmpp.NewIQ("get", client.GenId(), "", &core.RosterQuery{}))
	if err != nil {
		log.Fatal(err)
	}

	_, err = talk.SendIQ(xmpp.NewIQ("get", client.GenId(), "", &xep.DiscoItemsQuery{}))
	if err != nil {
		log.Println(err)
	}

	_, err = talk.SendIQ(xmpp.NewIQ("get", client.GenId(), "", &xep.DiscoInfoQuery{}))
	if err != nil {
		log.Println(err)
	}
	/*
		vcard := &xep.VCard{}
		_, err := talk.SendIQ(xmpp.NewIQ("get", client.GenId(), "", vcard))
		if err != nil {
			log.Println(err)
		}

		if vcard.Photo != nil {
			fmt.Println(vcard.Photo.Type)
			data, err := base64.StdEncoding.DecodeString(vcard.Photo.BinVal)
			if err != nil {
				log.Println(err)
			} else {
				ioutil.WriteFile("photo.jpg", data, os.ModePerm)
			}
		}
	*/

	<-exit
}
