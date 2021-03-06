// stream feature
package core

import (
	"encoding/xml"
)

type StreamFeatures struct {
	XMLName xml.Name `xml:"http://etherx.jabber.org/streams features"`

	StartTLS   *TlsStartTLS
	Mechanisms *SaslMechanisms
	Compress   *FeatureCompress
	Bind       *FeatureBind
	Session    *FeatureSession
}

func (_ StreamFeatures) Name() string {
	return "features"
}

func (_ StreamFeatures) FullName() string {
	return "http://etherx.jabber.org/streams features"
}

type FeatureCompress struct {
	XMLName xml.Name `xml:"http://jabber.org/features/compress compression"`
	Method  []string `xml:"method"`
}

func (_ FeatureCompress) Name() string {
	return "compression"
}

func (_ FeatureCompress) FullName() string {
	return "http://jabber.org/features/compress compression"
}

type FeatureBind struct {
	XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
	Resource string   `xml:"resource,omitempty"`
	Jid      string   `xml:"jid,omitempty"`
}

func (_ FeatureBind) Name() string {
	return "bind"
}

func (_ FeatureBind) FullName() string {
	return "urn:ietf:params:xml:ns:xmpp-bind bind"
}

func (b FeatureBind) String() string {
	return "[bind] Jid:" + b.Jid + ", resouce:" + b.Resource
}

type FeatureSession struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-session session"`
}

func (_ FeatureSession) Name() string {
	return "session"
}

func (_ FeatureSession) FullName() string {
	return "urn:ietf:params:xml:ns:xmpp-session session"
}

func (_ FeatureSession) String() string {
	return "[session]"
}
