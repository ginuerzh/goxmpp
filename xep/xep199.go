// XEP-0199: XMPP Ping
// http://xmpp.org/extensions/xep-0199.html

package xep

import (
	"encoding/xml"
	"github.com/ginuerzh/goxmpp/core"
)

type IQPing struct {
	XMLName xml.Name `xml:"jabber:client iq"`
	core.Stanza
	P *Ping
}

func (_ IQPing) Name() string {
	return "iq"
}

func (e IQPing) Elem() core.Element {
	return e.P
}

type Ping struct {
	XMLName xml.Name `xml:"urn:xmpp:ping ping"`
}

func (_ Ping) Name() string {
	return "ping"
}
