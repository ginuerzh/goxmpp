// client
package client

import (
	"bufio"
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	xmpp "github.com/ginuerzh/goxmpp"
	"github.com/ginuerzh/goxmpp/core"
	"io"
	//"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Options struct {
	// Resource specifies an XMPP client resource, like "bot", instead of accepting one
	// from the server.  Use "" to let the server generate one for your client.
	Resource string

	// NoTLS disables TLS and specifies that a plain old unencrypted TCP connection should
	// be used.
	NoTLS bool

	TlsConfig *tls.Config

	// Debug output
	Debug bool
}

type HandlerFunc func(st xmpp.Stan)

type Client struct {
	// Host specifies what host to connect to, as either "hostname" or "hostname:port"
	// If host is not specified, the  DNS SRV should be used to find the host from the domainpart of the JID.
	// Default the port to 5222.
	Host string

	// User specifies what user to authenticate to the remote server.
	User string

	// Password supplies the password to use for authentication with the remote server.
	Password string

	// Jabber ID for our connection
	Jid xmpp.JID

	// connection to server
	conn *Conn

	dec *xml.Decoder
	enc *xml.Encoder

	Opts *Options

	sendChan chan xmpp.Element
	rt       *roundTrip
}

func NewClient(host, user, pwd string, opts *Options) *Client {
	ch := make(chan xmpp.Element, 10)

	return &Client{
		Host:     host,
		User:     user,
		Password: pwd,
		Opts:     opts,
		sendChan: ch,
		rt:       NewRoundTrip(ch),
	}
}

func (c *Client) Run(handler HandlerFunc) error {
	go func() {
		for {
			v := <-c.sendChan
			if err := c.enc.Encode(v); err != nil {
				return
			}
		}
	}()

	for {
		st, err := c.Recv()
		if err != nil {
			fmt.Println("Recv error:", err)
			return err
		}

		if handler != nil {
			handler(st)
		}
	}

	panic("unreachable")

}

func (c *Client) Init() error {

	if c.Opts == nil {
		c.Opts = &Options{}
	}

	host := c.Host
	conn, err := connect(host, c.User, c.Password)
	if err != nil {
		return err
	}

	if c.Opts.Debug {
		c.conn = NewConn(conn, os.Stdout)
	} else {
		c.conn = NewConn(conn, nil)
	}

	if !c.Opts.NoTLS {
		if c.conn.c, err = tlsHandShake(c.conn.c, c.Host, c.Opts.TlsConfig); err != nil {
			return err
		}
	}

	if err := c.init(); err != nil {
		//c.Close()
		return err
	}

	return nil
}
func connect(host, user, passwd string) (net.Conn, error) {
	addr := host

	if strings.TrimSpace(host) == "" {
		a := strings.SplitN(user, "@", 2)
		if len(a) == 2 {
			host = a[1]
		}
	}
	a := strings.SplitN(host, ":", 2)
	if len(a) == 1 {
		host += ":5222"
	}
	proxy := os.Getenv("HTTP_PROXY")
	if proxy == "" {
		proxy = os.Getenv("http_proxy")
	}
	if proxy != "" {
		url, err := url.Parse(proxy)
		if err == nil {
			addr = url.Host
		}
	}
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	if proxy != "" {
		fmt.Fprintf(c, "CONNECT %s HTTP/1.1\r\n", host)
		fmt.Fprintf(c, "Host: %s\r\n", host)
		fmt.Fprintf(c, "\r\n")
		br := bufio.NewReader(c)
		req, _ := http.NewRequest("CONNECT", host, nil)
		resp, err := http.ReadResponse(br, req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != 200 {
			f := strings.SplitN(resp.Status, " ", 2)
			return nil, errors.New(f[1])
		}
	}
	return c, nil
}

func streamElement(domain string) []byte {
	return []byte("<stream:stream " +
		"xmlns='jabber:client' " +
		"xmlns:stream='http://etherx.jabber.org/streams' " +
		"version='1.0'" +
		" to='" + domain + "'>")
}

func (c *Client) Send(st xmpp.Stan) error {
	c.sendChan <- st
	return nil
}

func (c *Client) SendIQ(iq xmpp.Stan) (xmpp.Stan, error) {
	return c.rt.Request(iq)
}

func (c *Client) send(e xmpp.Element) error {
	return c.enc.Encode(e)
}

func (c *Client) sendRaw(data []byte) error {
	_, err := c.conn.Write(data)
	return err
}

func (c *Client) openStream(domain string) (*core.StreamFeatures, error) {
	if err := c.sendRaw(streamElement(domain)); err != nil {
		return nil, err
	}
	if _, err := c.recv(); err != nil {
		return nil, err
	}

	f, err := c.recv()
	if err != nil || f.Name() != "features" {
		return nil, errors.New("unmarshal <features>: " + err.Error())
	}
	return f.(*core.StreamFeatures), nil
}

func (c *Client) Recv() (xmpp.Stan, error) {
	e, err := c.recv()
	if err != nil {
		return nil, err
	}

	st, ok := e.(*xmpp.Stanza)
	if !ok {
		return nil, errors.New("Not stanza: " + e.Name())
	}
	c.rt.Put(st)

	return st, nil
}

func (c *Client) recv() (xmpp.Element, error) {
	se, err := nextStart(c.dec)
	if err != nil {
		return nil, err
	}

	elemName := se.Name.Space + " " + se.Name.Local
	elem := xmpp.E(elemName)
	if elem == nil {
		fmt.Println("Unknown element:", elemName)
		return new(xmpp.NullElement), nil
	}
	switch elemName {
	// stream start element
	case xmpp.NSStream + " stream":
		return elem, nil
	case xmpp.NSClient + " iq", xmpp.NSClient + " message", xmpp.NSClient + " presence":
		return decodeStan(c.dec, &se)
	}

	//TODO: nil element handling

	if err := c.dec.DecodeElement(elem, &se); err != nil {
		return nil, err
	}

	return elem, nil
}

func (c *Client) request(req xmpp.Element) (xmpp.Element, error) {
	if err := c.send(req); err != nil {
		return nil, err
	}
	return c.recv()
}

func (c *Client) init() error {
	c.dec = xml.NewDecoder(c.conn)
	c.enc = xml.NewEncoder(c.conn)

	a := strings.SplitN(c.User, "@", 2)
	if len(a) != 2 {
		return errors.New("xmpp: invalid username (want user@domain): " + c.User)
	}
	user := a[0]
	domain := a[1]

	features, err := c.openStream(domain)
	if err != nil {
		return err
	}

	if features.StartTLS != nil && features.StartTLS.Require != nil {
		_, err := c.request(&core.TlsStartTLS{})
		if err != nil {
			return err
		}

		if c.conn.c, err = tlsHandShake(c.conn.c, domain, c.Opts.TlsConfig); err != nil {
			return err
		}

		features, err = c.openStream(domain)
		if err != nil {
			return err
		}
	}

	mechanism := ""
	for _, m := range features.Mechanisms.Mechanism {
		if m == "PLAIN" {
			mechanism = m
			if _, err := c.request(
				&core.SaslAuth{Mechanism: m,
					Value: saslAuthPlain(user, c.Password)}); err != nil {
				return err
			}
			break
		}
		if m == "DIGEST-MD5" {
			mechanism = m
			//var ch core.SaslChallenge
			ch, err := c.request(&core.SaslAuth{Mechanism: m})
			if err != nil || ch.Name() != "challenge" {
				return errors.New("unmarshal <challenge>: " + err.Error())
			}

			if _, err := c.request(
				&core.SaslResponse{
					Value: saslAuthDigestMd5(ch.(*core.SaslChallenge).Value, domain, user, c.Password)}); err != nil {
				return errors.New("unmarshal <success>: " + err.Error())
			}
			break
		}
	}
	if mechanism == "" {
		return errors.New(
			fmt.Sprintf("PLAIN authentication is not an option: %v",
				features.Mechanisms.Mechanism))
	}

	// Now that we're authenticated, we're supposed to start the stream over again.
	// Declare intent to be a jabber client.
	// Here comes another <stream> and <features>.
	features, err = c.openStream(domain)
	if err != nil {
		return err
	}

	// Send IQ message asking to bind to the local user name.
	iq, err := c.request(xmpp.NewIQ("set", GenId(), "",
		&core.FeatureBind{Resource: c.Opts.Resource}))
	if err != nil {
		return errors.New("bind: " + err.Error())
	}

	if err = iq.(*xmpp.Stanza).Error(); err != nil {
		fmt.Println(iq.(*xmpp.Stanza).Error())
		return errors.New("bind: " + err.Error())
	}
	bind := iq.(*xmpp.Stanza).Elements[0].(*core.FeatureBind)
	c.Jid = xmpp.NewJID(bind.Jid) // our local id
	fmt.Println("Jid:", c.Jid)

	// open session
	if features.Session != nil {
		iq, err = c.request(xmpp.NewIQ("set", GenId(), "",
			&core.FeatureSession{}))
		if err != nil {
			return errors.New("session: " + err.Error())
		}
		if err := iq.(*xmpp.Stanza).Error(); err != nil {
			return errors.New("session: " + err.Error())
		}
	}
	return nil
}

// Scan XML token stream to find next StartElement.
func nextStart(p *xml.Decoder) (xml.StartElement, error) {
	for {
		t, err := p.Token()
		if err != nil && err != io.EOF {
			fmt.Println(err)
			return xml.StartElement{}, err
		}
		switch t := t.(type) {
		case xml.StartElement:
			return t, nil
		case xml.EndElement:
			return xml.StartElement{}, errors.New("Unexpected end element: " + t.Name.Local)
		}
	}
	panic("unreachable")
}

func nextElement(p *xml.Decoder) (xml.Token, error) {
	for {
		t, err := p.Token()
		if err != nil && err != io.EOF {
			fmt.Println(err)
			return nil, err
		}
		switch t := t.(type) {
		case xml.StartElement, xml.EndElement:
			return t, nil
		default:
			fmt.Println("unknown element")
		}
	}
	panic("unreachable")
}

func decodeStan(p *xml.Decoder, start *xml.StartElement) (xmpp.Stan, error) {
	st := xmpp.NewStanza(start.Name.Local)

	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "id":
			st.Ids = attr.Value
		case "type":
			st.Types = attr.Value
		case "from":
			st.From = attr.Value
		case "to":
			st.To = attr.Value
		case "lang":
			st.Lang = attr.Value
		default:
			fmt.Println("unknown stanza attr:", attr.Name.Local)
		}
	}

	for {
		t, err := nextElement(p)
		if err != nil {
			return nil, err
		}

		var se xml.StartElement
		switch t := t.(type) {
		case xml.EndElement:
			if t.Name.Local != st.Name() {
				return nil, errors.New("Unexpected end element: " +
					t.Name.Local + ", should be '</" + st.Name() + ">'")
			}
			return st, nil

		case xml.StartElement:
			se = t
		}

		elem := xmpp.E(se.Name.Space + " " + se.Name.Local)
		if elem == nil {
			elem = &xmpp.NullElement{}
		}
		if err := p.DecodeElement(elem, &se); err != nil {
			return nil, err
		}
		st.AddElement(elem)
	}

	panic("unreachable")
}

type Conn struct {
	c net.Conn
	w io.Writer
}

func NewConn(conn net.Conn, logger io.Writer) *Conn {
	return &Conn{
		c: conn,
		w: logger,
	}
}

func (t *Conn) Read(p []byte) (n int, err error) {
	n, err = t.c.Read(p)
	if n > 0 && t.w != nil {
		t.w.Write([]byte(">>> "))
		t.w.Write(p[0:n])
		t.w.Write([]byte("\n"))
	}
	return
}

func (t *Conn) Write(p []byte) (n int, err error) {
	n, err = t.c.Write(p)
	if n > 0 && t.w != nil {
		t.w.Write([]byte("<<< "))
		t.w.Write(p[:n])
		t.w.Write([]byte("\n"))
	}
	return
}

func (t *Conn) Close() error {
	return t.c.Close()
}

type roundTrip struct {
	sendChan chan<- xmpp.Element
	timeout  time.Duration
	m        map[string]chan xmpp.Stan
}

func NewRoundTrip(sendChan chan<- xmpp.Element) *roundTrip {
	return &roundTrip{
		sendChan: sendChan,
		timeout:  3 * time.Second,
		m:        make(map[string]chan xmpp.Stan),
	}
}

func (this *roundTrip) Request(iq xmpp.Stan) (resp xmpp.Stan, err error) {
	ch := make(chan xmpp.Stan, 1)
	this.m[iq.Id()] = ch

	defer delete(this.m, iq.Id())

	for retry := 3; retry > 0; retry-- {
		this.sendChan <- iq

		select {
		case <-time.NewTimer(this.timeout).C:
			err = errors.New("Time-out")
			continue
		case v := <-ch:
			return v, nil
		}
	}
	return
}

func (this *roundTrip) Put(iq xmpp.Stan) bool {
	ch, ok := this.m[iq.Id()]
	if !ok {
		return false
	}
	ch <- iq

	return true
}
