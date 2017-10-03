package timer

import (
	"fmt"
	"sync"
	"time"
)

const (
	timeNearShift  uint32 = 8
	timeNear       uint32 = 1 << timeNearShift
	timeNearMask   uint32 = timeNear - 1
	timeLevelShift uint32 = 6
	timeLevel      uint32 = 1 << timeLevelShift
	timeLevelMask  uint32 = timeLevel - 1
)

type timerEvent interface {
	Action()
}

type timerNode struct {
	event  timerEvent
	next   *timerNode
	expire uint32
}

type linkList struct {
	head timerNode // reserve space
	tail *timerNode
}

type timer struct {
	near        [timeNear]linkList
	far         [4][timeNear]linkList // far[0][0], far[1][0], far[2][0] is spared, far[3][0] storage num(0--2^26)exceed uint32(point is over 2^32 when it(2^31--2^32) plus time(1--2^31))
	mu          sync.Mutex
	point       uint32 // centisecond: 1/100 second
	currentTime uint64 // current time is centisecond
}

var Timer *timer = nil

func linkReset(list *linkList) *timerNode {
	ret := list.head.next
	list.head.next = nil
	list.tail = &list.head
	return ret
}

func link(list *linkList, node *timerNode) {
	list.tail.next = node
	list.tail = node
	node.next = nil
}

func (t *timer) addNode(node *timerNode) {
	expire := node.expire
	point := t.point
	if point|timeNearMask == expire|timeNearMask {
		link(&t.near[expire&timeNearMask], node)
	} else {
		mask := timeNear << timeLevelShift
		var i uint32
		for i = 0; i < 3; i++ {
			mask <<= timeLevelShift
			if point|(mask-1) == expire|(mask-1) {
				break
			}
		}
		// expire > point | expire < point(only expire exceed uint32)
		link(&t.far[i][(expire>>(i*timeLevelShift+timeNearShift))&timeLevelMask], node)
	}
}

func (t *timer) timerAdd(event timerEvent, time uint32) {
	node := new(timerNode)
	node.event = event
	node.expire = time + t.point

	t.mu.Lock()
	t.addNode(node)
	t.mu.Unlock()
}

func (t *timer) moveList(level uint32, index uint32) {
	currentNode := linkReset(&t.far[level][index])
	for currentNode != nil {
		temp := currentNode.next
		t.addNode(currentNode)
		currentNode = temp
	}
}

func (t *timer) timerShift() {
	mask := timeNear
	t.point++
	ct := t.point
	if ct == 0 {
		t.moveList(3, 0)
	} else {
		var level uint32 = 0
		point := ct >> timeNearShift
		for (ct & (mask - 1)) == 0 {
			index := point & timeLevelMask
			if index != 0 {
				t.moveList(level, index)
				break
			}
			level++
			mask <<= timeLevelShift
			point >>= timeLevelShift
		}
	}
}

func (t *timer) dispatchList(currentNode *timerNode) {
	for currentNode != nil {
		currentNode.event.Action()
		currentNode = currentNode.next
	}
}

func (t *timer) timerExecute() {
	index := t.point & timeNearMask
	for t.near[index].head.next != nil { // todo think if replace for
		currentNode := linkReset(&t.near[index])
		t.mu.Unlock()
		// dispatch don't need lock t
		t.dispatchList(currentNode)
		t.mu.Lock()
	}
}

func (t *timer) timerUpdate() {
	t.mu.Lock()

	// try to dispatch 0 (rare condition)
	t.timerExecute()

	t.timerShift()

	t.timerExecute()

	t.mu.Unlock()
}

func (t *timer) updateTime() {
	ct := getTime()
	if ct < t.currentTime {
		fmt.Printf("error time, ct:%v, t.currentTime:=%v\n", ct, t.currentTime)
	} else {
		diff := ct - t.currentTime
		t.currentTime = ct
		var i uint64
		for i = 0; i < diff; i++ {
			t.timerUpdate()
		}
	}
}

func timerCreate() *timer {
	t := new(timer)
	var i, j uint32
	for i = 0; i < timeNear; i++ {
		linkReset(&t.near[i])
	}

	for i = 0; i < 4; i++ {
		for j = 0; j < timeLevel; j++ {
			linkReset(&t.far[i][j])
		}
	}
	return t
}

func getTime() uint64 {
	return uint64(time.Now().UnixNano() / 1e7)
}

// time is sentisecond: 1/100 second
func AddTimer(event timerEvent, time int32) {
	if time <= 0 {
		event.Action()
	} else {
		Timer.timerAdd(event, uint32(time))
	}
}

func TimerInit() {
	Timer = timerCreate()
	Timer.currentTime = getTime()

	for {
		time.Sleep(25 * 1e5)
		Timer.updateTime()
	}
}
