package shotcutProjectBuilder

import (
	"time"
)

type XMLConvertible interface {
	ToXMLNode() XMLNode
}

type MLT struct {
	LC_NUMERIC string
	Version    string
	Title      string
	Producer   string
	Children   []XMLConvertible
}

func MLTNew(title, producer string) MLT {
	return MLT{
		LC_NUMERIC: MLT_LC_NUMERIC,
		Version:    MLT_VERSION,
		Title:      title,
		Producer:   producer,
		Children:   []XMLConvertible{},
	}
}

func (m MLT) ToXMLNode() XMLNode {
	childNodes := make([]XMLNode, len(m.Children))
	for i, child := range m.Children {
		childNodes[i] = child.ToXMLNode()
	}

	attrs := []XMLNodeAttr{
		{"LC_NUMERIC", m.LC_NUMERIC},
		{"version", m.Version},
		{"title", m.Title},
		{"producer", m.Producer},
	}

	return XMLNode{"mlt", attrs, childNodes, ""}
}

type Profile struct {
	Description      string
	Width            string
	Height           string
	Progressive      string
	SampleAspectNum  string
	SampleAspectDen  string
	DisplayAspectNum string
	DisplayAspectDen string
	FrameRateNum     string
	FrameRateDen     string
	Colorspace       string
}

func (p Profile) ToXMLNode() XMLNode {
	attrs := []XMLNodeAttr{
		{"description", p.Description},
		{"width", p.Width},
		{"height", p.Height},
		{"progressive", p.Progressive},
		{"sample_aspect_num", p.SampleAspectNum},
		{"sample_aspect_den", p.SampleAspectDen},
		{"display_aspect_num", p.DisplayAspectNum},
		{"display_aspect_den", p.DisplayAspectDen},
		{"frame_rate_num", p.FrameRateNum},
		{"frame_rate_den", p.FrameRateDen},
		{"colorspace", p.Colorspace},
	}

	return XMLNode{"profile", attrs, nil, ""}
}

type Playlist struct {
	Id       string
	Title    string
	Children []XMLConvertible
}

func (p Playlist) ToXMLNode() XMLNode {
	childNodes := make([]XMLNode, len(p.Children))
	for i, child := range p.Children {
		childNodes[i] = child.ToXMLNode()
	}

	attrs := []XMLNodeAttr{
		{"id", p.Id},
		{"title", p.Title},
	}

	return XMLNode{"playlist", attrs, childNodes, ""}
}

type PlaylistEntry struct {
	Producer string
	In       string
	Out      string
}

func (pe PlaylistEntry) ToXMLNode() XMLNode {
	attrs := []XMLNodeAttr{
		{"producer", pe.Producer},
		{"in", pe.In},
		{"out", pe.Out},
	}

	return XMLNode{"entry", attrs, nil, ""}
}

type Producer struct {
	Id       string
	In       string
	Out      string
	Children []XMLConvertible
}

func (p Producer) ToXMLNode() XMLNode {
	childNodes := make([]XMLNode, len(p.Children))
	for i, child := range p.Children {
		childNodes[i] = child.ToXMLNode()
	}

	attrs := []XMLNodeAttr{
		{"id", p.Id},
		{"in", p.In},
		{"out", p.Out},
	}

	return XMLNode{"producer", attrs, childNodes, ""}
}

type Chain struct {
	DefinedDuration time.Duration
	Id              string
	Out             string
	Children        []XMLConvertible
}

func (c Chain) ToXMLNode() XMLNode {
	childNodes := make([]XMLNode, len(c.Children))
	for i, child := range c.Children {
		childNodes[i] = child.ToXMLNode()
	}

	attrs := []XMLNodeAttr{
		{"id", c.Id},
		{"out", c.Out},
	}

	return XMLNode{"chain", attrs, childNodes, ""}
}

type Filter struct {
	Id       string
	Out      string
	Children []XMLConvertible
}

func (f Filter) ToXMLNode() XMLNode {
	attrs := []XMLNodeAttr{
		{"id", f.Id},
		{"out", f.Out},
	}

	childNodes := make([]XMLNode, len(f.Children))
	for i, child := range f.Children {
		childNodes[i] = child.ToXMLNode()
	}

	return XMLNode{"filter", attrs, childNodes, ""}
}

type Property struct {
	Name        string
	TextContent string
}

func (p Property) ToXMLNode() XMLNode {
	attrs := []XMLNodeAttr{
		{"name", p.Name},
	}
	return XMLNode{"property", attrs, nil, p.TextContent}
}

type Blank struct {
	Length string
}

func (p Blank) ToXMLNode() XMLNode {
	attrs := []XMLNodeAttr{
		{"length", p.Length},
	}
	return XMLNode{"blank", attrs, nil, ""}
}

type Tractor struct {
	Id       string
	Title    string
	Version  string
	In       string
	Out      string
	Children []XMLConvertible
}

func (t Tractor) ToXMLNode() XMLNode {
	attrs := []XMLNodeAttr{
		{"id", t.Id},
		{"title", t.Title},
		{"version", t.Version},
		{"in", t.In},
		{"out", t.Out},
	}

	childNodes := make([]XMLNode, len(t.Children))
	for i, child := range t.Children {
		childNodes[i] = child.ToXMLNode()
	}

	return XMLNode{"tractor", attrs, childNodes, ""}
}

type XMLTagEntity struct {
	Attrs    map[string]string
	Children []XMLTagEntity
}
