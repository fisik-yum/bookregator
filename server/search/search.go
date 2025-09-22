package search

import (
	"context"
	"database/sql"
	"errors"

	"github.com/blevesearch/bleve/v2"
)

type SearchMachine struct {
	index bleve.Index
}

type BookModel struct {
	Title string `json:"title"`
	Olid  string `json:"olid"`
}

func NewSearchMachine(loc string) (*SearchMachine, error) {
	index, err := bleve.Open(loc)
	// If bleve file doesn't exist, create it
	if errors.Is(err, bleve.ErrorIndexPathDoesNotExist) {
		mappng := bleve.NewIndexMapping()
		index, err = bleve.New(loc, mappng)
		if err != nil {
			return nil, err
		}
		return &SearchMachine{index}, nil
	} else if err != nil {
		return nil, err
	} else {
		return &SearchMachine{index}, nil
	}
}

// Index a book item. Wrapper around index.Index()
func (s *SearchMachine) AddItem(b BookModel) error {
	/*
		NOTE: Is confusing to read, but we want the OLID from any of the fields
		in our BookModel
	*/
	return s.index.Index(b.Olid, b)
}

// Search with book field. Wrapper around index.SearchInContext()
func (s *SearchMachine) SearchItem(ident string, ctx context.Context) ([]string, error) {
	q := bleve.NewMatchPhraseQuery(ident)
	q.SetFuzziness(2)
	r, err := s.index.SearchInContext(ctx, bleve.NewSearchRequest(q))
	if err != nil {
		return nil, err
	}
	results := make([]string, r.Hits.Len())
	for i, v := range r.Hits {
		results[i] = v.ID
	}
	return results, nil
}
func (s *SearchMachine) Close() {
	s.index.Close()
}

func (s *SearchMachine) Refresh(D *sql.DB) error {
	rows, err := D.Query("SELECT olid,title from works;")
	if err != nil {
		return err
	}
	for rows.Next() {
		val := BookModel{}
		err = rows.Scan(&val.Olid, &val.Title)
		if err != nil {
			return err
		}
		err = s.AddItem(val)
		if err != nil {
			return err
		}

	}
	return nil
}
