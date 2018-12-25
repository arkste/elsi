# elsi

elsi (Elasticsearch Indexer) is a simple multi-threaded command line tool written in Go (Golang) to help you quickly populate some data into Elasticsearch from different data sources.

Usually you'll have to import a third-party Elasticsearch client library into your project before you're able to get data into Elasticsearch, but sometimes third-party libraries don't exist, are outdated or there are other reasons why you can't import a client library. This is the problem elsi tries to solve.

elsi currently doesn't provide any syncing of the data, you'll have to reindex the data if it changes, but elsi will always populate a new index and then create an alias, the old data will be present while re-indexing until the new index is fully populated. A future update might include syncing.

Since Elasticsearch exposes a REST-API on Port 9200, there's no need for elsi providing a REST-API itself.

**elsi is compatible with Elasticsearch 5.x.x & 6.x.x**

## Install

+  Install Go (1.9+) and set your [GOPATH](https://golang.org/doc/code.html#GOPATH)
+ `go get -u github.com/arkste/elsi`
+ cd `$GOPATH/src/github.com/arkste/elsi`
+ `make`

or use Docker:

+ `docker pull arkste/elsi`

or get a pre-compiled binary from the [Releases](/releases)-Page.

## How to use

### Config

You need to create a config (default path `$HOME/.elsi.yml`) before running elsi:

    # $HOME/.elsi.yml

    host: http://127.0.0.1:9200
    index: elsi
    type: documents
    gzip: false
    mapping:
        settings:
            number_of_shards: 1
            number_of_replicas: 0
    pipeline:
        description: "extract attachment information"
        processors:
            - attachment:
                field: "data"

If you dont want to create a config file in your Home Directory, you can provide a custom config path with `--config`:

    $ elsi fs --config /path/to/config.yml --source /path/to/folder

### Data Sources

elsi supports the following data sources:

#### Filesystem

elsi will walk through the folder and index the path & name of every file it finds, the Elasticsearch document ID is the MD5-Hash of the filepath. It won't index the contents of the files per default, because you'll need the [Elasticsearch Ingest Attachment-Plugin](https://www.elastic.co/guide/en/elasticsearch/plugins/master/ingest-attachment.html).

    $ elsi fs --source /path/to/folder --exclude ".git,*.log"

It's recommended to install the [Elasticsearch Ingest Attachment-Plugin](https://www.elastic.co/guide/en/elasticsearch/plugins/master/ingest-attachment.html), which will provide meta informations (like the MIME-Type) about the files and will index contents of *.pdf, *.doc, etc. 
You can limit the filesize (in MB) of the files elsi will index with `--limit`. 
If you've installed the Ingest Attachment-Plugin, add the `--pipeline` flag:

    $ elsi fs --source /path/to/folder --exclude ".git,*.log" --limit 10 --pipeline 1

#### JSON

elsi will walk the folder and index the raw contents of `*.json` files. The Elasticsearch document ID is the MD5-Hash of the filepath.

    $ elsi json --source /path/to/folder/with/json/files

#### MySQL

elsi will index the response of a MySQL query. elsi will try to find an `ID` column and use it as the Elasticsearch document ID, otherwise it'll be auto-generated by Elasticsearch.

    $ elsi mysql --dsn "user:password@tcp(127.0.0.1:3306)/database?charset=utf8mb4&collation=utf8mb4_unicode_ci" --query "SELECT * FROM table"

## Custom-Mapping

You can provide a custom mapping in the config file, elsi will convert the yaml mapping 1:1 to json:

    # $HOME/.elsi.yml

    host: http://127.0.0.1:9200
    index: CUSTOM_INDEX_NAME
    type: CUSTOM_TYPE_NAME
    gzip: false
    mapping:
        settings:
            number_of_shards: 1
            number_of_replicas: 0
            analysis:
                filter:
                    autocomplete_filter:
                        type: edge_ngram
                        min_gram: 2
                        max_gram: 20
                analyzer:
                    custom_standard:
                        tokenizer: standard
                        filter: [lowercase,asciifolding,elision]
                    autocomplete:
                        tokenizer: standard
                        filter: [lowercase,asciifolding,elision,autocomplete_filter]
        mappings:
            CUSTOM_TYPE_NAME:
                dynamic_templates:
                    - default_integer:
                        match_mapping_type: long
                        mapping:
                            type: integer
                    - default_string:
                        match: "*"
                        match_mapping_type: string
                        mapping:
                            type: keyword
                            ignore_above: 256
                            fields:
                                search:
                                    type: text
                                    analyzer: custom_standard
                properties:
                    CUSTOM_FIELD_NAME:
                        type: text
                        fields:
                            keyword:
                                type: keyword
                                ignore_above: 256
                            search:
                                type: text
                                analyzer: custom_standard
                            ac:
                                type: text
                                analyzer: autocomplete
    pipeline:
        description: "extract attachment information"
        processors:
            - attachment:
                field: "data"
                ignore_missing: true

## License

elsi is released under the
[MIT License](http://www.opensource.org/licenses/MIT).