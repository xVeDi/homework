package main

import (
	"fmt"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// channel for control running of DataSignerMd5
var quotaChan chan struct{}

type singleHashStruct struct {
	wg        sync.WaitGroup // wait group for control calculating crc32data and crcmd5
	crc32data string         // result of crc32(data)
	crcmd5    string         // result of crc32(md5(data))
}

func (sh *singleHashStruct) calcSingleHashCrc32(data string) {
	defer sh.wg.Done()
	sh.crc32data = DataSignerCrc32(data)
	fmt.Println(data, " SingleHash crc32(data) ", sh.crc32data)
}

func (sh *singleHashStruct) calcSingleHashMd5(data string, qch chan struct{}) {
	defer sh.wg.Done()

	// control accessibility of DataSignerMd5
	qch <- struct{}{}
	md5Hash := DataSignerMd5(data)
	<-qch

	fmt.Println(data, " SingleHash md5(data) ", md5Hash)
	sh.crcmd5 = DataSignerCrc32(md5Hash)
	fmt.Println(data, " crc32(md5(data)) ", sh.crcmd5)
}

func singleHashCalsulate(data string, out chan interface{}, wg *sync.WaitGroup, qchan chan struct{}) {
	defer wg.Done()

	fmt.Println(data, " SingleHash data ", data)

	shObject := new(singleHashStruct)
	shObject.wg.Add(1)
	go shObject.calcSingleHashMd5(data, qchan)
	shObject.wg.Add(1)
	go shObject.calcSingleHashCrc32(data)
	shObject.wg.Wait()

	a := shObject.crc32data + "~" + shObject.crcmd5
	out <- a
	fmt.Println(data, " result ", a)
}

func SingleHash(in, out chan interface{}) {
	waitGroupSH := &sync.WaitGroup{}
	quotaChan = make(chan struct{}, 1)
	for i := range in {
		data := strconv.Itoa(i.(int))
		waitGroupSH.Add(1)
		go singleHashCalsulate(data, out, waitGroupSH, quotaChan)
		runtime.Gosched()
	}
	waitGroupSH.Wait()
}

func MultiHash(in, out chan interface{}) {

	waitGroupMH := &sync.WaitGroup{}

	for i := range in {
		data := i.(string)

		waitGroupMH.Add(1)
		go multiHashCalculator(data, out, waitGroupMH)
		runtime.Gosched()
	}
	waitGroupMH.Wait()
}

func multiHashCalculator(data string, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	wgForMHcalc := &sync.WaitGroup{}
	var mp [6]string

	for i := 0; i < 6; i++ {
		wgForMHcalc.Add(1)

		go func(data string, i int, wg *sync.WaitGroup) {
			defer wg.Done()
			mp[i] = DataSignerCrc32(data)
			runtime.Gosched()
		}(strconv.Itoa(i)+data, i, wgForMHcalc)
	}

	wgForMHcalc.Wait()

	a := ""
	for i, v := range mp {
		fmt.Println(data, " MultiHash: crc32(th+step1)) ", i, " ", v)
		a += v
	}
	// fmt.Println(a)
	out <- a
}

func CombineResults(in, out chan interface{}) {

	var inn []string

	for i := range in {
		inn = append(inn, i.(string))
	}
	sort.Strings(inn)

	a := ""
	for _, v := range inn {
		a = a + "_" + v

	}
	a = strings.TrimLeft(a, "_")
	fmt.Println("CombineResults ", a)
	out <- a
}

func ExecutePipeline(jbs ...job) {
	mainWaitGroup := &sync.WaitGroup{}
	in := make(chan interface{})
	out := make(chan interface{})
	for i, currentJob := range jbs {
		mainWaitGroup.Add(1)
		go starter(currentJob, in, out, mainWaitGroup)
		in = out
		if i != len(jbs)-1 {
			out = make(chan interface{})
		}
	}
	mainWaitGroup.Wait()
}

func starter(currentJob job, in, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(out)
	currentJob(in, out)
}
