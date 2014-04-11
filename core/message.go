// message
package core

import (
	"encoding/xml"
)

type StanMsg struct {
	XMLName xml.Name `xml:"jabber:client message"`
	Stanza
	Body    string `xml:"body"`
	Subject string `xml:"subject"`
	Thread  *MsgThread
}

func (_ StanMsg) Name() string {
	return "message"
}

type MsgThread struct {
	XMLName xml.Name `xml:"thread"`
	Parent  string   `xml:"parent,attr,omitempty"`
	Value   string   `xml:",chardata"`
}

func (_ MsgThread) Name() string {
	return "thread"
}
