package backend

import (
    "context"
    "fmt"

    "around/constants"

    "github.com/olivere/elastic/v7"
)

//client object（Session Factory）要保证只能有一个，只生成一次
//这行代码就是帮助我们定义 elasticSearch class的object，对应了struct的object
//*指针：直接指向struct的地址。保证唯一性。用的时候永远指向struct的client object
//之后对数据库做操作只需要调用 ESBackend（Session Factory）即可，global pointer
var (
	//sessionFactory vs SessionFactory
    ESBackend *ElasticsearchBackend
)
//对应java的class
type ElasticsearchBackend struct {
	//SecssionFactory object 对应go里面的client object
    client *elastic.Client
}

//go 和java不同，func要定义在class 的外部
//这些mehtod的作用是与数据库进行交互(Search API - ReadES)
//backend ElasticsearchBackend 作为 receiver，可以让程序知道属于那个struct的method
//input 为数据库index 和 search query（search by user or search by keywird）
func (backend *ElasticsearchBackend) ReadFromES(query elastic.Query, index string) (*elastic.SearchResult, error) {
    searchResult, err := backend.client.Search().
        Index(index).				//search in index
        Query(query).				//specify the query
        Pretty(true).				//Pretty print request and respons JSON
        Do(context.Background())    //execute
    //handle error
		if err != nil {
        return nil, err
    }
    return searchResult, nil
}

//类似于constructor 构造器，new elasticSearchBack object
//初始化生成 client object （EsBackend）
func InitElasticsearchBackend(){
	//创建client
	client, err := elastic.NewClient(elastic.SetURL(constants.ES_URL),
	elastic.SetBasicAuth(constants.ES_USERNAME, constants.ES_PASSWORD))
	//go 语言要check err
	if err != nil {
        panic(err)
    }
	//.Do(context.Background() context的用处，hold请求，设置ddl，如果超时，就取消请求
	exists, err := client.IndexExists(constants.POST_INDEX).Do(context.Background())
	//先check是否已经存在index
    if err != nil {
        panic(err)
    }
	//如果不存在index，创建index ---> mapping
    if !exists {
		//Schma : 5个column
		//keyword，text type都是字符串：不同点：keyword：搜索的时候完全匹配 = ； text：模糊匹配：约等于号
        mapping := `{
            "mappings": {
                "properties": {
                    "id":       { "type": "keyword" },
                    "user":     { "type": "keyword" },
                    "message":  { "type": "text" },
                    "url":      { "type": "keyword", "index": false },
                    "type":     { "type": "keyword", "index": false }
                }
            }
        }`
		//如果index默认为true的话搜索更快logn；为false则是线形搜索
        _, err := client.CreateIndex(constants.POST_INDEX).Body(mapping).Do(context.Background())
        if err != nil {
            panic(err)
        }
    }

    exists, err = client.IndexExists(constants.USER_INDEX).Do(context.Background())
    if err != nil {
        panic(err)
    }

    if !exists {
        mapping := `{
                        "mappings": {
                                "properties": {
                                        "username": {"type": "keyword"},
                                        "password": {"type": "keyword"},
                                        "age":      {"type": "long", "index": false},
                                        "gender":   {"type": "keyword", "index": false}
                                }
                        }
                }`
        _, err = client.CreateIndex(constants.USER_INDEX).Body(mapping).Do(context.Background())
        if err != nil {
            panic(err)
        }
    }
    fmt.Println("Indexes are created.")
	//EsBacked为指针类型，所以赋值用-取地址符&
    ESBackend = &ElasticsearchBackend{client: client}
}

//i interface{} : 空的interface，可以指代any type的父类：代表任何数据都可以存--->这个method可以被reuse到很多地方
//input：i，index，id
//output：判断错误
func (backend *ElasticsearchBackend) SaveToES(i interface{}, index string, id string) error {
    _, err := backend.client.Index().
        Index(index).
        Id(id).
        BodyJson(i).
        Do(context.Background())
    return err
}
