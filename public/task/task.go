package task

import (
	"errors"
	"log"
	"sync"
)

type Task struct {
	ElevatorID string `json:"elevator_id"`
	TaskID     string `json:"task_id"`
	Start      int64  `json:"start"`
	End        int64  `json:"end"`
}

type Tasks map[string]*Task

var m sync.Mutex
var Ts Tasks

func init()  {
	Ts =make(map[string]*Task)
}

func (tasks Tasks)AddTask(t *Task) error {
	m.Lock()
	defer m.Unlock()
	if task,ok:=tasks[t.TaskID];ok{
		log.Println("此任务已存在ID：",task.TaskID)
		return errors.New("此任务已存在")
	}
	tasks[t.TaskID] = t
	return nil
}

func (tasks Tasks)Get(taskID string) (*Task,error){
	task,ok:=tasks[taskID]
	if!ok{
		return nil,errors.New("任务不存在")
	}
	return task,nil
}

func (tasks Tasks)Delete(taskID string){
	m.Lock()
	delete(tasks, taskID) //删除任务
	m.Unlock()
}

