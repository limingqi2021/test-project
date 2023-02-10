package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"time"
)

func main() {
	//testMysqlAndCache()
	//time.Sleep(7 * time.Second)
	//testMysqlAndCache()
	//go InitTicker()
	//time.Sleep(2 * time.Minute)
	TestMap()
	// testMysqlAndCache()
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

var relations []string

func testMysqlAndCache() {
	// 读写sync.map
	nowUnix := time.Now().Unix()
	fmt.Println("----", len(relations))
	if relations != nil {
		if nowUnix-InMemoryTime < 5 {
			fmt.Println("cache", relations)
			return
		}
	}
	//var rows []*RecallUsers
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
	//memberId := 1234
	//InMemoryTime = time.Now().Unix()
	//err = db.Debug().Table("recall_users").Where("member_id=?", memberId).Scan(&rows).Error
	//if err != nil {
	//	fmt.Println("-----", err)
	//}
	//for _, v := range rows {
	//	relation := fmt.Sprintf("%d_", v.TargetId)
	//	relations = append(relations, relation)
	//}

	// 测试sql注入
	tt := "'mmmm'-- &@#'"
	//sql := fmt.Sprintf("update recall_users set nickname = %s where id = 1", "'ffff'-- &@#'")
	sql := "update recall_users set nickname = ? where id = 1"
	err = db.Debug().Exec(sql, tt).Error
	fmt.Println("第一次连接数据库", err)
	//sqlDb, _ := db.DB()
	//defer sqlDb.Close()
}

// 在未导出的顶级vars和consts， 前面加上前缀_，以使它们在使用时明确表示它们是全局符号。
const (
	_defaultPort = 8080
	DefaultUser  = "user"
)

func InitTicker() {
	clearTimer := time.NewTicker(3 * time.Second)
	defer clearTimer.Stop()

	for {
		<-clearTimer.C
		printT()
	}
}
func printT() {
	fmt.Println("-------", time.Now().Unix())
}

type relationData struct {
	RelationIds []string
	StoreTime   int64
}

func TestMap() {
	var a sync.Map
	t1 := &relationData{
		RelationIds: []string{"5018700_3_6", "5017728_2_3", "5018701_0_0"},
		StoreTime:   1675777933,
	}
	a.Store(124, t1)
	a.Store(123, t1)
	//a.Range(func(key,value interface{})) bool {
	//
	//	return true
	//}
	a.Range(func(key, value interface{}) bool {
		if v, ok := value.(*relationData); ok {
			fmt.Println(v.StoreTime)
		}
		return true
	})
	fmt.Printf("%.4f%%", float64(34/10000)*100)
	//t2, _ := a.Load(124)
	//t3, ok := t2.(*relationData)
	//if !ok {
	//	fmt.Println("不匹配")
	//}
	//fmt.Println(t3.RelationIds)

	//var res relationData
	//tTemp1, _ := json.Marshal(ttt)
	//_ = json.Unmarshal(tTemp1, &res)
}
