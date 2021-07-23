package cache

// ICache 缓存接口
type ICache interface {

	//KeyName key
	KeyName() string

	// PrimaryKey 返回主键名称，
	PrimaryKey(model interface{}) string

	// GetAllData 获取所有数据
	GetAllData() (data interface{}, err error)
}
