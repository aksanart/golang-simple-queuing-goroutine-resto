package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

var seatCapacity = 4
var openTime = 10 * time.Second

func main() {
	color.Green("| Resto terbuka |")
	color.Green("| ------------- |")

	customerChan := make(chan *CustomerStruct, seatCapacity)
	closedResto := make(chan bool)
	exit := make(chan bool)
	pesananChan := make(chan *PesananStruct)
	waiterChan := make(chan *WaiterStruct)
	chefChan := make(chan *ChefStruct)

	resto := RestoStruct{
		Name:             "Resto dubidubi dam dam",
		IsOpen:           true,
		NumberOfStaff:    0,
		NumberOfCustomer: 0,
		StaffDoneChan:    make(chan bool),
		CustDoneChan:     make(chan bool),
		PesananChan:      pesananChan,
		CustomerChan:     customerChan,
		WaiterChan:       waiterChan,
		ChefChan:         chefChan,
	}

	go func() {
		<-time.After(openTime)
		closedResto <- true
		resto.closing()
		exit <- true
	}()

	// new waiter
	waiters := []*WaiterStruct{
		{Name: "waiter-1"},
		{Name: "waiter-2"},
	}
	for _, waiter := range waiters {
		go resto.workingWaiter(waiter)
	}

	// new Chef
	chefs := []*ChefStruct{
		{Name: "chef-1"},
		{Name: "chef-2"},
	}
	for _, chef := range chefs {
		go resto.workingChef(chef)
	}

	// add customer
	go func() {
		i := 1
		for {
			select {
			case <-closedResto:
				return
			case <-time.After(time.Second * time.Duration(interval(1, 3))):
				resto.AddCustomer(fmt.Sprintf("customer-%d", i))
				i++
			}
		}
	}()

	// wait until resto is closed, and all people leave
	<-exit
}
