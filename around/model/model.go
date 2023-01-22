// like entity package 
package model
//1.public 和 private
// 大小写区分public和private：大写public小写private

//2.json标签的作用
// Id : json:'id' ---> 自带 json convertor  ---> 写json（标签）用来配对和转换go文件和json文件
// json数据：key-value pair，
//	 'id' : 1
//	 'user' : v
// 转译为post.go文件
// java中通过jakson library帮助我们将java object 转化成json文件，需要annotation

// 3.反引号的作用
// `的作用：在string已经自带双引号时，string的开头结尾用反引号表示，这样字符串中有双引号时，默认其为普通字符

// 4.ID通过uuid unique的生成
// 5.GCS返回url存到ES中 ：调用SaveToEs Method

// 6.前端通过 url 访问GCS

type Post struct {
    Id		string `json:"id"`
    User    string `json:"user"`
    Message string `json:"message"`
    Url     string `json:"url"`
    Type    string `json:"type"`
}

type User struct {
    Username string `json:"username"`
    Password string `json:"password"`
    Age      int64  `json:"age"`
    Gender   string `json:"gender"`
}