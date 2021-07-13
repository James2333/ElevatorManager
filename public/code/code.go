package code

//初次解包先解开byte[0:2]前两个字节，转换成code，之后switch 这个code 跳转进不同的方法。
//电梯管理需要向调度发送的信息:
//向调度发送电梯已经抵达起点楼层 ROBOT_START=2004
//向调度发送电梯已经抵达终点楼层 ROBOT_END=2005
//电梯管理需要向调度和电梯发送的信息:
//请求电梯到起点楼层 ELE_TO_START=2001
//请求电剃到终点楼层 ELE_TO_END=2002
//向电梯发送任务结束，进入空闲状态 ELE_TO_FREE=2003
//电梯管理需要向调度接受的信息:
//CHOOSE_ELE    = 1002	调度请求最优电梯
//ROBOT_In_Floor     = 1005 机器人是否进电梯里面
//ROBOT_OUT_Floor     = 1006	机器人是否已经出电梯里面
//电梯管理需要向电梯接受的信息:
//UPDATE_ELE    = 1001	更新电梯信息
//ARRIVED_START = 1003	电梯抵达起点楼层
//ARRIVED_END   = 1004	电梯抵达终点楼层
//请求一个电梯到结束任务的流程
//调度请求最优电梯>请求电梯到起点楼层>电梯抵达起点楼层>向调度发送电梯已经抵达起点楼层>机器人是否在电梯里面\n
//>请求电剃到终点楼层>电梯抵达终点楼层>向调度发送电梯已经抵达终点楼层>机器人是否在电梯里面>向电梯发送任务结束，进入空闲状态
const (
	REGISTER_SCHEDULE = 1000
	REGISTER_ELE      = 1010
	UPDATE_ELE        = 1001
	CHOOSE_ELE        = 1002
	ARRIVED_START     = 1003
	ARRIVED_END       = 1004
	ROBOT_In_Floor    = 1005
	ROBOT_OUT_Floor   = 1006
	ERROR             = 1024

	ELE_TO_START = 2001
	ELE_TO_END   = 2002
	ELE_TO_FREE  = 2003
	ROBOT_START  = 2004
	ROBOT_END    = 2005
)
