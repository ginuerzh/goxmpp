// tls
package core

import (
	"encoding/xml"
)

// STARTTLS Negotiation
type TlsStartTLS struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls starttls"`
	Require *string  `xml:"required"`
}

func (_ TlsStartTLS) Name() string {
	return "starttls"
}

func (_ TlsStartTLS) FullName() string {
	return "urn:ietf:params:xml:ns:xmpp-tls starttls"
}

type TlsFailure struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls failure"`
}

func (_ TlsFailure) Name() string {
	return "failure"
}

func (_ TlsFailure) FullName() string {
	return "urn:ietf:params:xml:ns:xmpp-tls failure"
}

func (_ TlsFailure) Error() string {
	return "tls failure"
}

type TlsProceed struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls proceed"`
}

func (_ TlsProceed) Name() string {
	return "proceed"
}

func (_ TlsProceed) FullName() string {
	return "urn:ietf:params:xml:ns:xmpp-tls proceed"
}
