package elasticsearch

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/olivere/elastic"
)

// Client represents the Client-Config
type Client struct {
	Host          string `yaml:"host"`
	Index         string `yaml:"index"`
	Type          string `yaml:"type"`
	Gzip          bool   `yaml:"gzip"`
	Mapping       string `yaml:"mapping,omitempty"`
	Pipeline      string `yaml:"pipeline,omitempty"`
	PipelineName  string `yaml:"pipeline_name,omitempty"`
	client        *elastic.Client
	processor     *elastic.BulkProcessor
	indexName     string
	tmpOldIndices []string
	UsePipeline   bool
}

// Init initializes the Elasticsearch Client, called at the beginning of a command
func (es *Client) Init() {
	if es.PipelineName == "" {
		es.PipelineName = "elsi_attachment"
	}
	es.createClient()
	es.pingClient()
	es.aliasedIndices()
	es.createIndex()
	if es.UsePipeline {
		es.createPipeline()
	}
	es.createProcessor()
	log.Println("Indexing ...")
}

// Flush flushes the Elasticsearch Client, called at the end of a command
func (es *Client) Flush() {
	log.Println("... done")
	es.flushProcessor()
	es.createAlias()
	es.deleteOldIndices()
}

// AddDocument adds a Document to the Elasticsearch Bulk Processor
func (es *Client) AddDocument(id string, document interface{}, pipeline string) {
	// index path to elasticsearch
	es.processor.Add(elastic.
		NewBulkIndexRequest().
		Index(es.indexName).
		Type(es.Type).
		Id(id).
		Doc(document).
		Pipeline(pipeline))
}

func (es *Client) createClient() {
	client, err := elastic.NewClient(
		elastic.SetURL(es.Host),
		elastic.SetGzip(es.Gzip),
	)
	if err != nil {
		log.Fatalln(err)
	}

	es.client = client
}

func (es *Client) pingClient() {
	info, code, err := es.client.Ping(es.Host).Do(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Elasticsearch returned code \"%d\" and version \"%s\"\n", code, info.Version.Number)
}

func (es *Client) aliasedIndices() {
	res, err := es.client.Aliases().Index("_all").Do(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	es.tmpOldIndices = res.IndicesByAlias(es.Index)
}

func (es *Client) createIndex() {
	// Create Index
	t := time.Now()
	es.indexName = fmt.Sprintf("%s_%s", es.Index, t.Format("2006-01-02-150405"))
	_, err := es.client.CreateIndex(es.indexName).BodyString(es.Mapping).Do(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Elasticsearch Index \"%s\" created\n", es.indexName)
}

func (es *Client) createPipeline() {
	// Check if Pipeline exists, if yes delete
	_, err := es.client.IngestGetPipeline(es.PipelineName).Do(context.Background())
	if err == nil {
		es.client.IngestDeletePipeline(es.PipelineName).Do(context.Background())
	}
	// Create a new Pipeline (in case the mapping changed)
	_, err = es.client.IngestPutPipeline(es.PipelineName).BodyJson(es.Pipeline).Do(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Elasticsearch Pipeline \"%s\" created\n", es.PipelineName)
}

func (es *Client) createProcessor() {
	workers := runtime.NumCPU()

	// Create BulkProcessor
	processor, err := es.client.BulkProcessor().
		Name("Indexer").
		Workers(workers).
		BulkActions(1000).
		Do(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	es.processor = processor
	log.Printf("Elasticsearch Bulk Processor with %d Workers created\n", workers)
}

func (es *Client) flushProcessor() {
	// Flush BulkProcessor
	err := es.processor.Flush()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Elasticsearch Bulk Processor flushed")

	// Close BulkProcessor
	es.processor.Close()
	log.Println("Elasticsearch Bulk Processor closed")
}

func (es *Client) createAlias() {
	// Create/Replace alias
	_, err := es.client.Alias().Add(es.indexName, es.Index).Do(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Elasticsearch Alias \"%s\" mapped to Index \"%s\"\n", es.Index, es.indexName)
}

func (es *Client) deleteOldIndices() {
	// Delete Old Indices and ignore all errors
	for _, oldIndex := range es.tmpOldIndices {
		es.client.DeleteIndex(oldIndex).Do(context.Background())
		log.Printf("Old Elasticsearch Index \"%s\" deleted\n", oldIndex)
	}
}
