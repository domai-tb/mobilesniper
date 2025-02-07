package models

import "encoding/xml"

type WsaEndpointReference struct {
	XMLName    xml.Name
	WsaAddress string `xml:"Address"`
}
