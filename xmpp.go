// xmpp
package xmpp

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/ginuerzh/goxmpp/core"
	"github.com/ginuerzh/goxmpp/xep"
	"strings"
)

const (
	NSClient       = "jabber:client"
	NSStream       = "http://etherx.jabber.org/streams"
	NSTLS          = "urn:ietf:params:xml:ns:xmpp-tls"
	NSSASL         = "urn:ietf:params:xml:ns:xmpp-sasl"
	NSBind         = "urn:ietf:params:xml:ns:xmpp-bind"
	NSSession      = "urn:ietf:params:xml:ns:xmpp-session"
	NSStanza       = "urn:ietf:params:xml:ns:xmpp-stanzas"
	NSRoster       = "jabber:iq:roster"
	NSDiscoInfo    = "http://jabber.org/protocol/disco#info"
	NSDiscoItems   = "http://jabber.org/protocol/disco#items"
	NSVcardTemp    = "vcard-temp"
	NSVcardUpdate  = "vcard-temp:x:update"
	NSPing         = "urn:xmpp:ping"
	NSSI           = "http://jabber.org/protocol/si"
	NSHtml         = "http://jabber.org/protocol/xhtml-im"
	NSChatState    = "http://jabber.org/protocol/chatstates"
	NSCaps         = "http://jabber.org/protocol/caps"
	NSFileTransfer = "http://jabber.org/protocol/si/profile/file-transfer"
	NSByteStreams  = "http://jabber.org/protocol/bytestreams"
	NSIBB          = "http://jabber.org/protocol/ibb"
	NSMUC          = "http://jabber.org/protocol/muc"
)

var clientFeatures = []string{
	"http://jabber.org/protocol/bytestreams",
}

func DiscInfoResult() *xep.DiscoInfoQuery {
	query := &xep.DiscoInfoQuery{
		Identities: []*xep.InfoIdentity{&xep.InfoIdentity{Category: "client", Type: "pc", Name: "goxmpp"}},
	}
	for _, f := range clientFeatures {
		query.Features = append(query.Features, &xep.InfoFeature{Var: f})
	}

	return query
}

type NewFunc func() Element

var elements = make(map[string]NewFunc)

func Register(elemName string, newFunc NewFunc) {
	elements[elemName] = newFunc
}

func E(elemName string) Element {
	newFunc, ok := elements[elemName]
	if !ok {
		return nil
	}

	return newFunc()
}

type Element interface {
	Name() string
	FullName() string
}

type Stan interface {
	Element
	Id() string
	Type() string
	Error() error
	E() []Element
}

type Stanza struct {
	XMLName xml.Name
	core.StanzaHeader
	Err      error
	Elements []Element
}

func (st Stanza) Name() string {
	return st.XMLName.Local
}

func (st Stanza) FullName() string {
	return st.XMLName.Space + " " + st.XMLName.Local
}

func (st Stanza) Id() string {
	return st.Ids
}

func (st Stanza) Type() string {
	return st.Types
}

func (st Stanza) E() []Element {
	return st.Elements
}

func (st *Stanza) AddE(elements ...Element) {
	st.Elements = append(st.Elements, elements...)
}

func (st *Stanza) Error() error {
	return st.Err
}

func (st *Stanza) String() string {
	b := &bytes.Buffer{}
	b.WriteString("[" + st.Name() + "] " + st.Type() + " " + st.Id() + "\n")
	for _, e := range st.Elements {
		fmt.Fprintln(b, e)
	}
	return b.String()
}

type JID string

func ToJID(s string) JID {
	return JID(s)
}

func (jid JID) Bare() string {
	a := strings.SplitN(string(jid), "/", 2)
	return a[0]
}

func (jid JID) Domain() string {
	a := strings.SplitN(jid.Bare(), "@", 2)
	return a[1]
}

func (jid JID) Split() (local string, domain string, resource string) {
	a := strings.SplitN(string(jid), "@", 2)
	if len(a) != 2 {
		return
	}
	local = a[0]

	a = strings.SplitN(a[1], "/", 2)
	domain = a[0]
	if len(a) != 2 {
		return
	}

	resource = a[1]
	return
}

func (jid JID) String() string {
	return string(jid)
}

type NullElement struct {
	XMLName xml.Name
}

func (e NullElement) Name() string {
	return e.XMLName.Local
}

func (e NullElement) FullName() string {
	return e.XMLName.Space + " " + e.XMLName.Local
}

func (e NullElement) String() string {
	return "[unknown] " + e.XMLName.Space + " " + e.XMLName.Local
}

func NewStanza(name string, elems ...Element) *Stanza {
	st := new(Stanza)
	st.XMLName.Space = "jabber:client"
	st.XMLName.Local = name
	st.Elements = elems

	return st
}

func NewIQ(typ string, id string, to string, payload Element) *Stanza {
	iq := NewStanza("iq", payload)
	iq.Types = typ
	iq.Ids = id
	iq.To = to

	return iq
}

func NewMessage(typ string, to string, body string, subject string) *Stanza {
	msg := NewStanza("message")
	msg.Types = typ
	msg.To = to

	if len(body) > 0 {
		msg.AddE(&core.MsgBody{Body: body})
	}
	if len(subject) > 0 {
		msg.AddE(&core.MsgSubject{Subject: subject})
	}

	return msg
}

func NewPresence(typ string, id string, to string) *Stanza {
	presence := NewStanza("presence")
	presence.Types = typ
	presence.Ids = id
	presence.To = to

	return presence
}
