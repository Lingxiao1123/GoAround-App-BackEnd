package main

import (
    "fmt"
    "log"
    "net/http" 

	"around/backend"
    "around/handler"   
)
func main() {
    fmt.Println("started-service")
	backend.InitElasticsearchBackend()
    backend.InitGCSBackend()
	//启动go的内置http surver ListenAndServer监听8080端口，router负责请求分发模块
    log.Fatal(http.ListenAndServe(":8080", handler.InitRouter()))
}
