package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/rs/usbpower"
)

func main() {
	duration := flag.Duration("duration", 0, "If non-zero, sampling is stopped after specified duration")
	output := flag.String("output", "", "Raw output format, no raw output if not specified")
	flag.Parse()

	// Open the device
	device, err := usbpower.OpenDevice()
	if err != nil {
		fmt.Printf("Error opening device: %v\n", err)
		os.Exit(1)
	}
	defer device.Close()

	voltages := make([]float64, 0)
	intencities := make([]float64, 0)

	// Read samples in a loop
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "")
	start := time.Now()
	for {
		sample, err := device.Read()
		if err != nil {
			panic(err)
		}
		if *output == "json" {
			if err := enc.Encode(sample); err != nil {
				panic(err)
			}
		}
		voltages = append(voltages, float64(sample.Voltage))
		intencities = append(intencities, float64(sample.Current))
		if *duration > 0 && time.Since(start) >= *duration {
			break
		}
	}

	if *output == "" {
		fmt.Println("Statistics:")
		fmt.Printf("Voltage: min=%.2fV, max=%.2fV, avg=%.2fV, p50=%.2fV, p90=%.2fV\n",
			min(voltages), max(voltages), average(voltages), percentile(voltages, 50), percentile(voltages, 90))
		fmt.Printf("Current: min=%.2fA, max=%.2fA, avg=%.2fA, p50=%.2fA, p90=%.2fA\n",
			min(intencities), max(intencities), average(intencities), percentile(intencities, 50), percentile(intencities, 90))
	}
}

func min(numbers []float64) float64 {
	if len(numbers) == 0 {
		return 0
	}
	min := numbers[0]
	for _, n := range numbers {
		if n < min {
			min = n
		}
	}
	return min
}

func max(numbers []float64) float64 {
	if len(numbers) == 0 {
		return 0
	}
	max := numbers[0]
	for _, n := range numbers {
		if n > max {
			max = n
		}
	}
	return max
}

func average(numbers []float64) float64 {
	if len(numbers) == 0 {
		return 0
	}
	var sum float64
	for _, n := range numbers {
		sum += n
	}
	return sum / float64(len(numbers))
}

func percentile(numbers []float64, p int) float64 {
	if len(numbers) == 0 || p < 0 || p > 100 {
		return 0
	}
	numbers = append([]float64{}, numbers...)
	sort.Float64s(numbers)
	index := (p * (len(numbers) - 1)) / 100
	return numbers[index]
}
