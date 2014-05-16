// XEP-0203: Delayed Delivery
// http://xmpp.org/extensions/xep-0203.html
package xep

import (
	"encoding/xml"
)

type Delay struct {
	XMLName xml.Name `xml:"urn:xmpp:delay delay"`
	From    string   `xml:"from,attr,omitemtpy"`
	Stamp   string   `xml:"stamp,attr"`
}

func (_ Delay) Name() string {
	return "delay"
}

func (_ Delay) FullName() string {
	return "urn:xmpp:delay delay"
}
