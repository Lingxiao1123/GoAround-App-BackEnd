package service

//user service,import package
import (
    "fmt"
    "reflect"

    "around/backend"
    "around/constants"
    "around/model"

    "github.com/olivere/elastic/v7"
)
//API addUser() : signup
//bool：添加用户是否成功：or 用户名存不存在
func AddUser(user *model.User) (bool, error) {
    query := elastic.NewTermQuery("username", user.Username) //基于user的query
    searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
    if err != nil {
        return false, err
    }

    if searchResult.TotalHits() > 0 {
        return false, nil
    }

    err = backend.ESBackend.SaveToES(user, constants.USER_INDEX, user.Username)
    if err != nil {
        return false, err
    }
    fmt.Printf("User is added: %s\n", user.Username)
    return true, nil
}

//API checkUser(): login
//怎么验证用户：
//1.Read ES based on username and compare the given password with the password returned from ES
//2.Read ES based on both username and password,check if there is a hit
//bool 返回用户是否验证成功
func CheckUser(username, password string) (bool, error) {
    query := elastic.NewBoolQuery()
    query.Must(elastic.NewTermQuery("username", username)) //query.Must : 必须满足
    query.Must(elastic.NewTermQuery("password", password))
    searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
    if err != nil {
        return false, err
    }

    var utype model.User //utpye为user object
    for _, item := range searchResult.Each(reflect.TypeOf(utype)) {
        u := item.(model.User)
        if u.Password == password {
            fmt.Printf("Login as %s\n", username)
            return true, nil
        }
    }
    return false, nil
}
