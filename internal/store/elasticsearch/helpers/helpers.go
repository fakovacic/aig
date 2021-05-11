package helpers

import (
	"fmt"

	"github.com/fakovacic/aig/internal/meta"
	elastic "github.com/olivere/elastic/v7"
)

// SetSearch set search
func SetSearch(field, search string) elastic.Query {
	return elastic.NewWildcardQuery(
		field,
		fmt.Sprintf("*%s*", search),
	)
}

// SetQueryFilters  set query filters for elasticsearch
func SetQueryFilters(filters []meta.Filters) []elastic.Query {
	var query []elastic.Query
	for i := range filters {
		switch len(filters[i].Values) {
		case 0:
			continue
		case 1:
			if filters[i].Values[0] == "" {
				continue
			}

			query = append(query, elastic.NewTermQuery(
				filters[i].Field,
				filters[i].Values[0],
			))
		default:
			query = append(query, elastic.NewTermsQuery(
				filters[i].Field,
				filters[i].Values))
		}
	}

	return query
}

// InterfaceToSlice convert interface to slice
func InterfaceToSlice(x interface{}) []interface{} {
	sl := x.([]string)
	isl := make([]interface{}, len(sl))
	for i := range sl {
		isl[i] = sl[i]
	}
	return isl
}
