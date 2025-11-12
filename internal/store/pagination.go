package store

import (
	"net/http"
	"strconv"
	"time"
)

type PaginatedFeedQuery struct {
	Limit  int      `json:"limit" validate:"gte=1,lte=20"`
	Offset int      `json:"offset" validate:"gte=0"`
	Sort   string   `json:"sort" validate:"oneof=asc desc"`
	Search string   `json:"search" validate:"max=100"`
	Tags   []string `json:"tags" validate:"max=5"`
	Since  string   `json:"since"`
	Until  string   `json:"until"`
}

func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	qs := r.URL.Query()

	// Parse limit with default value of 20
	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, err
		}
		fq.Limit = l
	} else {
		fq.Limit = 20
	}

	// Parse offset with default value of 0
	offset := qs.Get("offset")
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return fq, err
		}
		fq.Offset = o
	} else {
		fq.Offset = 0
	}

	// Parse sort with default value of "desc"
	sort := qs.Get("sort")
	if sort != "" {
		fq.Sort = sort
	} else {
		fq.Sort = "desc"
	}

	// Parse search filter (optional)
	fq.Search = qs.Get("search")

	// Parse tags filter (optional, comma-separated)
	if tagsParam := qs.Get("tags"); tagsParam != "" {
		// Split by comma, supporting multiple tags
		tags := qs["tags"]
		if len(tags) > 0 {
			fq.Tags = tags
		}
	}

	// Parse since filter (optional, RFC3339 format)
	if since := qs.Get("since"); since != "" {
		// Validate date format
		if _, err := time.Parse(time.RFC3339, since); err != nil {
			return fq, err
		}
		fq.Since = since
	}

	// Parse until filter (optional, RFC3339 format)
	if until := qs.Get("until"); until != "" {
		// Validate date format
		if _, err := time.Parse(time.RFC3339, until); err != nil {
			return fq, err
		}
		fq.Until = until
	}

	return fq, nil
}
