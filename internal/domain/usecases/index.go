package usecases

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"time"
)

type St struct {
}

func New() *St {
	return &St{}
}

func (u *St) Run(uri string, reqCount, workerCount int) {
	if workerCount > reqCount {
		workerCount = reqCount
	}

	workerReqCount := int(math.Ceil(float64(reqCount) / float64(workerCount)))

	var sentCount int

	done := make(chan time.Duration, reqCount)

	startTime := time.Now()

	for i := 0; i < workerCount; i++ {
		if sentCount+workerReqCount > reqCount {
			workerReqCount = reqCount - sentCount
		}

		go u.runRoutine(done, uri, workerReqCount)

		sentCount += workerReqCount
	}

	var maxDur time.Duration
	var failCount int
	var successTotalDur time.Duration

	for i := 0; i < sentCount; i++ {
		dur := <-done

		if dur < 0 {
			failCount++
		} else {
			successTotalDur += dur

			if dur > maxDur {
				maxDur = dur
			}
		}
	}

	totalDur := time.Now().Sub(startTime)

	avgDur := time.Duration(int64(successTotalDur) / int64(sentCount - failCount))

	fmt.Printf("Sent count:\t\t%2d\n", sentCount)
	fmt.Printf("Worker count:\t\t%2d\n", workerCount)
	fmt.Printf("Success count:\t\t%2d\n", sentCount-failCount)
	fmt.Printf("Fail count:\t\t%2d\n", failCount)
	fmt.Printf("Total duration:\t\t%s\n", totalDur.String())
	fmt.Printf("Avg duration:\t\t%s\n", avgDur.String())
	fmt.Printf("Max duration:\t\t%s\n", maxDur.String())
}

func (u *St) runRoutine(done chan<- time.Duration, uri string, reqCount int) {
	for i := 0; i < reqCount; i++ {
		startTime := time.Now()

		if u.sendRequest(uri) {
			done <- time.Now().Sub(startTime)
		} else {
			done <- -1
		}
	}
}

func (u *St) sendRequest(uri string) bool {
	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		fmt.Println("Fail to create request", err)
		return false
	}

	client := &http.Client{Timeout: 20 * time.Second}

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("Fail to do request", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		fmt.Println("Bad response code:", resp.StatusCode)
		return false
	}

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Fail to read response-body", err)
		return false
	}

	return true
}
