package opds2

import (
	"time"
)

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
	Language        StringOrArray `json:"language,omitempty"`
	Modified        *time.Time    `json:"modified,omitempty"`
	PublicationDate *time.Time    `json:"published,omitempty"`
	Description     string        `json:"description,omitempty"`
	Source          string        `json:"source,omitempty"`
	Rights          string        `json:"rights,omitempty"`
	Subject         Subjects      `json:"subject,omitempty"`
	BelongsTo       *BelongsTo    `json:"belongsTo,omitempty"`
	Duration        int           `json:"duration,omitempty"`
}

func NewPublication(meta any, links ...*Link) Publication {
	pub := Publication{}
	for _, l := range links {
		pub.Links = append(pub.Links, l)
	}

	if d, ok := meta.(string); ok {
		parsePublicationMetadata(d, &pub.Metadata)
		return pub
	}

	parsePublication(meta, &pub)
	return pub
}

func NewPublicationMetadata(data any) PublicationMetadata {
	if d, ok := data.(string); ok {
		return PublicationMetadata{
			Title: parseMultiLanguage(d),
		}
	}
	var m PublicationMetadata
	parsePublicationMetadata(data, &m)
	return m
}

// AddLink add a link to Publication
func (publication *Publication) AddLink(data any) *Link {
	i := NewLink(data)
	publication.Links = append(publication.Links, i)
	return i
}

// AddImage add a image link to Publication
func (publication *Publication) AddImage(data any) *Link {
	i := NewLink(data)
	publication.Images = append(publication.Images, i)
	return i
}

func (publication *Publication) BelongsToSeries(data any) *Collection {
	col := parseCollection(data)
	publication.Metadata.BelongsTo.Series = append(publication.Metadata.BelongsTo.Series, col)
	return col
}

func (publication *Publication) BelongsToCollection(data any) *Collection {
	col := parseCollection(data)
	publication.Metadata.BelongsTo.Collection = append(publication.Metadata.BelongsTo.Collection, col)
	return col
}

func (publication *Publication) FindFirstImageByRel(rel string) *Link {
	return publication.Images.FindFirstLinkByRel(rel)
}

func (publication *Publication) FindFirstLinkByRel(rel string) *Link {
	return publication.Links.FindFirstLinkByRel(rel)
}

func (publication *Publication) FindFirstLinkByType(mt string) *Link {
	return publication.Links.FindFirstLinkByType(mt)
}
