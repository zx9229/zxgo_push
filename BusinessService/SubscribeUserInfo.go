package businessservice

import (
	"strings"
)

//SubscribeUserInfo 用户订阅管理器
type SubscribeUserInfo struct {
	UserID      int64
	AllID       map[int64]bool  //接收指定UserID的消息.
	AllCategory map[string]bool //接收指定类别的消息.
}

//New_SubscribeUserInfo omit
func New_SubscribeUserInfo() *SubscribeUserInfo {
	curData := new(SubscribeUserInfo)
	curData.AllID = make(map[int64]bool)
	curData.AllCategory = make(map[string]bool)
	return curData
}

//SubUser omit
func (thls *SubscribeUserInfo) SubUser(userID int64) {
	thls.AllID[userID] = true
}

//UnsubUser omit
func (thls *SubscribeUserInfo) UnsubUser(userID int64) {
	delete(thls.AllID, userID)
}

//SubCategory omit
func (thls *SubscribeUserInfo) SubCategory(data string) {
	dataPrefix := data + string(CagegorySep)
	dictDel := map[string]bool{}
	for key := range thls.AllCategory {
		if key == data {
			return
		}
		if strings.HasPrefix(key, dataPrefix) {
			dictDel[key] = true
		}
	}
	for key := range dictDel {
		delete(thls.AllCategory, key)
	}
	thls.AllCategory[data] = true
}

//UnsubCategory omit
func (thls *SubscribeUserInfo) UnsubCategory(data string) {
	delete(thls.AllCategory, data)
}

//ShouldSend 应当发送数据给这个用户,临时函数,待优化订阅管理器
func (thls *SubscribeUserInfo) ShouldSend(userID int64, data string) bool {
	if _, ok := thls.AllID[userID]; ok {
		return true
	}

	for key := range thls.AllCategory {
		if strings.HasPrefix(data, key) {
			keyLen := len(key)
			if len(data) == keyLen {
				return true
			} else if data[keyLen] == CagegorySep {
				//key  => AA|B
				//key  => AA|BB
				//data => AA|BB|CC
				return true
			}
		}
	}

	return false
}
