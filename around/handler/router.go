package handler

//处理请求的url和entity的关系
import (
    "net/http" 

    jwtmiddleware "github.com/auth0/go-jwt-middleware" //middleware：token验证
    jwt "github.com/form3tech-oss/jwt-go"  //library的别名：jwt（类似于object），调用时直接jwt.mtehod

    "github.com/gorilla/mux"
    "github.com/gorilla/handlers"   
)

func InitRouter() http.Handler{
    // update Initrouter（）function
    //router := mux.NewRouter()
    // router.Handle("/upload", http.HandlerFunc(uploadHandler)).Methods("POST")
	// //Search Engine realize ; router(Controller) --> service(by user? by keyword) ---> backend(ElasticSearch.go)(DAO) --->ReadFromES
 	// router.Handle("/search", http.HandlerFunc(searchHandler)).Methods("GET")

    jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
        ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) { //添加token用于验证
            return []byte(mySigningKey), nil      //加密密钥
        },
        SigningMethod: jwt.SigningMethodHS256,    //加密算法
    })

    router := mux.NewRouter()

    router.Handle("/upload", jwtMiddleware.Handler(http.HandlerFunc(uploadHandler))).Methods("POST") //验证token
    router.Handle("/search", jwtMiddleware.Handler(http.HandlerFunc(searchHandler))).Methods("GET")

    router.Handle("/signup", http.HandlerFunc(signupHandler)).Methods("POST")
    router.Handle("/signin", http.HandlerFunc(signinHandler)).Methods("POST")

    //return router


    //跨域访问：前端和后端在不同的平台跑
    //支持前端跨域访问后端时的规范程序

    //allowedOrigin（*）：允许支持哪些地址的前端请求，给后端发送的跨域请求。string("*")默认任何请求都接受，*为通配符
    originsOk := handlers.AllowedOrigins([]string{"*"})
    //allowedHeaders：后段规范可以接受的前端请求的Header类型：eg.这个例子中支持前端发送的request包含authorization，content-type header
    //支持Content-Type header；规范前端发送的请求的body中的形式。eg.如果是POST请求，前端数据应该支持json格式
    //service中signup成功会返回前端一个token信息
    //支持authorization header：前端发送的token信息。将token decode加入request中。通过decode token判断是否登陆成功
    headersOk := handlers.AllowedHeaders([]string{"Authorization", "Content-Type"})
    //“get - 读操作：search”，“Post - 写操作：upload，signup，saveES”
    methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "DELETE"})


    return handlers.CORS(originsOk, headersOk, methodsOk)(router)
}
