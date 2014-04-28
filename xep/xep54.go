// XEP-0054: vcard-temp
// http://xmpp.org/extensions/xep-0054.html

package xep

import (
	"encoding/xml"
)

type VCard struct {
	XMLName    xml.Name     `xml:"vcard-temp vCard"`
	FName      string       `xml:"FN,omitempty"`
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
	Photo      *VCardPhoto  `xml:"PHOTO"`
}

func (_ VCard) Name() string {
	return "vCard"
}

func (_ VCard) FullName() string {
	return "vcard-temp vCard"
}

func (vc VCard) String() string {
	return "vcard"
}

type VCardPhoto struct {
	Type   string `xml:"TYPE"`
	BinVal string `xml:"BINVAL"`
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
