package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"testing"
	"time"
)

type Person struct {
	Name string
	Age  int64
	Sex  bool
}

// 参数t用于报告测试失败和附加的日志信息
func TestStruct(t *testing.T) {
	var a Person
	var b *Person
	var c = new(Person)
	var d = &Person{}
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
	fmt.Println(d)
	e := ""
	fmt.Println(e == a.Name)
	//fmt.Fprintf(os.Stdout, "an %s\n", "error")
	//os.Create("ts.txt")
}

func TestTime(t *testing.T) {
	// time.Parse UTC 标准时区
	// time.ParseInLocation CST 北京时区（time.Local）
	begin1 := "2023-01-18"
	begin, _ := time.Parse("2006-01-02", begin1)

	end1 := "2023-01-29"
	end, _ := time.Parse("2006-01-02", end1)
	fmt.Println("---***---")
	if time.Now().Before(begin) || time.Now().AddDate(0, 0, 11).After(end) {
		fmt.Println(end)
		return
	}
	fmt.Println("------")
}

type RecallUsers struct {
	Id        int64     `json:"id" form:"id" validate:"required"` // id
	TargetId  int64     `json:"target_id" form:"target_id"`       // 被召回人id
	MemberId  int64     `json:"member_id" form:"member_id"`       // 用户id
	CreatedAt time.Time `json:"created_at" form:"created_at"`     // 创建时间
	UpdatedAt time.Time `json:"updated_at" form:"updated_at"`     // 更新时间
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

	//rows := make([]*RecallUsers, 0, 0)
	var rows *RecallUsers
	memberId := 1234
	err = db.Debug().Table("recall_users").Where("member_id=?", memberId).Scan(&rows).Error
	if err != nil {
		fmt.Println("-----", err)
	}

	fmt.Println("连接一次数据库", rows)
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

func TestEnum(t *testing.T) {

}
