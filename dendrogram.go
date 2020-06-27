package hdbscan

import "log"

type node struct {
	key             int
	parentKey       int
	parent          *node
	children        []*node
	descendantCount int
}

type link struct {
	id              int
	parent          *link
	children        []*link
	points          []int
	descendantCount int
}

func (c *Clustering) buildDendrogram(baseEdges edges) []*link {
	if c.verbose {
		log.Println("starting dendrogram")
	}

	var links []*link
	for _, e := range baseEdges {
		var p1TopLink *link
		var p2TopLink *link

		for _, link := range links {
			if containsInt(link.points, e.p1) && link.parent == nil {
				p1TopLink = link
			}

			if containsInt(link.points, e.p2) && link.parent == nil {
				p2TopLink = link
			}
		}

		uniquePoints := make(map[int]bool)
		if p1TopLink != nil && p2TopLink != nil {
			for _, p := range p1TopLink.points {
				uniquePoints[p] = true
			}
			for _, p := range p2TopLink.points {
				uniquePoints[p] = true
			}
			var points []int
			for p, ok := range uniquePoints {
				if ok {
					points = append(points, p)
				}
			}

			newLink := link{
				id:       len(links),
				children: []*link{p1TopLink, p2TopLink},
				points:   points,
			}

			p1TopLink.parent = &newLink
			p2TopLink.parent = &newLink

			links = append(links, &newLink)
		} else if p1TopLink != nil && p2TopLink == nil {
			uniquePoints[e.p2] = true
			for _, p := range p1TopLink.points {
				uniquePoints[p] = true
			}
			var points []int
			for p, ok := range uniquePoints {
				if ok {
					points = append(points, p)
				}
			}

			newLink := link{
				id:       len(links),
				children: []*link{p1TopLink},
				points:   points,
			}

			p1TopLink.parent = &newLink

			links = append(links, &newLink)
		} else if p2TopLink != nil && p1TopLink == nil {
			uniquePoints[e.p1] = true
			for _, p := range p2TopLink.points {
				uniquePoints[p] = true
			}
			var points []int
			for p, ok := range uniquePoints {
				if ok {
					points = append(points, p)
				}
			}

			newLink := link{
				id:       len(links),
				children: []*link{p2TopLink},
				points:   points,
			}

			p2TopLink.parent = &newLink

			links = append(links, &newLink)
		} else {
			newLink := link{
				id:     len(links),
				points: []int{e.p1, e.p2},
			}

			links = append(links, &newLink)
		}
	}

	if c.verbose {
		log.Println("finished dendrogram")
	}

	return links
}
