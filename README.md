# Benchmarking Redis vs NATS (JetStream)

> [!IMPORTANT]
> Disclaimer: This is just a quick benchmark and should not be taken as the definitive truth. Numerous factors such as network latency, system hardware, configuration settings, and environmental conditions can significantly affect the results. This benchmark is simplified for the sake of a basic comparison between Redis and NATS in a Pub/Sub use case. While Redis is not primarily designed for Pub/Sub at scale, this test is intended as a basic performance indicator under stress conditions. For real-world use cases, further tuning and testing are necessary.

## Table of Contents
- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [How to Benchmark](#how-to-benchmark)
- [Benchmarking Methodology](#benchmarking-methodology)
- [Results](#results)

## Overview
This project benchmarks **Redis Pub/Sub** against **NATS (JetStream)** in a high-traffic microservice scenario. The goal is to compare **latency, throughput and resource usage** under stress conditions.

## Prerequisites
Before running the benchmarks, ensure that you have the following software installed:

- **Docker**: A platform for developing, shipping, and running applications in containers.
  - Install Docker from [here](https://docs.docker.com/get-docker/).
  
- **Docker Compose**: A tool for defining and running multi-container Docker applications.
  - Install Docker Compose from [here](https://docs.docker.com/compose/install/).

## How to Benchmark
### **1. Clone the Repository**
```sh
git clone https://github.com/Phillezi/redis-vs-nats.git
cd redis-vs-nats
```

### **2. Run Benchmarks**
```sh
make compose/bench
```

## Benchmarking Methodology
Each benchmark measures:
- **Latency:** Time taken per message (`time.Since(start)`).
- **Throughput:** Messages processed per second.

## Results

> [!IMPORTANT]
> Measurements are only from the client.

## Redis Benchmarks

| Test Type   | Total Time  | Throughput (msg/sec) | Avg Latency | Samples | Avg Goroutines (Min/Max) | Avg Heap (Min/Max) | Avg Stack (Min/Max) | Total GC Cycles | Avg Sys Mem (Min/Max) |
|------------|------------|----------------------|-------------|---------|-------------------------|--------------------|--------------------|----------------|--------------------|
| **Sequential** | 33.45s      | 29,892.89            | 16.69s       | 3,345   | 5 (5/5)                 | 12,063 KB (8,155 KB / 16,172 KB) | 626 KB (512 KB / 640 KB) | 111,975          | 24,883 KB (15,511 KB / 24,983 KB) |
| **Parallel**   | 7.76s       | 128,927.66           | 3.88s        | 775     | 105 (70/106)            | 13,317 KB (9,084 KB / 17,632 KB) | 1,754 KB (1,632 KB / 1,856 KB) | 74,316           | 33,175 KB (33,175 KB / 33,175 KB) |

## NATS Benchmarks

| Test Type   | Total Time  | Throughput (msg/sec) | Avg Latency | Samples | Avg Goroutines (Min/Max) | Avg Heap (Min/Max) | Avg Stack (Min/Max) | Total GC Cycles | Avg Sys Mem (Min/Max) |
|------------|------------|----------------------|-------------|---------|-------------------------|--------------------|--------------------|----------------|--------------------|
| **Sequential** | 493.56ms    | 2,026,108.39         | 250.08ms     | 49      | 6 (6/6)                 | 12,601 KB (8,963 KB / 16,277 KB) | 506 KB (480 KB / 512 KB) | 535              | 24,596 KB (19,607 KB / 24,727 KB) |
| **Parallel**   | 516.61ms    | 1,935,680.65         | 260.08ms     | 51      | 63 (6/105)              | 12,690 KB (8,760 KB / 16,449 KB) | 1,099 KB (768 KB / 1,216 KB) | 1,608            | 33,420 KB (33,175 KB / 33,431 KB) |

---

### **Notes**
- **Throughput** represents the number of messages processed per second.
- **Latency** is the average time per message.
- **GC Cycles** shows the total number of garbage collection events during the benchmark.
- **Memory Metrics** (Heap, Stack, and System Memory) are in **KB**
