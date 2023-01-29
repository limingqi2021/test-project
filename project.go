package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func main() {
	testMysqlAndCache()
	time.Sleep(7 * time.Second)
	testMysqlAndCache()
}

type Person struct {
	Name string
	Age  int64
	Sex  bool
}

type RecallUsers struct {
	Id        int64     `json:"id" form:"id" validate:"required"` // id
	TargetId  int64     `json:"target_id" form:"target_id"`       // 被召回人id
	MemberId  int64     `json:"member_id" form:"member_id"`       // 用户id
	CreatedAt time.Time `json:"created_at" form:"created_at"`     // 创建时间
	UpdatedAt time.Time `json:"updated_at" form:"updated_at"`     // 更新时间
}

var InMemoryTime int64 //缓存时间
var rows *RecallUsers

func testMysqlAndCache() {
	nowUnix := time.Now().Unix()
	if rows != nil {
		if nowUnix-InMemoryTime < 5 {
			fmt.Println("cache", rows)
			return
		}
	}
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
	memberId := 1234
	err = db.Debug().Table("recall_users").Where("member_id=?", memberId).Scan(&rows).Error
	if err != nil {
		fmt.Println("-----", err)
	}
	InMemoryTime = time.Now().Unix()
	fmt.Println("第一次连接数据库", rows)
	//sqlDb, _ := db.DB()
	//defer sqlDb.Close()
}
