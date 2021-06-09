package utils

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type Entry interface {
	GetId() int
}

type ElasticHelper struct {
	client *elastic.Client
	index  string
}

func NewEsClient(host string, port int, user string, pass string) (*elastic.Client, error) {

	url := fmt.Sprintf("http://%s:%d", host, port)

	client, err := elastic.NewClient(
		// elastic 服务地址
		elastic.SetURL(url),
		elastic.SetBasicAuth(user, pass),
		// 设置错误日志输出
		elastic.SetErrorLog(log.New(os.Stderr, "", log.LstdFlags)),
		// 设置info日志输出
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (es *ElasticHelper) CreateMapping(index string, mappingFile string) error {
	mapping, err := ioutil.ReadFile(mappingFile)

	if err != nil {
		return err
	}

	mp := string(mapping)

	_, err = es.client.CreateIndex(index).Body(mp).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func NewElasticHelper(host string, port int, user, pass string, index string, batchSize int) (*ElasticHelper, error) {

	client, err := NewEsClient(host, port, user, pass)
	if err != nil {
		return nil, err
	}

	return &ElasticHelper{
		client: client,
		index:  index,
		//batchSize: batchSize,
	}, nil
}

func (es *ElasticHelper) add(ctx context.Context, entries []Entry) error {
	req := es.client.Bulk().Index(es.index)
	for _, entry := range entries {
		doc := elastic.NewBulkIndexRequest().Id(strconv.Itoa(entry.GetId())).Doc(entry)
		req.Add(doc)
	}
	if req.NumberOfActions() < 0 {
		return nil
	}
	if _, err := req.Do(ctx); err != nil {
		return err
	}
	return nil
}

func (es *ElasticHelper) Add(entries []Entry, batchSize int) {
	bs := batchSize

	var batch []Entry

	for {
		if len(entries) > bs {
			batch = entries[:bs]
			entries = entries[bs:]
		} else {
			break
		}
		err := es.add(context.Background(), batch)

		if err != nil {
			fmt.Println(err)
		}
	}

	// last batch
	err := es.add(context.Background(), entries)

	if err != nil {
		fmt.Println(err)
	}
}
