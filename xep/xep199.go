// XEP-0199: XMPP Ping
// http://xmpp.org/extensions/xep-0199.html

package xep

import (
	"encoding/xml"
)

type Ping struct {
	XMLName xml.Name `xml:"urn:xmpp:ping ping"`
}

func (_ Ping) Name() string {
	return "ping"
}

func (_ Ping) FullName() string {
	return "urn:xmpp:ping ping"
}

func (_ Ping) String() string {
	return "[ping]"
}
