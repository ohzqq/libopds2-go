package opds2

import "strings"

// Subject Slice
type Subjects []*Subject

// Subject as based on EPUB 3.1 and WebPub
type Subject struct {
	Name   string `json:"name"`
	SortAs string `json:"sort_as,omitempty"`
	Scheme string `json:"scheme,omitempty"`
	Code   string `json:"code,omitempty"`
}

func NewSubject(con any) Subjects {
	return parseSubs(con)
}

func (s Subjects) StringSlice() []string {
	var subs []string
	for _, sub := range s {
		subs = append(subs, sub.Name)
	}
	return subs
}

func (s Subjects) String() string {
	return strings.Join(s.StringSlice(), ", ")
}
