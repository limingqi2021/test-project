package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
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

func TestSinceTime(t *testing.T) {
	start := time.Now()
	time.Sleep(2 * time.Second)
	fmt.Println(time.Since(start))
}
func TestMysql(t *testing.T) {
	//配置MySQL连接参数
	username := "root"       //账号
	password := ""           //密码
	host := "localhost"      //数据库地址，可以是Ip或者域名
	port := 3306             //数据库端口
	Dbname := "project_test" //数据库名
	//timeout := "10s"         //连接超时，10秒

	//拼接下dsn参数, dsn格式可以参考上面的语法，这里使用Sprintf动态拼接dsn参数，因为一般数据库连接参数，我们都是保存在配置文件里面，需要从配置文件加载参数，然后拼接dsn。
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, Dbname)
	//连接MYSQL, 获得DB类型实例，用于后面的数据库读写操作。
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
	var res []Users
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
	fmt.Printf("%.4f%%", float64(345)/float64(10000))
	fmt.Println(float64(3432543545) / float64(1000450))
	//var _ http.Handler = (*Handler)(nil)
}

func TestAtomic(t *testing.T) {
	for i := 0; i < 100; i++ {
		go HitCount()
	}
	time.Sleep(1 * time.Minute)
}

func TestArr(t *testing.T) {
	//ids := []int64{10112, 9714, 9187, 9987, 10147, 9109, 10022, 10070, 9799, 10043, 9291, 9723}
	//id := 10112
	//ids = RemoveIdFromIds(int64(id), ids)
	//fmt.Println(ids)
	//id1 := reflect.ValueOf([]int64{1, 2, 3})
	//id2 := reflect.ValueOf([]int64{3, 2, 3})
	//id3 := reflect.AppendSlice(id1, id2)

	// type rune = int32；官方对它的解释是：rune是类型int32的别名
	// rune跟byte是 Go 语言中仅有的两个类型别名，专门用来处理字符
	// 在 Go 语言中，字符可以被分成两种类型处理：对占 1 个字节的英文类字符，可以使用byte（或者unit8）
	// 对占 1 ~ 4 个字节的其他字符，可以使用rune（或者int32），如中文、特殊符号等。
	tt := "你好google"
	fmt.Println(string([]rune(tt)[:3]))
}
func RemoveIdFromIds(id int64, ids []int64) []int64 {
	var arr []int64
	for _, item := range ids {
		if item != id {
			arr = append(arr, item)
		}
	}
	fmt.Println(arr)
	return arr
}

func TestDefer(t *testing.T) {
	sumFunc := lazySum([]int{1, 2, 3, 4, 5})
	fmt.Println("等待一会")
	fmt.Println("结果：", sumFunc())

}
func lazySum(arr []int) func() int {
	fmt.Println("先获取函数，不求结果")
	var sum = func() int {
		fmt.Println("求结果...")
		result := 0
		for _, v := range arr {
			result = result + v
		}
		return result
	}
	return sum
}

type MsgTimestamp struct {
	ExpirationTimestamp int64 `json:"expiration_timestamp"`
}

func TestJsonMarshal(t *testing.T) {
	content := make(map[string]interface{})
	content["msgType"] = "GET_FLOW_CARD"
	content["flow_card_desc"] = fmt.Sprintf("使用%d分钟内获得额外人气", 1800/60)
	content["expiration_timestamp"] = 100000000 + int64(3*3600)
	msgContent, err := json.Marshal(content)
	if err != nil {
		fmt.Println(err)
	}
	tt := string(msgContent)
	var msgRes MsgTimestamp
	err = json.Unmarshal([]byte(tt), &msgRes)
	if err != nil {
		fmt.Println(err)
	}
	now := time.Now().Unix()
	fmt.Println(msgRes.ExpirationTimestamp)
	if msgRes.ExpirationTimestamp < now {
		fmt.Println(msgRes.ExpirationTimestamp)
		fmt.Println(now)
	}
}
func TestRedis(t *testing.T) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // 没有密码，默认值
		DB:       0,  // 默认DB 0
	})
	content := make(map[string]interface{})
	content["msgType"] = "GET_FLOW_CARD"
	content["flow_card_desc"] = fmt.Sprintf("使用%d分钟内获得额外人气", 1800/60)
	content["expiration_timestamp"] = 100000000 + int64(3*3600)
	msgContent, _ := json.Marshal(content)
	err := rdb.HSet(ctx, "fixcard", 1234, 2345, string(msgContent))
	if err != nil {

	}
	result := rdb.HGetAll(ctx, "1234")
	fmt.Println(result)
	for _, v := range result.Val() {
		fmt.Println(v)
	}

}
