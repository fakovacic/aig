package meta

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/fakovacic/aig/internal/errors"
)

// Model struct
type Model struct {
	Page          int               `json:"page"`
	Pages         int               `json:"pages"`
	Count         int64             `json:"count"`
	Limit         int               `json:"limit"`
	Search        string            `json:"search"`
	Sort          Sort              `json:"sort"`
	Filters       []Filters         `json:"filters"`
	SortFields    map[string]string `json:"-"`
	FilterFields  map[string]string `json:"-"`
	FiltersValues map[string]string `json:"-"`
}

// Sort struct
type Sort struct {
	Field string   `json:"field"`
	Value SortType `json:"value"`
}

// Filters struct
type Filters struct {
	Field     string   `json:"field"`
	Condition string   `json:"condition"`
	Statement string   `json:"statement"`
	Values    []string `json:"values"`
}

// New create new meta model
func New() *Model {
	return &Model{
		Page:  1,
		Limit: 10,
		Sort: Sort{
			Field: "created",
			Value: SortDESC,
		},
		FiltersValues: make(map[string]string),
	}
}

// ParseRequest parse request from http
func (d *Model) ParseRequest(r *http.Request) error {
	vals := r.URL.Query()

	search := vals.Get("search")
	if search != "" {
		d.Search = search
	}

	pager := vals.Get("pager")
	page := vals.Get("page")
	if page != "" {
		iPage, err := strconv.Atoi(page)
		if err != nil {
			return err
		}

		if pager == "next" {
			iPage++
			if iPage == d.Pages {
				iPage--
			}
		}

		if pager == "prev" {
			iPage--
			if iPage == d.Pages {
				iPage++
			}
		}

		d.Page = iPage
	}

	limit := vals.Get("limit")
	if limit != "" {
		iLimit, err := strconv.Atoi(limit)
		if err != nil {
			return err
		}
		d.Limit = iLimit
	}

	d.Sort.Field = vals.Get("sort.field")
	if d.Sort.Field == "" {
		d.Sort.Field = "created"
	}
	err := d.Sort.Value.Parse(vals.Get("sort.value"))
	if err != nil {
		return err
	}
	if d.Sort.Value == "" {
		d.Sort.Value = SortDESC
	}

	for key, values := range vals {
		if strings.HasPrefix(key, "filters.") {
			key = strings.ReplaceAll(key, "filters.", "")
			if len(values) == 0 {
				continue
			}
			d.Filters = append(d.Filters, Filters{
				Field:  key,
				Values: values,
			})
			d.FiltersValues[key] = strings.Join(values, ",")
		}
	}

	return nil
}

// validate

// CalcPages calculate pages
func (d *Model) CalcPages() {
	d.Pages = (int(d.Count) / d.Limit) + 1
}

// PrevPage return previous page
func PrevPage(m *Model) string {
	if m.Page == 1 {
		return fmt.Sprintf("%d", m.Page)
	}
	return fmt.Sprintf("%d", (m.Page - 1))
}

// NextPage return next page
func NextPage(m *Model) string {
	if m.Page == m.Pages {
		return fmt.Sprintf("%d", m.Page)
	}
	return fmt.Sprintf("%d", (m.Page + 1))
}

// SortType sort type
type SortType string

const (
	// SortASC cost
	SortASC SortType = "ASC"
	// SortDESC cost
	SortDESC SortType = "DESC"
)

// Parse source from text
func (t *SortType) Parse(s string) error {
	s = strings.Trim(s, "\"")
	switch s {
	case "ASC":
		*t = SortASC
	case "DESC":
		*t = SortDESC
	case "":
	default:
		return errors.BadRequest("invalid sort '%s'", s)
	}

	return nil
}

func (t SortType) String() string {
	return string(t)
}

// UnmarshalJSON implements json.Marshaler
func (t *SortType) UnmarshalJSON(b []byte) error {
	return t.Parse(string(b))
}
