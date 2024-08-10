package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"time"
)

type Percentiles struct {
	percentiles []float64
	defaulted   bool
}

func (percentile *Percentiles) Set(input string) error {
	val, err := strconv.ParseFloat(input, 32)
	if err != nil {
		return err
	}

	if percentile.defaulted {
		percentile.percentiles = make([]float64, 0)
		percentile.defaulted = false
	}

	percentile.percentiles = append(percentile.percentiles, val)
	return nil
}

func (percentile *Percentiles) String() string {
	output := ""

	for _, v := range percentile.percentiles {
		output += strconv.FormatFloat(v, 'f', -1, 64) + ","
	}

	return output
}

func main() {
	percentiles := &Percentiles{
		defaulted:   true,
		percentiles: []float64{50, 90, 99},
	}

	flag.Var(percentiles, "p", "the percentile you want. Defaults to 50/90/99 if not specified")
	mean := flag.Bool("mean", false, "If the mean should be shown")
	max := flag.Bool("max", false, "If the max should be shown")
	refreshRate := flag.Int64("refreshRate", 1, "The rate at which to refresh the calculated values in seconds")

	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)

	values := make([]float64, 0)
	valuesChan := make(chan float64, 100)

	for scanner.Scan() {
		input := scanner.Text()
		val, err := strconv.ParseFloat(input, 64)
		if err != nil {
			log.Fatal("An error occured parsing user input", err)
		}
		valuesChan <- val
	}

	ticker := time.NewTicker(time.Duration(*refreshRate) * time.Second)

	for {
		select {

		case <-ticker.C:
			sort.Float64s(values)
			fmt.Print("|")

			for _, percentile := range percentiles.percentiles {
				index := int(math.Floor(float64(len(values)) * percentile / 100.0))
				if index < len(values) {
					fmt.Printf(" p%.0f = %.2f      |", percentile, values[index])
				}
			}

			if *mean && len(values) > 0 {
				sum := 0.0
				for _, value := range values {
					sum += value
				}
				mean := sum / float64(len(values))
				fmt.Printf(" Mean = %.2f      |", mean)
			}

			if *max && len(values) > 0 {
				fmt.Printf(" Max = %.2f      |", values[len(values)-1])
			}

			fmt.Println()
			break
		case val := <-valuesChan:
			values = append(values, val)
			break

		}
	}

}
