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

	talk.HandleFunc(xmpp.NSClient+" message", func(header *core.StanzaHeader, e xmpp.Element) {
		msg := e.(*xmpp.Stanza)
		body := ""
		subject := ""
		for _, e := range msg.E() {
			if e.Name() == "body" {
				body = e.(*core.MsgBody).Body
				break
			}
			if e.Name() == "subject" {
				subject = e.(*core.MsgSubject).Subject
			}
		}
		if len(body) > 0 {
			talk.Send(xmpp.NewMessage(header.Types, header.From, body, subject))
		}
	})

	talk.HandleFunc(xmpp.NSClient+" presence", func(header *core.StanzaHeader, e xmpp.Element) {
		if header.Types == "subscribe" {
			talk.Send(xmpp.NewPresence("subscribed", header.Ids, header.From))
			talk.Send(xmpp.NewPresence("subscribe", client.GenId(), header.From))
		}

	})

	talk.HandleFunc(xmpp.NSPing+" ping", func(header *core.StanzaHeader, e xmpp.Element) {
		talk.Send(xmpp.NewIQ("result", header.Ids, header.From, nil))
	})

	filename := ""
	talk.HandleFunc(xmpp.NSSI+" si", func(header *core.StanzaHeader, e xmpp.Element) {
		si := e.(*xep.SI)
		filename = si.File.Name
		submit := xep.NewFormData("submit", "", "",
			xep.NewFormField("", "", "stream-method", "",
				[]string{xmpp.NSByteStreams}, false, nil))

		talk.Send(xmpp.NewIQ("result", header.Ids, header.From,
			xep.NewSI("", "", "", nil, xep.NewFeature(submit))))
	})

	talk.HandleFunc(xmpp.NSRoster+" query", func(header *core.StanzaHeader, e xmpp.Element) {
		fmt.Println(e)
	})
	talk.HandleFunc(xmpp.NSDiscoItems+" query", func(header *core.StanzaHeader, e xmpp.Element) {
		fmt.Println(e)
	})
	talk.HandleFunc(xmpp.NSDiscoInfo+" query", func(header *core.StanzaHeader, e xmpp.Element) {
		if header.Types == "result" {
			fmt.Println(e)
			return
		}
		talk.Send(xmpp.NewIQ("result", header.Ids, header.From, xmpp.DiscInfoResult()))
	})

	talk.HandleFunc(xmpp.NSByteStreams+" query", func(header *core.StanzaHeader, e xmpp.Element) {
		query := e.(*xep.ByteStreamsQuery)
		addr := query.Hosts[1].Host + ":" + query.Hosts[1].Port

		c, err := net.Dial("tcp", addr)
		if err != nil {
			log.Println(err)
			return
		}
		defer c.Close()
		if err := socks5(c, xep.Sha1Addr(query.Sid, header.From, header.To)); err != nil {
			log.Println(err)
			return
		}
		talk.Send(xmpp.NewIQ("result", header.Ids, header.From,
			xep.NewByteStreamQuery(query.Sid, "",
				xep.NewStreamHostUsed(query.Hosts[1].Jid), nil)))

		file, err := os.Create(filename)
		if err != nil {
			log.Println(err)
			return
		}
		io.Copy(file, c)
	})

	talk.HandleFunc(xmpp.NSRoster+" query", func(header *core.StanzaHeader, e xmpp.Element) {
		fmt.Println(e)
	})

	talk.OnLogined(func(err error) {
		if err != nil {
			return
		}

		run(talk)
	})

	log.Fatal(talk.Run())
}

func run(talk *client.Client) {
	talk.Send(xmpp.NewStanza("presence"))
	talk.Send(xmpp.NewIQ("get", client.GenId(), "", &core.RosterQuery{}))
	talk.Send(xmpp.NewIQ("get", client.GenId(), "", &xep.DiscoInfoQuery{}))

	iq, err := talk.SendIQ(xmpp.NewIQ("get", client.GenId(), "", &xep.DiscoItemsQuery{}))
	s5b := ""
	for _, item := range iq.E()[0].(*xep.DiscoItemsQuery).Items {
		iq, err = talk.SendIQ(xmpp.NewIQ("get", client.GenId(), item.Jid, &xep.DiscoInfoQuery{}))
		log.Println(iq)
		result := iq.E()[0].(*xep.DiscoInfoQuery)
		if result.Identities[0].Category == "proxy" && result.Identities[0].Type == "bytestreams" {
			s5b = item.Jid
		}
	}

	iq, err = talk.SendIQ(xmpp.NewIQ("get", client.GenId(), s5b, &xep.ByteStreamsQuery{}))
	if err != nil {
		log.Println(err)
	}
	fmt.Println(iq)
	streamHost := iq.E()[0].(*xep.ByteStreamsQuery).Hosts[0]

	form := xep.NewFormData("form", "", "",
		xep.NewFormField(xep.FieldSList, "", "stream-method", "", nil, false,
			xep.NewFormOption("", xmpp.NSByteStreams)))
	si := xep.NewSI(client.GenId(), "image/jpeg", xmpp.NSFileTransfer,
		xep.NewFileTransfer("001.jpg", "1024", "", "", ""),
		xep.NewFeature(form))
	iq, err = talk.SendIQ(xmpp.NewIQ("set", client.GenId(), "user002@gerry-ubuntu-work/Spark 2.6.3", si))
	fmt.Println(iq)

	iq, err = talk.SendIQ(xmpp.NewIQ("get", client.GenId(),
		"user002@gerry-ubuntu-work/Spark 2.6.3", &xep.DiscoInfoQuery{}))
	fmt.Println(iq)

	//iq, err = talk.SendIQ(xmpp.NewIQ("get", client.GenId(), "user002@gerry-ubuntu-work/Spark 2.6.3",
	//	&xep.ByteStreamsQuery{}))
	//fmt.Println(iq)

	iq, err = talk.SendIQ(xmpp.NewIQ("set", client.GenId(),
		"user002@gerry-ubuntu-work/Spark 2.6.3", xep.NewByteStreamQuery(client.GenId(), "tcp", nil, streamHost)))
	fmt.Println(iq)

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
