package orm

import (
	"fmt"
	"github.com/zituocn/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var (
	//dbs
	dbs map[string]*gorm.DB

	//defaultDBName
	defaultDBName string
)

// DBConfig mysql配置文件
type DBConfig struct {
	Name            string //库名
	User            string //登录名
	Password        string //密码
	Host            string //主机
	Port            int    //port
	Debug           bool   //是否debug
	DisablePrepared bool   //是否disable prepared
}

// InitDefaultDB init single db config
func InitDefaultDB(db *DBConfig) (err error) {
	if db == nil {
		err = fmt.Errorf("没有需要init的DB")
		return
	}
	defaultDBName = db.Name
	dbs = make(map[string]*gorm.DB, 1)
	newORM(db)
	return
}

// InitDB init multiple db
func InitDB(list []*DBConfig) (err error) {
	if len(list) == 0 {
		err = fmt.Errorf("没有需要init的DB")
		return
	}
	dbs = make(map[string]*gorm.DB, len(list))
	for _, item := range list {
		newORM(item)
	}
	return
}

// newORM a new ORM
func newORM(db *DBConfig) {
	var (
		orm *gorm.DB
		err error
	)
	if db.User == "" || db.Password == "" || db.Host == "" || db.Port == 0 {
		logx.Panicf("[DB]-[%s] 数据库配置信息获取失败", db.Name)
		return
	}
	config := &gorm.Config{}
	if db.Debug {
		config.Logger = logger.Default.LogMode(logger.Info)
	}
	str := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", db.User, db.Password, db.Host, db.Port, db.Name) + "?charset=utf8mb4&parseTime=true&loc=Local"
	if db.DisablePrepared {
		str = str + "&interpolateParams=true"
	}
	for orm, err = gorm.Open(mysql.Open(str), config); err != nil; {
		logx.Errorf("[DB]-[%v] 连接异常:%v，正在重试: %v", db.Name, err, str)
		time.Sleep(5 * time.Second)
		orm, err = gorm.Open(mysql.Open(str), config)
	}
	dbs[db.Name] = orm
}

// GetORM return default *gorm.DB
func GetORM() *gorm.DB {
	m, ok := dbs[defaultDBName]
	if !ok {
		logx.Panic("[DB] 没有初始化mysql连接，请参考github.com/zituocn/gow/lib/orm 的配置")
	}
	return m
}

// GetORMByName get orm by name
func GetORMByName(name string) *gorm.DB {
	m, ok := dbs[name]
	if !ok {
		logx.Panic("[DB] 没有初始化mysql连接，请参考github.com/zituocn/gow/lib/orm 的配置")
	}
	return m
}
