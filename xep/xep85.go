// XEP-0085: Chat State Notifications
// http://xmpp.org/extensions/xep-0085.html

package xep

import (
	"encoding/xml"
)

type ChatStateActive struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/chatstates active"`
}

func (_ ChatStateActive) String() string {
	return "active"
}

type ChatStateComposing struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/chatstates composing"`
}

func (_ ChatStateComposing) String() string {
	return "composing"
}

type ChatStatePaused struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/chatstates paused"`
}

func (_ ChatStatePaused) String() string {
	return "paused"
}

type ChatStateInactive struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/chatstates inactive"`
}

func (_ ChatStateInactive) String() string {
	return "inactive"
}

type ChatStateGone struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/chatstates gone"`
}

func (_ ChatStateGone) String() string {
	return "gone"
}
