// xmpp
package xmpp

import (
	"github.com/ginuerzh/goxmpp/core"
	"strings"
)

const (
	NSStream = "http://etherx.jabber.org/streams"
	NSTLS    = "urn:ietf:params:xml:ns:xmpp-tls"
	NSSASL   = "urn:ietf:params:xml:ns:xmpp-sasl"
	NSBind   = "urn:ietf:params:xml:ns:xmpp-bind"
	NSStanza = "urn:ietf:params:xml:ns:xmpp-stanzas"
	NSClient = "jabber:client"
	NSRoster = "jabber:iq:roster"
)

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

func NewIQ(typ, id, to string, e core.Element) core.IQ {
	var iq core.IQ

	st := core.Stanza{Types: typ, To: to, Ids: id}

	switch v := e.(type) {
	case *core.FeatureBind:
		iq = &core.IQBind{Stanza: st, Bind: v}

	case *core.FeatureSession:
		iq = &core.IQSession{Stanza: st, Session: v}
	case *core.RosterQuery:
		iq = &core.IQRosterQuery{Stanza: st, Query: v}
	default:
		iq = &core.IQEmpty{Stanza: st}
	}

	return iq
}
