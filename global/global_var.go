package global

import (
	"sync"
)

// Greeting is an exported global variable, accessible from other packages.
var Filename string = "Hello, world!"
var NetworkQueueType_global = -1
var NetworkDealingSpeed = 1000.0

const Cores int = 4

// var MutexChanLock sync.Mutex
var WriteRate float64 = 0.1
var MutexForGenerator [Cores]sync.Mutex

//var ServiceTimeDistriType = -1
