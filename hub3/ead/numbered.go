package ead

import (
	"encoding/xml"
)

type Cc01 struct {
	XMLName xml.Name `xml:"c01,omitempty"`
	Cc
	Numbered []*Cc02 `xml:"c02,omitempty"`
}

func (c *Cc01) GetNested() []CLevel {
	levels := []CLevel{}
	for _, cc := range c.Numbered {
		levels = append(levels, cc)
	}

	return levels
}

func (c *Cc01) GetCc() *Cc {
	return &c.Cc
}

type Cc02 struct {
	XMLName xml.Name `xml:"c02,omitempty"`
	Cc
	Numbered []*Cc03 `xml:"c03,omitempty"`
}

func (c *Cc02) GetNested() []CLevel {
	levels := []CLevel{}
	for _, cc := range c.Numbered {
		levels = append(levels, cc)
	}

	return levels
}

func (c *Cc02) GetCc() *Cc {
	return &c.Cc
}

type Cc03 struct {
	XMLName xml.Name `xml:"c03,omitempty"`
	Cc
	Numbered []*Cc04 `xml:"c04,omitempty"`
}

func (c *Cc03) GetNested() []CLevel {
	levels := []CLevel{}
	for _, cc := range c.Numbered {
		levels = append(levels, cc)
	}

	return levels
}

func (c *Cc03) GetCc() *Cc {
	return &c.Cc
}

type Cc04 struct {
	XMLName xml.Name `xml:"c04,omitempty"`
	Cc
	Numbered []*Cc05 `xml:"c05,omitempty"`
}

func (c *Cc04) GetNested() []CLevel {
	levels := []CLevel{}
	for _, cc := range c.Numbered {
		levels = append(levels, cc)
	}

	return levels
}

func (c *Cc04) GetCc() *Cc {
	return &c.Cc
}

type Cc05 struct {
	XMLName xml.Name `xml:"c05,omitempty"`
	Cc
	Numbered []*Cc06 `xml:"c06,omitempty"`
}

func (c *Cc05) GetNested() []CLevel {
	levels := []CLevel{}
	for _, cc := range c.Numbered {
		levels = append(levels, cc.GetCc())
	}

	return levels
}

func (c *Cc05) GetCc() *Cc {
	return &c.Cc
}

type Cc06 struct {
	XMLName xml.Name `xml:"c06,omitempty"`
	Cc
	Numbered []*Cc07 `xml:"c07,omitempty"`
}

func (c *Cc06) GetNested() []CLevel {
	levels := []CLevel{}
	for _, cc := range c.Numbered {
		levels = append(levels, cc.GetCc())
	}

	return levels
}

func (c *Cc06) GetCc() *Cc {
	return &c.Cc
}

type Cc07 struct {
	XMLName xml.Name `xml:"c07,omitempty"`
	Cc
	Numbered []*Cc08 `xml:"c08,omitempty"`
}

func (c *Cc07) GetNested() []CLevel {
	levels := []CLevel{}
	for _, cc := range c.Numbered {
		levels = append(levels, cc.GetCc())
	}

	return levels
}

func (c *Cc07) GetCc() *Cc {
	return &c.Cc
}

type Cc08 struct {
	XMLName xml.Name `xml:"c08,omitempty"`
	Cc
	Numbered []*Cc09 `xml:"c09,omitempty"`
}

func (c *Cc08) GetNested() []CLevel {
	levels := []CLevel{}
	for _, cc := range c.Numbered {
		levels = append(levels, cc.GetCc())
	}

	return levels
}

func (c *Cc08) GetCc() *Cc {
	return &c.Cc
}

type Cc09 struct {
	XMLName xml.Name `xml:"c09,omitempty"`
	Cc
	Numbered []*Cc10 `xml:"c10,omitempty"`
}

func (c *Cc09) GetNested() []CLevel {
	levels := []CLevel{}
	for _, cc := range c.Numbered {
		levels = append(levels, cc.GetCc())
	}

	return levels
}

func (c *Cc09) GetCc() *Cc {
	return &c.Cc
}

type Cc10 struct {
	XMLName xml.Name `xml:"c10,omitempty"`
	Cc
	Numbered []*Cc11 `xml:"c11,omitempty"`
}

func (c *Cc10) GetNested() []CLevel {
	levels := []CLevel{}
	for _, cc := range c.Numbered {
		levels = append(levels, cc.GetCc())
	}

	return levels
}

func (c *Cc10) GetCc() *Cc {
	return &c.Cc
}

type Cc11 struct {
	XMLName xml.Name `xml:"c11,omitempty"`
	Cc
	Numbered []*Cc12 `xml:"c12,omitempty"`
}

func (c *Cc11) GetNested() []CLevel {
	levels := []CLevel{}
	for _, cc := range c.Numbered {
		levels = append(levels, cc.GetCc())
	}

	return levels
}

func (c *Cc11) GetCc() *Cc {
	return &c.Cc
}

type Cc12 struct {
	XMLName xml.Name `xml:"c12,omitempty"`
	Cc
	Numbered []*Cc13 `xml:"c13,omitempty"`
}

func (c *Cc12) GetNested() []CLevel {
	levels := []CLevel{}
	for _, cc := range c.Numbered {
		levels = append(levels, cc.GetCc())
	}

	return levels
}

func (c *Cc12) GetCc() *Cc {
	return &c.Cc
}

type Cc13 struct {
	XMLName xml.Name `xml:"c13,omitempty"`
	Cc
	Numbered []*Cc14 `xml:"c14,omitempty"`
}

func (c *Cc13) GetNested() []CLevel {
	levels := []CLevel{}
	for _, cc := range c.Numbered {
		levels = append(levels, cc.GetCc())
	}

	return levels
}

func (c *Cc13) GetCc() *Cc {
	return &c.Cc
}

type Cc14 struct {
	XMLName xml.Name `xml:"c14,omitempty"`
	Cc
	Numbered []*Cc15 `xml:"c15,omitempty"`
}

func (c *Cc14) GetNested() []CLevel {
	levels := []CLevel{}
	for _, cc := range c.Numbered {
		levels = append(levels, cc.GetCc())
	}

	return levels
}

func (c *Cc14) GetCc() *Cc {
	return &c.Cc
}

type Cc15 struct {
	XMLName xml.Name `xml:"c15,omitempty"`
	Cc
	Numbered []*Cc16 `xml:"c16,omitempty"`
}

func (c *Cc15) GetNested() []CLevel {
	levels := []CLevel{}
	for _, cc := range c.Numbered {
		levels = append(levels, cc.GetCc())
	}

	return levels
}

func (c *Cc15) GetCc() *Cc {
	return &c.Cc
}

type Cc16 struct {
	XMLName xml.Name `xml:"c16,omitempty"`
	Cc
	Numbered []*Cc17 `xml:"c17,omitempty"`
}

func (c *Cc16) GetNested() []CLevel {
	levels := []CLevel{}
	for _, cc := range c.Numbered {
		levels = append(levels, cc.GetCc())
	}

	return levels
}

func (c *Cc16) GetCc() *Cc {
	return &c.Cc
}

type Cc17 struct {
	XMLName xml.Name `xml:"c17,omitempty"`
	Cc
	Numbered []*Cc18 `xml:"c18,omitempty"`
}

func (c *Cc17) GetNested() []CLevel {
	levels := []CLevel{}
	for _, cc := range c.Numbered {
		levels = append(levels, cc.GetCc())
	}

	return levels
}

func (c *Cc17) GetCc() *Cc {
	return &c.Cc
}

type Cc18 struct {
	XMLName xml.Name `xml:"c18,omitempty"`
	Cc
	Numbered []*Cc19 `xml:"c19,omitempty"`
}

func (c *Cc18) GetNested() []CLevel {
	levels := []CLevel{}
	for _, cc := range c.Numbered {
		levels = append(levels, cc.GetCc())
	}

	return levels
}

func (c *Cc18) GetCc() *Cc {
	return &c.Cc
}

type Cc19 struct {
	XMLName xml.Name `xml:"c19,omitempty"`
	Cc
	Numbered []*Cc20 `xml:"c20,omitempty"`
}

func (c *Cc19) GetNested() []CLevel {
	levels := []CLevel{}
	for _, cc := range c.Numbered {
		levels = append(levels, cc.GetCc())
	}

	return levels
}

func (c *Cc19) GetCc() *Cc {
	return &c.Cc
}

type Cc20 struct {
	XMLName xml.Name `xml:"c20,omitempty"`
	Cc
}

func (c *Cc20) GetNested() []CLevel {
	return []CLevel{}
}

func (c *Cc20) GetCc() *Cc {
	return &c.Cc
}
