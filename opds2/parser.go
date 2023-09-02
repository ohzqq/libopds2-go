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

	f, err := os.Open(filePath)
	if err != nil {
		return &Feed{}, err
	}
	buff, errRead := io.ReadAll(f)
	if err != nil {
		return &Feed{}, errRead
	}

	feed, errParse := ParseBuffer(buff)
	if errParse != nil {
		return &Feed{}, errParse
	}

	return feed, nil
}

// ParseBuffer parse opds2 feed from a buffer of byte usually get
// from a file or url
func ParseBuffer(buff []byte) (*Feed, error) {
	var feed Feed

	errParse := json.Unmarshal(buff, &feed)

	if errParse != nil {
		fmt.Println(errParse)
	}

	return &feed, nil
}

// UnmarshalJSON make all unmarshalling by hand to handle all case
func (feed *Feed) UnmarshalJSON(data []byte) error {
	var info map[string]any

	json.Unmarshal(data, &info)

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
			parseMetadata(&feed.Metadata, v)
		case "links":
			parseLinks(feed, v)
		case "facets":
			parseFacets(feed, v)
		case "publications":
			parsePublications(feed, v)
		case "navigation":
			parseNavigation(feed, v)
		case "groups":
			parseGroups(feed, v)
		}
	}

	return nil
}

func parseMetadata(m *Metadata, data any) {
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
}

func parseLinks(feed *Feed, data any) {
	infoA := cast.ToSlice(data)
	for _, vA := range infoA {
		l := parseLink(vA)
		feed.Links = append(feed.Links, l)
	}
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

func parseFacets(feed *Feed, data any) {
	info := cast.ToSlice(data)
	f := Facet{}
	for _, fa := range info {
		infoA := cast.ToStringMap(fa)
		for k, v := range infoA {
			switch k {
			case "metadata":
				parseMetadata(&f.Metadata, v)
			case "links":
				infoAL := cast.ToSlice(v)
				for _, vA := range infoAL {
					l := parseLink(vA)
					f.Links = append(f.Links, l)
				}
			}
		}
		feed.Facets = append(feed.Facets, f)
	}
}

func parseGroups(feed *Feed, data any) {
	info := cast.ToSlice(data)
	for _, ga := range info {
		g := Group{}
		infoA := cast.ToStringMap(ga)
		for k, v := range infoA {
			switch k {
			case "metadata":
				parseMetadata(&g.Metadata, v)
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
					p := parsePublication(vP)
					g.Publications = append(g.Publications, p)
				}
			}
		}
		feed.Groups = append(feed.Groups, g)
	}
}

func parsePublications(feed *Feed, data any) {
	info := cast.ToSlice(data)
	for _, fa := range info {
		p := parsePublication(fa)
		feed.Publications = append(feed.Publications, p)
	}
}

func parsePublication(data any) Publication {
	var p Publication

	infoA := cast.ToStringMap(data)
	for k, v := range infoA {
		switch k {
		case "metadata":
			parsePublicationMetadata(&p.Metadata, v)
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

	return p
}

func parsePublicationMeta(data any) PublicationMetadata {
	metadata := PublicationMetadata{}
	parsePublicationMetadata(&metadata, data)
	return metadata
}

func parsePublicationMetadata(metadata *PublicationMetadata, data any) {
	info := cast.ToStringMap(data)
	for k, v := range info {
		switch k {
		case "title": // handle multistring
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
			metadata.Author = Author.New(v)
		case "translator":
			metadata.Translator = Translator.New(v)
		case "editor":
			metadata.Editor = Editor.New(v)
		case "artist":
			metadata.Artist = Artist.New(v)
		case "illustrator":
			metadata.Illustrator = Illustrator.New(v)
		case "letterer":
			metadata.Letterer = Letterer.New(v)
		case "penciler":
			metadata.Penciler = Penciler.New(v)
		case "colorist":
			metadata.Colorist = Colorist.New(v)
		case "inker":
			metadata.Inker = Inker.New(v)
		case "narrator":
			metadata.Narrator = Narrator.New(v)
		case "contributor":
			metadata.Contributor = parseCons(v)
		case "publisher":
			metadata.Publisher = Publisher.New(v)
		case "imprint":
			metadata.Imprint = Imprint.New(v)
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
			metadata.Subject = parseSubs(v)
		case "belongs_to", "belongsTo":
			belong := BelongsTo{}
			fmt.Printf("%T\n", v)
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

func parseSub(data any) *Subject {
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

func parseSubs(data any) Subjects {
	var cons Subjects
	switch d := data.(type) {
	case string:
		c := parseSub(d)
		cons = append(cons, c)
		return cons
	case map[string]any:
		c := parseSub(d)
		cons = append(cons, c)
		return cons
	case []any:
		for _, con := range d {
			cons = append(cons, parseSub(con))
		}
		return cons
	case []map[string]any:
		for _, con := range d {
			cons = append(cons, parseSub(con))
		}
		return cons
	}
	return cons
}

func parseCollection(data any) *Collection {
	collection := &Collection{
		Contributor: parseCon(data),
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
	lang := MultiLanguage{
		MultiString: make(map[string]string),
	}
	switch d := data.(type) {
	case string:
		lang.SingleString = d
		return lang
	case map[string]string:
		for k, v := range d {
			lang.MultiString[k] = v
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

func parseCon(data any) *Contributor {
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

func parseCons(data any) Contributors {
	var cons Contributors
	switch d := data.(type) {
	case string:
		c := parseCon(d)
		cons = append(cons, c)
	case map[string]any:
		c := parseCon(d)
		cons = append(cons, c)
	case []any:
		for _, con := range d {
			cons = append(cons, parseCon(con))
		}
	case []map[string]any:
		for _, con := range d {
			cons = append(cons, parseCon(con))
		}
	}
	return cons
}

func parseNavigation(feed *Feed, data any) {
	infoA := data.([]interface{})
	for _, vA := range infoA {
		l := parseLink(vA)
		feed.Navigation = append(feed.Navigation, l)
	}
}

// UnmarshalJSON overwrite json unmarshalling for Rel for handling
// when we have a array of a string
// func (r *StringOrArray) UnmarshalJSON(data []byte) error {
// 	var relAr []string
//
// 	if data[0] == '[' {
// 		err := json.Unmarshal(data, &relAr)
// 		if err != nil {
// 			return err
// 		}
// 		for _, ra := range relAr {
// 			*r = append(*r, ra)
// 		}
// 	} else {
// 		*r = append(*r, string(data))
// 	}
//
// 	return nil
// }

// UnmarshalJSON overwrite json unmarshalling for MultiLanguage
// when we have an entry in the Multi fields we use it
// otherwise we use the single string
// func (m *MultiLanguage) UnmarshalJSON(data []byte) error {
// 	var mParse map[string]string
//
// 	if data[0] == '{' {
// 		json.Unmarshal(data, &mParse)
// 		m.MultiString = mParse
// 	} else {
// 		m.SingleString = string(data)
// 	}
//
// 	return nil
// }
