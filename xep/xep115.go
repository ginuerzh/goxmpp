// XEP-0115: Entity Capabilities
// http://xmpp.org/extensions/xep-0115.html

package xep

import (
	"encoding/xml"
)

type EntityCaps struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/caps c"`
	Hash    string   `xml:"hash,attr"`
	Node    string   `xml:"node,attr"`
	Ver     string   `xml:"ver,attr"`
	Ext     string   `xml:"ext,attr"`
}

func (_ EntityCaps) Name() string {
	return "caps"
}

func (c EntityCaps) String() string {
	return "[caps] " + c.Node + "," + c.Ver + "," + c.Ext
}
