package park

import (
	"fmt"
	"log"
	"math/rand"
	"sync/atomic"
	"time"
)

type Park struct {
	Money atomic.Int32
	queue *BlockQueue
	cap   int
}

func InitPark(cap int) *Park {
	ret := &Park{
		Money: atomic.Int32{},
		queue: InitBlockQueue(int64(cap)),
		cap:   cap,
	}
	ret.Money.Store(0)
	return ret
}

// 10秒钟，每秒产生0-5汽车，每辆车独立协程
func (p *Park) GenerateCar() {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < 10; i++ {
		numCar := r.Intn(5)
		fmt.Println(i, "generate", numCar, "cars")
		for j := 0; j < numCar; j++ {
			go p.CarRun()
		}
		time.Sleep(time.Second)
	}
}

// 每个车辆独立协程运行该函数，尝试占用停车场，每辆车最多在停车场停留10秒
func (p *Park) CarRun() {
	p.queue.Put(1) // 车辆进场时，put占用

	r := rand.New(rand.NewSource(time.Now().Unix()))
	seconds := 0
	for seconds == 0 {
		seconds = r.Intn(11)
	}
	fmt.Println("car get in, park", seconds, "seconds")
	time.Sleep(time.Duration(seconds) * time.Second)

	_, err := p.queue.Get() // 车辆出场时，get解除占用
	fmt.Println("car get out")
	p.Money.Add(int32(seconds) * 2)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func ParkTest() {
	p := InitPark(10)
	p.GenerateCar()
	for p.queue.Size != 0 {
		time.Sleep(time.Second)
	}
	fmt.Println("money:", p.Money.Load())
}
