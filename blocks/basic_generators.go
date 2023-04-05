package blocks

import (
	"math/rand"
	"time"

	"github.com/epfl-dcsl/schedsim/engine"
	"github.com/epfl-dcsl/schedsim/global"
)

// Generator interface describes how a generator behaves when describing
// a topology
type Generator interface {
	engine.ActorInterface
	SetCreator(ReqCreator)
}

type genericGenerator struct {
	engine.Actor
	Creator           ReqCreator
	ServiceTime       randDist
	WaitTime          randDist
	TargetApplication int
}

func (g *genericGenerator) SetCreator(rc ReqCreator) {
	g.Creator = rc
}

type randGenerator struct {
	genericGenerator
}

/* // for one generator
func (g *randGenerator) Run() {
	//var test_count = 0

	var test_List []int

	test_List = append(test_List, 0)
	test_List = append(test_List, 1)
	test_List = append(test_List, 2)
	test_List = append(test_List, 3)

	for {
		tarAppli := rand.Intn(global.Cores)

		//tarAppli := test_List[test_count]

		service_time := g.ServiceTime.getRand()
		//fmt.Println("test_ service time:", service_time)
		serviceTypeRandom := rand.Float64()
		serviceType := -1
		if serviceTypeRandom < global.WriteRate {
			serviceType = 1 //Write Operation

			minRand := 5
			maxRand := 15
			randomInt := rand.Intn(maxRand-minRand+1) + minRand
			service_time = service_time * float64(randomInt)
		} else {
			serviceType = 0 //Read Operation
		}
		req := g.Creator.NewRequest(service_time, tarAppli, serviceType)
		//qIdx := rand.Intn(g.GetOutQueueCount())
		qIdx := tarAppli % g.GetOutQueueCount()
		if monitorReq, ok := req.(*MonitorReq); ok {
			monitorReq.initLength = g.GetAllOutQueueLens()[qIdx]
		}
		global.MutexForGenerator[qIdx].Lock()
		g.WriteOutQueueI(req, qIdx)
		global.MutexForGenerator[qIdx].Unlock()
		waitTime := g.WaitTime.getRand()
		//fmt.Println("test_ WaitTime time:", waitTime)

			//test_count += 1
			//if test_count >= 4 {
			//	break
			//}
		g.Wait(waitTime)

	}

}*/

// the generator we use
func (g *randGenerator) Run() {
	//var test_count = 0
	/*
		var test_List []int

		test_List = append(test_List, 0)
		test_List = append(test_List, 1)
		test_List = append(test_List, 2)
		test_List = append(test_List, 3)
	*/
	for {
		tarAppli := g.TargetApplication

		//tarAppli := test_List[test_count]

		service_time := g.ServiceTime.getRand()
		//fmt.Println("test_ service time:", service_time)
		serviceTypeRandom := rand.Float64()
		serviceType := -1
		if serviceTypeRandom < global.WriteRate {
			serviceType = 1 //Write Operation

			minRand := 5
			maxRand := 15
			randomInt := rand.Intn(maxRand-minRand+1) + minRand
			service_time = service_time * float64(randomInt)
		} else {
			serviceType = 0 //Read Operation
		}
		req := g.Creator.NewRequest(service_time, tarAppli, serviceType)
		//qIdx := rand.Intn(g.GetOutQueueCount())
		qIdx := 0 //
		if monitorReq, ok := req.(*MonitorReq); ok {
			monitorReq.initLength = g.GetAllOutQueueLens()[qIdx]
		}
		global.MutexForGenerator[qIdx].Lock()
		g.WriteOutQueueI(req, 0)
		global.MutexForGenerator[qIdx].Unlock()
		waitTime := g.WaitTime.getRand()
		//fmt.Println("test_ WaitTime time:", waitTime)

		//test_count += 1
		//if test_count >= 4 {
		//	break
		//}
		g.Wait(waitTime)

	}

}

type rRGenerator struct {
	genericGenerator
}

func (g *rRGenerator) Run() {
	for count := 0; ; count++ {
		tarAppli := rand.Intn(global.Cores)
		service_time := g.ServiceTime.getRand()
		serviceTypeRandom := rand.Float64()
		serviceType := -1
		if serviceTypeRandom < global.WriteRate {
			serviceType = 1 //Write Operation

			minRand := 5
			maxRand := 15
			randomInt := rand.Intn(maxRand-minRand+1) + minRand
			service_time = service_time * float64(randomInt)
		} else {
			serviceType = 0 //Read Operation
		}
		req := g.Creator.NewRequest(service_time, tarAppli, serviceType)
		g.WriteOutQueueI(req, count%g.GetOutQueueCount())
		g.Wait(g.WaitTime.getRand())
	}
}

// DDGenerator is a fixed waiting time generator that produces fixed service time requests
type DDGenerator struct {
	rRGenerator
}

// NewDDGenerator returns a DDGenerator
func NewDDGenerator(waitTime, serviceTime float64) *DDGenerator {
	g := &DDGenerator{}
	g.ServiceTime = newDeterministicDistr(serviceTime)
	g.WaitTime = newDeterministicDistr(waitTime)
	return g
}

// MDGenerator is a exponential waiting time generator that produces fixed service time requests
// If multiple queues they are fed round robin
type MDGenerator struct {
	rRGenerator
}

// NewMDGenerator returns a MDGenerator
func NewMDGenerator(waitLambda float64, serviceTime float64) *MDGenerator {
	// Seed with time
	rand.Seed(time.Now().UTC().UnixNano())

	g := &MDGenerator{}
	g.ServiceTime = newDeterministicDistr(serviceTime)
	g.WaitTime = newExponDistr(waitLambda)
	return g
}

// MDRandGenerator is a exponential waiting time generator that produces fixed service time requests
// If multiple queues they are fed randomly
type MDRandGenerator struct {
	randGenerator
}

// NewMDRandGenerator returns a MDRandGenerator
func NewMDRandGenerator(waitLambda float64, serviceTime float64, targetApplication int) *MDRandGenerator {
	// Seed with time
	rand.Seed(time.Now().UTC().UnixNano())

	g := &MDRandGenerator{}
	g.WaitTime = newExponDistr(waitLambda)

	g.ServiceTime = newDeterministicDistr(serviceTime)
	g.TargetApplication = targetApplication
	return g
}

// MMGenerator is a exponential waiting time generator that produces exponential service time requests
// If multiple queues they are fed round robin
type MMGenerator struct {
	rRGenerator
}

// NewMMGenerator returns a MMGenerator
func NewMMGenerator(waitLambda float64, serviceMu float64) *MMGenerator {
	// Seed with time
	rand.Seed(time.Now().UTC().UnixNano())

	g := &MMGenerator{}
	g.ServiceTime = newExponDistr(serviceMu)
	g.WaitTime = newExponDistr(waitLambda)
	return g
}

// MMRandGenerator is a exponential waiting time generator that produces exponential service time requests
// If multiple queues they are fed randomly
type MMRandGenerator struct {
	randGenerator
}

// NewMMRandGenerator returns a MMRandGenerator
func NewMMRandGenerator(waitLambda float64, serviceMu float64, targetApplication int) *MMRandGenerator {
	// Seed with time
	rand.Seed(time.Now().UTC().UnixNano())

	g := &MMRandGenerator{}
	//fmt.Print("using NewMMRandGenerator")
	g.ServiceTime = newExponDistr(serviceMu)
	g.WaitTime = newExponDistr(waitLambda)
	g.TargetApplication = targetApplication
	return g
}

// MLNGenerator is exponential waiting time lognormal service time generator
// If multiple queues they are fed round robin
type MLNGenerator struct {
	rRGenerator
}

// NewMLNGenerator returns an MLNGenerator
func NewMLNGenerator(waitLambda, mu, sigma float64) *MLNGenerator {
	// Seed with time
	rand.Seed(time.Now().UTC().UnixNano())

	g := &MLNGenerator{}
	g.ServiceTime = newLGDistr(mu, sigma)
	g.WaitTime = newExponDistr(waitLambda)
	return g
}

// MBGenerator is a poisson interarrival generator with
// requests with bimodal service times (2 values)
// If multiple queues they are fed roundrobin
type MBGenerator struct {
	rRGenerator
}

// NewMBGenerator returns a MBGenerator
func NewMBGenerator(waitLambda, peak1, peak2, ratio float64) *MBGenerator {
	// Seed with time
	rand.Seed(time.Now().UTC().UnixNano())

	g := &MBGenerator{}
	g.ServiceTime = newBiDistr(peak1, peak2, ratio)
	g.WaitTime = newExponDistr(waitLambda)
	return g
}

// MBRandGenerator is a poisson interarrival generator with
// requests with bimodal service times (2 values)
// If multiple queues they are fed randomly
type MBRandGenerator struct {
	randGenerator
}

// NewMBRandGenerator returns a new MBRandGenerator
func NewMBRandGenerator(waitLambda, peak1, peak2, ratio float64) *MBRandGenerator {
	// Seed with time
	rand.Seed(time.Now().UTC().UnixNano())

	g := &MBRandGenerator{}
	g.ServiceTime = newBiDistr(peak1, peak2, ratio)
	g.WaitTime = newExponDistr(waitLambda)
	return g
}
