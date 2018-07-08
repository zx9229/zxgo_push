package businessservice

import (
	"time"
)

//TotalUserManager 所有用户的管理器
type TotalUserManager struct {
	lastUserID int64 //最后创建的用户的ID
	allUser    []*UserSummary
}

func New_TotalUserManager() *TotalUserManager {
	curData := new(TotalUserManager)
	curData.lastUserID = 0
	curData.allUser = make([]*UserSummary, 0)
	return curData
}

//IsValidPassword omit
func (thls *TotalUserManager) IsValidPassword(password string) bool {
	if len(password) == 0 {
		return false
	}
	for _, b := range password {
		if (' ' <= b && b <= '~') == false {
			return false
		}
	}
	return true
}

//CreateUser omit
func (thls *TotalUserManager) CreateUser(password string) int64 {
	if !thls.IsValidPassword(password) {
		return -1
	}
	if int64(len(thls.allUser)) != thls.lastUserID { //违反了逻辑规则.
		return -1
	}
	for idx, us := range thls.allUser {
		if us == nil {
			continue
		}
		if int64(idx) != us.Base.UserID { //违反了逻辑规则.
			return -1
		}
	}
	userData := New_UserSummary()
	userData.Base.UserID = thls.lastUserID + 1
	userData.Base.Password = password
	userData.Base.CreateTime = time.Now()
	userData.Base.UpdateTime = userData.Base.CreateTime
	userData.Base.Memo = ""

	thls.lastUserID = userData.Base.UserID
	thls.allUser = append(thls.allUser, userData)

	for idx, us := range thls.allUser {
		if us == nil {
			continue
		}
		if int64(idx) != us.Base.UserID {
			panic("违反了逻辑规则")
		}
	}

	return userData.Base.UserID
}

//DeleteUser omit
func (thls *TotalUserManager) DeleteUser(userID int64, password string) bool {
	if (0 < userID && userID <= int64(len(thls.allUser))) == false {
		return false
	}
	userData := thls.allUser[userID-1]
	if userData == nil {
		return false
	}
	if userData.Base.UserID != userID { //违反了逻辑规则.
		return false
	}
	if userData.Base.Password != password {
		return false
	}
	for _, state := range userData.State {
		if state.conn == nil {
			continue
		}
		//TODO:关闭之前,发送注销的消息.
		state.conn.Close()
	}

	thls.allUser[userID-1] = nil
	//TODO:更新到数据库中.

	return true
}
