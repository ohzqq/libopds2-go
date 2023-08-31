package opds2

import "strings"

type Role string

const (
	Author             Role = "author"
	Translator              = "translator"
	Editor                  = "editor"
	Artist                  = "artist"
	Illustrator             = "illustrator"
	Letterer                = "letterer"
	Penciler                = "penciler"
	Colorist                = "colorist"
	Inker                   = "inker"
	Narrator                = "narrator"
	Publisher               = "publisher"
	Imprint                 = "imprint"
	GenericContributor      = "contributor"
)

// Contributor Slice
type Contributors []*Contributor

// Contributor construct used internally for all contributors
type Contributor struct {
	Name       MultiLanguage `json:"name,omitempty"`
	SortAs     string        `json:"sort_as,omitempty"`
	Identifier string        `json:"identifier,omitempty"`
	Role       string        `json:"role,omitempty"`
	Links      Links         `json:"links,omitempty"`
}

func NewContributor(name string) *Contributor {
	return &Contributor{
		Name: MultiLanguage{
			SingleString: name,
		},
	}
}

func (c Contributors) StringSlice() []string {
	var cons []string
	for _, con := range c {
		cons = append(cons, con.Name.String())
	}
	return cons
}

func (c Contributors) String() string {
	return strings.Join(c.StringSlice(), " & ")
}

func (r Role) New(name string) *Contributor {
	c := NewContributor(name)
	c.Role = string(r)
	return c
}
