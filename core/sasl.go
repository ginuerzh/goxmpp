// sasl
package core

import (
	"encoding/xml"
)

//SASL Negotiation
type SaslAuth struct {
	XMLName   xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl auth"`
	Mechanism string   `xml:"mechanism,attr"`
	Value     string   `xml:",chardata"`
}

func (_ SaslAuth) Name() string {
	return "auth"
}

func (_ SaslAuth) FullName() string {
	return "urn:ietf:params:xml:ns:xmpp-sasl auth"
}

type SaslChallenge struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl challenge"`
	Value   string   `xml:",chardata"`
}

func (_ SaslChallenge) Name() string {
	return "challenge"
}

func (_ SaslChallenge) FullName() string {
	return "urn:ietf:params:xml:ns:xmpp-sasl challenge"
}

type SaslAbort struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl abort"`
}

func (_ SaslAbort) Name() string {
	return "abort"
}

func (_ SaslAbort) FullName() string {
	return "urn:ietf:params:xml:ns:xmpp-sasl abort"
}

func (_ SaslAbort) Error() string {
	return "sasl abort"
}

type SaslFailure struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl failure"`
	Reason  xml.Name `xml:",any"`
	Text    string   `xml:"text"`
}

func (_ SaslFailure) Name() string {
	return "failure"
}

func (_ SaslFailure) FullName() string {
	return "urn:ietf:params:xml:ns:xmpp-sasl failure"
}

func (e SaslFailure) Error() string {
	return e.Reason.Local + ": " + e.Text
}

type SaslResponse struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl response"`
	Value   string   `xml:",chardata"`
}

func (_ SaslResponse) Name() string {
	return "response"
}

func (_ SaslResponse) FullName() string {
	return "urn:ietf:params:xml:ns:xmpp-sasl response"
}

type SaslSuccess struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl success"`
	Value   string   `xml:",chardata"`
}

func (_ SaslSuccess) Name() string {
	return "success"
}

func (_ SaslSuccess) FullName() string {
	return "urn:ietf:params:xml:ns:xmpp-sasl success"
}

type SaslMechanisms struct {
	XMLName   xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl mechanisms"`
	Mechanism []string `xml:"mechanism"`
}

func (_ SaslMechanisms) Name() string {
	return "mechanisms"
}

func (_ SaslMechanisms) FullName() string {
	return "urn:ietf:params:xml:ns:xmpp-sasl mechanisms"
}
