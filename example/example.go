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
	"github.com/ginuerzh/gosocks5"
	"io"
	"net"
)

var server = flag.String("server", "talk.google.com:443", "server")
var username = flag.String("username", "", "username")
var password = flag.String("password", "", "password")
var notls = flag.Bool("notls", false, "No TLS")
var debug = flag.Bool("debug", false, "debug output")

var fileName = "recvFile"

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
			stanza := st.(*xmpp.Stanza)
			if st.Name() == "message" {
				msg := stanza
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
			if st.Name() == "iq" {
				iq := stanza
				for _, e := range iq.Elements {
					if e.Name() == "ping" {
						talk.Send(xmpp.NewIQ("result", iq.Id(), "", nil))
					}
					if e.Name() == "si" {
						fileName = e.(*xep.SI).File.Name

						submit := xep.NewFormData("submit", "", "",
							xep.NewFormField("", "", "stream-method", "",
								[]string{xmpp.NSByteStreams}, false, nil))
						si := xep.NewSI("", "", "",
							nil,
							xep.NewFeatureNeg(submit))
						talk.Send(xmpp.NewIQ("result", iq.Id(), iq.From, si))
					}

					if e.Name() == xmpp.NSDiscoInfo+" query" {
						talk.Send(xmpp.NewIQ("result", iq.Id(), iq.From, xmpp.DiscInfoResult()))
					}
					if e.Name() == xmpp.NSByteStreams+" query" {
						query := e.(*xep.ByteStreamsQuery)
						addr := query.Hosts[0].Host + ":" + query.Hosts[0].Port

						c, err := net.Dial("tcp", addr)
						if err != nil {
							log.Println(err)
							continue
						}
						defer c.Close()
						if err := socks5(c, xep.Sha1Addr(query.Sid, iq.From, iq.To)); err != nil {
							log.Println(err)
							continue
						}
						talk.Send(xmpp.NewIQ("result", iq.Id(), iq.From,
							xep.NewByteStreamQuery(query.Sid, "",
								xep.NewStreamHostUsed(query.Hosts[0].Jid), nil)))

						file, err := os.Create(fileName)
						if err != nil {
							log.Println(err)
							continue
						}
						io.Copy(file, c)
					}
				}

			}

			if st.Name() == "presence" {
				presence := stanza
				if presence.Type() == "subscribe" {
					talk.Send(xmpp.NewPresence("subscribed", presence.Id(), presence.From))
					talk.Send(xmpp.NewPresence("subscribe", client.GenId(), presence.From))
				}
			}
		})
		exit <- 1
	}()

	talk.Send(xmpp.NewStanza("presence"))

	iq, err := talk.SendIQ(xmpp.NewIQ("get", client.GenId(), "", &core.RosterQuery{}))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(iq)

	iq, err = talk.SendIQ(xmpp.NewIQ("get", client.GenId(), "", &xep.DiscoItemsQuery{}))
	if err != nil {
		log.Println(err)
	}
	fmt.Println(iq)

	iq, err = talk.SendIQ(xmpp.NewIQ("get", client.GenId(), "", &xep.DiscoInfoQuery{}))
	if err != nil {
		log.Println(err)
	}
	fmt.Println(iq)

	iq, err = talk.SendIQ(xmpp.NewIQ("get", client.GenId(), "proxy.gerry-ubuntu-work", &xep.ByteStreamsQuery{}))
	if err != nil {
		log.Println(err)
	}
	fmt.Println(iq)

	/*
		form := xep.NewFormData("form", "", "",
			xep.NewFormField(xep.FieldSList, "", "stream-method", "", nil, false,
				xep.NewFormOption("", xmpp.NSByteStreams), xep.NewFormOption("", xmpp.NSIBB)))
		si := xep.NewSI(client.GenId(), "image/jpeg", xmpp.NSFileTransfer,
			xep.NewFileTransfer("001.jpg", "1024", "", "", ""),
			xep.NewFeatureNeg(form))
		iq, err = talk.SendIQ(xmpp.NewIQ("set", client.GenId(), "user002@gerry-ubuntu-work/Spark 2.6.3", si))
		fmt.Println(iq)

		iq, err = talk.SendIQ(xmpp.NewIQ("get", client.GenId(), "user002@gerry-ubuntu-work/Spark 2.6.3",
			&xep.ByteStreamsQuery{}))
		fmt.Println(iq)

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

func socks5(c net.Conn, addr string) error {
	b := make([]byte, 128)
	methods := gosocks5.Methods([]gosocks5.MethodType{gosocks5.MethodNoAuth})

	if _, err := c.Write(methods.Encode()); err != nil {
		return err
	}

	n, err := c.Read(b)
	if err != nil {
		return err
	}
	log.Println(b[:n])

	cmd := gosocks5.NewCMD(gosocks5.CmdConnect, gosocks5.AddrDomainName, addr, 0)
	if _, err := c.Write(cmd.Encode()); err != nil {
		return err
	}
	fmt.Printf("%d, %x\n", len(cmd.Encode()), cmd.Encode())
	n, err = c.Read(b)
	if err != nil {
		return err
	}
	log.Println(b[:n])

	return nil
}
