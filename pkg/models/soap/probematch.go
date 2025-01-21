package models

import (
	"encoding/xml"
	"fmt"
)

type ProbeMatchMsg struct {
	XMLName          xml.Name
	ProbeMatchHeader ProbeMatchHeader `xml:"Header"`
	ProbeMatchBody   ProbeMatchBody   `xml:"Body"`
}

type ProbeMatchHeader struct {
	XMLName          xml.Name
	WsaAction        string `xml:"Action"`
	WsaTo            string `xml:"To"`
	WsaRelatesTo     string `xml:"RelatesTo"`
	WsaMessageId     string `xml:"MessageID"`
	WsdAppSequence   string `xml:"AppSequence"`
	WsdInstanceId    string `xml:"InstanceId,attr"`
	WsdMessageNumber string `xml:"MessageNumber,attr"`
}

type ProbeMatchBody struct {
	XMLName         xml.Name
	WsdProbeMatches WsdProbeMatches `xml:"ProbeMatches"`
}

type WsdProbeMatches struct {
	XMLName       xml.Name
	WsdProbeMatch WsdProbeMatch `xml:"ProbeMatch"`
}

type WsdProbeMatch struct {
	XMLName              xml.Name
	WsdEndpointReference WsdEndpointReference `xml:"EndpointReference"`
	WsdTypes             string               `xml:"Types"`
	WsdScopes            string               `xml:"Scopes"`
	WsdXAddrs            string               `xml:"XAddrs"`
	WsdMetadataVersion   string               `xml:"MetadataVersion"`
}

type WsdEndpointReference struct {
	XMLName    xml.Name
	WsaAddress string `xml:"Address"`
}

// Getter of ProbeMatchs Header

func (p *ProbeMatchMsg) GetAction() string {
	return p.ProbeMatchHeader.WsaAction
}

func (p *ProbeMatchMsg) GetTo() string {
	return p.ProbeMatchHeader.WsaTo
}

func (p *ProbeMatchMsg) GetRelatesTo() string {
	return p.ProbeMatchHeader.WsaRelatesTo
}

func (p *ProbeMatchMsg) GetMessageId() string {
	return p.ProbeMatchHeader.WsaMessageId
}

func (p *ProbeMatchMsg) GetAppSequence() string {
	return fmt.Sprintf("%s (InstanceId=%s MessageNumber=%s)", p.ProbeMatchHeader.WsdAppSequence, p.ProbeMatchHeader.WsdInstanceId, p.ProbeMatchHeader.WsdMessageNumber)
}

// Getter of ProbeMatchs Body

func (p *ProbeMatchMsg) GetXAddrs() string {
	return p.ProbeMatchBody.WsdProbeMatches.WsdProbeMatch.WsdXAddrs
}

func (p *ProbeMatchMsg) GetTypes() string {
	return p.ProbeMatchBody.WsdProbeMatches.WsdProbeMatch.WsdTypes
}

func (p *ProbeMatchMsg) GetScopes() string {
	return p.ProbeMatchBody.WsdProbeMatches.WsdProbeMatch.WsdScopes
}

func (p *ProbeMatchMsg) GetMetadataVersion() string {
	return p.ProbeMatchBody.WsdProbeMatches.WsdProbeMatch.WsdMetadataVersion
}

func (p *ProbeMatchMsg) GetEndpointReferenceAddress() string {
	return p.ProbeMatchBody.WsdProbeMatches.WsdProbeMatch.WsdEndpointReference.WsaAddress
}
