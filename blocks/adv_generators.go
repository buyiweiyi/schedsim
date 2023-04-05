package blocks

import (
	"bufio"
	"math/rand"
	"os"
	"strconv"

	"github.com/epfl-dcsl/schedsim/global"
)

// PBGenerator implements a playback generator for given service times.
// The interarrival distribution is exponential
type PBGenerator struct {
	genericGenerator
	sTimes   [][]int
	cpuCount int
	WaitTime randDist
}

// NewPBGenerator returns a PBGenerator
// Parameters: lambda for the exponential interarrival and the filenames
// with the service times
func NewPBGenerator(lambda float64, paths []string) *PBGenerator {
	g := PBGenerator{}

	for _, p := range paths {
		/* Read service times */
		inFile, _ := os.Open(p)
		defer inFile.Close()
		scanner := bufio.NewScanner(inFile)
		scanner.Split(bufio.ScanLines)

		newTimes := make([]int, 0)
		for scanner.Scan() {
			n, _ := strconv.Atoi(scanner.Text())
			newTimes = append(newTimes, n)
		}
		g.sTimes = append(g.sTimes, newTimes)
	}
	g.cpuCount = len(paths)
	g.WaitTime = newExponDistr(lambda)
	return &g
}

// Run is the main loop of the generator
func (g *PBGenerator) Run() {
	for {
		i := rand.Intn(g.cpuCount)
		j := rand.Intn(len(g.sTimes[i]))
		serviceTime := g.sTimes[i][j]
		tarAppli := rand.Intn(global.Cores)
		//fmt.Println("test_ service time:", service_time)
		serviceTypeRandom := rand.Float64()
		serviceType := -1
		if serviceTypeRandom < global.WriteRate {
			serviceType = 1 //Write Operation

			minRand := 5
			maxRand := 15
			randomInt := rand.Intn(maxRand-minRand+1) + minRand
			serviceTime = serviceTime * randomInt
		} else {
			serviceType = 0 //Read Operation
		}
		req := g.Creator.NewRequest(float64(serviceTime), tarAppli, serviceType)
		g.WriteOutQueueI(req, i)
		g.Wait(g.WaitTime.getRand())
	}
}
