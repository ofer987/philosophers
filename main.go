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
	leftFork  *fork
	rightFork *fork
	isEating  bool
}

type fork struct {
	mut        *sync.Mutex
	isPickedUp bool
}

const FORK_COUNT = PHILOSOPHER_COUNT + 1

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
	}

	// Create Philosophers
	i = 0
	for i < PHILOSOPHER_COUNT {
		philosophers[i] = philosopher{
			leftFork:  &forks[i],
			rightFork: &forks[i+1],
			isEating:  false,
		}
		i += 1
	}

	for {
		go startEating()
		duration := 6 * time.Second

		ticker := time.NewTicker(duration)

		<-ticker.C

		dropAllForks()
		startEating()
	}

	fmt.Println("vim-go")
}

func incrementEating() {
	eatingMutex.Lock()
	defer eatingMutex.Unlock()

	eating += 1
}

func (f *fork) pickUp() {
	f.mut.Lock()
	defer f.mut.Unlock()

	f.isPickedUp = true
}

func (f *fork) drop() {
	f.mut.Lock()
	defer f.mut.Unlock()

	f.isPickedUp = false
}

func startEating() {
	i := 0
	for i < PHILOSOPHER_COUNT {
		go philosophers[i].startToEat()
	}
}

func (phil *philosopher) startToEat() {
	if phil.isEating {
		return
	}

	phil.leftFork.pickUp()
	phil.rightFork.pickUp()

	if phil.leftFork.isPickedUp && phil.rightFork.isPickedUp {
		phil.isEating = true
	}

	phil.stopEating()
}

func (phil *philosopher) stopEating() {
	duration := 2 * time.Second

	ticker := time.NewTicker(duration)

	for {
		<-ticker.C

		phil.isEating = false
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

	i := 0
	for i < FORK_COUNT {
		forks[i].mut.Unlock()
		forks[i].isPickedUp = false
	}
}
