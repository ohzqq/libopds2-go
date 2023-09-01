package opds2

import "strings"

type Role string

const (
	Author             Role = "author"
	Translator         Role = "translator"
	Editor             Role = "editor"
	Artist             Role = "artist"
	Illustrator        Role = "illustrator"
	Letterer           Role = "letterer"
	Penciler           Role = "penciler"
	Colorist           Role = "colorist"
	Inker              Role = "inker"
	Narrator           Role = "narrator"
	Publisher          Role = "publisher"
	Imprint            Role = "imprint"
	GenericContributor Role = "contributor"
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

func NewContributor(con any) Contributors {
	return parseCons(con)
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

func (r Role) New(con any) Contributors {
	c := parseCons(con)
	for i := range c {
		c[i].Role = string(r)
	}
	return c
}
