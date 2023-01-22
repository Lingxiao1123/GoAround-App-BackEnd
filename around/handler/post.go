package handler

import (
    "encoding/json"
    "fmt"
    "net/http"
   	//前三个都是go的标准库自带的
    "path/filepath"
    
	"around/service"
    "around/model"
	//指的是刚才创建的 package/src/around/model
    "github.com/pborman/uuid"
)

//judge type
var (
    mediaTypes = map[string]string{
        ".jpeg": "image",
        ".jpg":  "image",
        ".gif":  "image",
        ".png":  "image",
        ".mov":  "video",
        ".mp4":  "video",
        ".avi":  "video",
        ".flv":  "video",
        ".wmv":  "video",
    }
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one upload request")

    p := model.Post{
        Id:      uuid.New(),
        User:    r.FormValue("user"),
        Message: r.FormValue("message"),
    }

    //header拿到文件名，来判断type是什么（from map）
    file, header, err := r.FormFile("media_file")
    if err != nil {
        http.Error(w, "Media file is not available", http.StatusBadRequest)
        fmt.Printf("Media file is not available %v\n", err)
        return
    }

    suffix := filepath.Ext(header.Filename)
    if t, ok := mediaTypes[suffix]; ok {
        p.Type = t
    } else {
        p.Type = "unknown"
    }

    err = service.SavePost(&p, file)
    if err != nil {
        http.Error(w, "Failed to save post to backend", http.StatusInternalServerError)
        fmt.Printf("Failed to save post to backend %v\n", err)
        return
    }

    fmt.Println("Post is saved successfully.")
}

//*http.Tequest: *是什么? pointer
//传入指向request对象的指针（传入地址）

//func uploadHandler(w http.ResponseWriter, r *http.Request) {
    // Parse from body of request to get a json object.
//    fmt.Println("Received one post request")
//    decoder := json.NewDecoder(r.Body)   //request body的json格式数据，变成Post类型对象
//    var p model.Post
//    if err := decoder.Decode(&p); err != nil { //传入reference地址，来修改P的内容，如果传入P不是&P，则deep copy，P不会被修改
//        panic(err)
//    }

//    fmt.Fprintf(w, "Post received: %s\n", p.Message)
//}

// req := http.Request{...}     -----> Request req = new Request(...)
// resp := http.ResponseWriter{...}
// req_ptr := &req               //2.1 先取地址
// // uploadHandler(resp,&req)   //1.  直接传入地址
// uploadHandler(resp,*req_ptr)  //2.2 然后在传入指向地址的pointer ----> shallow copy
// shallow copy能够提升运行效率，deep copy占资源率高，每次传入reques 时，不需要copy本身，只需要copy地址就行了

//加*与不加*的区别
// req:= http.Request{method:"POST"}
// resp := http.ResponseWriter{...}
// uploadHandler(resp,req)  ------->deep copy
// req.method == FETCH

//为什么reponsewriter不需要*，request 需要*

//通过ResponseWriter 写回前端
func searchHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one request for search")
    w.Header().Set("Content-Type", "application/json")

    user := r.URL.Query().Get("user")
    keywords := r.URL.Query().Get("keywords")

	//返回结果可能是一堆post，用array装起来
    var posts []model.Post
    var err error
	//either search user 或者 search keyword
    if user != "" {
        posts, err = service.SearchPostsByUser(user)
    } else {
        posts, err = service.SearchPostsByKeywords(keywords)
    }

	//定义StatusInternalSeverError：常用：200成功；404not found；500server error, etc
    if err != nil {
        http.Error(w, "Failed to read post from backend", http.StatusInternalServerError)
        fmt.Printf("Failed to read post from backend %v.\n", err)
        return
    }

    js, err := json.Marshal(posts)
    if err != nil {
        http.Error(w, "Failed to parse posts into JSON format", http.StatusInternalServerError)
        fmt.Printf("Failed to parse posts into JSON format %v.\n", err)
        return
    }
	
    w.Write(js)
}


