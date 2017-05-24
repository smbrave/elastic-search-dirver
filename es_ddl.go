package es_driver

import "encoding/json"

type SearchHitHighlight map[string][]string
type Aggregations map[string]*json.RawMessage
type SearchSuggest map[string][]SearchSuggestion

// shardsInfo represents information from a shard.
type shardsInfo struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Failed     int `json:"failed"`
}

type ErrorDetails struct {
	Type         string                   `json:"type"`
	Reason       string                   `json:"reason"`
	ResourceType string                   `json:"resource.type,omitempty"`
	ResourceId   string                   `json:"resource.id,omitempty"`
	Index        string                   `json:"index,omitempty"`
	Phase        string                   `json:"phase,omitempty"`
	Grouped      bool                     `json:"grouped,omitempty"`
	CausedBy     map[string]interface{}   `json:"caused_by,omitempty"`
	RootCause    []*ErrorDetails          `json:"root_cause,omitempty"`
	FailedShards []map[string]interface{} `json:"failed_shards,omitempty"`
}

type SearchSuggestion struct {
	Text    string                   `json:"text"`
	Offset  int                      `json:"offset"`
	Length  int                      `json:"length"`
	Options []SearchSuggestionOption `json:"options"`
}

type SearchSuggestionOption struct {
	Text         string      `json:"text"`
	Highlighted  string      `json:"highlighted"`
	Score        float64     `json:"score"`
	CollateMatch bool        `json:"collate_match"`
	Freq         int         `json:"freq"` // deprecated in 2.x
	Payload      interface{} `json:"payload"`
}

type SearchExplanation struct {
	Value       float64             `json:"value"`             // e.g. 1.0
	Description string              `json:"description"`       // e.g. "boost" or "ConstantScore(*:*), product of:"
	Details     []SearchExplanation `json:"details,omitempty"` // recursive details
}

type SearchHitInnerHits struct {
	Hits *SearchHits `json:"hits"`
}

// SearchHit is a single hit.
type SearchHit struct {
	Score          *float64                       `json:"_score"`          // computed score
	Index          string                         `json:"_index"`          // index name
	Type           string                         `json:"_type"`           // type meta field
	Id             string                         `json:"_id"`             // external or internal
	Uid            string                         `json:"_uid"`            // uid meta field (see MapperService.java for all meta fields)
	Timestamp      int64                          `json:"_timestamp"`      // timestamp meta field
	TTL            int64                          `json:"_ttl"`            // ttl meta field
	Routing        string                         `json:"_routing"`        // routing meta field
	Parent         string                         `json:"_parent"`         // parent meta field
	Version        *int64                         `json:"_version"`        // version number, when Version is set to true in SearchService
	Sort           []interface{}                  `json:"sort"`            // sort information
	Highlight      SearchHitHighlight             `json:"highlight"`       // highlighter information
	Source         *json.RawMessage               `json:"_source"`         // stored document source
	Fields         map[string]interface{}         `json:"fields"`          // returned fields
	Explanation    *SearchExplanation             `json:"_explanation"`    // explains how the score was computed
	MatchedQueries []string                       `json:"matched_queries"` // matched queries
	InnerHits      map[string]*SearchHitInnerHits `json:"inner_hits"`      // inner hits with ES >= 1.5.0

	// Shard
	// HighlightFields
	// SortValues
	// MatchedFilters
}

// SearchHits specifies the list of search hits.
type SearchHits struct {
	TotalHits int64        `json:"total"`     // total number of hits found
	MaxScore  *float64     `json:"max_score"` // maximum score of all hits
	Hits      []*SearchHit `json:"hits"`      // the actual hits returned
}

// SearchResult is the result of a search in Elasticsearch.
type SearchResult struct {
	TookInMillis    int64         `json:"took"`             // search time in milliseconds
	ScrollId        string        `json:"_scroll_id"`       // only used with Scroll and Scan operations
	Hits            *SearchHits   `json:"hits"`             // the actual search hits
	Suggest         SearchSuggest `json:"suggest"`          // results from suggesters
	Aggregations    Aggregations  `json:"aggregations"`     // results from aggregations
	TimedOut        bool          `json:"timed_out"`        // true if the search timed out
	TerminatedEarly bool          `json:"terminated_early"` // true if the operation has terminated before e.g. an expiration was reached
	//Error        string        `json:"error,omitempty"` // used in MultiSearch only
	// TODO double-check that MultiGet now returns details error information
	Error  *ErrorDetails `json:"error,omitempty"`   // only used in MultiGet
	Shards *shardsInfo   `json:"_shards,omitempty"` // shard information
}

// {"from":0,"post_filter":{"bool":{"must":[{"term":{"order_type":"2"}}]}},"size":10,"sort":[{"order_time":{"order":"desc"}}]}
type SearchRequest struct {
	From   int         `json:"from"`
	Size   int         `json:"size"`
	Filter interface{} `json:"post_filter"`
	Sort   interface{} `json:"sort"`
}
