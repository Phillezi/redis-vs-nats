package bench

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/Phillezi/redis-vs-nats/pkg/messaging"
)

const (
	numMessages = 1000000
	topic       = "benchmark"
	publishers  = 100
	subscribers = 100
)

type ResourceSample struct {
	Timestamp  time.Time
	Goroutines int
	HeapAlloc  uint64
	StackInuse uint64
	GCCount    uint32
	SystemMem  uint64
}

// Run all scenarios
func RunBenchmarks(brokerName string, connect func() messaging.Broker) {
	fmt.Printf("\nStarting benchmarks for %s...\n", brokerName)

	runBenchmarkWithTracking("Sequential", func() { BenchmarkSequential(connect(), brokerName) })
	runBenchmarkWithTracking("Parallel", func() { BenchmarkParallel(connect(), brokerName) })
	//BenchmarkFanOut(connect(), brokerName)
	//BenchmarkFullMesh(connect(), brokerName)
}

// Scenario 1: Sequential (1 publisher, 1 subscriber)
func BenchmarkSequential(broker messaging.Broker, brokerName string) {
	fmt.Printf("\n🔹 Running Sequential Benchmark (%s)...\n", brokerName)

	var wg sync.WaitGroup
	latencies := make([]time.Duration, numMessages)
	receivedMessages := 0

	startTime := time.Now()

	wg.Add(1)
	err := broker.Subscribe(topic, func(msg []byte) {
		latencies[receivedMessages] = time.Since(startTime)
		receivedMessages++

		if receivedMessages == numMessages {
			wg.Done()
		}
	})
	if err != nil {
		log.Fatalf("[%s] Subscription error: %v\n", brokerName, err)
	}

	for i := range numMessages {
		broker.Publish(topic, fmt.Appendf(nil, "Message %d", i))
	}

	wg.Wait()
	broker.Close()

	printResults(brokerName, "Sequential", startTime, latencies)
}

// Scenario 2: Parallel Publishers (Multiple publishers, 1 subscriber)
func BenchmarkParallel(broker messaging.Broker, brokerName string) {
	fmt.Printf("\n🔹 Running Parallel Benchmark (%s)...\n", brokerName)

	var wg sync.WaitGroup
	latencies := make([]time.Duration, numMessages)
	receivedMessages := 0
	var mu sync.Mutex

	startTime := time.Now()

	wg.Add(1)
	err := broker.Subscribe(topic, func(msg []byte) {

		mu.Lock()
		if receivedMessages >= numMessages {
			log.Println("[ERR] Test is bad")
			mu.Unlock()
			return
		}

		latencies[receivedMessages] = time.Since(startTime)
		receivedMessages++
		mu.Unlock()

		if receivedMessages >= numMessages {
			wg.Done()
		}
	})
	if err != nil {
		log.Fatalf("[%s] Subscription error: %v\n", brokerName, err)
	}

	startPublishing(broker, publishers)

	wg.Wait()
	broker.Close()

	printResults(brokerName, "Parallel", startTime, latencies)
}

// Scenario 3: Fan-out (1 publisher, multiple subscribers)
func BenchmarkFanOut(broker messaging.Broker, brokerName string) {
	fmt.Printf("\n🔹 Running Fan-out Benchmark (%s)...\n", brokerName)

	var wg sync.WaitGroup
	latencies := make([]time.Duration, numMessages*subscribers)
	receivedMessages := 0
	var mu sync.Mutex

	startTime := time.Now()

	for range subscribers {
		wg.Add(1)
		go func() {
			err := broker.Subscribe(topic, func(msg []byte) {
				mu.Lock()
				latencies[receivedMessages] = time.Since(startTime)
				receivedMessages++
				mu.Unlock()

				if receivedMessages == numMessages*subscribers {
					wg.Done()
				}
			})
			if err != nil {
				log.Fatalf("[%s] Subscription error: %v\n", brokerName, err)
			}
		}()
	}

	startPublishing(broker, 1)

	wg.Wait()
	broker.Close()

	printResults(brokerName, "Fan-out", startTime, latencies)
}

// Scenario 4: Full Mesh (Multiple publishers, multiple subscribers)
func BenchmarkFullMesh(broker messaging.Broker, brokerName string) {
	fmt.Printf("\n🔹 Running Full Mesh Benchmark (%s)...\n", brokerName)

	var wg sync.WaitGroup
	latencies := make([]time.Duration, numMessages*subscribers)
	receivedMessages := 0
	var mu sync.Mutex

	startTime := time.Now()

	for range subscribers {
		wg.Add(1)
		go func() {
			err := broker.Subscribe(topic, func(msg []byte) {
				mu.Lock()
				latencies[receivedMessages] = time.Since(startTime)
				receivedMessages++
				mu.Unlock()

				if receivedMessages == numMessages*subscribers {
					wg.Done()
				}
			})
			if err != nil {
				log.Fatalf("[%s] Subscription error: %v\n", brokerName, err)
			}
		}()
	}

	startPublishing(broker, publishers)

	wg.Wait()
	broker.Close()

	printResults(brokerName, "Full Mesh", startTime, latencies)
}

func startPublishing(broker messaging.Broker, pubCount int) {
	var pubWG sync.WaitGroup
	pubWG.Add(pubCount)

	batchSize := numMessages / pubCount

	for i := range pubCount {
		start := i * batchSize
		go func(start int) {
			defer pubWG.Done()
			for j := start; j < start+batchSize; j++ {
				err := broker.Publish(topic, fmt.Appendf(nil, "Message %d", j))
				if err != nil {
					log.Printf("[Publisher] Failed to send message %d: %v", j, err)
				}
			}
		}(start)
	}

	pubWG.Wait()
}

func calculateAvgLatency(latencies []time.Duration) time.Duration {
	var total time.Duration
	for _, l := range latencies {
		total += l
	}
	return total / time.Duration(len(latencies))
}

func printResults(brokerName, scenario string, startTime time.Time, latencies []time.Duration) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	elapsed := time.Since(startTime)
	avgLatency := calculateAvgLatency(latencies)

	fmt.Printf("\n[%s - %s] Benchmark Results:\n", brokerName, scenario)
	fmt.Printf("🕒 Total Time: %v\n", elapsed)
	fmt.Printf("⚡ Throughput: %.2f msg/sec\n", float64(numMessages)/elapsed.Seconds())
	fmt.Printf("⏳ Avg Latency: %v\n", avgLatency)
}

func trackResourceUsage(stopChan chan struct{}, interval time.Duration, samples *[]ResourceSample, mu *sync.Mutex) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)

			sample := ResourceSample{
				Timestamp:  time.Now(),
				Goroutines: runtime.NumGoroutine(),
				HeapAlloc:  memStats.HeapAlloc,
				StackInuse: memStats.StackInuse,
				GCCount:    memStats.NumGC,
				SystemMem:  memStats.Sys,
			}

			mu.Lock()
			*samples = append(*samples, sample)
			mu.Unlock()
		case <-stopChan:
			fmt.Println("\n✅ [Tracking] Stopped resource monitoring.")
			return
		}
	}
}

func runBenchmarkWithTracking(name string, benchmarkFunc func()) {
	var samples []ResourceSample
	var mu sync.Mutex
	stopChan := make(chan struct{})

	go trackResourceUsage(stopChan, 10*time.Millisecond, &samples, &mu) // Track every 10ms

	startTime := time.Now()
	benchmarkFunc()
	elapsed := time.Since(startTime)

	close(stopChan)

	analyzeResourceUsage(samples)

	fmt.Printf("\n🏁 Benchmark '%s' completed in %v\n", name, elapsed)
}

func analyzeResourceUsage(samples []ResourceSample) {
	if len(samples) == 0 {
		fmt.Println("⚠️ No resource samples recorded.")
		return
	}

	var totalGoroutines int
	var totalHeap uint64
	var totalStack uint64
	var totalGC uint32
	var totalSysMem uint64

	minGoroutines, maxGoroutines := samples[0].Goroutines, samples[0].Goroutines
	minHeap, maxHeap := samples[0].HeapAlloc, samples[0].HeapAlloc
	minStack, maxStack := samples[0].StackInuse, samples[0].StackInuse
	minSysMem, maxSysMem := samples[0].SystemMem, samples[0].SystemMem

	for _, sample := range samples {
		totalGoroutines += sample.Goroutines
		totalHeap += sample.HeapAlloc
		totalStack += sample.StackInuse
		totalGC += sample.GCCount
		totalSysMem += sample.SystemMem

		if sample.Goroutines < minGoroutines {
			minGoroutines = sample.Goroutines
		}
		if sample.Goroutines > maxGoroutines {
			maxGoroutines = sample.Goroutines
		}

		if sample.HeapAlloc < minHeap {
			minHeap = sample.HeapAlloc
		}
		if sample.HeapAlloc > maxHeap {
			maxHeap = sample.HeapAlloc
		}

		if sample.StackInuse < minStack {
			minStack = sample.StackInuse
		}
		if sample.StackInuse > maxStack {
			maxStack = sample.StackInuse
		}

		if sample.SystemMem < minSysMem {
			minSysMem = sample.SystemMem
		}
		if sample.SystemMem > maxSysMem {
			maxSysMem = sample.SystemMem
		}
	}

	n := uint64(len(samples))

	fmt.Println("\n📊 Resource Usage Summary:")
	fmt.Printf("🔹 Total Samples: %d\n", n)
	fmt.Printf("🔹 Avg Goroutines: %d (Min: %d, Max: %d)\n", totalGoroutines/int(n), minGoroutines, maxGoroutines)
	fmt.Printf("🔹 Avg Heap Usage: %d KB (Min: %d KB, Max: %d KB)\n", (totalHeap/n)/1024, minHeap/1024, maxHeap/1024)
	fmt.Printf("🔹 Avg Stack Usage: %d KB (Min: %d KB, Max: %d KB)\n", (totalStack/n)/1024, minStack/1024, maxStack/1024)
	fmt.Printf("🔹 Total GC Cycles: %d\n", totalGC)
	fmt.Printf("🔹 Avg System Memory: %d KB (Min: %d KB, Max: %d KB)\n", (totalSysMem/n)/1024, minSysMem/1024, maxSysMem/1024)
}
