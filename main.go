package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Comment struct {
	ID      int
	Name    string `json:"name"`
	Content string `json:"content"`
}

var db *gorm.DB
var err error

// 连接数据库
func initSQL() {
	//配置MySQL连接参数
	host := viper.GetString(`mysql.url`)          //数据库地址，可以是Ip或者域名
	username := viper.GetString(`mysql.username`) //账号
	password := viper.GetString(`mysql.password`) //密码
	port := viper.GetInt(`mysql.port`)            //数据库端口
	Dbname := viper.GetString(`mysql.Dbname`)     //数据库名
	timeout := viper.GetString(`mysql.timeout`)   //连接超时，10秒

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	} else {
		fmt.Println("Connected")
	}
}

func main() {
	initSQL()
	fmt.Println("Server Start")
	http.HandleFunc("/comment/get", getComment)
	http.HandleFunc("/comment/add", addComment)
	http.HandleFunc("/comment/delete", deleteComment)

	fs := http.FileServer(http.Dir('.'))
	http.Handle("/", fs)

	_ = http.ListenAndServe(":8080", nil)

}

// 获取评论
func getComment(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	currentPage, _ := strconv.Atoi(r.Form.Get("page"))
	size, _ := strconv.Atoi(r.Form.Get("size"))

	fmt.Println("GO:getComment ", currentPage, size)
	// 处理GET请求，返回评论列表
	// 查询数据库获取评论列表
	var comments []Comment

	db.Order("ID desc").Limit(size).Offset((currentPage - 1) * size).Find(&comments)

	// 返回评论列表
	var count int64
	var tmp_comments []Comment
	db.Find(&tmp_comments).Count(&count)

	fmt.Println("Got:", comments, "|total:", count)
	data := map[string]interface{}{
		"total":    count,
		"comments": comments,
	}

	response := map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": data,
	}
	response_json, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.Write(response_json)
}

// 添加评论
func addComment(w http.ResponseWriter, r *http.Request) {
	var comment Comment
	_ = json.NewDecoder(r.Body).Decode(&comment)
	var topComment Comment
	db.Last(&topComment)
	// 新评论的id是最大id+1
	comment.ID = topComment.ID + 1
	fmt.Println(comment)
	db.Create(&comment)

	// 返回添加成功的提示信息

	data := map[string]interface{}{
		"id":      comment.ID,
		"name":    comment.Name,
		"content": comment.Content,
	}

	response := map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": data,
	}
	response_json, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response_json)
}

// 按照id删除评论
func deleteComment(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var comment Comment
	db.First(&comment, r.Form.Get("id"))
	fmt.Println("deleting:", comment, "|id:", r.Form.Get("id"))
	db.Delete(&comment)

	// 返回删除成功的提示信息
	response := map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": nil,
	}
	response_json, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response_json)
}
