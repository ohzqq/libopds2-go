package opds2

import "strings"

//go:generate stringer -type Role -linecomment
type Role int

const (
	Author      Role = iota + 1 // author
	Translator                  // translator
	Editor                      // editor
	Artist                      // artist
	Illustrator                 // illustrator
	Letterer                    // letterer
	Penciler                    // penciler
	Colorist                    // colorist
	Inker                       // inker
	Narrator                    // narrator
	Publisher                   // publisher
	Imprint                     // imprint
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
	return parseContributors(con)
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

func (role Role) New(con any) Contributors {
	c := parseContributors(con)
	for i := range c {
		if c[i].Role == "" {
			c[i].Role = role.String()
		}
	}
	return c
}
