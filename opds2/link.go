package opds2

import "strings"

// Links used in collections and links
type Links []*Link

// Link object used in collections and links
type Link struct {
	Href       string        `json:"href"`
	TypeLink   string        `json:"type,omitempty"`
	Rel        StringOrArray `json:"rel,omitempty"`
	Height     int           `json:"height,omitempty"`
	Width      int           `json:"width,omitempty"`
	Title      string        `json:"title,omitempty"`
	Properties *Properties   `json:"properties,omitempty"`
	Duration   string        `json:"duration,omitempty"`
	Templated  bool          `json:"templated,omitempty"`
	Children   Links         `json:"children,omitempty"`
	Bitrate    int           `json:"bitrate,omitempty"`
}

func NewLink(data any) *Link {
	return parseLink(data)
}

func (l Links) StringSlice() []string {
	var links []string
	for _, link := range l {
		links = append(links, link.Href)
	}
	return links
}

func (l Links) String() string {
	return strings.Join(l.StringSlice(), ", ")
}
