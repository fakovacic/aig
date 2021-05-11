package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"

	"github.com/fakovacic/aig/internal/errors"
	"github.com/fakovacic/aig/internal/log"
	"github.com/fakovacic/aig/internal/meta"
	"github.com/fakovacic/aig/internal/store"
	"github.com/fakovacic/aig/internal/store/elasticsearch/helpers"
	elastic "github.com/olivere/elastic/v7"
)

// Store struct
type Store struct {
	store *store.Store
	model store.Model
}

// New create new model store
func New(ss *store.Store, model store.Model) store.ModelStore {
	s := Store{
		store: ss,
		model: model,
	}
	return &s
}

// GetModel return model
func (ms Store) GetModel() store.Model {
	return ms.model
}

// Init create, truncate, drop model table
func (ms Store) Init(ctx context.Context) {
	var err error
	if os.Getenv("DROP") == "true" {
		log.Infof("- drop index %s\n", ms.model.Table())
		_, err = ms.deleteIndex(ctx)
		if err != nil {
			log.Errorf("delete index:%v", err)
		}
	}

	_, err = ms.createIndex(ctx)
	if err != nil {
		log.Errorf("create index:%v", err)
	}
}

func (ms Store) checkIndex(ctx context.Context) (bool, error) {
	return ms.store.ES().IndexExists(ms.model.Table()).Do(ctx)
}

func (ms Store) createIndex(ctx context.Context) (*elastic.IndicesCreateResult, error) {
	exist, err := ms.checkIndex(ctx)
	if err != nil {
		return nil, err
	}

	if !exist {
		log.Infof("- create index %s\n", ms.model.Table())
		mapping, err := ms.model.Mapping()
		if err != nil {
			return nil, err
		}
		return ms.store.ES().CreateIndex(ms.model.Table()).Body(mapping).Do(ctx)
	}

	return nil, nil
}

func (ms Store) deleteIndex(ctx context.Context) (*elastic.IndicesDeleteResponse, error) {
	exist, err := ms.checkIndex(ctx)
	if err != nil {
		return nil, err
	}

	if exist {
		return ms.store.ES().DeleteIndex(ms.model.Table()).Do(ctx)
	}

	return nil, nil
}

// Create add row in database
func (ms Store) Create(ctx context.Context, model store.Model) (store.Model, error) {
	dataJSON, err := json.Marshal(model)
	if err != nil {
		return nil, errors.InternalWrap(err, "marshal model:")
	}
	js := string(dataJSON)

	_, err = ms.store.ES().Index().
		Index(ms.model.Table()).
		Id(model.GetID()).
		BodyJson(js).
		Do(ctx)

	if err != nil {
		return model, errors.InternalWrap(err, "create:")
	}

	return model, nil
}

// ReadBy read first row from database by given field and value
func (ms Store) ReadBy(ctx context.Context, field string, value interface{}) (store.Model, error) {
	var source json.RawMessage
	switch field {
	case "id":
		get, err := ms.store.ES().Get().
			FetchSource(true).
			Index(ms.model.Table()).
			Id(fmt.Sprintf("%v", value)).
			Do(ctx)
		if err != nil {
			val, ok := err.(*elastic.Error)
			if !ok {
				return nil, errors.InternalWrap(err, "read:")
			}

			switch val.Status {
			case http.StatusNotFound:
				return nil, errors.NotFoundWrap(err, "read:")
			default:
				return nil, errors.InternalWrap(err, "read:")
			}
		}
		source = get.Source

	default:
		search, err := ms.store.ES().Search().
			Index(ms.model.Table()).SearchSource(
			elastic.NewSearchSource().
				Query(elastic.NewMatchQuery(field, fmt.Sprintf("%v", value))),
		).Do(ctx)
		if err != nil {
			val, ok := err.(*elastic.Error)
			if !ok {
				return nil, errors.InternalWrap(err, "read:")
			}

			switch val.Status {
			case http.StatusNotFound:
				return nil, errors.NotFoundWrap(err, "read:")
			default:
				return nil, errors.InternalWrap(err, "read:")
			}
		}

		if search.TotalHits() == 0 {
			return nil, errors.NotFound("read")
		}

		source = search.Hits.Hits[0].Source
	}

	if err := json.Unmarshal(source, ms.model); err != nil {
		return nil, errors.InternalWrap(err, "unmarshal document")
	}

	return ms.model, nil
}

// Update update row in database
func (ms Store) Update(ctx context.Context, model store.Model, fields []string) (store.Model, error) {
	_, err := ms.store.ES().Update().
		Index(ms.model.Table()).
		Id(model.GetID()).
		Doc(model).
		Do(ctx)
	if err != nil {
		return nil, errors.InternalWrap(err, "update:")
	}

	return model, nil
}

// Delete deletes row from database
func (ms Store) Delete(ctx context.Context, model store.Model) error {
	_, err := ms.store.ES().Delete().
		Index(ms.model.Table()).
		Id(model.GetID()).
		Do(ctx)
	if err != nil {
		return errors.InternalWrap(err, "delete:")
	}

	return nil
}

// List return list of Users by parsed meta in request
func (ms Store) List(ctx context.Context, m *meta.Model) ([]json.RawMessage, error) {
	defer m.CalcPages()

	s := ms.store.ES().Search().
		Index(ms.model.Table()).
		TrackTotalHits(true).
		From((m.Page * m.Limit) - m.Limit).
		Size(m.Limit)

	var q []elastic.Query
	if len(m.Filters) != 0 {
		q = append(q, helpers.SetQueryFilters(m.Filters)...)
	}

	if m.Search != "" {
		q = append(q, helpers.SetSearch(ms.model.SearchField(), m.Search))
	}

	if len(q) != 0 {
		s.Query(elastic.NewBoolQuery().Must(q...))
	}

	// switch m.Sort.Value {
	// case meta.SortASC:
	// 	s.Sort(m.Sort.Field, true)
	// case meta.SortDESC:
	// 	s.Sort(m.Sort.Field, false)
	// }

	result, err := s.Do(ctx)
	if err != nil {
		return nil, errors.InternalWrap(err, "list:")
	}

	m.Count = result.TotalHits()

	list := make([]json.RawMessage, len(result.Hits.Hits))
	for i := range result.Hits.Hits {
		list[i] = result.Hits.Hits[i].Source
	}
	return list, nil
}

// ListByID return model list by ids
func (ms Store) ListByID(ctx context.Context, idList ...store.IDf) (map[string]map[string]json.RawMessage, error) {
	s := ms.store.ES().Mget()
	for i := range idList {
		mgi := elastic.NewMultiGetItem()
		mgi.Id(idList[i].ID).Index(idList[i].Table)
		s.Add(mgi)
	}

	result, err := s.Do(ctx)
	if err != nil {
		return nil, errors.InternalWrap(err, "list by ids:")
	}

	list := make(map[string]map[string]json.RawMessage)
	for i := range result.Docs {
		if !result.Docs[i].Found {
			continue
		}

		_, ok := list[result.Docs[i].Index]
		if !ok {
			list[result.Docs[i].Index] = make(map[string]json.RawMessage)
		}

		list[result.Docs[i].Index][result.Docs[i].Id] = result.Docs[i].Source
	}

	return list, nil
}

// ListIDs return id list
func (ms Store) ListIDs(ctx context.Context, querys ...store.IDs) ([]store.IDf, error) {
	var list []store.IDf
	s := ms.store.ES().MultiSearch()
	var srList []*elastic.SearchRequest
	for i := range querys {
		xType := reflect.TypeOf(querys[i].Value)
		switch xType.Kind() {
		case reflect.Slice:
			sr := elastic.NewSearchRequest().
				Index(querys[i].Table).
				Query(elastic.NewTermsQuery(
					querys[i].Field,
					helpers.InterfaceToSlice(querys[i].Value)...,
				)).Size(1000)
			srList = append(srList, sr)
		default:
			sr := elastic.NewSearchRequest().
				Index(querys[i].Table).
				Query(elastic.NewTermQuery(
					querys[i].Field,
					querys[i].Value,
				)).Size(1000)
			srList = append(srList, sr)
		}
	}

	if len(srList) == 0 {
		return list, nil
	}
	s.Add(srList...)

	result, err := s.Do(ctx)
	if err != nil {
		return nil, errors.InternalWrap(err, "list by ids:")
	}

	for i := range result.Responses {
		for x := range result.Responses[i].Hits.Hits {
			list = append(list, store.IDf{
				Table: result.Responses[i].Hits.Hits[x].Index,
				ID:    result.Responses[i].Hits.Hits[x].Id,
			})
		}
	}

	return list, nil
}
