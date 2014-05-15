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
	//"github.com/golang/glog"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Options struct {
	// Resource specifies an XMPP client resource, like "bot", instead of accepting one
	// from the server.  Use "" to let the server generate one for your client.
	Resource string

	// NoTLS disables TLS and specifies that a plain old unencrypted TCP connection should
	// be used.
	NoTLS bool

	Proxy string

	TlsConfig *tls.Config

	// Debug output
	Debug bool
}

type HandlerFunc func(stanza *core.StanzaHeader, e xmpp.Element)
type LoginFunc func(err error)
type ErrorFunc func(err error)

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
	recvChan chan xmpp.Stan
	rt       *roundTrip

	handlers     map[string]HandlerFunc
	loginHandler LoginFunc
	errorHandler ErrorFunc
}

func NewClient(host, user, pwd string, opts *Options) *Client {
	ch := make(chan xmpp.Element, 10)

	return &Client{
		Host:     host,
		User:     user,
		Password: pwd,
		Opts:     opts,
		sendChan: ch,
		recvChan: make(chan xmpp.Stan, 10),
		rt:       NewRoundTrip(ch),
		handlers: make(map[string]HandlerFunc),
	}
}

func (c *Client) HandleFunc(fullName string, handler HandlerFunc) {
	c.handlers[fullName] = handler
}

func (c *Client) OnLogined(loginFunc LoginFunc) {
	c.loginHandler = loginFunc
}

func (c *Client) OnError(errFunc ErrorFunc) {
	c.errorHandler = errFunc
}

func (c *Client) Run() error {
	exit := make(chan error, 1)

	err := c.Login()
	if c.loginHandler != nil {
		go c.loginHandler(err)
	}
	if err != nil {
		return err
	}

	go func() {
		for {
			_, err := c.Recv()
			if err != nil {
				exit <- err
				return
			}
		}
	}()

	for {
		select {
		case v := <-c.sendChan:
			if err := c.enc.Encode(v); err != nil {
				exit <- err
			}
		case <-c.recvChan:
		case err := <-exit:
			if c.errorHandler != nil {
				go c.errorHandler(err)
			}
			return err
		}
	}

	panic("unreachable")

}

func (c *Client) Login() error {

	if c.Opts == nil {
		c.Opts = &Options{}
	}

	host := c.Host
	conn, err := connect(host, c.User, c.Password, c.Opts.Proxy)
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

func (c *Client) Close() error {
	c.sendRaw([]byte("</stream:stream>"))
	return c.conn.Close()
}

func connect(host, user, passwd, proxy string) (net.Conn, error) {
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

	addr := host
	if len(proxy) > 0 {
		addr = proxy
	}

	c, err := net.DialTimeout("tcp", addr, time.Second*10)
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

	c.recvChan <- st

	if ok := c.rt.Put(st); ok {
		return st, nil
	}

	for _, e := range st.E() {
		if handler, ok := c.handlers[e.FullName()]; ok {
			go handler(&st.StanzaHeader, e)
		}
	}

	if handler, ok := c.handlers[st.FullName()]; ok {
		go handler(&st.StanzaHeader, e)
	}

	return st, nil
}

func (c *Client) recv() (xmpp.Element, error) {
	se, err := nextStart(c.dec)
	if err != nil {
		return nil, err
	}

	elem := xmpp.E(se.Name.Space + " " + se.Name.Local)
	if elem == nil {
		fmt.Println("Unknown element:", elem.FullName())
		return new(xmpp.NullElement), nil
	}
	switch elem.Name() {
	// stream start element
	case "stream":
		return elem, nil
	case "iq", "message", "presence":
		return c.decodeStan(&se)
	}

	//TODO: nil element handling

	if err := c.dec.DecodeElement(elem, &se); err != nil {
		return nil, err
	}

	if err, ok := checkError(elem); ok {
		return nil, err
	}

	return elem, nil
}

func checkError(e xmpp.Element) (error, bool) {
	switch v := e.(type) {
	case *core.StreamError:
		return v, true
	case *core.SaslAbort:
		return v, true
	case *core.SaslFailure:
		return v, true
	case *core.TlsFailure:
		return v, true
	}
	return nil, false
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
	c.Jid = xmpp.ToJID(bind.Jid) // our local id
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

func (c *Client) decodeStan(start *xml.StartElement) (xmpp.Stan, error) {
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
		t, err := nextElement(c.dec)
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
		if err := c.dec.DecodeElement(elem, &se); err != nil {
			return nil, err
		}
		if err, ok := elem.(*core.StanzaError); ok {
			st.Err = err
		} else {
			st.AddE(elem)
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

type roundTrip struct {
	sendChan chan<- xmpp.Element
	timeout  time.Duration
	m        map[string]chan xmpp.Stan
	lock     *sync.RWMutex
}

func NewRoundTrip(sendChan chan<- xmpp.Element) *roundTrip {
	return &roundTrip{
		sendChan: sendChan,
		timeout:  60 * time.Second,
		m:        make(map[string]chan xmpp.Stan),
		lock:     new(sync.RWMutex),
	}
}

func (this *roundTrip) Request(iq xmpp.Stan) (resp xmpp.Stan, err error) {
	ch := make(chan xmpp.Stan, 1)

	this.lock.Lock()
	this.m[iq.Id()] = ch
	this.lock.Unlock()

	defer func() {
		this.lock.Lock()
		delete(this.m, iq.Id())
		this.lock.Unlock()
	}()

	for retry := 1; retry > 0; retry-- {
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
	this.lock.RLock()
	ch, ok := this.m[iq.Id()]
	this.lock.RUnlock()
	if !ok {
		return false
	}
	ch <- iq

	return true
}
