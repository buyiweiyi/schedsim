package topologies

import (
	"fmt"

	"github.com/epfl-dcsl/schedsim/blocks"
	"github.com/epfl-dcsl/schedsim/engine"
	"github.com/epfl-dcsl/schedsim/global"
)

// MultiQueue describes a single-generator-multi-processor topology where every
// processor has its own incoming queue
func MultiQueue(lambda, mu, duration float64, genType, procType int) {

	engine.InitSim()

	//Init the statistics
	//stats := blocks.NewBookKeeper()
	stats_network := &blocks.AllKeeper{}
	stats_network.SetName("Network Stats")
	engine.InitStats(stats_network)
	stats := &blocks.AllKeeper{}
	stats.SetName("Main Stats")
	engine.InitStats(stats)

	// Add generator
	var g [global.Cores]blocks.Generator
	for i := 0; i < global.Cores; i++ {
		if i == 2 {
			lambdaHeavyLoad := lambda * 1
			if genType == 0 {
				g[i] = blocks.NewMMRandGenerator(lambdaHeavyLoad, mu, i)
			} else if genType == 1 {
				g[i] = blocks.NewMDRandGenerator(lambdaHeavyLoad, 1/mu, i)
			} else if genType == 2 {
				g[i] = blocks.NewMBRandGenerator(lambdaHeavyLoad, 1, 10*(1/mu-0.9), 0.9)
			} else if genType == 3 {
				g[i] = blocks.NewMBRandGenerator(lambdaHeavyLoad, 1, 1000*(1/mu-0.999), 0.999)
			}
		} else {

			if genType == 0 {
				g[i] = blocks.NewMMRandGenerator(lambda, mu, i)
			} else if genType == 1 {
				g[i] = blocks.NewMDRandGenerator(lambda, 1/mu, i)
			} else if genType == 2 {
				g[i] = blocks.NewMBRandGenerator(lambda, 1, 10*(1/mu-0.9), 0.9)
			} else if genType == 3 {
				g[i] = blocks.NewMBRandGenerator(lambda, 1, 1000*(1/mu-0.999), 0.999)
			}
		}
		g[i].SetCreator(&blocks.SimpleReqCreator{})
	}

	// Create queues
	fastQueues := make([]engine.QueueInterface, global.Cores)
	for i := range fastQueues {
		fastQueues[i] = blocks.NewQueue()
	}

	// Create processors
	processors := make([]blocks.Processor, global.Cores)

	// first the slow Cores
	for i := 0; i < global.Cores; i++ {
		if procType == 0 {
			processors[i] = &blocks.RTCProcessor{}
		} else if procType == 1 {
			processors[i] = blocks.NewPSProcessor()
		}
	}
	//var networkQueue engine.QueueInterface
	//var networkQueues []engine.QueueInterface
	networkQueue := blocks.NewQueue()
	networkQueues := make([]engine.QueueInterface, global.Cores)
	for i := range networkQueues {
		networkQueues[i] = blocks.NewQueue()
	}
	//pNetwork := blocks.NetworkManager(global.NetworkDealingSpeed) //network manager
	pNetwork := blocks.SimpleTestNetworkManager(global.NetworkDealingSpeed) //network manager
	if global.NetworkQueueType_global == 0 {
		// Connect the fast queues
		for i := 0; i < global.Cores; i++ {
			g[i].AddOutQueue(networkQueue)
		}

		pNetwork.AddInQueue(networkQueue)

		for i, q := range fastQueues {
			pNetwork.AddOutQueue(q)
			processors[i].AddInQueue(q)
		}

		// Add the stats and register processors
		pNetwork.SetReqDrain(stats_network)
		engine.RegisterActor(pNetwork)
		for _, p := range processors {
			p.SetReqDrain(stats)
			engine.RegisterActor(p)
		}
	} else {
		// Connect the fast queues

		//pNetwork := blocks.NetworkManager(global.NetworkDealingSpeed)

		for i, q := range networkQueues {

			g[i].AddOutQueue(q)
			pNetwork.AddInQueue(q)

		}

		for i, q := range fastQueues {

			//g.AddOutQueue(networkQueues[i])
			//pNetwork.AddInQueue(networkQueues[i])
			pNetwork.AddOutQueue(q)
			processors[i].AddInQueue(q)
		}
		//fmt.Printf("Cores:%d_%d\n", pNetwork.GetInQueueCount(), pNetwork.GetOutQueueCount())
		// Add the stats and register processors
		pNetwork.SetReqDrain(stats_network)
		engine.RegisterActor(pNetwork)
		for _, p := range processors {
			p.SetReqDrain(stats)
			engine.RegisterActor(p)
		}
	}
	for i := 0; i < global.Cores; i++ {
		// Register the generator
		engine.RegisterActor(g[i])
	}

	fmt.Printf("Cores:%v\tservice_rate:%v\tinterarrival_rate:%v\n", global.Cores, mu, lambda)
	engine.Run(duration)
}
