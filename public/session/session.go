package session

import (
	"errors"
	"log"
	"net"
)

type Session struct {
	Ch   chan []byte
	C    net.Conn
	SsiD string
}

var ssinfo map[string]*Session

//电梯的用电梯的id，调度统一用pss.
func init() {
	ssinfo = make(map[string]*Session)
}

//实现类似注册的功能，每条连接都会新生成一个session。
//想向其他的连接里面写值，就获取他的session，然后向channel里面写值
func NewSession(c net.Conn, ch chan []byte) *Session {
	s := &Session{}
	s.Ch = ch
	s.C = c
	return s
}

func AddEle(ssID string, session *Session) {
	session.SsiD = ssID
	ssinfo[ssID] = session
	//ssinfo[elID].ssiD=elID
}

func Delete(ssID string) {
	if _, ok := ssinfo[ssID]; !ok {
		log.Printf("删除连接失败")
	}
	log.Printf("删除连接成功%s", ssID)
	delete(ssinfo, ssID)
}

func Send(ssID string, b []byte) error {
	_, ok := ssinfo[ssID]
	if !ok {
		log.Printf("此连接不存在%s", ssID)
		return errors.New("此连接不存在")
	}
	if ssID == "pss" {
		log.Printf("向调度发送了一条信息,调度IP为%s", ssinfo[ssID].C.RemoteAddr().String())
	} else {
		log.Printf("向电梯%s发送了一条信息,电梯IP为%s", ssID, ssinfo[ssID].C.RemoteAddr().String())
	}
	ssinfo[ssID].Ch <- b
	return nil
}
