// XEP-0004: Data Forms
// http://xmpp.org/extensions/xep-0004.html
package xep

import (
	"bytes"
	"encoding/xml"
	"strings"
)

const (
	FormForm   = "form"
	FormSubmit = "submit"
	FormCancel = "cancel"
	FormResult = "result"

	FieldBool   = "boolean"
	FieldFixed  = "fixed"
	FieldHidden = "hidden"
	FieldMJid   = "jid-multi"
	FieldSJid   = "jid-single"
	FiledMList  = "list-multi"
	FieldSList  = "list-single"
	FieldMText  = "text-multi"
	FieldPText  = "text-private"
	FieldSText  = "text-single"
)

type XFormData struct {
	XMLName      xml.Name     `xml:"jabber:x:data x"`
	Type         string       `xml:"type,attr,omitempty"`
	Title        string       `xml:"title,omitempty"`
	Instructions string       `xml:"instructions,omitempty"`
	Fields       []*FormField `xml:"field"`
}

func NewFormData(typ string, title string, ins string,
	fields ...*FormField) *XFormData {
	return &XFormData{
		Type:         typ,
		Title:        title,
		Instructions: ins,
		Fields:       fields,
	}
}

func (form XFormData) Name() string {
	return "xdata"
}

func (form XFormData) String() string {
	b := &bytes.Buffer{}
	b.WriteString("[xdata " + form.Type + "] " +
		form.Title + " " + form.Instructions + "\n")
	for _, field := range form.Fields {
		b.WriteString(field.String() + "\n")
	}
	return b.String()
}

type FormField struct {
	Var   string `xml:"var,attr,omitempty"`
	Type  string `xml:"type,attr,omitempty"`
	Label string `xml:"label,attr,omitempty"`

	Desc     string        `xml:"desc,omitempty"`
	Required *string       `xml:"required"`
	Value    []string      `xml:"value,omitempty"`
	Options  []*FormOption `xml:"option"`
}

func NewFormField(typ, label, varAttr, desc string, value []string, required bool,
	opts ...*FormOption) *FormField {
	field := &FormField{
		Type:    typ,
		Label:   label,
		Var:     varAttr,
		Desc:    desc,
		Value:   value,
		Options: opts,
	}
	if required {
		field.Required = new(string)
	}

	return field
}

func (field FormField) String() string {
	b := &bytes.Buffer{}
	b.WriteString("[field " + field.Var + "] " + field.Type + " " + field.Label +
		" " + field.Desc + " " + strings.Join(field.Value, ",") + "\n")
	for _, opt := range field.Options {
		b.WriteString(opt.String() + "\n")
	}
	return b.String()
}

type FormOption struct {
	Label string `xml:"label,attr,omitempty"`
	Value string `xml:"value"`
}

func NewFormOption(label, value string) *FormOption {
	return &FormOption{
		Label: label,
		Value: value,
	}
}

func (opt FormOption) String() string {
	return "[option " + opt.Label + "] " + opt.Value
}
