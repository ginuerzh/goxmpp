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
	NSClient     = "jabber:client"
	NSStream     = "http://etherx.jabber.org/streams"
	NSTLS        = "urn:ietf:params:xml:ns:xmpp-tls"
	NSSASL       = "urn:ietf:params:xml:ns:xmpp-sasl"
	NSBind       = "urn:ietf:params:xml:ns:xmpp-bind"
	NSSession    = "urn:ietf:params:xml:ns:xmpp-session"
	NSStanza     = "urn:ietf:params:xml:ns:xmpp-stanzas"
	NSRoster     = "jabber:iq:roster"
	NSDiscoInfo  = "http://jabber.org/protocol/disco#info"
	NSDiscoItems = "http://jabber.org/protocol/disco#items"
	NSVcardTemp  = "vcard-temp"
	NSPing       = "urn:xmpp:ping"
	NSHtml       = "http://jabber.org/protocol/xhtml-im"
	NSChatState  = "http://jabber.org/protocol/chatstates"
	NSCaps       = "http://jabber.org/protocol/caps"
)

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

func init() {
	/* core elements */
	Register("http://etherx.jabber.org/streams stream",
		func() Element { return new(core.Stream) })
	Register("http://etherx.jabber.org/streams error",
		func() Element { return new(core.StreamError) })
	Register("http://etherx.jabber.org/streams features",
		func() Element { return new(core.StreamFeatures) })
	Register("http://jabber.org/features/compress compression",
		func() Element { return new(core.FeatureCompress) })
	Register("urn:ietf:params:xml:ns:xmpp-bind bind",
		func() Element { return new(core.FeatureBind) })
	Register("urn:ietf:params:xml:ns:xmpp-session session",
		func() Element { return new(core.FeatureSession) })
	Register("urn:ietf:params:xml:ns:xmpp-sasl auth",
		func() Element { return new(core.SaslAuth) })
	Register("urn:ietf:params:xml:ns:xmpp-sasl challenge",
		func() Element { return new(core.SaslChallenge) })
	Register("urn:ietf:params:xml:ns:xmpp-sasl abort",
		func() Element { return new(core.SaslAbort) })
	Register("urn:ietf:params:xml:ns:xmpp-sasl failure",
		func() Element { return new(core.SaslFailure) })
	Register("urn:ietf:params:xml:ns:xmpp-sasl response",
		func() Element { return new(core.SaslResponse) })
	Register("urn:ietf:params:xml:ns:xmpp-sasl success",
		func() Element { return new(core.SaslSuccess) })
	Register("urn:ietf:params:xml:ns:xmpp-sasl mechanisms",
		func() Element { return new(core.SaslMechanisms) })
	Register("urn:ietf:params:xml:ns:xmpp-tls starttls",
		func() Element { return new(core.TlsStartTLS) })
	Register("urn:ietf:params:xml:ns:xmpp-tls failure",
		func() Element { return new(core.TlsFailure) })
	Register("urn:ietf:params:xml:ns:xmpp-tls proceed",
		func() Element { return new(core.TlsProceed) })
	Register("jabber:client iq",
		func() Element { return NewStanza("iq") })
	Register("jabber:client message",
		func() Element { return NewStanza("message") })
	Register("jabber:client presence",
		func() Element { return NewStanza("presence") })
	Register("jabber:client error",
		func() Element { return new(core.StanzaError) })
	Register("jabber:iq:roster query",
		func() Element { return new(core.RosterQuery) })
	Register("jabber:client show",
		func() Element { return new(core.PresenceShow) })
	Register("jabber:client status",
		func() Element { return new(core.PresenceStatus) })
	Register("jabber:client priority",
		func() Element { return new(core.PresencePriority) })
	Register("jabber:client body",
		func() Element { return new(core.MsgBody) })
	Register("jabber:client subject",
		func() Element { return new(core.MsgSubject) })
	Register("jabber:client thread",
		func() Element { return new(core.MsgThread) })
	Register("http://jabber.org/protocol/xhtml-im html",
		func() Element { return new(core.MsgHtml) })

	/* XEP elements */

	// XEP30
	Register("http://jabber.org/protocol/disco#info query",
		func() Element { return new(xep.DiscoInfoQuery) })
	Register("http://jabber.org/protocol/disco#items query",
		func() Element { return new(xep.DiscoItemsQuery) })
	//XEP54
	Register("vcard-temp vCard",
		func() Element { return new(xep.VCard) })
	// XEP85
	Register("http://jabber.org/protocol/chatstates active",
		func() Element { return new(xep.ChatStateActive) })
	Register("http://jabber.org/protocol/chatstates composing",
		func() Element { return new(xep.ChatStateComposing) })
	Register("http://jabber.org/protocol/chatstates paused",
		func() Element { return new(xep.ChatStatePaused) })
	Register("http://jabber.org/protocol/chatstates inactive",
		func() Element { return new(xep.ChatStateInactive) })
	Register("http://jabber.org/protocol/chatstates gone",
		func() Element { return new(xep.ChatStateGone) })
	// XEP115
	Register("http://jabber.org/protocol/caps c",
		func() Element { return new(xep.EntityCaps) })
	// XEP166
	Register("urn:xmpp:jingle:1 jingle",
		func() Element { return new(xep.Jingle) })
	// XEP199
	Register("urn:xmpp:ping ping",
		func() Element { return new(xep.Ping) })
}

type Element interface {
	Name() string
}

type Stan interface {
	Element
	Id() string
	Type() string
	Error() error
}

type Stanza struct {
	XMLName xml.Name
	core.Stanza
	Elements []Element
}

func (st Stanza) Name() string {
	return st.XMLName.Local
}

func (st Stanza) Id() string {
	return st.Ids
}

func (st Stanza) Type() string {
	return st.Types
}

func (st *Stanza) Error() (err error) {
	if st.Types != "error" {
		return nil
	}
	for _, e := range st.Elements {
		if e.Name() == "error" {
			err = e.(*core.StanzaError)
			break
		}
	}
	return
}

func (st *Stanza) String() string {
	b := &bytes.Buffer{}
	b.WriteString("[" + st.Name() + "] " + st.Type() + " " + st.Id() + "\n")
	for _, e := range st.Elements {
		b.WriteString("\t")
		fmt.Fprintln(b, e)
	}
	return b.String()
}

func (st *Stanza) AddElement(elements ...Element) {
	st.Elements = append(st.Elements, elements...)
}

type JID string

func NewJID(s string) JID {
	return JID(s)
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
	return e.XMLName.Space + " " + e.XMLName.Local
}

func (e NullElement) String() string {
	return "[unknown] " + e.XMLName.Space + " " + e.XMLName.Local
}

func NewStanza(name string, elems ...Element) *Stanza {
	st := new(Stanza)
	st.XMLName.Space = "jabber:client"
	st.XMLName.Local = name
	st.Elements = append(st.Elements, elems...)

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
		msg.AddElement(&core.MsgBody{Body: body})
	}
	if len(subject) > 0 {
		msg.AddElement(&core.MsgSubject{Subject: subject})
	}

	return msg
}
