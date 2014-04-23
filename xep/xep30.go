// XEP-0030: Service Discovery
// http://xmpp.org/extensions/xep-0030.html

package xep

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

type DiscoInfoQuery struct {
	XMLName    xml.Name        `xml:"http://jabber.org/protocol/disco#info query"`
	Ver        string          `xml:"ver,attr,omitempty"`
	Node       string          `xml:"node,attr,omitempty"`
	Identities []*InfoIdentity `xml:"identity"`
	Features   []*InfoFeature  `xml:"feature"`
}

func (_ DiscoInfoQuery) Name() string {
	return "query"
}

func (e DiscoInfoQuery) String() string {
	b := &bytes.Buffer{}
	for _, id := range e.Identities {
		fmt.Fprintf(b, "%s/%s/%s\n", id.Name, id.Type, id.Category)
	}
	for _, f := range e.Features {
		fmt.Fprintln(b, f.Var)
	}

	return b.String()
}

type InfoIdentity struct {
	XMLName  xml.Name `xml:"identity"`
	Category string   `xml:"category,attr"`
	Type     string   `xml:"type,attr"`
	Name     string   `xml:"name,attr"`
}

type InfoFeature struct {
	XMLName xml.Name `xml:"feature"`
	Var     string   `xml:"var,attr"`
}

type DiscoItemsQuery struct {
	XMLName xml.Name     `xml:"http://jabber.org/protocol/disco#items query"`
	Ver     string       `xml:"ver,attr,omitempty"`
	Node    string       `xml:"node,attr,omitempty"`
	Items   []*DiscoItem `xml:"item"`
}

func (_ DiscoItemsQuery) Name() string {
	return "query"
}

func (e DiscoItemsQuery) String() string {
	b := &bytes.Buffer{}
	for _, item := range e.Items {
		fmt.Fprintf(b, "%s(%s) %s\n", item.Jid, item.Name, item.Node)
	}

	return b.String()
}

type DiscoItem struct {
	XMLName xml.Name `xml:"item"`
	Jid     string   `xml:"jid,attr,omitempty"`
	Node    string   `xml:"node,attr,omitempty"`
	Name    string   `xml:"name,attr,omitempty"`
}
