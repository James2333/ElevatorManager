package elevator

import (
	"errors"
	"log"
	"sync"
)

type (
	ElevatorM interface {
		Update(elevator *Elevator) error
		RightElevator(int64) (string, error)
	}
	Elevator struct {
		//电梯id    最终返回这个
		//当前楼层   从连接中获取 ，需要取最优解
		//当前状态   取空闲状态的电梯   繁忙/空闲/不可用
		//临时状态	选中这个电梯之后临时变为繁忙状态，后续状态由心跳更新。
		ElevatorId   string
		Floor        int64
		State        string
		CurrentState string
		IsInFloor    bool //电梯内是否有机器人 true 为在
	}
)

type Elevators map[string]*Elevator
var Els Elevators

//type ELs sync.Map
var wg sync.RWMutex

func init()  {
	Els= make(map[string]*Elevator)
}

func (els Elevators) Update(el *Elevator) error{
	wg.RLock()
	defer wg.RUnlock()
	//如果CurrentState不为0则不允许更新。
	if ele,ok:=els[el.ElevatorId];ok{
		if ele.CurrentState!="0"{
			log.Printf("电梯ID:%s当前不允许更新",ele.ElevatorId)
			return errors.New("当前状态不允许更新")
		}
	}
	//有电梯信息则覆盖，无则新增
	el.CurrentState="0"
	els[el.ElevatorId] = el
	return nil
}

func (els Elevators) RightElevator(start int64) (string, error) {
	wg.RLock()
	defer wg.RUnlock()
	freeEl := make(map[string]int64)
	for key, value := range els {
		if value.State != "0" || value.CurrentState != "0" {
			continue
		} else {
			freeEl[key] = value.Floor
			//取出空闲的电梯，id为key 所在楼层为value
		}
	}
	//判定：全都不在空闲中则返回“目前无可用电梯 稍后再试”
	//如果只有一个电梯则直接返回
	//如果多个电梯则设置一个中间数，看谁的距离绝对值最短。同最短则取最小的
	elId := ""
	if len(freeEl) == 0 {
		return "", errors.New("目前无可用电梯,稍后再试")
	}
	if len(freeEl) == 1 {
		for key, _ := range freeEl {
			//midEl := els[key]
			//midEl.CurrentState = "1"
			//els.Update(midEl)
			els[key].CurrentState = "1"
			return key, nil
		}
	} else {
		min := int64(99)
		for k, v := range freeEl {
			if Abs(v-start) < min {
				elId = k
				min = Abs(v - start)
			}
		}
	}
	//midEl := els[elId]
	//midEl.CurrentState = "1"
	//els.Update(midEl)
	els[elId].CurrentState = "1"
	return elId, nil
}
func Abs(n int64) int64 {
	//els["zbc"]=nil
	if n < 0 {
		return -n
	}
	return n
}