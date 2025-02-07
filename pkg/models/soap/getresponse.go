package models

import "encoding/xml"

type GetResponse struct {
	XMLName           xml.Name
	GetResponseHeader GetResponseHeader `xml:"Header"`
	GetResponseBody   GetResponseBody   `xml:"Body"`
}

type GetResponseHeader struct {
	XMLName      xml.Name
	WsaAction    string `xml:"Action"`
	WsaTo        string `xml:"To"`
	WsaRelatesTo string `xml:"RelatesTo"`
	WsaMessageId string `xml:"MessageID"`
}

type GetResponseBody struct {
	XMLName     xml.Name
	MexMetadata MexMetadata `xml:"Metadata"`
}

type MexMetadata struct {
	XMLName             xml.Name
	MexMetadataSections []MexMetadataSection `xml:"MetadataSection"`
}

type MexMetadataSection struct {
	XMLName                   xml.Name
	MexMetadataSectionDialect string               `xml:"Dialect,attr"`
	DpwsThisModel             DpwsThisModel        `xml:"ThisModel"`
	DpwsThisDevice            DpwsThisDevice       `xml:"ThisDevice"`
	DpwsThisRelationship      DpwsThisRelationship `xml:"Relationship"`
	MexLocation               string               `xml:"Location"`
}

type DpwsThisModel struct {
	XMLName             xml.Name
	DpwsManufacturer    string `xml:"Manufacturer"`
	DpwsManufacturerUrl string `xml:"ManufacturerUrl"`
	DpwsModelName       string `xml:"ModelName"`
	DpwsModelUrl        string `xml:"ModelNumber"`
	DpwsPresentationUrl string `xml:"PresentationUrl"`
}

type DpwsThisDevice struct {
	XMLName             xml.Name
	DpwsFriendlyName    string `xml:"FriendlyName"`
	DpwsFirmwareVersion string `xml:"FirmwareVersion"`
	DpwsSerialNumber    string `xml:"SerialNumber"`
}

type DpwsThisRelationship struct {
	XMLName    xml.Name
	DpwsHost   DpwsHost   `xml:"Host"`
	DpwsHosted DpwsHosted `xml:"Hosted"`
}

type DpwsHost struct {
	XMLName              xml.Name
	WsaEndpointReference WsaEndpointReference `xml:"EndpointReference"`
	DpwsTypes            string               `xml:"Types"`
}

type DpwsHosted struct {
	XMLName              xml.Name
	WsaEndpointReference WsaEndpointReference `xml:"EndpointReference"`
	DpwsTypes            string               `xml:"Types"`
	DpwsServiceID        string               `xml:"ServiceId"`
}
