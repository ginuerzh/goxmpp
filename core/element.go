// element
package core

import (
	"encoding/xml"
)

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

func (_ Stream) FullName() string {
	return "http://etherx.jabber.org/streams stream"
}

type StreamError struct {
	XMLName xml.Name `xml:"http://etherx.jabber.org/streams error"`
	Err     xml.Name `xml:",any"`
	Text    string   `xml:"text"`
}

func (_ StreamError) Name() string {
	return "error"
}

func (_ StreamError) FullName() string {
	return "http://etherx.jabber.org/streams error"
}
func (e *StreamError) Error() string {
	return e.Err.Local + ": " + e.Text
}
