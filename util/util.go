package util

import "flag"

type Command struct {
	Name     string
	Flags    *flag.FlagSet
	Commands []Command
	Run      func()
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}

type Scorable interface {
	Score() int
}

func MaxChannel(in <-chan Scorable) <-chan Scorable {
	result := make(chan Scorable)
	go func() {
		maxScore := 0
		var maxValue Scorable
		for v := range in {
			s := v.Score()
			if s > maxScore {
				maxScore = s
				maxValue = v
			}
		}
		result <- maxValue
	}()

	return result
}

func MaxArray(in []Scorable) Scorable {
	maxScore := 0
	var maxValue Scorable
	for _, v := range in {
		s := v.Score()
		if s > maxScore {
			maxScore = s
			maxValue = v
		}
	}

	return maxValue
}
