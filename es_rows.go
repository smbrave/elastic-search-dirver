package es_driver

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
)

type ESRows struct {
	searchResult SearchResult
	idx          int
}

func (er *ESRows) Close() error {
	return nil
}

func (er *ESRows) Columns() []string {
	columns := make([]string, 0)
	columns = append(columns, "aa")
	columns = append(columns, "bb")
	columns = append(columns, "cc")
	return columns
}

func (er *ESRows) Next(dest []driver.Value) error {

	if er.searchResult.Hits == nil {
		return io.EOF
	}

	source := make(map[string]interface{})
	if er.idx >= len(er.searchResult.Hits.Hits) {
		return io.EOF
	}

	r := er.searchResult.Hits.Hits[er.idx]
	err := json.Unmarshal(*r.Source, &source)
	if err != nil {
		return err
	}

	for _, v := range source {
		fmt.Printf("%v\n", v)
	}

	er.idx += 1
	return nil

}
