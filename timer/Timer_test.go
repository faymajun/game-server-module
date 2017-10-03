package timer

import (
	"fmt"
	"testing"
	"time"
)

//func TestTimerInit(t *testing.T) {
//	go TimerInit()
//	time.Sleep(1e9)
//	for i := 0; i < 100; i++ {
//		AddTimer(LoggingTimerEvent{i}, int32(i*100))
//		AddTimer(LoggingTimerEvent{i}, int32(i*100))
//		AddTimer(LoggingTimerEvent{i}, int32(i*100))
//	}
//	time.Sleep(100 * 1e9)
//	fmt.Println("input over")
//}

func BenchmarkTimerInit(b *testing.B) {
	go TimerInit()
	time.Sleep(1e9)
	for i := 0; i < 20; i++ {
		AddTimer(LoggingTimerEvent{i}, int32(i*1))
		AddTimer(LoggingTimerEvent{i}, int32(i*10))
		AddTimer(LoggingTimerEvent{i}, int32(i*100))
	}
	time.Sleep(100 * 1e9)
	fmt.Println("input over")
}

type LoggingTimerEvent struct {
	num int
}

func (l LoggingTimerEvent) Action() {
	fmt.Printf("time:%d, nanosecond:%d, num:%d\n", time.Now().Second(), time.Now().Nanosecond(), l.num)
}
