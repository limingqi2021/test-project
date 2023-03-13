package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	goredis "github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	s := time.Now().Format("2006-01-02 15:04:05")
	tt, _ := time.Parse("2006-01-02 15:04:05", s)
	fmt.Println(s)
	fmt.Println(tt)
	fmt.Println("******1*****")
	var ttt time.Time
	fmt.Println(ttt)
	fmt.Println(ttt.IsZero())

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

type ReceptionSetting struct {
	LocationIdStr  string `json:"location_id"`
	ReceptionScene int64  `json:"reception_scene"`
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
	//var res MsgTimestamp
	//ttt := ""
	//err = json.Unmarshal([]byte(ttt), &res)
	//if err != nil {
	//	panic(err)
	//}
	fmt.Println("---------1----------")

	var res ReceptionSetting
	value := `{"location_id":"11","reception_scene":1}`

	_ = json.Unmarshal([]byte(value), &res)
	locationIdStr := strings.Split(res.LocationIdStr, ",")
	var locationIdsTemp []int32
	for _, v := range locationIdStr {
		locationIdTemp, _ := strconv.Atoi(v)
		locationIdsTemp = append(locationIdsTemp, int32(locationIdTemp))
	}
	fmt.Println(locationIdsTemp)
	fmt.Println(res.ReceptionScene)
	fmt.Println(res)
	fmt.Println("---------2----------")
	fmt.Println(strings.Fields("hello widuu golang"))
}

// EX second：设置键的过期时间为 second 秒。SET key value EX second 效果等同于 SETEX key second value。
// PX millisecond：设置键的过期时间为毫秒。SET key value PX millisecond 效果等同于 PSETEX key millisecond value。
// NX：只在键不存在时，才对键进行设置操作。SET key value NX 效果等同于 SETNX key value。
// zrank:返回排名,索引值 zrange:入参索引  zrangebyscore:入参分数
func TestRedis(t *testing.T) {
	ctx := context.Background()
	rdb := goredis.NewClient(&goredis.Options{
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
func TestRedisGo(t *testing.T) {
	//1. 链接到 redis
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("redis.Dial err=", err)
		return
	}

	defer conn.Close() //关闭..
	res, err := conn.Do("SET", "time31key", 1, "EX", 86400, "NX")
	if err != nil {
		fmt.Println("err1=", err)
	}
	fmt.Println("err1=", res)
	fmt.Println(strings.Title("heer royal highness"))
}

func TestCtx(t *testing.T) {
	//ctx := context.Background().
}

func fib(n int) int {
	if n == 0 || n == 1 {
		return n
	}
	return fib(n-2) + fib(n-1)
}

func BenchmarkFib(b *testing.B) {
	//time.Sleep(time.Second * 3) // 模拟耗时准备任务
	for n := 0; n < b.N; n++ {
		fib(30) // run fib(30) b.N times
	}
}

// 拼接字符串
// + 和 fmt.Sprintf 的效率是最低的，和其余的方式相比，性能相差约 1000 倍，而且消耗了超过 1000 倍的内存。fmt.Sprintf 通常是用来格式化字符串的，一般不会用来拼接字符串。
// strings.Builder 和 + 性能和内存消耗差距如此巨大，是因为两者的内存分配方式不一样。
// + 拼接 2 个字符串时，生成一个新的字符串，那么就需要开辟一段新的空间，新空间的大小是原来两个字符串的大小之和
// strings.Builder，bytes.Buffer，包括切片 []byte 的内存是以倍数申请的。例如，初始大小为 0，当第一次写入大小为 10 byte 的字符串时，则会申请大小为 16 byte 的内存（恰好大于 10 byte 的 2 的指数），第二次写入 10 byte 时，内存不够，则申请 32 byte 的内存，第三次写入内存足够，则不申请新的，
func TestBuilderConcat(t *testing.T) {
	var str strings.Builder
	fmt.Println(str.Cap())
	str.WriteString("adf")
	fmt.Println(str.Cap())
	str.WriteString("efe")
	fmt.Println(str.Cap())
	fmt.Println(str.String())
}

func Increase() func() int {
	n := 0
	return func() int {
		n++
		return n
	}
}

func TestFunc(t *testing.T) {
	in := Increase()
	fmt.Println(in()) // 1
	fmt.Println(in()) // 2
	var s, sep string
	for i := 1; i < len(os.Args); i++ {
		s += sep + os.Args[i]
		fmt.Println(os.Args[i])
	}
	fmt.Println(s)
}

// TrimSpace只能去掉两边的空格
func TestStrconv(t *testing.T) {
	tt := " 8"
	tt = strings.TrimSpace(tt)
	id, _ := strconv.Atoi(tt)
	fmt.Println(id)
	var runes []rune
	for _, v := range "hello,世界" {
		runes = append(runes, v)
	}
	fmt.Println(runes)
}

func TestGoFunc(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(30)
	for i := 0; i < 10; i++ {
		// 每一步循环至少间隔一秒，而这一秒的时间足够启动一个goroutine了，因此这样可以输出正确的结果。
		// 在实际的工程中，不可能进行延时，这样就没有并发的优势，有两种方法：
		// time.Sleep(1 * time.Second)
		go func() {
			// 大部分都是10
			// 现象的原因在于闭包共享外部的变量i，注意到，每次调用go就会启动一个goroutine，这需要一定时间；
			// 但是，启动的goroutine与循环变量递增不是在同一个goroutine，可以把i认为处于主goroutine中。启动一个goroutine的速度远小于循环执行的速度
			// 所以即使是第一个goroutine刚起启动时，外层的循环也执行到了最后一步了。
			// 由于所有的goroutine共享i，而且这个i会在最后一个使用它的goroutine结束后被销毁，所以最后的输出结果都是最后一步的i==10。
			fmt.Println("A: ", i)
			wg.Done()
		}()
	}
	// 1:共享的环境变量作为函数参数传递
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println("B: ", i)
			wg.Done()
		}(i)
	}

	// 2:使用同名的变量保留当前的状态
	for i := 0; i < 10; i++ {
		i := i
		go func() {
			fmt.Println("C: ", i)
			wg.Done()
		}()
	}

	wg.Wait()
}

// 局部变量在delete所有元素后内存会释放，而全局变量只有在将map设置为nil后内存才会释放
// 如果删除的元素是值类型，如int，float，bool，string以及数组和struct，map的内存不会自动释放
// 如果删除的元素是引用类型，如指针，slice，map，chan等，map的内存会自动释放，但释放的内存是子元素应用类型的内存占用
// sync.map开箱即用，不能赋值为nil
// 初始值为:sync.Map{mu:sync.Mutex{state:0, sema:0x0}, read:atomic.Value{v:interface {}(nil)}, dirty:map[interface {}]*sync.entry(nil), misses:0}
// sync.map的read部分没有改变，dirty部分在执行完delete操作之后被回收
func TestMapDelete(t *testing.T) {

}
