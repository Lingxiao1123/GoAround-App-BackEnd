package service

import (
    "mime/multipart"
	//reflect读取数据做过滤
    "reflect"

    "around/backend"
    "around/constants"
    "around/model"
	//client 开发包
    "github.com/olivere/elastic/v7"
)

//search Post by user
func SearchPostsByUser(user string) ([]model.Post, error) {
    query := elastic.NewTermQuery("user", user)		//query type ：termQuery：精准匹配 ---> 按照field来匹配
    searchResult, err := backend.ESBackend.ReadFromES(query, constants.POST_INDEX) //call backend Func：backend.ESBackend.ReadFromEs
    if err != nil {
        return nil, err
    }
    return getPostFromSearchResult(searchResult), nil
}

func SearchPostsByKeywords(keywords string) ([]model.Post, error) {
    query := elastic.NewMatchQuery("message", keywords)
    query.Operator("AND")
    if keywords == "" {
        query.ZeroTermsQuery("all")
    }
    searchResult, err := backend.ESBackend.ReadFromES(query, constants.POST_INDEX)
    if err != nil {
        return nil, err
    }
    return getPostFromSearchResult(searchResult), nil
}

func getPostFromSearchResult(searchResult *elastic.SearchResult) []model.Post {
    var ptype model.Post
    var posts []model.Post

	//Reflect过滤，如果数据是ptype就拿出来
    for _, item := range searchResult.Each(reflect.TypeOf(ptype)) {
        p := item.(model.Post)
        posts = append(posts, p)
    }
    return posts
}

// func SavePost(post *model.Post,file multipart.File) error{
//     //SaveToGCS() //返回url
//     mediaLink,err := backend.GCSBackend.SaveToGCS(file,post.Id)
//     if err != nil{
//         return err
//     }
//     post.url = mediaLink

//     //SaveToES() ---> id,
//     err = backend.ESBackend.SaveToES(post,constants.POST_INDEX,post.Id)
//     if err != ril{
          //解决方法
//        //1. rollback ---- >delete from GCS 保持原子性：特点共存亡，要有一期有，要无一起无
//        //2. retry - call SaveToES again
//        //3. offline service
//     }
// }

func SavePost(post *model.Post, file multipart.File) error {
    medialink, err := backend.GCSBackend.SaveToGCS(file, post.Id)
    if err != nil {
        return err
    }
    post.Url = medialink

    return backend.ESBackend.SaveToES(post, constants.POST_INDEX, post.Id)
}