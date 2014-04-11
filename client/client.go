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
	"github.com/ginuerzh/goxmpp/xep"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
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

type HandlerFunc func(st core.Stan)

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

	// handlers for received stanzas
	//handlers map[string]HandlerFunc

	dec *xml.Decoder
	enc *xml.Encoder

	Opts *Options

	iqw *iqWait
}

func NewClient(host, user, pwd string, opts *Options) *Client {
	return &Client{
		Host:     host,
		User:     user,
		Password: pwd,
		//handlers: make(map[string]HandlerFunc),
		Opts: opts,
		iqw:  &iqWait{m: make(map[string]*iqResp)},
	}
}

func (c *Client) Run(handler HandlerFunc) error {
	for {
		st, err := c.Recv()
		if err != nil {
			fmt.Println(err)
			return err
		}

		if st.Name() == "iq" {
			if iq, ok := st.(*xmpp.IQDefault); ok {
				switch iq.Elem().(type) {
				case *xep.Ping:
					c.Send(xmpp.NewIQ("result", st.Id(), "", nil))
				default:
					fmt.Println("unknown iq:", iq.Name())
				}
				continue
			}
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

func (c *Client) Send(st core.Stan) error {
	return c.send(st)
}

func (c *Client) SendIQ(iq core.IQ) (core.IQ, error) {
	ch := c.iqw.WaitChan(iq.Id(), iq)
	defer c.iqw.Clean(iq.Id())

	if err := c.send(iq); err != nil {
		return nil, err
	}
	st := core.IQ(<-ch)
	if st.Error() != nil {
		return st, st.Error()
	}
	return st, nil
}

func (c *Client) send(e core.Element) error {
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
	if err := c.recv(&core.Stream{}); err != nil {
		return nil, err
	}

	f := &core.StreamFeatures{}
	if err := c.recv(f); err != nil {
		return nil, errors.New("unmarshal <features>: " + err.Error())
	}
	return f, nil
}

func (c *Client) Recv() (core.Stan, error) {
	se, err := nextStart(c.dec)
	if err != nil {
		return nil, err
	}
	var e core.Stan
	switch se.Name.Space + " " + se.Name.Local {
	case xmpp.NSClient + " iq":
		id := ""
		var iq core.IQ

		for _, attr := range se.Attr {
			if attr.Name.Local == "id" {
				id = attr.Value
				break
			}
		}
		resp := c.iqw.Get(id)
		if resp == nil {
			iq = &xmpp.IQDefault{}
		} else {
			iq = resp.st
		}

		if err := c.dec.DecodeElement(iq, &se); err != nil {
			return nil, err
		}
		if resp != nil {
			resp.ch <- iq
		}
		return iq, nil
	case xmpp.NSClient + " presence":
		e = &core.StanPresence{}
	case xmpp.NSClient + " message":
		e = &xmpp.StanMsg{}
	default:
		return nil, errors.New("unexpected XMPP message " +
			se.Name.Space + " <" + se.Name.Local + "/>")
	}

	// Unmarshal into that storage.
	if err = c.dec.DecodeElement(e, &se); err != nil {
		return nil, err
	}

	return e, nil
}

func (c *Client) recv(e core.Element) (err error) {
	se, err := nextStart(c.dec)
	if err != nil {
		return err
	}

	switch se.Name.Space + " " + se.Name.Local {
	// stream start element
	case xmpp.NSStream + " stream":
		if _, ok := e.(*core.Stream); !ok {
			return errors.New("xmpp: expected <stream> but got <" +
				se.Name.Local + "> in " + se.Name.Space)
		}
		return nil
		// stream error
	case xmpp.NSStream + " error":
		err = &core.StreamError{}
	case xmpp.NSSASL + " failure":
		err = &core.SaslFailure{}
		// sasl abort
	case xmpp.NSSASL + " abort":
		err = &core.SaslAbort{}
		// tls failture
	case xmpp.NSTLS + " failure":
		err = &core.TlsFailure{}
	}

	if err != nil {
		e = err.(core.Element)
	}

	if err := c.dec.DecodeElement(e, &se); err != nil {
		return err
	}

	return
}

func (c *Client) request(req core.Element, resp core.Element) error {
	if err := c.send(req); err != nil {
		return err
	}
	return c.recv(resp)
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
		if err := c.request(&core.TlsStartTLS{}, &core.TlsProceed{}); err != nil {
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
			if err := c.request(
				&core.SaslAuth{Mechanism: m, Value: saslAuthPlain(user, c.Password)},
				&core.SaslSuccess{}); err != nil {
				return err
			}
			break
		}
		if m == "DIGEST-MD5" {
			mechanism = m
			var ch core.SaslChallenge
			if err := c.request(&core.SaslAuth{Mechanism: m}, &ch); err != nil {
				return errors.New("unmarshal <challenge>: " + err.Error())
			}
			if err := c.request(
				&core.SaslResponse{Value: saslAuthDigestMd5(ch.Value, domain, user, c.Password)},
				&core.SaslSuccess{}); err != nil {
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
	bind := &core.FeatureBind{Resource: c.Opts.Resource}
	iq := xmpp.NewIQ("set", GenId(), "", bind)
	if err := c.request(iq, iq); err != nil {
		return errors.New("bind: " + err.Error())
	}

	c.Jid = xmpp.NewJID(bind.Jid) // our local id
	fmt.Println("Jid:", c.Jid)

	// open session
	if features.Session != nil {
		iq = xmpp.NewIQ("set", GenId(), "", &core.FeatureSession{})
		if err := c.request(iq, iq); err != nil {
			return errors.New("session: " + err.Error())
		}
	}

	//c.Send(&Presence{})
	//c.Send(&IQDiscoItems{})
	//c.Send(&IQDiscoInfo{})

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
			return xml.StartElement{}, errors.New("End element: " + t.Name.Local)
		}
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

type iqWait struct {
	m map[string]*iqResp
}

func (w *iqWait) Get(id string) *iqResp {
	return w.m[id]
}

func (w *iqWait) WaitChan(id string, stan core.IQ) <-chan core.IQ {
	resp := &iqResp{
		st: stan,
		ch: make(chan core.IQ, 1),
	}
	w.m[id] = resp
	return resp.ch
}

func (w *iqWait) Clean(id string) {
	delete(w.m, id)
}

type iqResp struct {
	st core.IQ
	ch chan core.IQ
}
