package businessservice

import (
	"time"

	wscmanager "github.com/zx9229/zxgo_push/WSCManager"
)

const (
	LoginTypeDEFAULT = iota //(登录方式)0,默认值,无效值
	LoginTypeWeb            //(登录方式)网页
	LoginTypeMobile         //(登录方式)手机
	LoginTypePC             //(登录方式)电脑
	LoginTypeEND            //(登录方式)结束值,最大值的下一个
)

//UserTempData 附加到socket上的用户临时数据
type UserTempData struct {
	tInfo *UserTotalInfo
	sInfo *UserSessionInfo
}

//UserBaseInfo 用户的基本信息
type UserBaseInfo struct {
	UserID     int64
	Password   string
	Memo       string    //备注.
	CreateTime time.Time //创建时刻.
	UpdateTime time.Time //更新时刻.
}

//UserSessionInfo 用户的会话所持有的信息
type UserSessionInfo struct {
	conn       *wscmanager.WSConnection //有值,表示在线.(这个字段不往外导出)
	LoginType  int                      //电脑登录,网页登录,APP登录,等.
	MaxPushID  int64                    //推送给它的最大的推送序号
	LastRecvID int64                    //它上报的自己接收到的最后一个序号
}

//UserTotalInfo 用户的汇总信息
type UserTotalInfo struct {
	Base    UserBaseInfo          //用户基础信息
	State   []*UserSessionInfo    //用户会话信息
	SubInfo *UserSubscriptionInfo //用户订阅信息
}

//New_UserTotalInfo omit
func New_UserTotalInfo() *UserTotalInfo {
	curData := new(UserTotalInfo)
	//
	curData.State = make([]*UserSessionInfo, 0)
	for i := LoginTypeDEFAULT + 1; i < LoginTypeEND; i++ {
		stateInfo := new(UserSessionInfo)
		stateInfo.LoginType = i
		curData.State = append(curData.State, stateInfo)
	}
	curData.SubInfo = New_UserSubscriptionInfo()
	//
	return curData
}
