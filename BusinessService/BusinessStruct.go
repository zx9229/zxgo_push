package businessservice

//go get -u -v golang.org/x/net/websocket
//响应消息,要带过来请求消息的结构体.
//校验用户名密码=>在内存数据库查count(id=?&&pwd=?)
//用[当前接收到哪个推送序号了]作为业务心跳消息.
//把所有用户数据都加载到内存中

import (
	"time"

	"github.com/zx9229/zxgo_push/TxStruct"

	wsconnectionmanager "github.com/zx9229/zxgo_push/WSConnectionManager"
)

const (
	LoginTypeNA     = iota //(登录方式)0,默认值,无效值
	LoginTypeWeb           //(登录方式)网页
	LoginTypeMobile        //(登录方式)手机
	LoginTypePC            //(登录方式)电脑
	LoginTypeEND           //(登录方式)结束值,最大值的下一个
)

//CacheData 所有的缓存信息
type CacheData struct {
	LastUserID  int64           //最后一个注册的用户ID
	LastPushID  int64           //最后一个推送消息的序号
	AllCategory map[string]bool //总共有哪些种类
	AllUser     []*UserSummary
}

func New_CacheData() *CacheData {
	curData := new(CacheData)
	return curData
}

//UserSummary 用户的汇总信息
type UserSummary struct {
	Base    UserBaseInfo
	State   []*LoginInfo
	SubInfo *SubscribeUserInfo
}

//UserBaseInfo 用户的基础信息
type UserBaseInfo struct {
	UserID     int64
	Password   string
	Memo       string    //备注.
	CreateTime time.Time //创建时刻.
	UpdateTime time.Time //更新时刻.
}

//LoginInfo 用户的登录信息
type LoginInfo struct {
	conn       *wsconnectionmanager.WSConnection //有值,表示在线.(这个字段不往外导出)
	LoginType  int                               //电脑登录,网页登录,APP登录,等.
	MaxPushID  int64                             //推送给它的最大的推送序号
	LastRecvID int64                             //它上报的自己接收到的最后一个序号
}

////////////////////////////////////////////////////////////////

//ReportData omit
type ReportData struct {
	ID   int64     `xorm:"notnull pk"` //数据库的递增序号.
	Time time.Time //插入数据库的时刻.
	//
	UserID      int64
	UserTagID   int64     //用户在维护这一条数据时,可能会给它贴上一个序号,这个序号的值
	UserTagTime time.Time //用户在维护这一条数据时,可能会给它贴上一个时刻,这个时刻的值
	AttachedI   int64     //随附的一个数字(或许你上报的消息,想表示(成功/失败)或(某某的个数)呢.)
	Message     string    //消息的详情(当然,这里面也可以是json,这样就不需要上面的字段了)
	Category    string    //消息所属的类别
}

func convertToReportData(src *txstruct.ReportReq) *ReportData {
	dst := new(ReportData)
	dst.ID = 0
	dst.Time = time.Now()
	dst.UserID = src.UserID
	dst.UserTagID = src.UserTagID
	dst.UserTagTime = src.UserTagTime
	dst.AttachedI = src.AttachedI
	dst.Message = src.Message
	dst.Category = src.Category
	return dst
}
