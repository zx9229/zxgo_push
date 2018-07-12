package businessservice

import (
	"time"

	txstruct "github.com/zx9229/zxgo_push/TxStruct"
	wscmanager "github.com/zx9229/zxgo_push/WSCManager"
)

//UserInfoManager 用户信息管理器
type UserInfoManager struct {
	lastUserID int64 //最后创建的用户的ID
	allUser    []*UserTotalInfo
}

//new_UserInfoManager omit
func new_UserInfoManager() *UserInfoManager {
	curData := new(UserInfoManager)
	curData.lastUserID = 0
	curData.allUser = make([]*UserTotalInfo, 0)
	return curData
}

//WellFormedPassword 格式正确的密码
func (thls *UserInfoManager) WellFormedPassword(password string) bool {
	if password == emptyString {
		return false
	}
	for _, b := range password {
		if (' ' <= b && b <= '~') == false {
			return false
		}
	}
	return true
}

//PasswordOk omit
func (thls *UserInfoManager) PasswordOk(userID int64, password string) bool {
	if (0 < userID && userID <= int64(len(thls.allUser))) == false {
		return false
	}
	userInfo := thls.allUser[userID-1]
	if userInfo == nil {
		return false
	}
	if userInfo.Base.UserID != userID { //违反了逻辑规则.
		return false
	}
	if userInfo.Base.Password != password {
		return false
	}
	return true
}

//GetUserTotalInfo omit
func (thls *UserInfoManager) GetUserTotalInfo(userID int64) *UserTotalInfo {
	return thls.allUser[userID-1]
}

//CreateUser omit
func (thls *UserInfoManager) CreateUser(password string) int64 {
	var userID int64
	if !thls.WellFormedPassword(password) {
		return userID
	}
	if int64(len(thls.allUser)) != thls.lastUserID { //违反了逻辑规则.
		return userID
	}
	for idx, uInfo := range thls.allUser {
		if uInfo == nil {
			continue
		}
		if int64(idx+1) != uInfo.Base.UserID { //违反了逻辑规则.
			return userID
		}
	}

	newUserInfo := new_UserTotalInfo()
	newUserInfo.Base.UserID = thls.lastUserID + 1
	newUserInfo.Base.Password = password
	newUserInfo.Base.CreateTime = time.Now()
	newUserInfo.Base.UpdateTime = newUserInfo.Base.CreateTime
	newUserInfo.Base.Memo = ""

	thls.lastUserID = newUserInfo.Base.UserID
	thls.allUser = append(thls.allUser, newUserInfo)

	for idx, uInfo := range thls.allUser {
		if uInfo == nil {
			continue
		}
		if int64(idx+1) != uInfo.Base.UserID {
			panic("违反了逻辑规则")
		}
	}

	return newUserInfo.Base.UserID
}

//DeleteUser omit
func (thls *UserInfoManager) DeleteUser(userID int64, password string) bool {
	if !thls.PasswordOk(userID, password) {
		return false
	}
	for _, sInfo := range thls.allUser[userID-1].State {
		if sInfo.conn == nil {
			continue
		}
		//TODO:关闭之前,发送注销的消息.
		sInfo.conn.Close()
	}

	thls.allUser[userID-1] = nil
	//TODO:更新到数据库中.

	return true
}

//LoginUser 用[string]代替[error]
func (thls *UserInfoManager) LoginUser(conn *wscmanager.WSConnection, req *txstruct.LoginReq) string {
	var msg string
	for range "1" {
		if conn.ExtraData != nil {
			msg = ErrConnectionIsLoggedOn
			break
		}
		if !thls.PasswordOk(req.UserID, req.Password) {
			msg = ErrIncorrectPassword
			break
		}
		if (LoginTypeDEFAULT < req.Way && req.Way < LoginTypeEND) == false {
			msg = ErrInvalidLoginType
			break
		}
		sInfo := thls.allUser[req.UserID-1].State[req.Way-1]
		if sInfo.conn != nil && !req.ForceLogin {
			msg = ErrUserHasLoggedIn
			break
		}
		if sInfo.conn != nil {
			//TODO:关闭之前,发送一个"您被踢下线了"的消息
			sInfo.conn.Close()
		}

		sInfo.conn = conn
		conn.ExtraData = &UserTempData{tInfo: thls.allUser[req.UserID-1], sInfo: sInfo}

		if 0 <= req.LastMsgID {
			sInfo.LastRecvID = req.LastMsgID
			//TODO:哪些消息尚未推送,把它们推送过去
		}
	}
	return msg
}

//PushData omit
func (thls *UserInfoManager) PushData(data *txstruct.ReportData) {
	jsonStr := data.ToJSON(true)
	for _, tInfo := range thls.allUser {
		if tInfo == nil {
			continue
		}
		if !tInfo.SubInfo.ShouldSend(data.UserID, data.Category) {
			continue
		}
		for _, sInfo := range tInfo.State {
			if sInfo.conn == nil {
				continue
			}
			sInfo.conn.Send(jsonStr)
		}
	}
}
