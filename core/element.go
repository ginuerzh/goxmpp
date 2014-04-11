// element
package core

import (
	"encoding/xml"
)

type Element interface {
	Name() string
}

type Stream struct {
	XMLName xml.Name `xml:"http://etherx.jabber.org/streams stream"`
	Id      string   `xml:"id,attr,omitempty"`
	From    string   `xml:"from,attr,omitempty"`
	To      string   `xml:"to,attr,omitempty"`
	Version string   `xml:"version,attr"`
	Lang    string   `xml:"lang,attr"`
}

func (_ Stream) Name() string {
	return "stream"
}

type StreamError struct {
	XMLName xml.Name `xml:"http://etherx.jabber.org/streams error"`
	Err     xml.Name `xml:",any"`
	Text    string   `xml:"text"`
}

func (_ StreamError) Name() string {
	return "error"
}

func (e StreamError) Error() string {
	return e.Err.Local + ": " + e.Text
}
