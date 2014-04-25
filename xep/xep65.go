// XEP-0065: SOCKS5 Bytestreams
// http://xmpp.org/extensions/xep-0065.html
package xep

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
)

type ByteStreamsQuery struct {
	XMLName  xml.Name        `xml:"http://jabber.org/protocol/bytestreams query"`
	Sid      string          `xml:"sid,attr,omitempty"`
	Mode     string          `xml:"mode,attr,omitempty"`
	Hosts    []*StreamHost   `xml:"streamhost"`
	HostUsed *StreamHostUsed `xml:"streamhost-used"`
}

func NewByteStreamQuery(sid, mode string, usedHost *StreamHostUsed, hosts ...*StreamHost) *ByteStreamsQuery {
	return &ByteStreamsQuery{
		Sid:      sid,
		Mode:     mode,
		Hosts:    hosts,
		HostUsed: usedHost,
	}
}

func (q ByteStreamsQuery) Name() string {
	return q.XMLName.Space + " " + q.XMLName.Local
}

func (q ByteStreamsQuery) String() string {
	s := "[bytestreams query] " + q.Sid + " " + q.Mode
	for _, host := range q.Hosts {
		s += "\n" + host.String()
	}
	return s
}

type StreamHost struct {
	Host string `xml:"host,attr,omitempty"`
	Jid  string `xml:"jid,attr,omitempty"`
	Port string `xml:"port,attr,omitempty"`
}

func NewStreamHost(host, port, jid string) *StreamHost {
	return &StreamHost{
		Host: host,
		Port: port,
		Jid:  jid,
	}
}
func (h StreamHost) String() string {
	return h.Host + ":" + h.Port + " (" + h.Jid + ")"
}

type StreamHostUsed struct {
	Jid string `xml:"jid,attr"`
}

func NewStreamHostUsed(jid string) *StreamHostUsed {
	return &StreamHostUsed{
		Jid: jid,
	}
}

// XEP0065 6.3.2
func Sha1Addr(sid, requester, target string) string {
	fmt.Println(sid, requester, target)
	h := sha1.New()
	h.Write([]byte(sid + requester + target))
	return fmt.Sprintf("%x", h.Sum(nil))
}
