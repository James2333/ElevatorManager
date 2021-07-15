package eleManager

import (
	"EleManager/public/code"
	"EleManager/public/elevator"
	"EleManager/public/packet"
	"EleManager/public/session"
	"EleManager/public/task"
	"encoding/json"
	"fmt"
	"log"
)

func ParseCode(cod uint16, s *session.Session) {
	log.Printf("收到%s的请求，请求头%d.", s.C.RemoteAddr().String(), cod)
	switch cod {
	//
	case code.UPDATE_ELE:
		ReplyUpdateElevator(s)
	case code.CHOOSE_ELE:
		ReplyRightElevator(s)
	case code.ARRIVED_START:
		ReplyElevatorArriveStart(s)
	case code.ARRIVED_END:
		ReplyElevatorArriveEnd(s)
	case code.ROBOT_In_Floor:
		ReplyRobotInFloor(s)
	case code.ROBOT_OUT_Floor:
		ReplyRobotOutFloor(s)
	case code.REGISTER_SCHEDULE:
		ReplyRegisterSchedule(s)
	case code.REGISTER_ELE:
		ReplyRegisterElevator(s)
	default:
		ReplyError(s)
	}
}

//注册调度信息
func ReplyRegisterSchedule(s *session.Session) {
	log.Println("触发了一次注册调度操作")
	session.AddEle("pss",s)
	b := packet.Packet("注册电梯调度信息成功", code.REGISTER_SCHEDULE)
	session.Send("pss", b)
}

//注册电梯信息
func ReplyRegisterElevator(s *session.Session) {
	log.Println("触发了一次注册电梯操作")
	q, err := packet.UnPacket(s.C)
	if err != nil {
		log.Println(err)
	}
	var ele elevator.Elevator
	err = json.Unmarshal(q.Content, &ele)
	if err != nil {
		log.Println("json unmarshal faild:", err)
	}
	log.Println("客户端传来的信息", ele)
	err = elevator.Els.Update(&ele)
	if err != nil {
		b := packet.Packet("", code.ERROR)
		session.Send(ele.ElevatorId,b)
		return
	}
	session.AddEle(ele.ElevatorId, s)
	b := packet.Packet("注册电梯信息成功", code.REGISTER_ELE)
	session.Send(ele.ElevatorId,b)
}

//更新电梯信息
func ReplyUpdateElevator(s *session.Session) {
	log.Println("触发了一次更新操作")
	q, err := packet.UnPacket(s.C)
	if err != nil {
		log.Println(err)
	}
	var ele elevator.Elevator
	err = json.Unmarshal(q.Content, &ele)
	if err != nil {
		log.Println("json unmarshal faild:", err)
	}
	log.Println("客户端传来的信息", ele)
	err = elevator.Els.Update(&ele)
	if err != nil {
		b := packet.Packet("", code.ERROR)
		session.Send(ele.ElevatorId,b)
		return
	}
	b := packet.Packet("更新电梯信息成功", code.UPDATE_ELE)
	session.Send(ele.ElevatorId,b)
}

//调度请求最优电梯
func ReplyRightElevator(s *session.Session) {
	q, err := packet.UnPacket(s.C)
	if err != nil {
		log.Println(err)
	}
	ts := new(task.Task)
	err = json.Unmarshal(q.Content, ts)
	if err != nil {
		log.Println("json unmarshal faild:", err)
	}

	log.Printf("客户端传来的信息ElevatorID:%s,TaskID:%s,Start:%d,End:%d", ts.ElevatorID, ts.TaskID, ts.Start, ts.End)
	elID, err := elevator.Els.RightElevator(ts.Start)
	if err != nil {
		b := packet.Packet("暂无空闲电梯", code.CHOOSE_ELE)
		session.Send("pss",b)
		return
	}
	ts.ElevatorID = elID
	task.Ts.AddTask(ts)
	log.Printf("任务选用电梯之后ElevatorID:%s,TaskID:%s,Start:%d,End:%d", ts.ElevatorID, ts.TaskID, ts.Start, ts.ElevatorID)
	elevator.Els[elID].CurrentState = "1" //电梯变为繁忙状态 这个后面任务结束才能更新成空闲。
	ReqElevatorToStart(ts)
	b := packet.Packet(fmt.Sprintf("选取了ID为%s的电梯", ts.ElevatorID), code.CHOOSE_ELE)
	session.Send("pss",b)
}

//电梯抵达起点楼层
func ReplyElevatorArriveStart(s *session.Session) {
	q, err := packet.UnPacket(s.C)
	if err != nil {
		log.Println(err)
	}
	ts := new(task.Task)
	err = json.Unmarshal(q.Content, ts)
	if err != nil {
		log.Println("json unmarshal faild:", err)
	}

	ts, err = task.Ts.Get(ts.TaskID)
	if err != nil {
		log.Printf("没有找到任务%s", ts.TaskID)
		return
	}
	log.Printf("任务%s电梯抵达起点楼层%d.", ts.TaskID, ts.Start)
	ReqElevatorArriveStart(ts)
	return
}

//机器人是否进电梯里面
func ReplyRobotInFloor(s *session.Session) {
	q, err := packet.UnPacket(s.C)
	if err != nil {
		log.Println(err)
	}
	ts := new(task.Task)
	_ = json.Unmarshal(q.Content, ts)

	elevator.Els[ts.ElevatorID].IsInFloor = true
	ReqElevatorToEnd(ts) //请求电梯到终点楼层
	return
}

//机器人是否已经出电梯里面
func ReplyRobotOutFloor(s *session.Session) {
	q, err := packet.UnPacket(s.C)
	if err != nil {
		log.Println(err)
	}
	ts := new(task.Task)
	_ = json.Unmarshal(q.Content, ts)

	elevator.Els[ts.ElevatorID].IsInFloor = false
	ReqElevatorTaskEnd(ts) //此次任务结束
	return
}

//电梯抵达终点楼层
func ReplyElevatorArriveEnd(s *session.Session) {
	q, err := packet.UnPacket(s.C)
	if err != nil {
		log.Println(err)
	}
	ts := new(task.Task)
	_ = json.Unmarshal(q.Content, ts)

	ReqElevatorArriveEnd(ts) //向调度发送电梯已经抵达终点楼层
	return
}

//错误返回
func ReplyError(s *session.Session) {
	//buffer := packet.Packet(errors.New("404 not found"), ERROR)
	//b:=make([]byte,1024)
	////回收垃圾包
	//_,err:=s.C.Read(b)
	//if err != nil {
	//	log.Println("回收垃圾包出错")
	//	return
	//}
	log.Printf("收到%s无效的请求", s.C.RemoteAddr().String())
	return
}
