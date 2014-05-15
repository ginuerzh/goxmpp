// XEP-0163: Personal Eventing Protocol
// http://xmpp.org/extensions/xep-0163.html
package xep

import (
	"encoding/xml"
)

type Event struct {
	XMLName xml.Name    `xml:"http://jabber.org/protocol/pubsub#event event"`
	Items   *EventItems `xml:"items"`
}

type EventItems struct {
	Node  string
	Items []*EventItem
}

type EventItem struct {
	Id string
}
