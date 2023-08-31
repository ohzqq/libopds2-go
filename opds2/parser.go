package opds2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
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

	buff, errRead := ioutil.ReadAll(res.Body)
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
	buff, errRead := ioutil.ReadAll(f)
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
	var info map[string]interface{}

	json.Unmarshal(data, &info)

	for k, v := range info {
		switch k {
		case "@context":
			switch v.(type) {
			case string:
				feed.Context = append(feed.Context, v.(string))
			case []string:
				feed.Context = v.([]string)
			}
		case "metadata":
			ParseMetadata(&feed.Metadata, v)
		case "links":
			ParseLinks(feed, v)
		case "facets":
			ParseFacets(feed, v)
		case "publications":
			ParsePublications(feed, v)
		case "navigation":
			ParseNavigation(feed, v)
		case "groups":
			ParseGroups(feed, v)
		}
	}

	return nil
}

func ParseMetadata(m *Metadata, data interface{}) {

	info := data.(map[string]interface{})
	for k, v := range info {
		switch k {
		case "title":
			m.Title = v.(string)
		case "numberOfItems":
			m.NumberOfItems = int(v.(float64))
		case "itemsPerPage":
			m.ItemsPerPage = int(v.(float64))
		case "modified":
			t, err := time.Parse(time.RFC3339, v.(string))
			if err == nil {
				m.Modified = &t
			}
		case "type":
			m.RDFType = v.(string)
		case "currentPage":
			m.CurrentPage = int(v.(float64))
		}
	}
}

func ParseLinks(feed *Feed, data interface{}) {
	infoA := data.([]interface{})
	for _, vA := range infoA {
		l := ParseLink(vA)
		feed.Links = append(feed.Links, l)
	}
}

func ParseLink(data interface{}) Link {
	info := data.(map[string]interface{})
	l := Link{}
	for k, v := range info {
		switch k {
		case "title":
			l.Title = v.(string)
		case "href":
			l.Href = v.(string)
		case "type":
			l.TypeLink = v.(string)
		case "rel":
			switch v.(type) {
			case string:
				l.Rel = append(l.Rel, v.(string))
			case []string:
				l.Rel = v.([]string)
			}
		case "height":
			l.Height = int(v.(float64))
		case "width":
			l.Width = int(v.(float64))
		case "bitrate":
			l.Bitrate = int(v.(float64))
		case "duration":
			l.Duration = strconv.FormatFloat(v.(float64), 'f', -1, 64)
		case "templated":
			l.Templated = v.(bool)
		case "properties":
			p := Properties{}
			infoProp := v.(map[string]interface{})
			for kp, vp := range infoProp {
				switch kp {
				case "numberOfItems":
					p.NumberOfItems = int(vp.(float64))
				case "indirectAcquisition":
					infoIndir := vp.([]interface{})
					for _, in := range infoIndir {
						indir := ParseIndirectAcquisition(in)
						p.IndirectAcquisition = append(p.IndirectAcquisition, indir)
					}
				case "price":
					pr := Price{}
					infoPrice := vp.(map[string]interface{})
					for kpr, vpr := range infoPrice {
						switch kpr {
						case "currency":
							pr.Currency = vpr.(string)
						case "value":
							pr.Value = vpr.(float64)
						}
					}
					p.Price = &pr
				}
			}
			l.Properties = &p
		case "children":
			lc := ParseLink(v)
			l.Children = append(l.Children, lc)
		}
	}

	return l
}

func ParseIndirectAcquisition(data interface{}) IndirectAcquisition {
	var i IndirectAcquisition

	info := data.(map[string]interface{})
	for k, v := range info {
		switch k {
		case "type":
			i.TypeAcquisition = v.(string)
		case "child":
			infoA := v.([]interface{})
			for _, in := range infoA {
				indirect := ParseIndirectAcquisition(in)
				i.Child = append(i.Child, indirect)
			}
		}
	}

	return i
}

func ParseFacets(feed *Feed, data interface{}) {
	info := data.([]interface{})
	f := Facet{}
	for _, fa := range info {
		infoA := fa.(map[string]interface{})
		for k, v := range infoA {
			switch k {
			case "metadata":
				ParseMetadata(&f.Metadata, v)
			case "links":
				infoAL := v.([]interface{})
				for _, vA := range infoAL {
					l := ParseLink(vA)
					f.Links = append(f.Links, l)
				}
			}
		}
		feed.Facets = append(feed.Facets, f)
	}
}

func ParseGroups(feed *Feed, data interface{}) {
	info := data.([]interface{})
	for _, ga := range info {
		g := Group{}
		infoA := ga.(map[string]interface{})
		for k, v := range infoA {
			switch k {
			case "metadata":
				ParseMetadata(&g.Metadata, v)
			case "links":
				infoAL := v.([]interface{})
				for _, vA := range infoAL {
					l := ParseLink(vA)
					g.Links = append(g.Links, l)
				}
			case "navigation":
				infoAN := v.([]interface{})
				for _, vAN := range infoAN {
					l := ParseLink(vAN)
					g.Navigation = append(g.Navigation, l)
				}
			case "publications":
				infoP := v.([]interface{})
				for _, vP := range infoP {
					p := ParsePublication(vP)
					g.Publications = append(g.Publications, p)
				}
			}
		}
		feed.Groups = append(feed.Groups, g)
	}
}

func ParsePublications(feed *Feed, data interface{}) {
	info := data.([]interface{})
	for _, fa := range info {
		p := ParsePublication(fa)
		feed.Publications = append(feed.Publications, p)
	}
}

func ParsePublication(data interface{}) Publication {
	var p Publication

	infoA := data.(map[string]interface{})
	for k, v := range infoA {
		switch k {
		case "metadata":
			ParsePublicationMetadata(&p.Metadata, v)
		case "links":
			infoAL := v.([]interface{})
			for _, vA := range infoAL {
				l := ParseLink(vA)
				p.Links = append(p.Links, l)
			}
		case "images":
			infoAL := v.([]interface{})
			for _, vA := range infoAL {
				l := ParseLink(vA)
				p.Images = append(p.Images, l)
			}
		}
	}

	return p
}

func ParsePublicationMetadata(metadata *PublicationMetadata, data interface{}) {
	info := data.(map[string]interface{})
	for k, v := range info {
		switch k {
		case "title": // handle multistring
			metadata.Title.SingleString = v.(string)
		case "identifier":
			metadata.Identifier = v.(string)
		case "@type":
			metadata.RDFType = v.(string)
		case "modified":
			t, err := time.Parse(time.RFC3339, v.(string))
			if err == nil {
				metadata.Modified = &t
			}
		case "type":
			metadata.RDFType = v.(string)
		case "author":
			c := ParseContributors(v)
			for _, cont := range c {
				metadata.Author = append(metadata.Author, cont)
			}
		case "translator":
			c := ParseContributors(v)
			for _, cont := range c {
				metadata.Translator = append(metadata.Translator, cont)
			}
		case "editor":
			c := ParseContributors(v)
			for _, cont := range c {
				metadata.Editor = append(metadata.Editor, cont)
			}
		case "artist":
			c := ParseContributors(v)
			for _, cont := range c {
				metadata.Artist = append(metadata.Artist, cont)
			}
		case "illustrator":
			c := ParseContributors(v)
			for _, cont := range c {
				metadata.Illustrator = append(metadata.Illustrator, cont)
			}
		case "letterer":
			c := ParseContributors(v)
			for _, cont := range c {
				metadata.Letterer = append(metadata.Letterer, cont)
			}
		case "penciler":
			c := ParseContributors(v)
			for _, cont := range c {
				metadata.Penciler = append(metadata.Penciler, cont)
			}
		case "colorist":
			c := ParseContributors(v)
			for _, cont := range c {
				metadata.Colorist = append(metadata.Colorist, cont)
			}
		case "inker":
			c := ParseContributors(v)
			for _, cont := range c {
				metadata.Inker = append(metadata.Inker, cont)
			}
		case "narrator":
			c := ParseContributors(v)
			for _, cont := range c {
				metadata.Narrator = append(metadata.Narrator, cont)
			}
		case "contributor":
			c := ParseContributors(v)
			for _, cont := range c {
				metadata.Contributor = append(metadata.Contributor, cont)
			}
		case "publisher":
			c := ParseContributors(v)
			for _, cont := range c {
				metadata.Publisher = append(metadata.Publisher, cont)
			}
		case "imprint":
			c := ParseContributors(v)
			for _, cont := range c {
				metadata.Imprint = append(metadata.Imprint, cont)
			}
		case "language":
		case "published":
			t, err := time.Parse(time.RFC3339, v.(string))
			if err == nil {
				metadata.PublicationDate = &t
			}
		case "description":
			metadata.Description = v.(string)
		case "source":
			metadata.Source = v.(string)
		case "rights":
			metadata.Rights = v.(string)
		case "subject":
			metadata.Subject = ParseSubject(v)
		case "belongs_to":
			belong := BelongsTo{}
			infoB := v.(map[string]interface{})
			for kb, vb := range infoB {
				switch kb {
				case "series":
					switch vb.(type) {
					case string:
						belong.Series = append(belong.Series, Collection{Name: vb.(string)})
					case []interface{}:
						for _, colls := range vb.([]interface{}) {
							coll := ParseCollection(colls)
							belong.Series = append(belong.Series, coll)
						}
					case interface{}:
						coll := ParseCollection(vb)
						belong.Series = append(belong.Series, coll)
					}
				case "collection":
					switch vb.(type) {
					case string:
						belong.Collection = append(belong.Collection, Collection{Name: vb.(string)})
					case []interface{}:
						for _, colls := range vb.([]interface{}) {
							coll := ParseCollection(colls)
							belong.Collection = append(belong.Collection, coll)
						}
					case interface{}:
						coll := ParseCollection(vb)
						belong.Collection = append(belong.Collection, coll)
					}
				}
			}
			metadata.BelongsTo = &belong
		case "duration":
			metadata.Duration = int(v.(float64))
		}
	}
}

func ParseSubject(v any) []Subject {
	var subs []Subject
	switch data := v.(type) {
	case string:
		s := Subject{}
		s.Name = data
		subs = append(subs, s)
	case []any:
		for _, subject := range data {
			s := Subject{}
			switch sub := subject.(type) {
			case string:
				s.Name = sub
			case map[string]any:
				for ks, vs := range sub {
					switch ks {
					case "name":
						s.Name = vs.(string)
					case "sort_as":
						s.SortAs = vs.(string)
					case "scheme":
						s.Scheme = vs.(string)
					case "code":
						s.Code = vs.(string)
					}
				}
			}
			subs = append(subs, s)
		}
	}
	return subs
}

func ParseCollection(data interface{}) Collection {
	var collection Collection

	info := data.(map[string]interface{})
	for k, v := range info {
		switch k {
		case "name":
			collection.Name = v.(string)
		case "sort_as":
			collection.SortAs = v.(string)
		case "identifier":
			collection.Identifier = v.(string)
		case "position":
			collection.Position = float32(v.(float64))
		case "links":
			infoL := v.([]interface{})
			for _, l := range infoL {
				link := ParseLink(l)
				collection.Links = append(collection.Links, link)
			}
		}
	}

	return collection
}

func ParseContributors(data interface{}) []Contributor {
	var c []Contributor

	switch d := data.(type) {
	case string:
		cont := Contributor{}
		cont.Name.SingleString = d
		c = append(c, cont)
	case []string:
		for _, i := range d {
			cont := Contributor{}
			cont.Name.SingleString = d
			c = append(c, cont)
		}
	case []interface{}:
		for _, i := range d {
			cont := ParseContributor(i)
			c = append(c, cont)
		}
	case interface{}:
		cont := ParseContributor(d)
		c = append(c, cont)
	}
	return c
}

func ParseContributor(data interface{}) Contributor {
	var c Contributor

	info := data.(map[string]interface{})
	for k, v := range info {
		switch k {
		case "name":
			switch v.(type) {
			case string:
				c.Name.SingleString = v.(string)
			case map[string]interface{}:
				infoN := v.(map[string]interface{})
				c.Name.MultiString = make(map[string]string)
				for kn, vn := range infoN {
					c.Name.MultiString[kn] = vn.(string)
				}
			}
		case "identifier":
			c.Identifier = v.(string)
		case "sort_as":
			c.SortAs = v.(string)
		case "role":
			c.Role = v.(string)
		case "links":
			l := ParseLink(v)
			c.Links = append(c.Links, l)
		}
	}

	return c
}

func ParseNavigation(feed *Feed, data interface{}) {
	infoA := data.([]interface{})
	for _, vA := range infoA {
		l := ParseLink(vA)
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
