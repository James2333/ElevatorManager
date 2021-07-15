package elevator

import (
	"EleManager/public/code"
	"EleManager/public/elevator"
	"EleManager/public/packet"
	"EleManager/public/session"
	"EleManager/public/task"
	"encoding/json"
	"log"
	"time"
)

func ParseCodeElevator(cod uint16, s *session.Session) {
	log.Printf("收到%s的请求，请求头%d.",s.C.RemoteAddr().String(),cod)
	switch cod {
	case code.ELE_TO_START:
		ReqElevatorArriveStart(s)
	case code.ELE_TO_END:
		ReqElevatorArriveEnd(s)
	case code.UPDATE_ELE:
		ReplyUpdateEle(s)
	case code.ELE_TO_FREE:
		ReplyTaskEnd(s)
	case code.REGISTER_ELE:
		ReplyRegister(s)
	default:
		ReplyError(s)
	}
}
//模拟更新电梯信息
func ReqRegisterElevator(s *session.Session)  {
	ele:=elevator.Elevator{
		ElevatorId:   "111",
		Floor:        1,
		State:        "0",
	}
	b:=packet.Packet(ele,code.REGISTER_ELE)
	s.Ch<-b
}

//电梯去到起点楼层
func ReqElevatorArriveStart(s *session.Session){
	q, err := packet.UnPacket(s.C)
	if err != nil {
		log.Println(err)
	}
	ts:=new(task.Task)
	_ = json.Unmarshal(q.Content,ts)
	log.Printf("收到任务驶向起点楼层：taskID:%s,ElevatorID:%s,Start:%d,End:%d",ts.TaskID,ts.ElevatorID,ts.Start,ts.End)
	log.Println("电梯驶向起点楼层:",ts.Start)
	time.Sleep(time.Second*5)
	log.Println("电梯抵达起点楼层:",ts.Start)
	b:=packet.Packet(ts,code.ARRIVED_START)
	s.Ch<-b
}
//电梯去到终点楼层
func ReqElevatorArriveEnd(s *session.Session){
	q, err := packet.UnPacket(s.C)
	if err != nil {
		log.Println(err)
	}
	ts:=new(task.Task)
	_ = json.Unmarshal(q.Content,ts)
	log.Printf("收到任务驶向终点楼层：taskID:%s,ElevatorID:%s,Start:%d,End:%d",ts.TaskID,ts.ElevatorID,ts.Start,ts.End)
	log.Println("电梯驶向终点楼层",ts.End)
	time.Sleep(time.Second*5)
	log.Println("电梯抵达终点楼层",ts.End)
	b:=packet.Packet(ts,code.ARRIVED_END)
	s.Ch<-b
}


//错误返回
func ReplyError(s *session.Session) {
	//buffer := packet.Packet(errors.New("404 not found"), ERROR)
	b:=make([]byte,1024)
	//回收垃圾包
	_,err:=s.C.Read(b)
	if err != nil {
		log.Println("回收垃圾包出错")
		return
	}
	log.Printf("收到%s无效的请求", s.C.RemoteAddr().String())
	return
}
