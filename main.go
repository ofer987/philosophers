package main

import (
	"fmt"
	"sync"
	"time"
)

const PHILOSOPHER_COUNT = 5
const EATING_TIME_IN_SECONDS = 2

var eating = 0
var eatingMutex sync.Mutex

type philosopher struct {
	i         int
	leftFork  *fork
	rightFork *fork
	isEating  bool
}

type fork struct {
	mut        *sync.Mutex
	isPickedUp bool
	phil       *philosopher
}

const FORK_COUNT = PHILOSOPHER_COUNT

var forks = [FORK_COUNT]fork{}
var philosophers = [PHILOSOPHER_COUNT]philosopher{}

func main() {
	// Create Forks
	i := 0
	for i < FORK_COUNT {
		var mut sync.Mutex

		forks[i] = fork{
			mut:        &mut,
			isPickedUp: false,
		}

		i += 1
	}

	// Create Philosophers
	i = 0
	for i < PHILOSOPHER_COUNT {
		leftForkIndex := i
		rightForkIndex := (i + 1) % 5
		philosophers[i] = philosopher{
			i:         i + 1,
			leftFork:  &forks[leftForkIndex],
			rightFork: &forks[rightForkIndex],
			isEating:  false,
		}

		i += 1
	}

	fmt.Println(3)
	startEating()

	for {
		duration := 6 * time.Second
		ticker := time.NewTicker(duration)

		<-ticker.C

		dropAllForks()
	}
}

func (f *fork) pickUp(phil *philosopher) {
	f.mut.Lock()
	defer f.mut.Unlock()

	if f.isPickedUp {
		return
	}

	f.phil = phil
	f.isPickedUp = true
}

func (f *fork) drop() {
	f.mut.Lock()
	defer f.mut.Unlock()

	f.isPickedUp = false
	f.phil = nil
}

func startEating() {
	eatingMutex.Lock()
	defer eatingMutex.Unlock()

	fmt.Printf("Start dining\n")

	i := 0
	for i < PHILOSOPHER_COUNT {
		go philosophers[i].startToEat()

		i += 1
	}
}

func (phil *philosopher) startToEat() {
	for {
		if phil.isEating {
			continue
		}

		phil.leftFork.pickUp(phil)
		phil.rightFork.pickUp(phil)

		if phil.leftFork.isPickedUp && phil.leftFork.phil == phil && phil.rightFork.isPickedUp && phil.rightFork.phil == phil {
			incrementEating(phil)

			phil.stopEating()
		}
	}
}

func incrementEating(phil *philosopher) {
	eatingMutex.Lock()
	defer eatingMutex.Unlock()

	phil.isEating = true
	eating += 1

	i := 0
	eatingCount := 0

	fmt.Println("")
	fmt.Printf("A total of %d philosopher are supposed to be eating\n", eating)
	for i < PHILOSOPHER_COUNT {
		philOther := philosophers[i]
		if philOther.isEating {
			fmt.Printf("Philosopher %d is eating\n", philOther.i)

			eatingCount += 1
		}

		i += 1
	}

	fmt.Printf("A total of %d philosopher are eating\n\n", eatingCount)
}

func (phil *philosopher) stopEating() {
	duration := 2 * time.Second

	ticker := time.NewTicker(duration)

	for {
		<-ticker.C

		eatingMutex.Lock()
		defer eatingMutex.Unlock()

		eating -= 1
		phil.isEating = false
		fmt.Printf("Philosopher %d has stopped eating\n", phil.i)

		phil.leftFork.drop()
		phil.rightFork.drop()

		break
	}
}

func dropAllForks() {
	eatingMutex.Lock()
	defer eatingMutex.Unlock()

	if eating > 0 {
		return
	}

	fmt.Printf("START: Drop all the forks because there is a deadlock!\n")
	i := 0
	for i < FORK_COUNT {
		forks[i].drop()

		i += 1
	}

	fmt.Printf("END: Drop all the forks because there is a deadlock!\n")
}
