package models

import (
	"encoding/xml"
	"fmt"
)

type ProbeMatch struct {
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
	WsaEndpointReference WsaEndpointReference `xml:"EndpointReference"`
	WsdTypes             string               `xml:"Types"`
	WsdScopes            string               `xml:"Scopes"`
	WsdXAddrs            string               `xml:"XAddrs"`
	WsdMetadataVersion   string               `xml:"MetadataVersion"`
}

// Getter of ProbeMatchs Header

func (p *ProbeMatch) GetAction() string {
	return p.ProbeMatchHeader.WsaAction
}

func (p *ProbeMatch) GetTo() string {
	return p.ProbeMatchHeader.WsaTo
}

func (p *ProbeMatch) GetRelatesTo() string {
	return p.ProbeMatchHeader.WsaRelatesTo
}

func (p *ProbeMatch) GetMessageId() string {
	return p.ProbeMatchHeader.WsaMessageId
}

func (p *ProbeMatch) GetAppSequence() string {
	return fmt.Sprintf("%s (InstanceId=%s MessageNumber=%s)", p.ProbeMatchHeader.WsdAppSequence, p.ProbeMatchHeader.WsdInstanceId, p.ProbeMatchHeader.WsdMessageNumber)
}

// Getter of ProbeMatchs Body

func (p *ProbeMatch) GetXAddrs() string {
	return p.ProbeMatchBody.WsdProbeMatches.WsdProbeMatch.WsdXAddrs
}

func (p *ProbeMatch) GetTypes() string {
	return p.ProbeMatchBody.WsdProbeMatches.WsdProbeMatch.WsdTypes
}

func (p *ProbeMatch) GetScopes() string {
	return p.ProbeMatchBody.WsdProbeMatches.WsdProbeMatch.WsdScopes
}

func (p *ProbeMatch) GetMetadataVersion() string {
	return p.ProbeMatchBody.WsdProbeMatches.WsdProbeMatch.WsdMetadataVersion
}

func (p *ProbeMatch) GetEndpointReferenceAddress() string {
	return p.ProbeMatchBody.WsdProbeMatches.WsdProbeMatch.WsaEndpointReference.WsaAddress
}
