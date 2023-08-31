package opds2

import "time"

// Publication is a collection for a given publication
type Publication struct {
	Metadata PublicationMetadata `json:"metadata"`
	Links    Links               `json:"links"`
	Images   Links               `json:"images"`
}

// PublicationMetadata for the default context in WebPub
type PublicationMetadata struct {
	RDFType         string        `json:"@type,omitempty"` //Defaults to schema.org for EBook
	Title           MultiLanguage `json:"title"`
	Identifier      string        `json:"identifier"`
	Author          Contributors  `json:"author,omitempty"`
	Translator      Contributors  `json:"translator,omitempty"`
	Editor          Contributors  `json:"editor,omitempty"`
	Artist          Contributors  `json:"artist,omitempty"`
	Illustrator     Contributors  `json:"illustrator,omitempty"`
	Letterer        Contributors  `json:"letterer,omitempty"`
	Penciler        Contributors  `json:"penciler,omitempty"`
	Colorist        Contributors  `json:"colorist,omitempty"`
	Inker           Contributors  `json:"inker,omitempty"`
	Narrator        Contributors  `json:"narrator,omitempty"`
	Contributor     Contributors  `json:"contributor,omitempty"`
	Publisher       Contributors  `json:"publisher,omitempty"`
	Imprint         Contributors  `json:"imprint,omitempty"`
	Language        []string      `json:"language,omitempty"`
	Modified        *time.Time    `json:"modified,omitempty"`
	PublicationDate *time.Time    `json:"published,omitempty"`
	Description     string        `json:"description,omitempty"`
	Source          string        `json:"source,omitempty"`
	Rights          string        `json:"rights,omitempty"`
	Subject         Subjects      `json:"subject,omitempty"`
	BelongsTo       *BelongsTo    `json:"belongs_to,omitempty"`
	Duration        int           `json:"duration,omitempty"`
}

func NewPublication() Publication {
	return Publication{
		Metadata: NewPublicationMetadata(),
	}
}

func NewPublicationMetadata() PublicationMetadata {
	return PublicationMetadata{
		BelongsTo: &BelongsTo{},
	}
}

// AddImage add a image link to Publication
func (publication *Publication) AddImage(href string) *Link {
	i := NewLink(href)
	publication.Images = append(publication.Images, i)
	return &i
}

// AddLink add a new link to the publication
func (publication *Publication) AddLink(href string) *Link {
	l := NewLink(href, rel)
	publication.Links = append(publication.Links, l)
	return &l
}

// AddContributor adds a contributor to publication with all parameters mostly optional
func (publication *Publication) AddAuthor(name string, identifier string, sortAs string, href string, typeLink string) {
	var c Contributor
	l := NewLink(href)

	c.Name.SingleString = name
	if identifier != "" {
		c.Identifier = identifier
	}
	if sortAs != "" {
		c.SortAs = sortAs
	}
	if typeLink != "" {
		l.TypeLink = typeLink
	}

	if l.Href != "" {
		c.Links = append(c.Links, l)
	}

	publication.Metadata.Author = append(publication.Metadata.Author, c)
}

// AddAuthor add author to publication with all parameters mostly optional
func (publication *Publication) AddAuthor(name string, identifier string, sortAs string, href string, typeLink string) {
	var c Contributor
	l := NewLink(href)

	c.Name.SingleString = name
	if identifier != "" {
		c.Identifier = identifier
	}
	if sortAs != "" {
		c.SortAs = sortAs
	}
	if typeLink != "" {
		l.TypeLink = typeLink
	}

	if l.Href != "" {
		c.Links = append(c.Links, l)
	}

	publication.Metadata.Author = append(publication.Metadata.Author, c)
}

// AddSerie add serie to publication
func (publication *Publication) AddSerie(name string, position float32, href string, typeLink string) {
	var c Collection
	l := NewLink(href)

	c.Name = name
	c.Position = position

	if publication.Metadata.BelongsTo == nil {
		publication.Metadata.BelongsTo = &BelongsTo{}
	}

	if typeLink != "" {
		l.TypeLink = typeLink
	}

	if l.Href != "" {
		c.Links = append(c.Links, l)
	}

	publication.Metadata.BelongsTo.Series = append(publication.Metadata.BelongsTo.Series, c)
}

// AddPublisher add publisher to publication
func (publication *Publication) AddPublisher(name string, href string, typeLink string) {
	var c Contributor
	l := NewLink(href)

	c.Name.SingleString = name

	if typeLink != "" {
		l.TypeLink = typeLink
	}

	if l.Href != "" {
		c.Links = append(c.Links, l)
	}

	publication.Metadata.Publisher = append(publication.Metadata.Publisher, c)
}
