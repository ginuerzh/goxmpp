// util
package client

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"math/big"
	"net"
	"strconv"
	"strings"
	"time"
)

func tlsHandShake(conn net.Conn, host string, config *tls.Config) (net.Conn, error) {
	tlsconn := tls.Client(conn, config)
	if err := tlsconn.Handshake(); err != nil {
		return nil, err
	}

	if strings.LastIndex(host, ":") > 0 {
		host = host[:strings.LastIndex(host, ":")]
	}
	if err := tlsconn.VerifyHostname(host); err != nil {
		return nil, err
	}

	return tlsconn, nil
}

func saslAuthPlain(username, password string) string {
	raw := "\x00" + username + "\x00" + password
	enc := make([]byte, base64.StdEncoding.EncodedLen(len(raw)))
	base64.StdEncoding.Encode(enc, []byte(raw))

	return string(enc)
}

func saslAuthDigestMd5(change string, domain, username, password string) string {
	b, err := base64.StdEncoding.DecodeString(change)
	if err != nil {
		return ""
	}

	tokens := map[string]string{}
	for _, token := range strings.Split(string(b), ",") {
		kv := strings.SplitN(strings.TrimSpace(token), "=", 2)
		if len(kv) == 2 {
			if kv[1][0] == '"' && kv[1][len(kv[1])-1] == '"' {
				kv[1] = kv[1][1 : len(kv[1])-1]
			}
			tokens[kv[0]] = kv[1]
		}
	}
	realm, _ := tokens["realm"]
	nonce, _ := tokens["nonce"]
	qop, _ := tokens["qop"]
	charset, _ := tokens["charset"]
	cnonceStr := cnonce()
	digestUri := "xmpp/" + domain
	nonceCount := fmt.Sprintf("%08x", 1)
	digest := saslDigestResponse(username, realm, password,
		nonce, cnonceStr, "AUTHENTICATE", digestUri, nonceCount)
	message := "username=" + username +
		", realm=" + realm +
		", nonce=" + nonce +
		", cnonce=" + cnonceStr +
		", nc=" + nonceCount +
		", qop=" + qop +
		", digest-uri=" + digestUri +
		", response=" + digest +
		", charset=" + charset

	return base64.StdEncoding.EncodeToString([]byte(message))
}

func cnonce() string {
	randSize := big.NewInt(0)
	randSize.Lsh(big.NewInt(1), 64)
	cn, err := rand.Int(rand.Reader, randSize)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%016x", cn)
}

func saslDigestResponse(username, realm, passwd, nonce, cnonceStr,
	authenticate, digestUri, nonceCountStr string) string {
	h := func(text string) []byte {
		h := md5.New()
		h.Write([]byte(text))
		return h.Sum(nil)
	}
	hex := func(bytes []byte) string {
		return fmt.Sprintf("%x", bytes)
	}
	kd := func(secret, data string) []byte {
		return h(secret + ":" + data)
	}

	a1 := string(h(username+":"+realm+":"+passwd)) + ":" +
		nonce + ":" + cnonceStr
	a2 := authenticate + ":" + digestUri
	response := hex(kd(hex(h(a1)), nonce+":"+
		nonceCountStr+":"+cnonceStr+":auth:"+
		hex(h(a2))))
	return response
}

func GenId() string {
	return strconv.Itoa(time.Now().Nanosecond())
}
