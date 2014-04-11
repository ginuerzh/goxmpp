// xmpp
package xmpp

import (
	"encoding/xml"
	"github.com/ginuerzh/goxmpp/core"
	"github.com/ginuerzh/goxmpp/xep"
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

type IQDefault struct {
	XMLName xml.Name `xml:"jabber:client iq"`
	core.Stanza

	// XEP-0199
	Ping *xep.Ping
}

func (_ IQDefault) Name() string {
	return "iq"
}

func (e IQDefault) Elem() core.Element {
	if e.Ping != nil {
		return e.Ping
	}
	return nil
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
	case *xep.DiscoInfoQuery: // XEP-0030
		iq = &xep.IQDiscoInfoQuery{Stanza: st, Query: v}
	case *xep.DiscoItemsQuery: // XEP-0030
		iq = &xep.IQDiscoItemsQuery{Stanza: st, Query: v}
	case *xep.VCard: // XEP-0054
		iq = &xep.IQVCard{Stanza: st, Card: v}
	case *xep.Ping: // XEP-0199
		iq = &xep.IQPing{}
	default:
		iq = &IQDefault{Stanza: st}
	}

	return iq
}

type StanMsg struct {
	XMLName xml.Name `xml:"jabber:client message"`
	core.StanMsg

	// XEP-0085
	Active    *xep.ChatStateActive
	Composing *xep.ChatStateComposing
	Paused    *xep.ChatStatePaused
	Inactive  *xep.ChatStateInactive
	Gone      *xep.ChatStateGone
}

func (e *StanMsg) ChatState() string {
	var state string

	if e.Active != nil {
		state = e.Active.String()
	} else if e.Composing != nil {
		state = e.Composing.String()
	} else if e.Paused != nil {
		state = e.Paused.String()
	} else if e.Inactive != nil {
		state = e.Inactive.String()
	} else if e.Gone != nil {
		state = e.Gone.String()
	} else {
		state = ""
	}

	return state
}

type StanPresence struct {
	XMLName xml.Name `xml:"jabber:client presence"`
	core.StanPresence

	// XEP-0115
	Caps *xep.EntityCaps
}
