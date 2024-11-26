package consumer

import "sync/atomic"

type typeCounter struct {
	zero  atomic.Uint64
	one   atomic.Uint64
	two   atomic.Uint64
	three atomic.Uint64
	four  atomic.Uint64
	five  atomic.Uint64
	six   atomic.Uint64
	seven atomic.Uint64
	eight atomic.Uint64
	nine  atomic.Uint64
}

func newTypeCounter() *typeCounter {
	return &typeCounter{
		zero:  atomic.Uint64{},
		one:   atomic.Uint64{},
		two:   atomic.Uint64{},
		three: atomic.Uint64{},
		four:  atomic.Uint64{},
		five:  atomic.Uint64{},
		six:   atomic.Uint64{},
		seven: atomic.Uint64{},
		eight: atomic.Uint64{},
		nine:  atomic.Uint64{},
	}
}

func (counter *typeCounter) getCounter(taskType int) *atomic.Uint64 {
	switch taskType {
	case 0:
		return &counter.zero
	case 1:
		return &counter.one
	case 2:
		return &counter.two
	case 3:
		return &counter.three
	case 4:
		return &counter.four
	case 5:
		return &counter.five
	case 6:
		return &counter.six
	case 7:
		return &counter.seven
	case 8:
		return &counter.eight
	case 9:
		return &counter.nine
	default:
		panic("unknown task type")
	}
}

func (counter *typeCounter) Inc(taskType int) {
	counter.getCounter(taskType).Add(1)
}

func (counter *typeCounter) Get(taskType int) int {
	return int(counter.getCounter(taskType).Load())
}

func (counter *typeCounter) Total() int {
	result := 0
	for i := range 10 {
		result += counter.Get(i)
	}

	return result
}
