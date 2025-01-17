// Package opds2 provide parsing and generation method for an OPDS2 feed
// https://github.com/opds-community/opds-revision/blob/master/opds-2.0.md
package opds2

import (
	"encoding/json"
	"time"
)

// Feed is a collection as defined in Readium Web Publication Manifest
type Feed struct {
	Context      []string      `json:"@context,omitempty"`
	Metadata     Metadata      `json:"metadata"`
	Links        Links         `json:"links"`
	Facets       []Facet       `json:"facets,omitempty"`
	Groups       []Group       `json:"groups,omitempty"`
	Publications []Publication `json:"publications,omitempty"`
	Navigation   Links         `json:"navigation,omitempty"`
}

// Metadata has a limited subset of metadata compared to a publication
type Metadata struct {
	RDFType       string     `json:"@type,omitempty"`
	Title         string     `json:"title"`
	NumberOfItems int        `json:"numberOfItems,omitempty"`
	ItemsPerPage  int        `json:"itemsPerPage,omitempty"`
	CurrentPage   int        `json:"currentPage,omitempty"`
	Modified      *time.Time `json:"modified,omitempty"`
}

// Facet is a collection that contains a facet group
type Facet struct {
	Metadata Metadata `json:"metadata"`
	Links    Links    `json:"links"`
}

// Group is a group collection that must contain publications
type Group struct {
	Metadata     Metadata      `json:"metadata"`
	Links        Links         `json:"links,omitempty"`
	Publications []Publication `json:"publications,omitempty"`
	Navigation   Links         `json:"navigation,omitempty"`
}

// Properties object use to link properties
// Use also in Rendition for fxl
type Properties struct {
	NumberOfItems       int                   `json:"numberOfItems,omitempty"`
	Price               *Price                `json:"price,omitempty"`
	IndirectAcquisition []IndirectAcquisition `json:"indirectAcquisition,omitempty"`
}

// IndirectAcquisition store
type IndirectAcquisition struct {
	TypeAcquisition string                `json:"type"`
	Child           []IndirectAcquisition `json:"child,omitempty"`
}

// Price price information
type Price struct {
	Currency string  `json:"currency"`
	Value    float64 `json:"value"`
}

// MultiLanguage store a basic string when we only have one lang
// Store in a hash by language for multiple string representation
type MultiLanguage struct {
	SingleString string
	MultiString  map[string]string
}

// StringOrArray is a array of string redefine for overriding json
// marshalling and unmarshalling
type StringOrArray []string

// New create a new feed structure
func New(title string) Feed {
	var feed Feed

	feed.Metadata.Title = title
	t := time.Now()
	feed.Metadata.Modified = &t

	return feed
}

// MarshalJSON overwrite json marshalling for MultiLanguage
// when we have an entry in the Multi fields we use it
// otherwise we use the single string
func (m MultiLanguage) MarshalJSON() ([]byte, error) {
	if len(m.MultiString) > 0 {
		return json.Marshal(m.MultiString)
	}
	return json.Marshal(m.SingleString)
}

func (m MultiLanguage) String() string {
	if len(m.MultiString) > 0 {
		for _, s := range m.MultiString {
			return s
		}
	}
	return m.SingleString
}

// MarshalJSON overwrite json marshalling for handling string or array
func (r StringOrArray) MarshalJSON() ([]byte, error) {
	if len(r) == 1 {
		return json.Marshal(r[0])
	}
	return json.Marshal(r)
}
