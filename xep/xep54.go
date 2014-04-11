// XEP-0054: vcard-temp
// http://xmpp.org/extensions/xep-0054.html

package xep

import (
	"encoding/xml"
	"github.com/ginuerzh/goxmpp/core"
)

type IQVCard struct {
	XMLName xml.Name `xml:"jabber:client iq"`
	core.Stanza
	Card *VCard
}

func (_ IQVCard) Name() string {
	return "iq"
}

func (e IQVCard) Elem() core.Element {
	return e.Card
}

type VCard struct {
	XMLName    xml.Name     `xml:"vcard-temp vCard"`
	FullName   string       `xml:"FN,omitempty"`
	FamilyName string       `xml:"N>FAMILY,omitempty"`
	GivenName  string       `xml:"N>GIVEN,omitempty"`
	MidName    string       `xml:"N>MIDDLE,omitempty"`
	NickName   string       `xml:"NICKNAME,omitempty"`
	Url        string       `xml:"URL,omitempty"`
	Birthday   string       `xml:"BDAY,omitempty"`
	OrgName    string       `xml:"ORG>ORGNAME,omitempty"`
	OrgUnit    string       `xml:"ORG>ORGUNIT,omitempty"`
	Title      string       `xml:"TITLE,omitempty"`
	Role       string       `xml:"ROLE,omitempty"`
	Addr       []*VCardAddr `xml:"ADR,omitempty"`
	JabberId   string       `xml:"JABBERID,omitempty"`
	Desc       string       `xml:"DESC,omitempty"`
}

func (_ VCard) Name() string {
	return "vCard"
}

type VCardAddr struct {
	Work     string `xml:"WORK"`
	TexAdd   string `xml:"EXTADD"`
	Street   string `xml:"STREET"`
	Locality string `xml:"LOCALITY"`
	Region   string `xml:"REGION"`
	Pcode    string `xml:"PCODE"`
	Country  string `xml:"CTRY"`
}
