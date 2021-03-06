package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Unknwon/GoConfig"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/jiangmitiao/cali/app/models"
	"github.com/jiangmitiao/cali/app/rcali"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/sha3"
	"io"
	"os"
	"path"
	"time"
)

var engineTest *xorm.Engine

func init() {
	fmt.Println("dbService ok")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func GetSqliteDbPath() string {
	c, err := goconfig.LoadConfigFile("conf/app.conf")
	if err != nil {
		fmt.Println("读取失败")
		return ""
	} else {
		str, _ := c.GetValue("", "sqlitedb.path")
		return str
	}
}

func ListTableContent(engine xorm.Engine) {
	authors := make([]models.Author, 0)
	engine.Limit(2, 0).Find(&authors)
	authorJosnByte, _ := json.Marshal(authors)
	fmt.Println(string(authorJosnByte))

	books := make([]models.Book, 0)
	engine.Limit(2, 0).Find(&books)
	booksJosnByte, _ := json.Marshal(books)
	fmt.Println(string(booksJosnByte))

}

func DbInit() (bool, error) {
	fmt.Println(GetSqliteDbPath())

	var err error
	engineTest, err = xorm.NewEngine("sqlite3", GetSqliteDbPath())
	//dataSourceName := username + ":" + password + "@tcp(" + host + ")/" + database + "?charset=utf8"
	//engine, err = xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		return false, err
	}
	engineTest.ShowSQL(true)
	engineTest.Logger().SetLevel(core.LOG_DEBUG)
	err = engineTest.Ping()
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	engineTest.Logger().Info("----------创建表----------")

	engineTest.Logger().Info("----------创建表结束----------")

	engineTest.Logger().Info("----------插入默认数据----------")

	engineTest.Logger().Info("----------插入默认数据结束----------")

	//ListTableContent(*engine)
	return true, nil

}

func PathTest() {
	fmt.Println(path.Join("/home", "path1/", "path2", "/path3/", "cover.jpeg"))
}

func Sha3_256(in string) string {
	m := sha3.New256()
	io.WriteString(m, in)
	return hex.EncodeToString(m.Sum(nil))
}

func SearchBooks() {
	DbInit()
	bookVos := make([]models.BookVo, 0)
	if engineTest == nil {
		fmt.Println("error")
		return
	}
	err := engineTest.SQL("select books.* ,ratings.rating,authors.name  from books,authors,books_authors_link left join (select books_ratings_link.book,ratings.rating from ratings,books_ratings_link where ratings.id=books_ratings_link.rating) as ratings on books.id=ratings.book where books.id=books_authors_link.book and authors.id=books_authors_link.author and (books.title like ? ) limit ?,?", "%Quick%", 0, 10).Find(&bookVos)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	fmt.Printf("%d\n\n%+v\n\n", len(bookVos), bookVos)
	time.Sleep(time.Second * 2)
	engineTest.Close()
}

func mailtest() {
	user, _ := os.LookupEnv("CALIEMAIL")
	password, _ := os.LookupEnv("CALIEMAILPASSWORD")
	host, _ := os.LookupEnv("CALISMTP")
	to, _ := os.LookupEnv("CALIEMAILTESTTO")

	fmt.Println(user, password, host, to)
	if user == "" || password == "" || host == "" || to == "" {
		return
	}
	subject := "hello"
	body := "ceshi kfnlgnlrnglkr <a href='http://baidu.com'>baidu</a>"
	rcali.SendToMail(user, password, host, to, subject, body, "html")
}

func main() {
	//DbInit()
	//PathTest()
	//fmt.Println(Sha3_256("admininit"))

	//userhome, _ := rcali.Home()
	//fmt.Println(userhome)

	//SearchBooks()
	mailtest()
}
