// init
package xmpp

import (
	"github.com/ginuerzh/goxmpp/core"
	"github.com/ginuerzh/goxmpp/xep"
)

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

	// XEP04
	Register("jabber:x:data x",
		func() Element { return new(xep.XFormData) })
	// XEP20
	Register("http://jabber.org/protocol/feature-neg feature",
		func() Element { return new(xep.Feature) })
	// XEP30
	Register("http://jabber.org/protocol/disco#info query",
		func() Element { return new(xep.DiscoInfoQuery) })
	Register("http://jabber.org/protocol/disco#items query",
		func() Element { return new(xep.DiscoItemsQuery) })
	//XEP54
	Register("vcard-temp vCard",
		func() Element { return new(xep.VCard) })
	//XEP65
	Register("http://jabber.org/protocol/bytestreams query",
		func() Element { return new(xep.ByteStreamsQuery) })
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
	// XEP96
	Register("http://jabber.org/protocol/si si",
		func() Element { return new(xep.SI) })
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
