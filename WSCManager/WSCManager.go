package wscmanager

import (
	"errors"
	"fmt"
	"log"

	"golang.org/x/net/websocket"
)

var ErrClosedByUser = errors.New("Closed By User")

//WSConnection websocket的包装器.
type WSConnection struct {
	ws           *websocket.Conn
	closedByUser bool        //(用户调用了close函数)
	ExtraData    interface{} //附加信息
}

//Send omit
func (thls *WSConnection) Send(v interface{}) error {
	return websocket.Message.Send(thls.ws, v)
}

//Close omit
func (thls *WSConnection) Close() (err error) {
	thls.closedByUser = true
	err = thls.ws.Close()
	return
}

//WSConnectionManager websocket的连接管理器.
type WSConnectionManager struct {
	connSet        map[*WSConnection]bool //所有连接的集合.
	CbConnected    func(conn *WSConnection)
	CbDisconnected func(conn *WSConnection, err error)
	CbReceive      func(conn *WSConnection, bytes []byte)
}

//New_WSConnectionManager omit
func New_WSConnectionManager() *WSConnectionManager {
	curData := new(WSConnectionManager)
	curData.connSet = make(map[*WSConnection]bool)
	return curData
}

//HandleWebsocket omit
func (thls *WSConnectionManager) HandleWebsocket(wsConn *websocket.Conn) {
	wsc := &WSConnection{ws: wsConn, closedByUser: false, ExtraData: nil}

	var err error

	defer func() {
		if thls.CbDisconnected != nil {
			if wsc.closedByUser {
				err = ErrClosedByUser
			}
			thls.CbDisconnected(wsc, err)
		}

		delete(thls.connSet, wsc)

		if !wsc.closedByUser {
			if err = wsc.ws.Close(); err != nil {
				log.Println(fmt.Sprintf("ws=%p,调用Close失败,err=%v", wsc, err))
			}
		}
	}()

	thls.connSet[wsc] = true

	if thls.CbConnected != nil {
		thls.CbConnected(wsc)
	}

	var recvRawMessage []byte
	for {
		recvRawMessage = nil
		if err = websocket.Message.Receive(wsc.ws, &recvRawMessage); err != nil {
			break
		}
		if thls.CbReceive != nil {
			thls.CbReceive(wsc, recvRawMessage)
		}
	}
}

//Example_CbConnected omit
func Example_CbConnected(conn *WSConnection) {
	log.Println(fmt.Sprintf("[   Connected][%p]LocalAddr=%v,RemoteAddr=%v", conn, conn.ws.LocalAddr(), conn.ws.RemoteAddr()))
}

//Example_CbDisconnected omit
func Example_CbDisconnected(conn *WSConnection, err error) {
	log.Println(fmt.Sprintf("[Disconnected][%p]err=%v", conn, err))
}

//Example_CbReceive omit
func Example_CbReceive(conn *WSConnection, bytes []byte) {
	log.Println(fmt.Sprintf("[     Receive][%p]data=%v", conn, string(bytes)))
}

/* 一个测试的例子
func main() {
	shs := SimpleHttpServer.New_SimpleHttpServer("localhost:8080")
	mngr := wsconnectionmanager.New_WSConnectionManager()
	shs.GetHttpServeMux().Handle("/websocket", websocket.Handler(mngr.HandleWebsocket))
	mngr.CbConnected = wsconnectionmanager.Example_CbConnected
	mngr.CbDisconnected = wsconnectionmanager.Example_CbDisconnected
	mngr.CbReceive = wsconnectionmanager.Example_CbReceive
	shs.Run()
}
*/
