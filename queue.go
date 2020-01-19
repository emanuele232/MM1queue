package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"

	"github.com/wcharczuk/go-chart"
)

//system status
var custormers_served int
var IS_BUSY bool
var simulation_clock float64
var num_in_queue int

//events
var next_arrival_time, next_departure_time float64
var next_event_type int
var next_time float64
var event_interarrival_time float64
var last_event_time float64
var time_arrivals []float64
var interarrivals []float64

//statistical variables
var total_delays float64
var customers_delayed int
var area_num_in_queue float64
var area_server_status float64

var max_customers_served = 1000

var lambda_service float64 = 2.0
var lambda_interarrival float64 = 1.0

func main() {

	initialize()

	for custormers_served < max_customers_served {
		//determine the next event
		next_event()

		switch next_event_type {
		case 1:
			arrive()
		case 2:
			departure()
		}

		update_stats()

	}

	plot(interarrivals, "arr.png")

	report()

}

func initialize() {
	//first event must be an arrive
	next_arrival_time = exp(lambda_interarrival)
	next_departure_time = 1.0e10
	last_event_time = 0.0

	num_in_queue = 0

	interarrivals = append(interarrivals, next_arrival_time)
}

func next_event() {
	next_event_type = 0

	if next_arrival_time < next_departure_time {
		next_event_type = 1
		next_time = next_arrival_time
	} else {
		next_event_type = 2
		next_time = next_departure_time
	}

	simulation_clock = next_time

}

func arrive() {
	//we schedule the next arrival
	next_arrival_time = simulation_clock + exp(lambda_interarrival)
	interarrivals = append(interarrivals, next_arrival_time-simulation_clock)

	if IS_BUSY {
		//server busy
		num_in_queue = num_in_queue + 1
		time_arrivals = append(time_arrivals, simulation_clock)

	} else {
		IS_BUSY = true
		next_departure_time = simulation_clock + exp(lambda_service)
	}

}

func departure() {
	if num_in_queue == 0 {
		IS_BUSY = false
		next_departure_time = 1.0e10
	} else {
		//queue not empty
		num_in_queue = num_in_queue - 1

		//add the delay of this departure to the total delays

		total_delays += simulation_clock - time_arrivals[0]
		customers_delayed++

		next_departure_time = simulation_clock + exp(lambda_service)
		time_arrivals = time_arrivals[1:]

		//add served customers
		custormers_served = custormers_served + 1

	}

}

func update_stats() {
	event_interarrival_time = simulation_clock - last_event_time
	last_event_time = simulation_clock

	//adding the time in which we have customers in queue times
	area_num_in_queue += float64(num_in_queue) * event_interarrival_time

	//adding the time in which the server is busy
	var busy float64
	if IS_BUSY {
		busy = 0.0
	} else {
		busy = 1.0
	}
	area_server_status = area_server_status + (busy * event_interarrival_time)

}

func report() {
	fmt.Println("average delay in queue:")
	fmt.Println(total_delays / float64(customers_delayed))

	fmt.Println("average number in queue:")
	fmt.Println(area_num_in_queue / simulation_clock)

	fmt.Println("server utilization:")
	fmt.Println(area_server_status / simulation_clock)

}

//utility functions
func exp(lambda float64) float64 {
	//return an exponential distribuited variate
	return rand.ExpFloat64() / lambda

}

func float_to_string(number float64) string {
	return strconv.FormatFloat(number, 'f', 6, 64)
}

func plot(XAxis []float64, filename string) {
	var rangee []float64
	for i := 0.0; i < 100.0; i++ {
		rangee = append(rangee, i)
	}

	sort.Slice(XAxis, func(i, j int) bool {
		return XAxis[i] > XAxis[j]
	})

	graph := chart.Chart{
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: XAxis[0:100],
				YValues: rangee,
			},
		},
	}
	f, _ := os.Create(filename)
	defer f.Close()
	graph.Render(chart.PNG, f)
}
