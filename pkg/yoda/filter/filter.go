package filter

import (
	"strconv"

	v1 "k8s.io/api/core/v1"
)

func PodFitsNumber(pod *v1.Pod, node *v1.Node) (bool, uint) {
	if number, ok := pod.GetLabels()["scv/number"]; ok {
		return strToUint(number) < 1, strToUint(number)
	}
	return true, 1
}

func PodFitsMemory(number uint, pod *v1.Pod, node *v1.Node) (bool, uint64) {
	if memory, ok := pod.GetLabels()["scv/memory"]; ok {
		fitsCard := uint(0)
		m := StrToUint64(memory)
		if fitsCard >= number {
			return true, m
		}
		return false, m
	}
	return true, 0
}

func PodFitsClock(number uint, pod *v1.Pod, node *v1.Node) (bool, uint) {
	if clock, ok := pod.GetLabels()["scv/clock"]; ok {
		fitsCard := uint(0)
		c := strToUint(clock)

		if fitsCard >= number {
			return true, c
		}
		return false, c
	}
	return true, 0
}

func CardFitsMemory(memory uint64) bool {
	return true
}

func CardFitsClock(clock uint) bool {
	return true
}

func strToUint(str string) uint {
	if i, e := strconv.Atoi(str); e != nil {
		return 0
	} else {
		return uint(i)
	}
}

func StrToUint64(str string) uint64 {
	if i, e := strconv.Atoi(str); e != nil {
		return 0
	} else {
		return uint64(i)
	}
}

func StrToInt64(str string) int64 {
	if i, e := strconv.Atoi(str); e != nil {
		return 0
	} else {
		return int64(i)
	}
}

func Uint64ToInt64(intNum uint64) int64 {
	return StrToInt64(strconv.FormatUint(intNum, 10))
}
