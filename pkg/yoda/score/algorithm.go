package score

import (
	"errors"
	"github.com/NJUPT-ISL/Yoda-Scheduler/pkg/yoda/collection"
	"github.com/NJUPT-ISL/Yoda-Scheduler/pkg/yoda/filter"
	v1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

// Sum is from collection/collection.go
// var Sum = []string{"Cores","FreeMemory","Bandwidth","MemoryClock","MemorySum","Number","Memory"}

const (
	BandwidthWeight   = 1
	ClockWeight       = 1
	CoreWeight        = 1
	PowerWeight       = 1
	FreeMemoryWeight  = 2
	TotalMemoryWeight = 1
	ActualWeight      = 2

	AllocateWeight = 3
)

func CalculateScore(s *v1.Node, state *framework.CycleState, pod *v1.Pod, info *framework.NodeInfo) (uint64, error) {
	d, err := state.Read("Max")
	if err != nil {
		return 0, errors.New("Error Get CycleState Info Max Error: " + err.Error())
	}
	data, ok := d.(*collection.Data)
	if !ok {
		return 0, errors.New("The Type is not Data ")
	}
	return CalculateBasicScore(data.Value, s, pod) + CalculateAllocateScore(info, s) + CalculateActualScore(), nil
}

func CalculateBasicScore(value collection.MaxValue, node *v1.Node, pod *v1.Pod) uint64 {
	var cardScore uint64
	if ok, number := filter.PodFitsNumber(pod, node); ok {
		isFitsMemory, memory := filter.PodFitsMemory(number, pod, node)
		isFitsClock, _ := filter.PodFitsClock(number, pod, node)
		if isFitsClock && isFitsMemory {
			return memory
		}
	}
	return cardScore
}

func CalculateCardScore(value collection.MaxValue) uint64 {
	var (
		bandwidth   = value.MaxBandwidth
		clock       = value.MaxBandwidth
		core        = value.MaxCore
		power       = value.MaxPower
		freeMemory  = value.MaxFreeMemory
		totalMemory = value.MaxTotalMemory
	)
	return uint64(bandwidth*BandwidthWeight+clock*ClockWeight+core*CoreWeight+power*PowerWeight) +
		freeMemory*FreeMemoryWeight + totalMemory*TotalMemoryWeight
}

func CalculateActualScore() uint64 {
	return ActualWeight
}

func CalculateAllocateScore(info *framework.NodeInfo, scv *v1.Node) uint64 {
	allocateMemorySum := uint64(0)
	for _, pod := range info.Pods {
		if mem, ok := pod.Pod.GetLabels()["scv/memory"]; ok {
			allocateMemorySum += filter.StrToUint64(mem)
		}
	}
	return 0
}
