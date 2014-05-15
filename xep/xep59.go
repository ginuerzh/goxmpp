// XEP-0059: Result Set Management
// http://xmpp.org/extensions/xep-0059.html
package xep

import (
	"encoding/xml"
)

type Rsm struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/rsm set"`
	Max     int      `xml:"max,omitempty"`
	Before  *string  `xml:"before"`
	After   string   `xml:"after,omitempty"`
	Fist    *RSFirst `xml:"first"`
	Last    string   `xml:"last,omitempty"`
	Count   int      `xml:"count,omitempty"`
}

type RSFirst struct {
	Index int    `xml:"index,attr"`
	Value string `xml:",chardata"`
}
