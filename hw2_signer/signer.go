package main

import (
	"fmt"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var quotaCh chan struct{}

func SingleHash(in, out chan interface{}) {

	fmt.Println("SHin :")
	data := <-in
	fmt.Println("SHin :", data)
	// workerInput := make(chan interface{}, 1)
	workerOutputCrc32 := make(chan interface{}, 1)
	workerOutputMd5 := make(chan interface{}, 1)
	// workerOutput3 := make(chan interface{})
	go workerDataSignerCrc32(data.(string), workerOutputCrc32)
	runtime.Gosched()
	go workerDataSignerMd5(data.(string), workerOutputMd5)
	runtime.Gosched()
	// md := DataSignerMd5(data)
	crc32res := <-workerOutputCrc32
	runtime.Gosched()
	md5res := <-workerOutputMd5
	runtime.Gosched()

	a := crc32res.(string) + "~" + DataSignerCrc32(md5res.(string))
	//println(a)
	out <- a
	return
}

func MultiHash(in, out chan interface{}) {
	var mp [6]string
	// for i := range in {
	data := <-in
	//	fmt.Println(data)
	wg := &sync.WaitGroup{}

	for i := 0; i < 6; i++ {
		//	mp[i] = DataSignerCrc32(strconv.Itoa(i) + data.(string))
		newstr := strconv.Itoa(i) + data.(string)

		wg.Add(1)
		go func(data2 string, i int) {
			defer wg.Done()
			// fmt.Println("closure data2: ", data2)

			// ch1 := make(chan interface{}, 1)
			// ch2 := make(chan interface{}, 1)
			// //	fmt.Println("mp before", mp)
			// ch1 <- data2
			// //	fmt.Println("mp before", mp)
			// go DataSignerCrc32(ch1, ch2)
			// mu := &sync.Mutex{}
			// res := <-ch2
			// fmt.Println(res)
			// fmt.Println("mp after", mp)
			// mu.Lock()
			// mp[i] = res.(string)
			// mu.Unlock()
			// mu := &sync.Mutex{}
			// mu.Lock()
			mp[i] = DataSignerCrc32(data2)
			// mu.Unlock()

		}(newstr, i)

	}
	wg.Wait()

	a := ""
	for _, v := range mp {
		a += v
	}
	// fmt.Println(a)
	out <- a
}

func CombineResults(in, out chan interface{}) {

	var inn []string
	fmt.Println("cr start")
	for i := range in {
		fmt.Println("cr: ", i)
		inn = append(inn, i.(string))

	}
	sort.Strings(inn)
	fmt.Println(inn)
	a := ""
	for _, v := range inn {
		a = a + "_" + v

	}
	a = strings.TrimLeft(a, "_")
	fmt.Println(a)
	out <- a
}

func ExecutePipeline(jbs ...job) {
	quotaCh = make(chan struct{}, 1)
	die := make(chan interface{})
	chan1 := make(chan interface{}, len(jbs))
	chan2 := make(chan interface{}, len(jbs))
	// chan3 := make(chan interface{}, len(jbs))
	// chan4 := make(chan interface{})
	wg := &sync.WaitGroup{}
	// ctx, finish := context.WithCancel(context.Background())

	for _, jb := range jbs {
		go jb(chan1, chan1)
		go workerSingleHash(chan1, chan2)
		wg.Add(1)
		go workerSingleHash(chan1, chan2, die, wg)
		// wg.Add(1)
		// go workerMultiHash(chan2, chan3, wg)
	}
	wg.Wait()
	close(die)
	// finish()
	// CombineResults(chan3, chan4)
	time.Sleep(5 * time.Second)
	// close(chan1)
	close(chan2)
	// close(chan3)
	for ii := range chan2 {
		fmt.Println(ii, "ss")

	}

}

func workerSingleHash(in, out, die chan interface{}, wg *sync.WaitGroup) {

	workerInput := make(chan interface{}, 1)
	workerOutput := make(chan interface{}, 1)
	i := <-in
	//	for i := range in {
	fmt.Println(i)
	ii, ok := i.(int)

	if !ok {
		fmt.Println("Не удалось преобразовать к типу *int")
	}
	fmt.Println("worker in: ", ii)
	workerInput <- strconv.Itoa(ii)
	fmt.Println("worker in: ", ii)
	go SingleHash(workerInput, workerOutput)
	res := <-workerOutput
	wg.Done()
	fmt.Println("wokkerHS res: ", res)

	select {
	case <-die:
		return
	default:
		out <- workerOutput
	}

	//runtime.Gosched()
	//			fmt.Printf("%T\n", i)

}

// func workerMultiHash(in, out chan interface{}, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	workerInput := make(chan interface{}, MaxInputDataLen)
// 	//	workerOutput := make(chan interface{}, MaxInputDataLen)

// 	for i := range in {
// 		workerInput <- i
// 		fmt.Println("wmh ", i)
// 		go MultiHash(workerInput, out)
// 		runtime.Gosched()
// 		//			fmt.Printf("%T\n", i)
// 	}

// }

func workerMultiHash(in, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	workerInput := make(chan interface{}, MaxInputDataLen)
	strr := <-in
	workerInput <- strr.(string)
	go MultiHash(workerInput, out)
	runtime.Gosched()

}

func workerDataSignerCrc32(in string, out chan interface{}) {
	out <- DataSignerCrc32(in)
	return
}

func workerDataSignerMd5(in string, out chan interface{}) {
	// println("start")
	quotaCh <- struct{}{}
	// println("middle")
	out <- DataSignerMd5(in)
	<-quotaCh
	// println("stop")
	return
}

func main() {

	inputData := []int{0, 1}

	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				out <- fibNum
			}
		}),
		// job(func(in, out chan interface{}) {
		// 	for i := 50; i < 55; i++ {
		// 		out <- i
		// 	}
		// }),
	}

	ExecutePipeline(hashSignJobs...)

	//	fmt.Scanln()
	fmt.Println("done")

	// res := []string{"4958044192186797981418233587017209679042592862002427381542", "29568666068035183841425683795340791879727309630931025356555"}
	// res := []string{}
	// res = append(res, MultiHash(SingleHash("1")))
	// res = append(res, MultiHash(SingleHash("0")))
	// res = append(res, MultiHash(SingleHash("1")))
	// res = append(res, MultiHash(SingleHash("2")))
	// res = append(res, MultiHash(SingleHash("3")))
	// res = append(res, MultiHash(SingleHash("5")))
	// res = append(res, MultiHash(SingleHash("8")))

	// fmt.Println(CombineResults(res))
}
