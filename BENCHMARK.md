# SolVM Performance Benchmarks

This document details the performance characteristics of SolVM, a Lua runtime implemented in Go. The benchmarks aim to provide insights into the execution speed of common Lua operations within the SolVM environment. Understanding these metrics can help developers optimize their Lua scripts and anticipate SolVM's behavior under various workloads.

## Benchmark Suite Overview

The benchmark suite, located in `benchmark/benchmark.go`, comprises tests designed to stress different facets of Lua execution:

*   **Simple Arithmetic**: This test focuses on the raw computational throughput of SolVM for basic mathematical operations. It involves a loop summing integers up to one million, primarily testing CPU-bound performance for numerical tasks.
*   **Table Operations**: Lua tables are fundamental data structures. This benchmark evaluates the efficiency of creating tables, populating them with a large number of elements (100,000 key-value pairs), and then iterating through the table to sum its values. This tests memory allocation, hash table performance, and iteration speed.
*   **Function Calls**: The overhead associated with function invocation is critical for many applications. This benchmark uses a recursive Fibonacci function (`fib(20)`) to measure the cost of repeated function calls and Lua's call stack management.
*   **String Operations**: String manipulation is a common task. This test measures the performance of repeated string concatenation, a potentially expensive operation in many scripting languages, by building a long string from 10,000 smaller segments.

Each benchmark is executed five times to gather a range of timings, allowing for the calculation of average, minimum, and maximum execution durations.

## Running the Benchmarks

To execute the benchmark suite on your own system, navigate to the SolVM project root and run:

```bash
go run benchmark/benchmark.go
```

## Benchmark Results and Performance Analysis

The following results were obtained from a sample run. Timings can vary based on the underlying hardware and system load.

```
Running benchmark: Simple Arithmetic
  Iteration 1: 31.1518ms
  Iteration 2: 23.5224ms
  Iteration 3: 24.1932ms
  Iteration 4: 26.413ms
  Iteration 5: 25.0678ms

Running benchmark: Table Operations
  Iteration 1: 20.4479ms
  Iteration 2: 18.9732ms
  Iteration 3: 17.5482ms
  Iteration 4: 16.2188ms
  Iteration 5: 17.5624ms

Running benchmark: Function Calls
  Iteration 1: 2.035ms    
  Iteration 2: 1.0022ms
  Iteration 3: 1.5321ms
  Iteration 4: 1.9994ms
  Iteration 5: 1.11ms      

Running benchmark: String Operations
  Iteration 1: 111.8972ms
  Iteration 2: 100.1595ms
  Iteration 3: 96.002ms
  Iteration 4: 101.6598ms
  Iteration 5: 99.4323ms

Benchmark Results Summary:
=========================
Table Operations:
  Average: 18.1501ms
  Min: 16.2188ms
  Max: 20.4479ms
Function Calls:
  Average: 1.53918ms
  Min: 1.0022ms
  Max: 2.035ms
String Operations:
  Average: 101.83016ms
  Min: 96.002ms
  Max: 111.8972ms
Simple Arithmetic:
  Average: 26.06964ms
  Min: 23.5224ms
  Max: 31.1518ms
```

### Detailed Performance Discussion:

**Simple Arithmetic:**
The task of summing one million integers completed with an average time of **26.07ms**. The individual iterations ranged from a minimum of **23.52ms** to a maximum of **31.15ms**. This indicates reasonably consistent and efficient handling of CPU-intensive numerical loops within SolVM. The performance here is largely dependent on the underlying gopher-lua VM's ability to optimize tight loops and arithmetic operations.

**Table Operations:**
Manipulating a table with 100,000 entries (insertions followed by iteration and summation) averaged **18.15ms**. The execution times were quite stable, varying between **16.22ms** and **20.45ms**. This suggests that SolVM, via gopher-lua, provides efficient table implementation for moderately large datasets, including quick key hashing, value storage, and iteration. This is crucial for many Lua applications that heavily rely on tables for data structuring.

**Function Calls:**
Executing the recursive `fib(20)` function, which involves a significant number of function calls, took an average of **1.54ms**. The timings for this test showed more variability, ranging from **1.00ms** to **2.04ms**. This relatively fast execution for a recursive task demonstrates that Lua function call overhead within SolVM is low. Efficient function dispatch is key for structuring complex applications and for algorithms that naturally lend themselves to recursion.


**String Operations:**
Concatenating "test" with an iterator 10,000 times to form a long string averaged **101.83ms**, with a range from **96.00ms** to **111.90ms**. This operation is notably slower than the others, which is typical for string concatenation in many scripting languages that use immutable strings. Each concatenation likely creates a new string object, leading to memory allocation and copying overhead. For performance-critical sections involving extensive string building, alternative strategies (like using a table of strings and `table.concat` at the end, if available and idiomatic within SolVM's provided Lua environment) might be considered.

## Interpreting the Numbers

These benchmarks provide a snapshot of SolVM's performance characteristics:

*   **CPU-Bound Tasks (Arithmetic):** SolVM shows good performance for numerical computations.
*   **Data Structures (Tables):** Table operations are efficient for moderately sized collections.
*   **Control Flow (Function Calls):** Lua function call overhead is minimal, supporting modular and recursive programming styles.
*   **String Manipulation:** Repeated string concatenation is relatively expensive, a common trait that developers should be mindful of when dealing with large string-building operations.

The performance is fundamentally tied to the `gopher-lua` library that SolVM uses as its Lua VM. The Go runtime itself also plays a role in how efficiently the overall SolVM process executes.

## SolVM Configuration for Benchmarks

The benchmarks were run with the following SolVM configuration:

*   **Timeout**: 30 seconds
*   **Debug**: false
*   **Trace**: false
*   **MemoryLimit**: 100MB (1024 * 1024 * 100 bytes)
*   **MaxGoroutines**: 1000
*   **WorkingDir**: "."

These settings ensure that the benchmarks are not artificially constrained by very low resource limits but also operate within a defined boundary.

This detailed analysis should help in understanding SolVM's strengths and areas where script optimization might be beneficial.
