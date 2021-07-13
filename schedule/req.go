package schedule

import (
	"EleManager/public/code"
	"EleManager/public/packet"
	"EleManager/public/session"
	"EleManager/public/task"
	"encoding/json"
	"log"
	"time"
)

func ParseCodeScheduling(cod uint16, s *session.Session) {
	log.Printf("收到%s的请求，请求头%d.",s.C.RemoteAddr().String(), cod)
	switch cod {
	case code.ROBOT_START:
		ReqRobotInStart(s)
	case code.ROBOT_END:
		ReqRobotOutEnd(s)
	case code.REGISTER_SCHEDULE:
		ReplyRegister(s)
	case code.CHOOSE_ELE:
		ReplyTask(s)
	default:
		ReplyError(s)
	}
}

func ReqRegister(s *session.Session)  {
	b := packet.Packet("", code.REGISTER_SCHEDULE)
	s.Ch<-b
}

func ReqRobotInStart(s *session.Session) {
	q, err := packet.UnPacket(s.C)
	if err != nil {
		log.Println(err)
	}
	ts:=new(task.Task)
	_ = json.Unmarshal(q.Content,ts)
	log.Printf("收到任务进入电梯：taskID:%s,ElevatorID:%s,Start:%d,End:%d",ts.TaskID,ts.ElevatorID,ts.Start,ts.End)
	time.Sleep(time.Second * 5)
	log.Println("机器人已经进入电梯了")
	b := packet.Packet(ts, code.ROBOT_In_Floor)
	s.Ch<-b
}


func ReqRobotOutEnd(s *session.Session) {
	q, err := packet.UnPacket(s.C)
	if err != nil {
		log.Println(err)
	}
	ts:=new(task.Task)
	_ = json.Unmarshal(q.Content,ts)
	log.Printf("收到任务驶出电梯：taskID:%s,ElevatorID:%s,Start:%d,End:%d",ts.TaskID,ts.ElevatorID,ts.Start,ts.End)
	time.Sleep(time.Second * 5)
	log.Println("机器人已经出电梯了")
	b := packet.Packet(ts, code.ROBOT_OUT_Floor)
	s.Ch<-b
}

//错误返回
func ReplyError(s *session.Session) {
	//buffer := packet.Packet(errors.New("404 not found"), ERROR)
	log.Printf("收到%s无效的请求",s.C.RemoteAddr().String())
	return
}
