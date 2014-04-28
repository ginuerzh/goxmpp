// XEP-0085: Chat State Notifications
// http://xmpp.org/extensions/xep-0085.html

package xep

import (
	"encoding/xml"
)

type ChatStateActive struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/chatstates active"`
}

func (_ ChatStateActive) Name() string {
	return "active"
}

func (_ ChatStateActive) FullName() string {
	return "http://jabber.org/protocol/chatstates active"
}

func (_ ChatStateActive) String() string {
	return "[chatstate] active"
}

type ChatStateComposing struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/chatstates composing"`
}

func (_ ChatStateComposing) Name() string {
	return "composing"
}

func (_ ChatStateComposing) FullName() string {
	return "http://jabber.org/protocol/chatstates composing"
}

func (_ ChatStateComposing) String() string {
	return "[chatstate] composing"
}

type ChatStatePaused struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/chatstates paused"`
}

func (_ ChatStatePaused) Name() string {
	return "paused"
}

func (_ ChatStatePaused) FullName() string {
	return "http://jabber.org/protocol/chatstates paused"
}

func (_ ChatStatePaused) String() string {
	return "[chatstate] paused"
}

type ChatStateInactive struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/chatstates inactive"`
}

func (_ ChatStateInactive) Name() string {
	return "inactive"
}

func (_ ChatStateInactive) FullName() string {
	return "http://jabber.org/protocol/chatstates inactive"
}

func (_ ChatStateInactive) String() string {
	return "[chatstate] inactive"
}

type ChatStateGone struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/chatstates gone"`
}

func (_ ChatStateGone) Name() string {
	return "gone"
}

func (_ ChatStateGone) FullName() string {
	return "http://jabber.org/protocol/chatstates gone"
}

func (_ ChatStateGone) String() string {
	return "[chatstate] gone"
}
