package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	queue "unicorn/queue"
	unicorns "unicorn/unicorns"

	"github.com/julienschmidt/httprouter"
)

var CPU_CORES_WORKERS = 2*runtime.NumCPU() + 1

type response struct {
	Message string `json:"message"`
	Stack   string `json:"stack"`
}

var Queue = &queue.QueueList{}
var petNames = unicorns.PetNames{}
var adjectives = unicorns.Adjectives{}

func setUnicornsProduction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("processing new request...")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	rb := struct {
		Amount int `json:"amount"`
	}{}
	err := unmarshal(r.Body, &rb)
	if err != nil {
		errorHandler(w, r, 500, response{
			Stack: fmt.Sprintf("%v", err),
		})
		return
	}

	requestId := rand.Intn(99999)
	unicornList := new(unicorns.UnicornList)
	_queue := new(queue.QueueElement)
	_queue = Queue.Enqueue(_queue)

	_queue.RequestId = requestId
	status := queue.PRODUCTION_STATUS_IN_PROGRESS
	_queue.Status = &status
	_queue.Unicorns = &unicorns.UnicornList{}
	fmt.Printf("Request_id '%v' enqueued on position '%v'\n", requestId, len(*Queue))

	// Concurrency
	ch := make(chan *unicorns.UnicornElement, CPU_CORES_WORKERS)
	for j := 0; j < rb.Amount; j++ {

		name := adjectives[rand.Intn(len(adjectives)-1)] + "-" + petNames[rand.Intn(len(petNames)-1)]
		// Check if unicorn collection doesn't contain name already
		if unicornList.GetUnicornByName(name) {
			fmt.Printf("Unicorn name '%v' repeated\n", name)
			j--
			continue
		}

		go processUnicorn(name, ch)
	}

	// Pushing Unicorns to the LIFO stack
	for i := 0; i < rb.Amount; i++ {
		_unicorn := *<-ch
		unicornList.Push(&_unicorn)
		fmt.Printf("Unicorn '%v' pushed to stack\n", _unicorn.Name)
	}
	fmt.Printf("All the Unicorns pushed to the stack:\n%v", unicornList)

	// Popping Unicorns from the stack into the FIFO Queue
	for i := 0; i < rb.Amount; i++ {
		_unicorn := unicornList.Pop()
		_queue.Unicorns.Push(_unicorn)
		fmt.Printf("Unicorn '%v' popped from stack\n", _unicorn.Name)
	}
	fmt.Printf("All the Unicorns popped:\n%v", _queue.Unicorns)

	status = queue.PRODUCTION_STATUS_READY
	_queue.Status = &status

	fmt.Printf("Unicorns production, request_id '%v' ready...\n", _queue.RequestId)
	d, _ := json.Marshal(_queue)
	w.Write(d)
}

func processUnicorn(name string, ch chan *unicorns.UnicornElement) {
	item := new(unicorns.UnicornElement)
	item.Name = name

	for i := 0; i < 3; i++ {
		idx := rand.Intn(len(unicorns.Capabilities) - 1)
		cap := unicorns.Capabilities[idx]
		if unicorns.GetCapabilityByName(item.Capabilities, cap) {
			fmt.Printf("Capability '%v' repeated in the Unicorn '%v'\n", cap, name)
			i--
			continue
		}
		item.Capabilities = append(item.Capabilities, cap)
	}

	prodTime := time.Duration(rand.Intn(1000)) * time.Millisecond
	fmt.Printf("Start producing Unicorn '%v'\n", name)
	time.Sleep(prodTime)
	fmt.Printf("Unicorn '%v' produced in %v\n", name, prodTime)
	ch <- item
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

func getRequestDetail(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println("retrieving request detail by requestId")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	requestId, _ := strconv.Atoi(params.ByName("request_id"))

	_queue := Queue.FindQueueByRequestId(requestId)
	if _queue == nil {
		errorHandler(w, r, 404, response{
			Message: "Request not found",
		})
		return
	}
	d, _ := json.Marshal(_queue)
	w.Write(d)
}

func getAllRequest(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("retrieving request detail by requestId")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	d, _ := json.Marshal(*Queue)
	w.Write(d)
}

func deliveryPackage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println("deliverying unicorns package by requestId")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	requestId, _ := strconv.Atoi(params.ByName("request_id"))

	_queue := Queue.FindQueueFirstPosition(requestId)
	if _queue == nil || *_queue.Status != queue.PRODUCTION_STATUS_READY {
		errorHandler(w, r, 500, response{
			Message: fmt.Sprintf("Package request_id '%v' not available yet", requestId),
		})
		return
	}

	d, _ := json.Marshal(struct {
		Message string             `json:"message"`
		Package queue.QueueElement `json:"package"`
	}{
		Message: "Package deliveried sucessfully",
		Package: *Queue.Dequeue(),
	})
	w.Write(d)
}

func cleanQueue(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("deliverying unicorns package by requestId")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	Queue = &queue.QueueList{}
	d, _ := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: "Queue cleaned sucessfully",
	})
	w.Write(d)
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int, msg response) {
	w.WriteHeader(status)
	d, _ := json.Marshal(msg)
	w.Write(d)
}

func main() {
	var err error
	petNames, err = unicorns.GetPetNames()
	if err != nil {
		os.Exit(1)
	}

	adjectives, err = unicorns.GetAdjectives()
	if err != nil {
		os.Exit(1)
	}

	router := httprouter.New()
	router.POST("/api/set-unicorns-production", setUnicornsProduction)
	router.GET("/api/get-request-detail/:request_id", getRequestDetail)
	router.GET("/api/get-all-request", getAllRequest)
	router.PUT("/api/delivery-package/:request_id", deliveryPackage)
	router.DELETE("/api/clean-queue", cleanQueue)

	PORT := ":8888"
	if os.Getenv("API_PORT") != "" {
		PORT = os.Getenv("API_PORT")
	}

	fmt.Printf("Unicorn Factory is firing this up on port '%v'\n", PORT)
	err = http.ListenAndServe(PORT, router)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
