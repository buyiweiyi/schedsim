package engine

import (
	"container/heap"
	"container/list"
	"fmt"
)

var mdl *model

// ActorInterface is the main interface to be used in main package.
// Every element of the topology should implement this interface.
// Init, AddInQueuem AddOutQueue are provided by the Actor nested struct and
// only the Run() function needs to be implemented
type ActorInterface interface {
	Run()
	AddInQueue(q QueueInterface)
	AddOutQueue(q QueueInterface)
	init(ch chan interface{})
}

// ReqInterface describes what a basic request should look like
type ReqInterface interface {
	GetDelay() float64
	GetServiceTime() float64
	GetTargetAppli() int
	SubServiceTime(t float64)
}

// QueueInterface describe basic queue functionality
type QueueInterface interface {
	Enqueue(ReqInterface)
	Dequeue() ReqInterface
	Len() int
}

// Stats is an interface that is called at the end of the simulation and
// prints the collected statistics
type Stats interface {
	PrintStats()
}

type timerEventInterface interface {
	getTime() float64
	setIdx(idx int)
	getChannel() chan int
	getEventType() int
	setEventType(EventType int)
}

type timerEvent struct {
	time      float64
	wakeUpCh  chan int
	idx       int
	eventType int
}

func (te *timerEvent) getTime() float64 {
	return te.time
}

func (te *timerEvent) setIdx(idx int) {
	te.idx = idx
}

func (te *timerEvent) getChannel() chan int {
	return te.wakeUpCh
}

func (te *timerEvent) getEventType() int {
	return te.eventType
}

func (te *timerEvent) setEventType(EventType int) {
	te.eventType = EventType
}

type blockEventInterface interface {
	getChannel() chan int
	getQueues() []QueueInterface
	deactivateReplicas()
	addReplica(pair listElPair)
}

type listElPair struct {
	el *list.Element
	l  *list.List
}

type blockEvent struct {
	wakeUpCh chan int
	queues   []QueueInterface
	replicas []listElPair
}

func (be *blockEvent) getChannel() chan int {
	return be.wakeUpCh
}

func (be *blockEvent) getQueues() []QueueInterface {
	return be.queues
}

func (be *blockEvent) deactivateReplicas() {
	for _, pair := range be.replicas {
		pair.l.Remove(pair.el)
	}
}

func (be *blockEvent) addReplica(pair listElPair) {
	be.replicas = append(be.replicas, pair)
}

type linkedEvent struct {
	timerEvent
	blockEvent
}

func (le *linkedEvent) getChannel() chan int {
	return le.blockEvent.wakeUpCh
}

type model struct {
	time            float64
	actorCount      int
	pq              priorityQueue
	eventChan       chan interface{}
	blockedInQueues map[QueueInterface]*list.List
	queues          map[QueueInterface]bool
	bookkeeping     []Stats
}

func newModel() *model {
	m := &model{}
	m.eventChan = make(chan interface{}, 1000)
	m.pq = make(priorityQueue, 0)
	m.queues = make(map[QueueInterface]bool)
	m.blockedInQueues = make(map[QueueInterface]*list.List)
	heap.Init(&m.pq)
	return m
}

func (m *model) registerActor(a ActorInterface) {
	a.init(m.eventChan)
	m.actorCount++

	go a.Run()
}

func (m *model) registerBlockEvent(e blockEventInterface) {
	for _, q := range e.getQueues() {
		if _, ok := m.blockedInQueues[q]; !ok {
			m.blockedInQueues[q] = list.New()
		}
		el := m.blockedInQueues[q].PushBack(e)
		e.addReplica(listElPair{el, m.blockedInQueues[q]})
	}
}

func (m *model) getTime() float64 {
	return m.time
}

func (m *model) waitActor() {
	newEvent := <-m.eventChan
	if timerE, ok := newEvent.(timerEvent); ok {
		heap.Push(&m.pq, &timerE)
		return
	}
	if blockE, ok := newEvent.(blockEvent); ok {
		m.registerBlockEvent(&blockE)
		return
	}
	if linkedE, ok := newEvent.(linkedEvent); ok {
		heap.Push(&m.pq, &linkedE)
		m.registerBlockEvent(&linkedE)
		return
	}
}

func (m *model) run(threshold float64) {
	////wait for all actors to start and add an event or block on a queue
	for i := 0; i < m.actorCount; i++ {
		m.waitActor()
	}
	//var test_count = 0

	var log_ServiceTime []int

	//all actors started
	for m.time < threshold {

		for q := range m.queues {
			if q.Len() == 0 {
				continue
			}

			// Check if none is waiting for this active queue
			if val, ok := m.blockedInQueues[q]; ok {
				if val.Len() == 0 {
					continue
				}
			} else {
				continue
			}

			for e := m.blockedInQueues[q].Front(); e != nil && q.Len() > 0; e = e.Next() {
				be := e.Value.(blockEventInterface)
				// Remove the blockEvents for the rest of the queues if any
				be.deactivateReplicas()

				if linkedE, ok := e.Value.(*linkedEvent); ok {
					heap.Remove(&m.pq, linkedE.timerEvent.idx)
				}
				be.getChannel() <- 1 // try to unblock
				m.waitActor()
				//m.blockedInQueues[q].Remove(e)
			}
		}

		// pick event and wake up process
		e := heap.Pop(&m.pq).(timerEventInterface)
		m.time = e.getTime()

		// if it's linked deactivate the blocked requests
		if _, ok := e.(*timerEvent); ok {
			//fmt.Println("event type:", e.getEventType())
			if e.getEventType() == 10 {
				log_ServiceTime = append(log_ServiceTime, 1)
				//fmt.Println("event resolved:", e.getTime())

				//test_count += 1
			}

		}
		// if it's linked deactivate the blocked requests
		if linkedE, ok := e.(*linkedEvent); ok {
			linkedE.blockEvent.deactivateReplicas()
		}
		e.getChannel() <- 1
		//
		/*
			fmt.Println("Preparing to output result 0:", test_count)
			if test_count >= 4 {
				fmt.Println("Preparing to output result 1:")
				break
			}*/
		// wait till process adds event or blocks in queue
		m.waitActor()

	}
	fmt.Println("Preparing to output result:")
	for _, s := range m.bookkeeping {
		s.PrintStats()
	}
	//os.Exit(0)
}

// InitSim initialises the simulation
func InitSim() {
	mdl = newModel()
}

// GetTime returns the current simulation time
func GetTime() float64 {
	return mdl.getTime()
}

// RegisterActor registers a specific simulation element.
// All actors should be registered
func RegisterActor(a ActorInterface) {
	mdl.registerActor(a)
}

// Run runs the simulation for till the given threshold time
func Run(threshold float64) {
	mdl.run(threshold)
}

// InitStats sets the interface in charge of collecting statistics.
// This is interface is called at the end of the simulation to print the
// collected statistics
func InitStats(s Stats) {
	mdl.bookkeeping = append(mdl.bookkeeping, s)
}
