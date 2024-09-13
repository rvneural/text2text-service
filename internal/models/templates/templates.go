package templates

import "encoding/xml"

type Template struct {
	XMLName     xml.Name `xml:"template" json:"-"`
	Name        string   `xml:"name" json:"name"`
	Value       string   `xml:"value" json:"-"`
	Temperature string   `xml:"temperature" json:"-"`
	Role        string   `xml:"role,attr" json:"-"`
	Description string   `xml:"description" json:"description"`
}

type Templates struct {
	XMLName   xml.Name   `xml:"templates" json:"-"`
	Templates []Template `xml:"template" json:"templates"`
}
