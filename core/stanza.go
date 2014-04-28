// stanza
package core

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

type Stanza struct {
	Ids   string `xml:"id,attr,omitempty"`
	Types string `xml:"type,attr,omitempty"`
	From  string `xml:"from,attr,omitempty"`
	To    string `xml:"to,attr,omitempty"`
	Lang  string `xml:"lang,attr,omitempty"`
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

func (_ StanzaError) FullName() string {
	return "jabber:client error"
}

func (e *StanzaError) Error() string {
	return e.Code + ": " + e.Reason.Local
}

type RosterQuery struct {
	XMLName xml.Name      `xml:"jabber:iq:roster query"`
	Ver     string        `xml:"ver,attr,omitempty"`
	Items   []*RosterItem `xml:"item"`
}

func (_ *RosterQuery) Name() string {
	return "query"
}

func (_ *RosterQuery) FullName() string {
	return "jabber:iq:roster query"
}

func (e RosterQuery) String() string {
	b := &bytes.Buffer{}
	for _, item := range e.Items {
		fmt.Fprintf(b, "%s(%s) %s\n", item.Jid, item.Name, item.Group)
	}

	return b.String()
}

type RosterItem struct {
	//XMLName      xml.Name `xml:"item"`
	Jid          string   `xml:"jid,attr,omitempty"`
	Name         string   `xml:"name,attr,omitempty"`
	Subscription string   `xml:"subscription,attr,omitempty"`
	Approved     bool     `xml:"approved,attr,omitempty"`
	Ask          string   `xml:"ask,attr,omitempty"`
	Group        []string `xml:"group,omitempty"`
}

type MsgBody struct {
	XMLName xml.Name `xml:"jabber:client body"`
	Lang    string   `xml:"lang,attr,omitempty"`
	Body    string   `xml:",chardata"`
}

func (_ MsgBody) Name() string {
	return "body"
}

func (_ MsgBody) FullName() string {
	return "jabber:client body"
}

func (mb MsgBody) String() string {
	return "[body] " + mb.Body
}

type MsgSubject struct {
	XMLName xml.Name `xml:"jabber:client subject"`
	Lang    string   `xml:"lang,attr,omitempty"`
	Subject string   `xml:",chardata"`
}

func (_ MsgSubject) Name() string {
	return "subject"
}

func (_ MsgSubject) FullName() string {
	return "jabber:client subject"
}

func (ms MsgSubject) String() string {
	return "[subject] " + ms.Subject
}

type MsgThread struct {
	XMLName xml.Name `xml:"jabber:client thread"`
	Parent  string   `xml:"parent,attr,omitempty"`
	Value   string   `xml:",chardata"`
}

func (_ MsgThread) Name() string {
	return "thread"
}

func (_ MsgThread) FullName() string {
	return "jabber:client thread"
}

func (t MsgThread) String() string {
	return "[subject] " + t.Value
}

type MsgHtml struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/xhtml-im html"`
	Body    string   `xml:",chardata"`
}

func (_ MsgHtml) Name() string {
	return "html"
}

func (_ MsgHtml) FullName() string {
	return "http://jabber.org/protocol/xhtml-im html"
}

func (h MsgHtml) String() string {
	return "[html] " + h.Body
}

type PresenceShow struct {
	XMLName xml.Name `xml:"jabber:client show"`
	Show    string   `xml:",chardata"`
}

func (_ PresenceShow) Name() string {
	return "show"
}

func (_ PresenceShow) FullName() string {
	return "jabber:client show"
}

func (p PresenceShow) String() string {
	return "[show] " + p.Show
}

type PresenceStatus struct {
	XMLName xml.Name `xml:"jabber:client status"`
	Status  string   `xml:",chardata"`
}

func (_ PresenceStatus) Name() string {
	return "status"
}

func (_ PresenceStatus) FullName() string {
	return "jabber:client status"
}

func (p PresenceStatus) String() string {
	return "[status] " + p.Status
}

type PresencePriority struct {
	XMLName  xml.Name `xml:"jabber:client priority"`
	Priority string   `xml:",chardata"`
}

func (_ PresencePriority) Name() string {
	return "priority"
}

func (_ PresencePriority) FullName() string {
	return "jabber:client priority"
}

func (p PresencePriority) String() string {
	return "[priority] " + p.Priority
}
