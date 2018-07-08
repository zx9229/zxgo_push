package businessservice

import (
	"strings"
)

//CagegorySep 父类别/子类别,之间的分隔符
const CagegorySep = byte('|')

//SubscribeBaseInfo 订阅的基本信息
type SubscribeBaseInfo struct {
	AllCategory map[string]bool
}

//New_SubscribeBaseInfo omit
func New_SubscribeBaseInfo() *SubscribeBaseInfo {
	curData := new(SubscribeBaseInfo)
	curData.AllCategory = make(map[string]bool)
	return curData
}

//AddCategory 添加类别
func (thls *SubscribeBaseInfo) AddCategory(data string) bool {
	SEP := string(CagegorySep)

	for _, field := range strings.Split(data, SEP) {
		if len(field) == 0 { //合法性校验
			return false
		}
	}

	parent := map[string]bool{} //data的父亲的集合

	for key := range thls.AllCategory {
		if key == data {
			return false
		}
		if strings.HasPrefix(key, data+SEP) {
			//key  => A|B|C
			//data => A|B
			return false
		}
		if strings.HasPrefix(data, key+SEP) {
			//key  => A|B
			//data => A|B|C
			parent[key] = true
		}
	}

	for key := range parent {
		delete(thls.AllCategory, key)
	}
	thls.AllCategory[data] = true

	return true
}

//DelCategory 删除类别和所有的子类别
func (thls *SubscribeBaseInfo) DelCategory(data string) {
	dataPrefix := data + string(CagegorySep)
	dictDel := map[string]bool{}
	for key := range thls.AllCategory {
		if key == data {
			dictDel[key] = true
		}
		if strings.HasPrefix(key, dataPrefix) {
			dictDel[key] = true
		}
	}
	for key := range dictDel {
		delete(thls.AllCategory, key)
	}
}

//IsRegistered 送进去的是不是一个注册过的类别
func (thls *SubscribeBaseInfo) IsRegistered(data string) bool {
	dataPrefix := data + string(CagegorySep)
	for key := range thls.AllCategory {
		if key == data {
			return true
		}
		if strings.HasPrefix(key, dataPrefix) {
			//key  => A|B|C
			//data => A|B
			return true
		}
	}
	return false
}
