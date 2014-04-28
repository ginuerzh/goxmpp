// XEP-0166: Jingle
// http://xmpp.org/extensions/xep-0166.html

package xep

import (
	"encoding/xml"
)

type Jingle struct {
	XMLName   xml.Name `xml:"urn:xmpp:jingle:1 jingle"`
	Action    string   `xml:"action,attr"`
	Initiator string   `xml:"initiator,attr"`
	Responder string   `xml:"responder,attr"`
	Sid       string   `xml:"sid,attr"`

	Contents []JingleContent `xml:"content"`
	Ringing  *JingleRinging
}

func (_ Jingle) Name() string {
	return "jingle"
}

func (_ Jingle) FullName() string {
	return "urn:xmpp:jingle:1 jingle"
}

func (j Jingle) String() string {
	return ""
}

type JingleContent struct {
	Description *JingleDesc      `xml:"description"`
	Transport   *JingleTransport `xml:"transport"`
}

type JingleDesc struct {
	Media    string        `xml:"media,attr"`
	Payloads []PayloadType `xml:"payload-type"`
}

type PayloadType struct {
	Id        int    `xml:"id,attr"`
	Name      string `xml:"name,attr"`
	ClockRate int    `xml:"clockrate"`
	Channels  int    `xml:"channels,attr"`
}

type JingleTransport struct {
	Pwd        string      `xml:"pwd,attr"`
	Ufrag      string      `xml:"ufrag"`
	Candidates []Candidate `xml:"candidate"`
}

type Candidate struct {
	Id         string `xml:"id,attr"`
	Component  int    `xml:"component,attr"`
	Foundation int    `xml:"foundation,attr"`
	Generation int    `xml:"generation,attr"`
	Ip         string `xml:"ip,attr"`
	Network    int    `xml:"ip,attr"`
	Port       uint16 `xml:"port,attr"`
	Priority   int    `xml:"priority,attr"`
	Protocol   string `xml:"protocol,attr"`
	RelAddr    string `xml:"rel-addr,attr"`
	RelPort    uint16 `xml:"rel-port,attr"`
	Type       string `xml:"type,attr"`
}

type JingleRinging struct {
	XMLName xml.Name `xml:"urn:xmpp:jingle:apps:rtp:info:1 ringing"`
}
