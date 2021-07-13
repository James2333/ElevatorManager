package eleManager

import (
	"EleManager/public/code"
	"EleManager/public/elevator"
	"EleManager/public/packet"
	"EleManager/public/session"
	"EleManager/public/task"
	"log"
)

//请求电梯到起点楼层
func ReqElevatorToStart(ts *task.Task) {
	log.Printf("正在请求电梯去往起点楼层")
	b := packet.Packet(ts, code.ELE_TO_START)
	session.Send(ts.ElevatorID, b)
}

//请求电梯到终点楼层
func ReqElevatorToEnd(ts *task.Task) {
	b := packet.Packet(ts, code.ELE_TO_END)
	session.Send(ts.ElevatorID, b)
	return
}

//向调度发送电梯已经抵达起点楼层
func ReqElevatorArriveStart(ts *task.Task) {
	//time.Sleep(time.Second * 5)
	//b := packet.Packet("", code.ERROR)
	//session.Send("pss", b)
	//b = packet.Packet("", code.ERROR)
	//session.Send("pss", b)
	//b = packet.Packet("", code.ERROR)
	//session.Send("pss", b)
	//b = packet.Packet("", code.ERROR)
	//session.Send("pss", b)
	b := packet.Packet(ts, code.ROBOT_START)
	session.Send("pss", b)
	log.Printf("发送任务ID为%s，已经抵达起点楼层%d.", ts.TaskID, ts.Start)
	return
}

//向调度发送电梯已经抵达终点楼层
func ReqElevatorArriveEnd(ts *task.Task) {
	b := packet.Packet(ts, code.ROBOT_END)
	session.Send("pss", b)
	log.Printf("发送任务%s已经抵达终点楼层%d.", ts.TaskID, ts.End)
	return
}

//向电梯发送任务结束，进入空闲状态
func ReqElevatorTaskEnd(ts *task.Task) {
	//更新电梯信息，防止心跳不及时。
	elevator.Els[ts.ElevatorID].IsInFloor = false
	elevator.Els[ts.ElevatorID].CurrentState = "0"
	elevator.Els[ts.ElevatorID].Floor = ts.End
	//s := els[tasks[taskid].ElevatorID].Sess
	task.Ts.Delete(ts.TaskID)
	log.Printf("任务%s结束，电梯%s进入空闲状态。",ts.TaskID,ts.ElevatorID)
	b := packet.Packet("任务结束，进入空闲状态", code.ELE_TO_FREE)
	session.Send(ts.ElevatorID, b)
	return
}
