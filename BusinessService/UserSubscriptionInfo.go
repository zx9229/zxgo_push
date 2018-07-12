package businessservice

import (
	"strings"
)

//UserSubscriptionInfo 用户订阅管理器
type UserSubscriptionInfo struct {
	AllID       map[int64]bool  //接收指定UserID的消息.
	AllCategory map[string]bool //接收指定类别的消息.
}

//new_UserSubscriptionInfo omit
func new_UserSubscriptionInfo() *UserSubscriptionInfo {
	curData := new(UserSubscriptionInfo)
	curData.AllID = make(map[int64]bool)
	curData.AllCategory = make(map[string]bool)
	return curData
}

//SubUser omit
func (thls *UserSubscriptionInfo) SubUser(userID int64) {
	thls.AllID[userID] = true
}

//UnsubUser omit
func (thls *UserSubscriptionInfo) UnsubUser(userID int64) {
	delete(thls.AllID, userID)
}

//SubCategory omit
func (thls *UserSubscriptionInfo) SubCategory(data string) {
	if data == emptyString {
		return
	}
	dataPrefix := data + string(CategorySep)
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
func (thls *UserSubscriptionInfo) UnsubCategory(data string) {
	delete(thls.AllCategory, data)
}

//ShouldSend 应当发送数据给这个用户,临时函数,待优化订阅管理器
func (thls *UserSubscriptionInfo) ShouldSend(userID int64, data string) bool {
	if _, ok := thls.AllID[userID]; ok {
		return true
	}

	for key := range thls.AllCategory {
		if strings.HasPrefix(data, key) {
			keyLen := len(key)
			if len(data) == keyLen {
				return true
			} else if data[keyLen] == CategorySep {
				//key  => AA|B
				//key  => AA|BB
				//data => AA|BB|CC
				return true
			}
		}
	}

	return false
}
