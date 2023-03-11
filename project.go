package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"time"
)

//func init() {
//	fmt.Println("init")
//}

func main() {
	//testMysqlAndCache()

	//go InitTicker()
	//time.Sleep(2 * time.Minute)
	//TestMap()
	//testMysqlAndCache()
	//for i := 0; i < 100; i++ {
	//	go HitCount()
	//}
	testMysqlAndCache()
	//fmt.Println("main")
}

var count int64

func HitCount() {
	count++
	if count%10 == 0 {
		fmt.Println("--", count)
		count = 0
	}
	//atomic.AddInt64(&count, 1)
	//value := atomic.LoadInt64(&count)
	//fmt.Println("**", value)
	//if value%10 == 0 {
	//	fmt.Println("--", count)
	//	count = 0
	//}
}

type Person struct {
	Name string
	Age  int64
	Sex  bool
}

//CREATE TABLE `users` (
//`id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
//`member_id` int NOT NULL DEFAULT 0 COMMENT '用户',
//`no` int NOT NULL DEFAULT 0 COMMENT '奖励编号',
//`status` int NOT NULL DEFAULT 0 COMMENT '状态1待发放2已发放',
//`created_at` datetime NOT NULL DEFAULT '1980-01-01 00:00:00' COMMENT '创建时间',
//`updated_at` datetime NOT NULL DEFAULT '1980-01-01 00:00:00' COMMENT '更新时间',
//PRIMARY KEY (`id`),
//UNIQUE KEY `idx_unique_member_id_no` (`member_id`,`no`)
//) ENGINE=InnoDB COMMENT='用户表';

//CreatedAt  time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;NOT NULL" json:"created_at"` // created_at
//UpdatedAt  time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated_at"` // updated_at
//
//CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
//UpdateAt  time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`

type Users struct {
	Id        int64     `json:"id" form:"id" validate:"required"`             // id
	MemberId  int64     `json:"member_id" form:"member_id"`                   // 用户id
	StartAt   time.Time `gorm:"column:start_at;NOT NULL" json:"start_at"`     // 开始时间
	CreatedAt time.Time `gorm:"column:created_at;NOT NULL" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;NOT NULL" json:"updated_at"` // 更新时间
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
		fmt.Println("连接数据库失败, error=" + err.Error())
	}
	//db.exec("sql语句")		//执行插入删除等操作使用
	//db.raw("sql语句")		//执行查询操作时使用

	// 测试sql注入
	//tt := "'mmmm'-- &@#'"
	//sql := fmt.Sprintf("update recall_users set nickname = %s where id = 1", "'ffff'-- &@#'")
	//sql := "update users set member_id = ? where id = 1"
	//不要更新某个字段
	//db.Model(&Users{}).Omit("name").Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
	//更新操作会自动运行 model 的 BeforeUpdate, AfterUpdate 方法，更新 UpdatedAt 时间戳,
	// 在更新时保存其 Associations, 如果你不想调用这些方法，你可以使用 UpdateColumn， UpdateColumns
	// 如果一个 model 有 DeletedAt 字段，他将自动获得软删除的功能！
	//当调用 Delete 方法时， 记录不会真正的从数据库中被删除， 只会将DeletedAt 字段的值会被设置为当前时间
	// err = db.Debug().Model(&Users{}).Table("users").Where("id = 1").Updates(map[string]interface{}{"member_id": 8080}).Error
	var res Users
	err = db.Debug().Table("users").Where("id = 1").Scan(&res).Error
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
	fmt.Println(res.StartAt)
	if res.StartAt.IsZero() {
		fmt.Println(res.StartAt.IsZero())
	}
	//errors.Is(err, gorm.ErrRecordNotFound)
	//fmt.Println("第一次连接数据库", err)
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

// sdfjsf
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
