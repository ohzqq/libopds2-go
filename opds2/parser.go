package opds2

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cast"
)

// ParseURL parse the opds2 feed from an url
func ParseURL(url string) (*Feed, error) {

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, errReq := http.DefaultClient.Do(request)
	if errReq != nil {
		return nil, errReq
	}

	buff, errRead := io.ReadAll(res.Body)
	if errRead != nil {
		return nil, errRead
	}

	feed, errParse := ParseBuffer(buff)
	if errParse != nil {
		return &Feed{}, errParse
	}

	return feed, nil
}

// ParseFile parse opds2 from a file on filesystem
func ParseFile(filePath string) (*Feed, error) {

	f, err := os.ReadFile(filePath)
	if err != nil {
		return &Feed{}, err
	}

	feed, errParse := ParseBuffer(f)
	if errParse != nil {
		return &Feed{}, errParse
	}

	return feed, nil
}

// ParseBuffer parse opds2 feed from a buffer of byte usually get
// from a file or url
func ParseBuffer(buff []byte) (*Feed, error) {
	feed := &Feed{}

	errParse := json.Unmarshal(buff, feed)

	if errParse != nil {
		fmt.Println(errParse)
		return feed, errParse
	}

	return feed, nil
}

// UnmarshalJSON make all unmarshalling by hand to handle all case
func (feed *Feed) UnmarshalJSON(data []byte) error {
	info, castErr := cast.ToStringMapE(string(data))
	if castErr != nil {
		return castErr
	}

	for k, v := range info {
		switch k {
		case "@context":
			switch v.(type) {
			case string:
				feed.Context = append(feed.Context, cast.ToString(v))
			case []string:
				feed.Context = cast.ToStringSlice(v)
			}
		case "metadata":
			feed.Metadata = parseMetadata(v)
		case "links":
			feed.Links = parseLinks(v)
		case "facets":
			feed.Facets = parseFacets(v)
		case "publications":
			feed.Publications = parsePublications(v)
		case "navigation":
			feed.Navigation = parseLinks(v)
		case "groups":
			feed.Groups = parseGroups(v)
		}
	}

	return nil
}

func parseMetadata(data any) Metadata {
	m := Metadata{}
	info := cast.ToStringMap(data)
	for k, v := range info {
		switch k {
		case "title":
			m.Title = cast.ToString(v)
		case "numberOfItems":
			m.NumberOfItems = cast.ToInt(v)
		case "itemsPerPage":
			m.ItemsPerPage = cast.ToInt(v)
		case "modified":
			m.Modified = parseDate(v)
		case "type":
			m.RDFType = cast.ToString(v)
		case "currentPage":
			m.CurrentPage = cast.ToInt(v)
		}
	}
	return m
}

func parseLinks(data any) Links {
	var links Links
	infoA := cast.ToSlice(data)
	for _, vA := range infoA {
		l := parseLink(vA)
		links = append(links, l)
	}
	return links
}

func parseLink(data any) *Link {
	info := cast.ToStringMap(data)
	l := Link{}
	for k, v := range info {
		switch k {
		case "title":
			l.Title = cast.ToString(v)
		case "href":
			l.Href = cast.ToString(v)
		case "type":
			l.TypeLink = cast.ToString(v)
		case "rel":
			switch v.(type) {
			case string:
				l.Rel = append(l.Rel, cast.ToString(v))
			case []string:
				l.Rel = cast.ToStringSlice(v)
			}
		case "height":
			l.Height = cast.ToInt(v)
		case "width":
			l.Width = cast.ToInt(v)
		case "bitrate":
			l.Bitrate = cast.ToInt(v)
		case "duration":
			l.Duration = cast.ToString(v)
		case "templated":
			l.Templated = cast.ToBool(v)
		case "properties":
			p := Properties{}
			infoProp := cast.ToStringMap(v)
			for kp, vp := range infoProp {
				switch kp {
				case "numberOfItems":
					p.NumberOfItems = cast.ToInt(vp)
				case "indirectAcquisition":
					infoIndir := cast.ToSlice(vp)
					for _, in := range infoIndir {
						indir := parseIndirectAcquisition(in)
						p.IndirectAcquisition = append(p.IndirectAcquisition, indir)
					}
				case "price":
					pr := Price{}
					infoPrice := cast.ToStringMap(vp)
					for kpr, vpr := range infoPrice {
						switch kpr {
						case "currency":
							pr.Currency = cast.ToString(vpr)
						case "value":
							pr.Value = cast.ToFloat64(vpr)
						}
					}
					p.Price = &pr
				}
			}
			l.Properties = &p
		case "children":
			lc := parseLink(v)
			l.Children = append(l.Children, lc)
		}
	}

	return &l
}

func parseIndirectAcquisition(data any) IndirectAcquisition {
	var i IndirectAcquisition

	info := cast.ToStringMap(data)
	for k, v := range info {
		switch k {
		case "type":
			i.TypeAcquisition = cast.ToString(v)
		case "child":
			infoA := cast.ToSlice(v)
			for _, in := range infoA {
				indirect := parseIndirectAcquisition(in)
				i.Child = append(i.Child, indirect)
			}
		}
	}

	return i
}

func parseFacets(data any) []Facet {
	var facets []Facet
	info := cast.ToSlice(data)
	f := Facet{}
	for _, fa := range info {
		infoA := cast.ToStringMap(fa)
		for k, v := range infoA {
			switch k {
			case "metadata":
				f.Metadata = parseMetadata(v)
			case "links":
				infoAL := cast.ToSlice(v)
				for _, vA := range infoAL {
					l := parseLink(vA)
					f.Links = append(f.Links, l)
				}
			}
		}
		facets = append(facets, f)
	}
	return facets
}

func parseGroups(data any) []Group {
	var groups []Group
	info := cast.ToSlice(data)
	for _, ga := range info {
		g := Group{}
		infoA := cast.ToStringMap(ga)
		for k, v := range infoA {
			switch k {
			case "metadata":
				g.Metadata = parseMetadata(v)
			case "links":
				infoAL := cast.ToSlice(v)
				for _, vA := range infoAL {
					l := parseLink(vA)
					g.Links = append(g.Links, l)
				}
			case "navigation":
				infoAN := cast.ToSlice(v)
				for _, vAN := range infoAN {
					l := parseLink(vAN)
					g.Navigation = append(g.Navigation, l)
				}
			case "publications":
				infoP := cast.ToSlice(v)
				for _, vP := range infoP {
					p := Publication{}
					parsePublication(vP, &p)
					g.Publications = append(g.Publications, p)
				}
			}
		}
		groups = append(groups, g)
	}
	return groups
}

func parsePublications(data any) []Publication {
	var pubs []Publication
	info := cast.ToSlice(data)
	for _, fa := range info {
		p := Publication{}
		parsePublication(fa, &p)
		pubs = append(pubs, p)
	}
	return pubs
}

func parsePublication(data any, p *Publication) {
	if d, ok := data.([]byte); ok {
		data = string(d)
	}
	infoA := cast.ToStringMap(data)
	for k, v := range infoA {
		switch k {
		case "metadata":
			parsePublicationMetadata(v, &p.Metadata)
		case "links":
			infoAL := cast.ToSlice(v)
			for _, vA := range infoAL {
				l := parseLink(vA)
				p.Links = append(p.Links, l)
			}
		case "images":
			infoAL := cast.ToSlice(v)
			for _, vA := range infoAL {
				l := parseLink(vA)
				p.Images = append(p.Images, l)
			}
		}
	}
}

func parsePublicationMetadata(data any, metadata *PublicationMetadata) {
	if d, ok := data.([]byte); ok {
		data = string(d)
	}
	info := cast.ToStringMap(data)
	for k, v := range info {
		switch k {
		case "title":
			metadata.Title = parseMultiLanguage(v)
		case "identifier":
			metadata.Identifier = cast.ToString(v)
		case "@type":
			metadata.RDFType = cast.ToString(v)
		case "modified":
			metadata.Modified = parseDate(v)
		case "type":
			metadata.RDFType = cast.ToString(v)
		case "author":
			metadata.Author = parseContributors(v)
		case "translator":
			metadata.Translator = parseContributors(v)
		case "editor":
			metadata.Editor = parseContributors(v)
		case "artist":
			metadata.Artist = parseContributors(v)
		case "illustrator":
			metadata.Illustrator = parseContributors(v)
		case "letterer":
			metadata.Letterer = parseContributors(v)
		case "penciler":
			metadata.Penciler = parseContributors(v)
		case "colorist":
			metadata.Colorist = parseContributors(v)
		case "inker":
			metadata.Inker = parseContributors(v)
		case "narrator":
			metadata.Narrator = parseContributors(v)
		case "contributor":
			metadata.Contributor = parseContributors(v)
		case "publisher":
			metadata.Publisher = parseContributors(v)
		case "imprint":
			metadata.Imprint = parseContributors(v)
		case "language":
			switch vb := v.(type) {
			case string:
				metadata.Language = append(metadata.Language, vb)
			case []any:
				for _, colls := range cast.ToStringSlice(vb) {
					metadata.Language = append(metadata.Language, colls)
				}
			}
		case "published":
			metadata.PublicationDate = parseDate(v)
		case "description":
			metadata.Description = cast.ToString(v)
		case "source":
			metadata.Source = cast.ToString(v)
		case "rights":
			metadata.Rights = cast.ToString(v)
		case "subject":
			metadata.Subject = parseSubjects(v)
		case "belongs_to", "belongsTo":
			belong := BelongsTo{}
			infoB := cast.ToStringMap(v)
			for kb, vb := range infoB {
				switch kb {
				case "series":
					belong.Series = parseCollections(vb)
				case "collection":
					belong.Collection = parseCollections(vb)
				}
			}
			metadata.BelongsTo = &belong
		case "duration":
			metadata.Duration = cast.ToInt(v)
		}
	}
}

func parseSubject(data any) *Subject {
	c := &Subject{}
	switch d := data.(type) {
	case string:
		c.Name = d
		return c
	case map[string]any:
		for ks, vs := range d {
			switch ks {
			case "name":
				c.Name = cast.ToString(vs)
			case "sort_as":
				c.SortAs = cast.ToString(vs)
			case "scheme":
				c.Scheme = cast.ToString(vs)
			case "code":
				c.Code = cast.ToString(vs)
			}
		}
	}
	return c
}

func parseSubjects(data any) Subjects {
	var cons Subjects
	switch d := data.(type) {
	case string:
		c := parseSubject(d)
		cons = append(cons, c)
		return cons
	case map[string]any:
		c := parseSubject(d)
		cons = append(cons, c)
		return cons
	case []any:
		for _, con := range d {
			cons = append(cons, parseSubject(con))
		}
		return cons
	case []map[string]any:
		for _, con := range d {
			cons = append(cons, parseSubject(con))
		}
		return cons
	}
	return cons
}

func parseCollection(data any) *Collection {
	collection := &Collection{
		Contributor: parseContributor(data),
	}

	info := cast.ToStringMap(data)
	if pos, ok := info["position"]; ok {
		collection.Position = cast.ToFloat64(pos)
	}
	return collection
}

func parseCollections(data any) Collections {
	var cons Collections
	switch d := data.(type) {
	case string:
		c := parseCollection(d)
		cons = append(cons, c)
		return cons
	case map[string]any:
		c := parseCollection(d)
		cons = append(cons, c)
		return cons
	case []any:
		for _, con := range d {
			cons = append(cons, parseCollection(con))
		}
		return cons
	case []map[string]any:
		for _, con := range d {
			cons = append(cons, parseCollection(con))
		}
		return cons
	}
	return cons
}

func parseMultiLanguage(data any) MultiLanguage {
	lang := MultiLanguage{}
	switch d := data.(type) {
	case string:
		lang.SingleString = d
		return lang
	case map[string]any:
		lang.MultiString = make(map[string]string)
		for k, v := range d {
			lang.MultiString[k] = cast.ToString(v)
		}
	}
	return lang
}

func parseDate(data any) *time.Time {
	t, err := time.Parse(time.RFC3339, cast.ToString(data))
	if err == nil {
		t = time.Now()
	}
	return &t
}

func parseContributor(data any) *Contributor {
	switch d := data.(type) {
	case string:
		c := &Contributor{}
		c.Name = parseMultiLanguage(d)
		return c
	case map[string]any:
		c := &Contributor{}
		for k, v := range d {
			switch k {
			case "name":
				c.Name = parseMultiLanguage(v)
			case "identifier":
				c.Identifier = cast.ToString(v)
			case "sort_as":
				c.SortAs = cast.ToString(v)
			case "role":
				c.Role = cast.ToString(v)
			case "links":
				l := parseLink(v)
				c.Links = append(c.Links, l)
			}
		}
		return c
	}
	return &Contributor{}
}

func parseContributors(data any) Contributors {
	var cons Contributors
	switch d := data.(type) {
	case string:
		c := parseContributor(d)
		cons = append(cons, c)
	case map[string]any:
		c := parseContributor(d)
		cons = append(cons, c)
	case []any:
		for _, con := range d {
			cons = append(cons, parseContributor(con))
		}
	case []map[string]any:
		for _, con := range d {
			cons = append(cons, parseContributor(con))
		}
	}
	return cons
}
