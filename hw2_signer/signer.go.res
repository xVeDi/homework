package main

import (
	"fmt"
	"strconv"
	"sync"
)

func SingleHash(in, out chan interface{}) {
	for i := range in {
		fmt.Println("SH", i)
		out <- strconv.Itoa(i.(int))
	}
}

func MultiHash(in, out chan interface{}) {
	for i := range in {
		fmt.Println("MH", i)
		out <- i
	}
}

func CombineResults(in, out chan interface{}) {
	res := ""
	for i := range in {
		res += i.(string)
	}
	out <- res
	fmt.Println("CH", res)

}

func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}
	in := make(chan interface{}, MaxInputDataLen)
	out := make(chan interface{}, MaxInputDataLen)

	for _, jbs := range jobs {
		wg.Add(1)
		go starter(in, out, wg, jbs)
		in = out
		out = make(chan interface{}, MaxInputDataLen)
	}
	wg.Wait()
}

func starter(in, out chan interface{}, wg *sync.WaitGroup, currentJob job) {
	defer wg.Done()
	defer close(out)
	currentJob(in, out)
}

func main() {

	inputData := []int{0, 1, 1, 2, 3, 5, 8}
	//inputData := []int{0, 1}
	testResult := "NOT_SET"

	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				out <- fibNum
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(func(in, out chan interface{}) {
			dataRaw := <-in
			data, ok := dataRaw.(string)
			if !ok {
				//t.Error("cant convert result data to string")
				panic(ok)
			}
			testResult = data
		}),
	}
	ExecutePipeline(hashSignJobs...)

	println(testResult)
}
