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

func (links Links) FindFirstLinkByRel(rel string) *Link {
	for _, l := range links {
		for _, r := range l.Rel {
			if r == rel {
				return l
			}
		}
	}
	return &Link{}
}

func (links Links) FindFirstLinkByType(mt string) *Link {
	for _, l := range links {
		if strings.Contains(l.TypeLink, mt) {
			return l
		}
	}
	return &Link{}
}
