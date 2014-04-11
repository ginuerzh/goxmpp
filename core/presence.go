// presence
package core

import (
	"encoding/xml"
)

type StanPresence struct {
	XMLName xml.Name `xml:"jabber:client presence"`
	Stanza
	Show     string `xml:"show,omitempty"`
	Status   string `xml:"status,omitempty"`
	Priority int    `xml:"priority,omitempty"`
}

func (_ StanPresence) Name() string {
	return "presence"
}
