package main

import (
	"fmt"
	"os"
	"time"

	"solvm/vm"
)

func benchmark() {
	config := vm.Config{
		Timeout:       time.Second * 30,
		Debug:         false,
		Trace:         false,
		MemoryLimit:   1024 * 1024 * 100,
		MaxGoroutines: 1000,
		WorkingDir:    ".",
	}
	vm := vm.NewSolVM(config)
	defer vm.Close()

	benchmarks := []struct {
		name       string
		code       string
		iterations int
	}{
		{
			name: "Simple Arithmetic",
			code: `
				local sum = 0
				for i = 1, 1000000 do
					sum = sum + i
				end
				return sum
			`,
			iterations: 5,
		},
		{
			name: "Table Operations",
			code: `
				local t = {}
				for i = 1, 100000 do
					t[i] = i * 2
				end
				local sum = 0
				for _, v in pairs(t) do
					sum = sum + v
				end
				return sum
			`,
			iterations: 5,
		},
		{
			name: "Function Calls",
			code: `
				local function fib(n)
					if n <= 1 then return n end
					return fib(n-1) + fib(n-2)
				end
				return fib(20)
			`,
			iterations: 5,
		},
		{
			name: "String Operations",
			code: `
				local s = ""
				for i = 1, 10000 do
					s = s .. "test" .. i
				end
				return #s
			`,
			iterations: 5,
		},
	}

	results := make(map[string][]time.Duration)

	for _, bench := range benchmarks {
		fmt.Printf("\nRunning benchmark: %s\n", bench.name)
		results[bench.name] = make([]time.Duration, bench.iterations)

		for i := 0; i < bench.iterations; i++ {
			start := time.Now()
			err := vm.LoadString(bench.code)
			duration := time.Since(start)
			results[bench.name][i] = duration

			if err != nil {
				fmt.Printf("Error in iteration %d: %v\n", i+1, err)
				os.Exit(1)
			}

			fmt.Printf("  Iteration %d: %v\n", i+1, duration)
		}
	}

	fmt.Println("\nBenchmark Results Summary:")
	fmt.Println("=========================")
	for name, durations := range results {
		var total time.Duration
		for _, d := range durations {
			total += d
		}
		avg := total / time.Duration(len(durations))
		fmt.Printf("%s:\n", name)
		fmt.Printf("  Average: %v\n", avg)
		fmt.Printf("  Min: %v\n", min(durations))
		fmt.Printf("  Max: %v\n", max(durations))
	}
}

func min(durations []time.Duration) time.Duration {
	min := durations[0]
	for _, d := range durations[1:] {
		if d < min {
			min = d
		}
	}
	return min
}

func max(durations []time.Duration) time.Duration {
	max := durations[0]
	for _, d := range durations[1:] {
		if d > max {
			max = d
		}
	}
	return max
}

func main() {
	benchmark()
}
