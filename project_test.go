package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"sync"
	"testing"
	"time"
	"unsafe"
)

// 参数t用于报告测试失败和附加的日志信息
func TestStruct(t *testing.T) {
	var a Person
	var b *Person
	var c = new(Person)
	var d = &Person{}
	e := Person{}
	fmt.Println("a", a)
	fmt.Println("b", b)
	fmt.Println("c", c)
	fmt.Println("d", d)
	fmt.Println("f", e)
	f := ""
	fmt.Println(f == a.Name)
	//fmt.Fprintf(os.Stdout, "an %s\n", "error")
	//os.Create("ts.txt")
}

func TestTime(t *testing.T) {
	// time.Parse UTC 标准时区 按照指定格式将字符串转换为日期
	// time.ParseInLocation CST 北京时区（time.Local）
	// time.Unix有两个参数，第一个时秒时间戳，另一个是纳秒时间戳
	begin1 := "2023-01-18"
	begin, _ := time.Parse("2006-01-02", begin1)
	timeParse, _ := time.Parse("2006-01-02 15:04:05", "2022-05-11 15:04:05")
	fmt.Println(timeParse)
	end1 := "2023-01-29"
	end, _ := time.Parse("2006-01-02", end1)
	if time.Now().Before(begin) || time.Now().AddDate(0, 0, 11).After(end) {
		fmt.Println(end)
	}
	// time.add最大是小时
	now := time.Now()
	nowDay := now.Add(2 * time.Hour)
	nowDay2 := now.AddDate(0, 0, 1)
	fmt.Println("----------")
	dd := time.Unix(time.Now().Unix()-2, 0).String()
	fmt.Println(dd)
	fmt.Println("----------")
	fmt.Println("yyyy-MM-dd HH:mm:ss", now.Format("2006-01-02 15:04:03"))
	fmt.Println("yyyy-MM-dd HH:mm:ss", now.Format("2006-01-02 15:04:05"))
	fmt.Println("yyyy-MM-dd HH:mm:ss", now.Format("2006-01-04 15:04:05"))
	fmt.Println(nowDay)
	fmt.Println(nowDay2)
}

func TestMysql(t *testing.T) {
	//配置MySQL连接参数
	username := "root"       //账号
	password := ""           //密码
	host := "localhost"      //数据库地址，可以是Ip或者域名
	port := 3306             //数据库端口
	Dbname := "project_test" //数据库名
	timeout := "10s"         //连接超时，10秒

	//拼接下dsn参数, dsn格式可以参考上面的语法，这里使用Sprintf动态拼接dsn参数，因为一般数据库连接参数，我们都是保存在配置文件里面，需要从配置文件加载参数，然后拼接dsn。
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)
	//连接MYSQL, 获得DB类型实例，用于后面的数据库读写操作。
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
	var res []RecallUsers
	memberId := 1234

	//err = db.Debug().Table("recall_users").Where("member_id=?", memberId).Order("target_id desc").Limit(1).Scan(&res).Error
	tt := time.Unix(time.Now().Unix()-2, 0)
	timestamp := time.Now().Add(-2 * time.Second)
	fmt.Println("timestamp", timestamp, "tt", tt)
	err = db.Debug().Raw("select * from recall_users where member_id = ? and created_at < ? limit 3", memberId, tt).Scan(&res).Error
	if err != nil {
		fmt.Println("-----", err)
	}
	InMemoryTime = time.Now().Unix()
	fmt.Println("第一次连接数据库", res)
	//sqlDb, _ := db.DB()
	//defer sqlDb.Close()
}
func TestGoroutine(t *testing.T) {
	var waitGroup sync.WaitGroup
	maxCount := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	num := 3
	goroutineNum := len(maxCount) / num
	for i := 0; i <= goroutineNum; i++ {
		waitGroup.Add(1)
		var testIds []int
		if i == goroutineNum {
			testIds = maxCount[(goroutineNum)*3:]
		} else {
			testIds = maxCount[i*3 : (i+1)*3]
		}
		go func() {
			defer waitGroup.Done()
			fmt.Println("--------", testIds)
			time.Sleep(2 * time.Second)
		}()
	}
	waitGroup.Wait()
}

// 在 Go 中引入枚举的标准方法是声明一个自定义类型和一个使用了 iota 的 const 组
// 由于变量的默认值为 0，因此通常应以非零值开头枚举
type Operation int

const (
	ADD Operation = iota + 1
	Subtract
	Multiply
)

// nil在概念上和其它语言的null、None、nil、NULL一样，都指代零值或空值。nil是预先说明的标识符，即通常意义上的关键字。
// 在Golang中，nil只能赋值给指针、channel、func、interface、map或slice类型的变量
// 不同类型的 nil 值占用的内存大小可能不一样
func TestNil(t *testing.T) {
	var p *struct{}
	fmt.Println(unsafe.Sizeof(p)) // 8

	var s []int64
	fmt.Println(unsafe.Sizeof(s)) // 24
	fmt.Println(s)
	fmt.Println(s == nil)
	var m map[int]bool
	fmt.Println(unsafe.Sizeof(m)) // 8

	var c chan string
	fmt.Println(unsafe.Sizeof(c)) // 8

	var f func()
	fmt.Println(unsafe.Sizeof(f)) // 8

	var i interface{}
	fmt.Println(unsafe.Sizeof(i)) // 16
}

type MyError struct {
}
type Handler struct {
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//TODO implement me
	panic("implement me")
}

func TestInterface(t *testing.T) {
	var _ error = error(nil)
	var _ *MyError = (*MyError)(nil)
	d := []int{1, 2, 3, 4, 5}
	d = d[1:len(d)]
	a, b := 12, 45
	if a-b < 9 {

	}
	//var _ http.Handler = (*Handler)(nil)
}
