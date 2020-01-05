package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var employeeArray []EmployeeFormat

var prizeArray []PrizeFormat

var availableNameList []EmployeeFormat

var prizeMap map[int][]EmployeeFormat

type ListPrizesResult struct {
	Action `json:"action"`
	Prizes []PrizeFormat `json:"prizes"`
}

type StartStopResult struct {
	Action  `json:"action"`
	Winners []EmployeeFormat `json:"winners"`
}

func SetUp() {
	employeeArray, prizeArray = getConfig()
	prizeMap = make(map[int][]EmployeeFormat)
	for i := 0; i < len(prizeArray); i++ {
		prizeMap[i] = nil
	}
	availableNameList = employeeArray
}

func prizes(c *Client, action Action) {
	response := ListPrizesResult{action, prizeArray}

	fmt.Printf("prizeArray: %v\n", prizeArray)
	message, _ := json.Marshal(response)
	fmt.Printf("prizeArray: %v\n", string(message))

	c.send <- message
	//	context.JSON(200, gin.H{"prizes": prizeArray})
}

func lottery(num int, available []EmployeeFormat) ([]EmployeeFormat, []EmployeeFormat) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var result []EmployeeFormat

	if (len(available)) < num || num < 1 {
		return nil, nil
	}

	arr := available

	for i := 0; i < num; i++ {
		n := r.Intn(len(arr))
		result = append(result, arr[n])
		arr = append(arr[:n], arr[n+1:]...)
	}

	return arr, result
}

func roll(ctx context.Context, client *Client, action Action, ch chan []EmployeeFormat, mutex *sync.Mutex) {

	mutex.Lock()

	defer func() {
		mutex.Unlock()
	}()

	var (
		//err    error
		result []EmployeeFormat
	)
	//	        defer func() {
	//              if err != nil {
	//                    fmt.Printf("roll() err: %v", err)
	//          }
	//        context.JSON(200, gin.H{"result": result})
	// }()
	prizeIndex := action.PrizeIndex
	var response StartStopResult

	for {
		select {

		case <-ctx.Done():
			ch <- result
			return
		default:
		}

		//str := message
		//fmt.Println("str:" + str)
		//prizeIndex, _ := strconv.Atoi(str)
		num := prizeArray[prizeIndex].Num
		//Already have result
		if prizeMap[prizeIndex] != nil {
			lastNameList := prizeMap[prizeIndex]
			availableNameList = append(availableNameList, lastNameList...)
			prizeMap[prizeIndex] = nil
			//err = errors.New("Duplicate roll")
			//return
		}
		availableNameList, result = lottery(num, availableNameList)
		if result != nil {
			prizeMap[prizeIndex] = result
		} else {
			//err = errors.New("error in lottrey()")
			return
		}

		response.Action = action
		response.Winners = result

		message, _ := json.Marshal(response)

		client.send <- []byte(message)
		time.Sleep(time.Millisecond * 100)
	}

}

func getWinners(client *Client, action Action, mutex *sync.Mutex) {
	mutex.Lock()
	defer mutex.Unlock()

	response := StartStopResult{action, prizeMap[action.PrizeIndex]}
	message, _ := json.Marshal(response)
	client.send <- []byte(message)
}
