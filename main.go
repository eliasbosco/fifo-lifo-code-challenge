package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

const (
	PRODUCTION_STATUS_IN_PROGRESS string = "in progress"
	PRODUCTION_STATUS_READY              = "ready"
)

var CPU_CORES_WORKERS = 2*runtime.NumCPU() + 1

type CapabilitiesList = []string

var Capabilities = CapabilitiesList{}

// LIFO stack
type UnicornElement struct {
	Name         string           `json:"name"`
	Capabilities CapabilitiesList `json:"capabilities"`
}

type UnicornList []UnicornElement

func (u *UnicornList) Push(unicorn *UnicornElement) {
	*u = append(*u, *unicorn)
}

func (u *UnicornList) Pop() *UnicornElement {
	uni := *u
	if len(*u) > 0 {
		res := uni[len(uni)-1]
		*u = uni[:len(uni)-1]
		return &res
	}
	return nil
}

type QueueElement struct {
	Unicorns  *UnicornList `json:"unicorns"`
	RequestId int          `json:"request_id"`
	Status    *string      `json:"status"`
}

type QueueList []QueueElement

var Queue = &QueueList{}

func (q *QueueList) Enqueue(value *QueueElement) *QueueElement {
	queue := *q
	queue = append(queue, *value)
	*q = queue
	return &queue[len(queue)-1]
}

func (q *QueueList) Dequeue() *QueueElement {
	queue := *q
	if len(*q) > 0 {
		dequeued := queue[0]
		*q = queue[1:]
		return &dequeued
	}
	return nil
}

func (q *QueueList) FindByRequestIdFirstPosition(requestId int) *QueueElement {
	queue := *q
	for i, item := range queue {
		if i == 0 && item.RequestId == requestId {
			return &item
		} else {
			return nil
		}
	}
	return nil
}

type response struct {
	Message string `json:"message"`
	Stack   string `json:"stack"`
}

var Names = []string{""}
var Adjectives = []string{""}

func setUnicornsProduction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("processing new request...")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	rb := struct {
		Amount int `json:"amount"`
	}{}
	err := unmarshal(r.Body, &rb)
	if err != nil {
		d, _ := json.Marshal(response{
			Stack: fmt.Sprintf("%v", err),
		})
		w.Write(d)
		return
	}

	requestId := rand.Intn(99999)
	unicornList := new(UnicornList)
	queue := new(QueueElement)
	fmt.Println("1 - %v", &queue)
	queue = Queue.Enqueue(queue)
	fmt.Println("3 - %v", &queue)

	queue.RequestId = requestId
	status := PRODUCTION_STATUS_IN_PROGRESS
	queue.Status = &status
	queue.Unicorns = &UnicornList{}

	fmt.Println("request_id:", requestId, " created.")

	// Concurrency
	ch := make(chan *UnicornElement, CPU_CORES_WORKERS)
	for j := 0; j < rb.Amount; j++ {

		name := Adjectives[rand.Intn(len(Adjectives)-1)] + "-" + Names[rand.Intn(len(Names)-1)]
		// Check if unicorn collection doesn't contain name already
		if unicornList.getUnicornByName(name) {
			j--
			continue
		}

		go processUnicorn(name, ch)
	}

	for i := 0; i < rb.Amount; i++ {
		unicornList.Push(<-ch)
	}
	fmt.Println(unicornList)

	for i := 0; i < rb.Amount; i++ {
		queue.Unicorns.Push(unicornList.Pop())
	}
	fmt.Println(queue.Unicorns)

	status = PRODUCTION_STATUS_READY
	queue.Status = &status

	fmt.Println(Queue)

	fmt.Println("Unicorns production ready...")
	d, _ := json.Marshal(queue)
	w.Write(d)
}

func processUnicorn(name string, ch chan *UnicornElement) {
	item := new(UnicornElement)
	item.Name = name

	for i := 0; i < 3; i++ {
		idx := rand.Intn(len(Capabilities) - 1)
		cap := Capabilities[idx]
		if getCapabilityByName(item.Capabilities, cap) {
			i--
			continue
		}
		item.Capabilities = append(item.Capabilities, cap)
	}

	fmt.Println("Unicorn", name, " is in production")
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	fmt.Println("Unicorn", name, " produced")
	ch <- item
}

func getCapabilityByName(capabilitiesList CapabilitiesList, name string) bool {
	for _, item := range capabilitiesList {
		if item == name {
			return true
		}
	}
	return false
}

func (u *UnicornList) getUnicornByName(name string) bool {
	for _, item := range *u {
		if item.Name == name {
			return true
		}
	}
	return false
}

func (q *QueueList) getQueueByRequestId(requestId int) *QueueElement {
	for _, item := range *q {
		if item.RequestId == requestId {
			return &item
		}
	}
	return nil
}

func getRequestDetail(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println("retrieving request detail by requestId")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	d, _ := json.Marshal(&Queue)
	w.Write(d)
}

func getAllRequest(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("retrieving request detail by requestId")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	d, _ := json.Marshal(*Queue)
	w.Write(d)
}

func unmarshal(input io.Reader, marshalStruct interface{}) error {
	dec := json.NewDecoder(input)
	dec.DisallowUnknownFields()

	err := dec.Decode(&marshalStruct)
	if err != nil {
		return err
	}
	return nil
}

func getArrayStringFromTextFile(fileName string) ([]string, error) {
	fn, err := os.Open(fileName)
	defer fn.Close()
	if err != nil {
		return nil, err
	}
	var names []string
	var scanner = bufio.NewScanner(fn)
	for scanner.Scan() {
		names = append(names, scanner.Text())
	}
	return names, nil
}

func deliveryPackage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println("deliverying unicorns package by requestId")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	requestId, _ := strconv.Atoi(params.ByName("request_id"))

	queue := Queue.FindByRequestIdFirstPosition(requestId)
	if queue == nil || *queue.Status != PRODUCTION_STATUS_READY {
		d, _ := json.Marshal(response{
			Message: fmt.Sprint("Package request_id: ", requestId, " not available."),
		})
		w.Write(d)
		return
	}

	d, _ := json.Marshal(struct {
		Message string       `json:"message"`
		Package QueueElement `json:"package"`
	}{
		Message: "Package deliveried sucessfully",
		Package: *Queue.Dequeue(),
	})
	w.Write(d)
}

func cleanQueue(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("deliverying unicorns package by requestId")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	Queue = &QueueList{}
	d, _ := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: "Queue cleaned sucessfully",
	})
	w.Write(d)
}

func main() {
	Capabilities = append(Capabilities, "super strong", "fullfill wishes", "fighting capabilities", "fly", "swim", "sing", "run", "cry", "change color", "talk", "dance", "code", "design", "drive", "walk", "talk chinese", "lazy")

	var err error
	Names, err = getArrayStringFromTextFile("petnames.txt")
	if err != nil {
		fmt.Println("Unicorn names not found")
		os.Exit(1)
	}

	Adjectives, err = getArrayStringFromTextFile("adj.txt")
	if err != nil {
		if err != nil {
			fmt.Println("Unicorn adjectives not found")
			os.Exit(1)
		}
	}

	router := httprouter.New()
	router.POST("/api/set-unicorns-production", setUnicornsProduction)
	router.GET("/api/get-request-detail/:request_id", getRequestDetail)
	router.GET("/api/get-all-request", getAllRequest)
	router.PUT("/api/delivery-package/:request_id", deliveryPackage)
	router.DELETE("/api/clean-queue", cleanQueue)

	http.ListenAndServe(":8888", router)
}
