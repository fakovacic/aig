package config

import (
	"net/http"
	"os"
	"time"

	logger "github.com/fakovacic/aig/internal/log"
	"github.com/olivere/elastic/v7"
	"go.elastic.co/apm/module/apmelasticsearch"
)

// Client new elasticsearch client
func Client(client *elastic.Client) (*elastic.Client, error) {
	if client != nil {
		return client, nil
	}

	client, err := elastic.NewClient(
		// elastic.SetErrorLog(log.New(os.Stderr, "ERROR:", log.LstdFlags)),
		// elastic.SetInfoLog(log.New(os.Stdout, "INFO:", log.LstdFlags)),
		// elastic.SetTraceLog(log.New(os.Stderr, "TRACE:", log.LstdFlags)),
		// elastic.SetSniff(true),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
		elastic.SetURL(os.Getenv("ELASTICSEARCH_CONN")),
		elastic.SetHttpClient(&http.Client{
			Transport: apmelasticsearch.WrapRoundTripper(http.DefaultTransport),
		}),
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// RefreshSession refresh client session
func RefreshSession(client *elastic.Client) (*elastic.Client, error) {
	var err error
	client, err = Client(client)
	if err == nil {
		return client, nil
	}

	logger.Errorf("elasticsearch error: %v", err.Error())
	logger.Info("sleeping 10 sec, than will try again connection 1/3")
	time.Sleep(10 * time.Second)

	client, err = Client(client)
	if err == nil {
		return client, nil
	}

	logger.Errorf("elasticsearch error: ", err.Error())
	logger.Info("sleeping 10 sec, than will try again connection 2/3")
	time.Sleep(10 * time.Second)

	client, err = Client(client)
	if err == nil {
		return client, nil
	}

	logger.Errorf("elasticsearch error: ", err.Error())
	logger.Info("sleeping 10 sec, than will try again connection 3/3")
	time.Sleep(10 * time.Second)

	client, err = Client(client)
	if err == nil {
		return client, nil
	}

	return client, nil
}
