// XEP-0096: SI File Transfer
// http://xmpp.org/extensions/xep-0096.html
package xep

import (
	"bytes"
	"encoding/xml"
)

type SI struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/si si"`
	Id      string   `xml:"id,attr,omitempty"`
	Mime    string   `xml:"mime-type,attr,omitempty"`
	Profile string   `xml:"profile,attr,omitempty"`
	File    *FileTransfer
	Feature *FeatureNeg
}

func NewSI(id, mime, profile string, file *FileTransfer, feature *FeatureNeg) *SI {
	return &SI{
		Id:      id,
		Mime:    mime,
		Profile: profile,
		File:    file,
		Feature: feature,
	}
}

func (si SI) Name() string {
	return "si"
}

func (si SI) String() string {
	b := &bytes.Buffer{}
	b.WriteString("[si] " + si.Id + " " + si.Mime + "\n")
	if si.File != nil {
		b.WriteString(si.File.String() + "\n")
	}
	if si.Feature != nil {
		b.WriteString(si.Feature.String() + "\n")
	}
	return b.String()
}

type FileTransfer struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/si/profile/file-transfer file"`
	Name    string   `xml:"name,attr,omitempty"`
	Size    string   `xml:"size,attr,omitempty"`
	Date    string   `xml:"date,attr,omitempty"`
	Hash    string   `xml:"hash,attr,omitempty"`
	Desc    string   `xml:"desc,omitempty"`
}

func NewFileTransfer(name, size, date, hash, desc string) *FileTransfer {
	return &FileTransfer{
		Name: name,
		Size: size,
		Date: date,
		Hash: hash,
		Desc: desc,
	}
}

func (ft FileTransfer) String() string {
	return "[file] " + ft.Name + " " + ft.Size + " " + ft.Desc
}

type FeatureNeg struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/feature-neg feature"`
	Form    *XFormData
}

func NewFeatureNeg(form *XFormData) *FeatureNeg {
	return &FeatureNeg{
		Form: form,
	}
}

func (f FeatureNeg) String() string {
	return "[feature]\n" + f.Form.String()
}
