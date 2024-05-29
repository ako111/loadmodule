package main

// import "fmt"

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	"math"
	"sync"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: loadmodule <url> <qps> <duration>")
		return
	}

	duration, err := time.ParseDuration(os.Args[3])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	url := os.Args[1] // replace with your URL
	qpsStr := os.Args[2] // replace with your QPS

	qps, err := strconv.Atoi(qpsStr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	sleepDuration := time.Second / time.Duration(qps)
	ticker := time.NewTicker(sleepDuration)
	defer ticker.Stop()

	endTime := time.Now().Add(duration)
	var responseTimes []time.Duration
	var num2xx, num4xx, num5xx int
	var wg sync.WaitGroup

	for i := 0; time.Now().Before(endTime); i++ {
		select {
		case <-ticker.C:
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				start := time.Now()
				resp, err := http.Get(url)
				if err != nil {
					fmt.Println("Error:", err)
					wg.Done()
					return
				}
				defer resp.Body.Close()

				elapsed := time.Since(start)
				responseTimes = append(responseTimes, elapsed)

				switch {
				case resp.StatusCode >= 200 && resp.StatusCode < 300:
					num2xx++
				case resp.StatusCode >= 400 && resp.StatusCode < 500:
					num4xx++
				case resp.StatusCode >= 500:
					num5xx++
				}

				wg.Done()
			}(&wg)
		}
	}

	wg.Wait()

	minTime := responseTimes[0]
	maxTime := responseTimes[0]
	totalTime := time.Duration(0)
	for _, t := range responseTimes {
		if t < minTime {
			minTime = t
		}
		if t > maxTime {
			maxTime = t
		}
		totalTime += t
	}
	avgTime := totalTime / time.Duration(len(responseTimes))

	var sdSquared time.Duration
	for _, t := range responseTimes {
		sdSquared += time.Duration(math.Pow(float64(t-avgTime), 2))
	}
	standardDeviation := time.Duration(math.Sqrt(float64(sdSquared / time.Duration(len(responseTimes)))))

	fmt.Println("Minimum response time:", minTime)
	fmt.Println("Average response time:", avgTime)
	fmt.Println("Maximum response time:", maxTime)
	fmt.Println("Standard deviation of response time:", standardDeviation)
	fmt.Println("2xx responses:", num2xx)
	fmt.Println("4xx responses:", num4xx)
	fmt.Println("5xx responses:", num5xx)
}
