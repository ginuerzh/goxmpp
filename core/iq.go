// iq
package core

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

type IQ interface {
	Stan
	Elem() Element
}

type IQEmpty struct {
	XMLName xml.Name `xml:"jabber:client iq"`
	Stanza
}

func (_ IQEmpty) Name() string {
	return "iq"
}

func (e IQEmpty) Elem() Element {
	return nil
}

type IQBind struct {
	XMLName xml.Name `xml:"jabber:client iq"`
	Stanza
	Bind *FeatureBind
}

func (_ IQBind) Name() string {
	return "iq"
}

func (e IQBind) Elem() Element {
	return e.Bind
}

type IQSession struct {
	XMLName xml.Name `xml:"jabber:client iq"`
	Stanza
	Session *FeatureSession
}

func (_ IQSession) Name() string {
	return "iq"
}

func (e IQSession) Elem() Element {
	return e.Session
}

type IQRosterQuery struct {
	XMLName xml.Name `xml:"jabber:client iq"`
	Stanza
	Query *RosterQuery
}

func (_ IQRosterQuery) Name() string {
	return "iq"
}

func (e IQRosterQuery) Elem() Element {
	return e.Query
}

type RosterQuery struct {
	XMLName xml.Name      `xml:"jabber:iq:roster query"`
	Ver     string        `xml:"ver,attr,omitempty"`
	Items   []*RosterItem `xml:"item"`
}

func (_ *RosterQuery) Name() string {
	return "query"
}

func (e RosterQuery) String() string {
	b := &bytes.Buffer{}
	for _, item := range e.Items {
		fmt.Fprintf(b, "%s(%s) %s\n", item.Jid, item.Name, item.Group)
	}

	return b.String()
}

type RosterItem struct {
	XMLName      xml.Name `xml:"item"`
	Jid          string   `xml:"jid,attr,omitempty"`
	Name         string   `xml:"name,attr,omitempty"`
	Subscription string   `xml:"subscription,attr,omitempty"`
	Approved     bool     `xml:"approved,attr,omitempty"`
	Ask          string   `xml:"ask,attr,omitempty"`
	Group        []string `xml:"group,omitempty"`
}
