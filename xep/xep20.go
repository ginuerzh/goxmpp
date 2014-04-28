// XEP-0020: Feature Negotiation
// http://xmpp.org/extensions/xep-0020.html
package xep

import (
	"encoding/xml"
)

type Feature struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/feature-neg feature"`
	Form    *XFormData
}

func NewFeature(form *XFormData) *Feature {
	return &Feature{
		Form: form,
	}
}

func (_ Feature) Name() string {
	return "feature"
}

func (_ Feature) FullName() string {
	return "http://jabber.org/protocol/feature-neg feature"
}

func (f Feature) String() string {
	return "[feature]\n" + f.Form.String()
}
