// XEP-0153: vCard-Based Avatars
// http://xmpp.org/extensions/xep-0153.html
package xep

import (
	"encoding/xml"
)

type VCardUpdate struct {
	XMLName xml.Name `xml:"vcard-temp:x:update x"`
	Photo   string   `xml:"photo"`
}

func (_ VCardUpdate) Name() string {
	return "x"
}

func (_ VCardUpdate) FullName() string {
	return "vcard-temp:x:update x"
}

func (vc VCardUpdate) String() string {
	return vc.Photo
}
