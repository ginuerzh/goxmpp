// stanza
package core

import (
	"encoding/xml"
)

type Stan interface {
	Element
	Id() string
	Type() string
	Error() *StanzaError
}

type Stanza struct {
	Ids   string `xml:"id,attr,omitempty"`
	Types string `xml:"type,attr,omitempty"`
	From  string `xml:"from,attr,omitempty"`
	To    string `xml:"to,attr,omitempty"`
	Lang  string `xml:"lang,attr,omitempty"`
	Err   *StanzaError
}

func (_ Stanza) Name() string {
	return "stanza"
}

func (st Stanza) Id() string {
	return st.Ids
}

func (st Stanza) Type() string {
	return st.Types
}

func (st Stanza) Error() *StanzaError {
	return st.Err
}

type StanzaError struct {
	XMLName xml.Name `xml:"jabber:client error"`
	Code    string   `xml:"code,attr"`
	Type    string   `xml:"type,attr"`
	Reason  xml.Name `xml:",any"`
	Text    string   `xml:"text,omitempty"`
}

func (_ StanzaError) Name() string {
	return "error"
}

func (e StanzaError) Error() string {
	return e.Code + ": " + e.Reason.Local
}
