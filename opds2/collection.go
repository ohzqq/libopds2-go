package opds2

import (
	"strings"
)

// BelongsTo is a list of collections/series that a publication belongs to
type BelongsTo struct {
	Series     Collections `json:"series,omitempty"`
	Collection Collections `json:"collection,omitempty"`
}

// Collections Slice
type Collections []*Collection

// Collection construct used for collection/serie metadata
type Collection struct {
	*Contributor
	Position float64 `json:"position,omitempty"`
}

func NewCollection(col any) Collections {
	return parseCollections(col)
}

func (c Collections) StringSlice() []string {
	var cols []string
	for _, col := range c {
		cols = append(cols, col.Name.SingleString)
	}
	return cols
}

func (c Collections) String() string {
	return strings.Join(c.StringSlice(), ", ")
}
