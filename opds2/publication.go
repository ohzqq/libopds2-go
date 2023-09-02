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

func NewPublication(meta any) Publication {
	return parsePublication(meta)
}

func NewPublicationMetadata(title any) PublicationMetadata {
	return parsePublicationMeta(title)
}

// AddLink add a link to Publication
func (publication *Publication) AddLink(href any) *Link {
	i := NewLink(href)
	publication.Links = append(publication.Links, i)
	return i
}

// AddImage add a image link to Publication
func (publication *Publication) AddImage(href any) *Link {
	i := NewLink(href)
	publication.Images = append(publication.Images, i)
	return i
}

func (publication *Publication) BelongsToSeries(name any) *Collection {
	col := parseCollection(name)
	publication.Metadata.BelongsTo.Series = append(publication.Metadata.BelongsTo.Series, col)
	return col
}
