package opds2

// AddLink add a new link in feed information
// at minimum the self link
func (feed *Feed) AddLink(href string, rel string, typeLink string, templated bool) {
	l := NewLink(href)
	l.TypeLink = typeLink
	if templated == true {
		l.Templated = true
	}

	feed.Links = append(feed.Links, l)
}

// AddNavigation add navigation element in feed
func (feed *Feed) AddNavigation(title string, href string, rel string, typeLink string) {
	l := NewLink(href)

	l.TypeLink = typeLink
	if title != "" {
		l.Title = title
	}

	feed.Navigation = append(feed.Navigation, l)
}

// AddPagination add pagination and link information in feed
func (feed *Feed) AddPagination(numberItems int, itemsPerPage int, currentPage int, nextLink string, prevLink string, firstLink string, lastLink string) {

	feed.Metadata.CurrentPage = currentPage
	feed.Metadata.ItemsPerPage = itemsPerPage
	feed.Metadata.NumberOfItems = numberItems

	if nextLink != "" {
		feed.AddLink(nextLink, "next", "application/opds+json", false)
	}
	if prevLink != "" {
		feed.AddLink(prevLink, "previous", "application/opds+json", false)
	}
	if firstLink != "" {
		feed.AddLink(firstLink, "first", "application/opds+json", false)
	}
	if lastLink != "" {
		feed.AddLink(lastLink, "last", "application/opds+json", false)
	}
}

// AddFacet add link to facet handler multiple add
func (feed *Feed) AddFacet(link *Link, group string) {
	var facet Facet

	for i, f := range feed.Facets {
		if f.Metadata.Title == group {
			feed.Facets[i].Links = append(feed.Facets[i].Links, link)
			return
		}
	}

	facet.Metadata.Title = group
	facet.Links = append(facet.Links, link)
	feed.Facets = append(feed.Facets, facet)
}

// AddPublicationInGroup smart adding of publication in Group
func (feed *Feed) AddPublicationInGroup(publication Publication, collLink *Link) {
	var group Group

	for i, g := range feed.Groups {
		for _, l := range g.Links {
			if l.Href == collLink.Href {
				feed.Groups[i].Publications = append(feed.Groups[i].Publications, publication)
				return
			}
		}
	}

	group.Metadata.Title = collLink.Title
	group.Publications = append(group.Publications, publication)
	group.Links = append(group.Links, &Link{Rel: []string{"self"}, Title: collLink.Title, Href: collLink.Href})
	feed.Groups = append(feed.Groups, group)
}

// AddNavigationInGroup add a navigation link to Group
func (feed *Feed) AddNavigationInGroup(link *Link, collLink *Link) {
	var group Group

	for i, g := range feed.Groups {
		for _, l := range g.Links {
			if l.Href == collLink.Href {
				feed.Groups[i].Navigation = append(feed.Groups[i].Navigation, link)
				return
			}
		}
	}

	group.Metadata.Title = collLink.Title
	group.Navigation = append(group.Navigation, link)
	group.Links = append(group.Links, &Link{Rel: []string{"self"}, Title: collLink.Title, Href: collLink.Href})
	feed.Groups = append(feed.Groups, group)
}
