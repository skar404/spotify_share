package spotify

import (
	"errors"
	spotify_type "github.com/skar404/spotify_share/spotify/type"
)

type Artists struct {
	ID   string
	Name string
	Href string
}

type Album struct {
	ID   string
	Name string
}

type History struct {
	ID         string
	URL        string
	Name       string
	PreviewURL string
	Img        string
	Album      Album
	Artists    []Artists
}

func makeItem(item *spotify_type.Item) History {
	h := History{
		ID:         item.ID,
		URL:        item.URI,
		Name:       item.Name,
		PreviewURL: item.PreviewURL,
		Img:        item.Album.Images[0].URL,
		Album: Album{
			ID:   item.Album.ID,
			Name: item.Album.Name,
		},
		Artists: make([]Artists, len(item.Artists)),
	}

	for i := range item.Artists {
		link := &item.Artists[i]
		h.Artists[i] = Artists{
			ID:   link.ID,
			Name: link.Name,
			Href: link.Href,
		}
	}

	return h
}

func makeItems(item *spotify_type.Items) History {
	track := &item.Track

	h := History{
		ID:         track.ID,
		URL:        track.URI,
		Name:       track.Name,
		PreviewURL: track.PreviewURL,
		Img:        track.Album.Images[0].URL,
		Album: Album{
			ID:   track.Album.ID,
			Name: track.Album.Name,
		},
		Artists: make([]Artists, len(track.Artists)),
	}

	for i := range track.Artists {
		link := &track.Artists[i]
		h.Artists[i] = Artists{
			ID:   link.ID,
			Name: link.Name,
			Href: link.Href,
		}
	}

	return h
}

func (c *ApiContext) GetAllHistory() ([]History, error) {
	// FIXME user make
	var r []History

	now, err := c.GetPlayNow()
	if err != nil && !errors.Is(err, NotFoundError) {
		return r, err
	} else if !errors.Is(err, NotFoundError) {
		r = append(r, makeItem(&now.Item))
	}

	history, err := c.GetHistory()
	if err != nil && !errors.Is(err, NotFoundError) {
		return r, err
	} else if !errors.Is(err, NotFoundError) {
		for i := range history.Items {
			r = append(r, makeItems(&history.Items[i]))
		}
	}
	return r, nil
}
