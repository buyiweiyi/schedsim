package main

import (
	"flag"
	"fmt"

	"github.com/epfl-dcsl/schedsim/global"
	"github.com/epfl-dcsl/schedsim/topologies"
)

func main() {
	var topo = flag.Int("topo", 0, "topology selector")
	var mu = flag.Float64("mu", 0.02, "mu service rate") // default 50usec
	var lambda = flag.Float64("lambda", 0.005, "lambda poisson interarrival")
	var genType = flag.Int("genType", 0, "type of generator")
	var procType = flag.Int("procType", 0, "type of processor")
	var networkQueueType = flag.Int("NetType", 0, "type of the network queue")
	var networkDealingSpe = flag.Float64("netSpeed", 0.02, "the speed of dealing network packages")
	var duration = flag.Float64("duration", 10000000, "experiment duration")
	var bufferSize = flag.Int("buffersize", 1, "size of the bounded buffer")
	//var ServiceTimeDistriType = flag.Int("ServiceTimeDistriType", 0, "The type of distribution of service time")
	flag.Parse()
	fmt.Printf("Selected topology: %v\n", *topo)
	global.Filename = fmt.Sprintf("./save/%d_%.4f_%.4f_%d_%d_%d_%.4f_%.2f_%d_%d.txt", *topo, *mu, *lambda, *genType, *procType, *networkQueueType, *networkDealingSpe, *duration, *bufferSize, global.Cores)
	global.NetworkQueueType_global = *networkQueueType
	global.NetworkDealingSpeed = *networkDealingSpe
	//global.ServiceTimeDistriType = *ServiceTimeDistriType
	if *topo == 0 {
		topologies.SingleQueue(*lambda, *mu, *duration, *genType, *procType)
	} else if *topo == 1 {
		topologies.MultiQueue(*lambda, *mu, *duration, *genType, *procType)
	} else if *topo == 2 {
		topologies.BoundedQueue(*lambda, *mu, *duration, *bufferSize)
	} else {
		panic("..Unknown topology..")
	}
}
