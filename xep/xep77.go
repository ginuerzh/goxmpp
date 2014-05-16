// XEP-0077: In-Band Registration
// http://xmpp.org/extensions/xep-0077.html
package xep

import (
	"encoding/xml"
)

type RegisterQuery struct {
	XMLName xml.Name `xml:"jabber:iq:register query"`
}

func (_ RegisterQuery) Name() string {
	return "query"
}

func (_ RegisterQuery) FullName() string {
	return "jabber:iq:register query"
}
