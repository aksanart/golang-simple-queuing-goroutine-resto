package main

import (
	"math/rand"
	"time"

	"github.com/fatih/color"
)

type RestoStruct struct {
	Name             string
	IsOpen           bool
	NumberOfStaff    int
	NumberOfCustomer int
	StaffDoneChan    chan bool
	CustDoneChan     chan bool
	PesananChan      chan *PesananStruct
	CustomerChan     chan *CustomerStruct
	WaiterChan       chan *WaiterStruct
	ChefChan         chan *ChefStruct
}

type CustomerStruct struct {
	Name                 string
	PesananDeliveredChan chan bool
}

type ChefStruct struct {
	Name string
}

type WaiterStruct struct {
	Name string
}

type PesananStruct struct {
	NamePelanggan   string
	MenuPesanan     string
	PesananDoneChan chan bool
}

func (resto *RestoStruct) workingWaiter(waiter *WaiterStruct) {
	resto.NumberOfStaff++
	go func() {
		color.Cyan("%s is a waiter, waiting for new order", waiter.Name)
		IsWorking := true
		for {
			if len(resto.CustomerChan) == 0 {
				color.Yellow("there is no new order, %s is take time to play game", waiter.Name)
				IsWorking = false
			}
			customer, ok := <-resto.CustomerChan
			// means resto is closed
			if !ok {
				resto.staffBackToHome(waiter.Name)
				return
			}
			if !IsWorking {
				color.Cyan("%s stop playing game and heading to customer", waiter.Name)
				IsWorking = true
			}
			color.Cyan("%s picking up new order from %s", waiter.Name, customer.Name)
			time.Sleep(2 * time.Second)
			color.Cyan("%s is sending new order to kitchen", waiter.Name)
			menus := []string{"nasi kuning", "coto makassar", "ikan bakar"}
			random := interval(0, 2)
			pesanan := PesananStruct{
				NamePelanggan:   customer.Name,
				MenuPesanan:     menus[random],
				PesananDoneChan: make(chan bool),
			}
			resto.sendOrderToKitchen(&pesanan)
			resto.waitingCustomer(customer)
			<-pesanan.PesananDoneChan
			color.Green("deliver the order %s, to %s", pesanan.MenuPesanan, pesanan.NamePelanggan)
			customer.PesananDeliveredChan <- true
			time.Sleep(1 * time.Second)

		}
	}()
}

func (resto *RestoStruct) sendOrderToKitchen(order *PesananStruct) {
	color.Yellow("waiter hanging note order %s from %s on kitchen", order.MenuPesanan, order.NamePelanggan)
	resto.PesananChan <- order
}

func (resto *RestoStruct) workingChef(chef *ChefStruct) {
	resto.NumberOfStaff++
	go func() {
		color.Cyan("%s is a chef, waiting for new order", chef.Name)
		IsWorking := true
		for {
			if len(resto.PesananChan) == 0 {
				color.Yellow("there is no new order, %s is take time to sleep", chef.Name)
				IsWorking = false
			}
			pesanan, ok := <-resto.PesananChan
			// means resto is closed
			if !ok {
				resto.staffBackToHome(chef.Name)
				return
			}

			if !IsWorking {
				color.Cyan("%s wakeup from ringing bell and picking up the order note ", chef.Name)
				IsWorking = true
			}
			color.Cyan("%s is cooking new order %s for %s", chef.Name, pesanan.MenuPesanan, pesanan.NamePelanggan)
			time.Sleep(3 * time.Second)
			color.Green("%s is done cooking new order %s for %s", chef.Name, pesanan.MenuPesanan, pesanan.NamePelanggan)
			pesanan.PesananDoneChan <- true

		}
	}()
}

func (resto *RestoStruct) staffBackToHome(name string) {
	color.Cyan("%s is back to home", name)
	resto.StaffDoneChan <- true
}

func (resto *RestoStruct) AddCustomer(name string) {
	color.Yellow("%s is new client", name)
	if !resto.IsOpen {
		color.Red("the resto already closed %s is leave", name)
		return
	}
	select {
	case resto.CustomerChan <- &CustomerStruct{
		Name:                 name,
		PesananDeliveredChan: make(chan bool),
	}:
		color.Yellow("%s take a seat and ordering", name)
		resto.NumberOfCustomer++
	default:
		color.White("the seat is full so %s is waiting for empty seat", name)
		resto.CustomerChan <- &CustomerStruct{
			Name:                 name,
			PesananDeliveredChan: make(chan bool),
		}
		color.Yellow("%s take a seat and ordering", name)
	}
}

func (resto *RestoStruct) waitingCustomer(customer *CustomerStruct) {
	go func() {
		<-customer.PesananDeliveredChan
		color.Cyan("%s is eating", customer.Name)
		time.Sleep(5 * time.Second)
		color.Green("%s is now leave", customer.Name)
		resto.CustDoneChan <- true
	}()
}

func interval(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func (resto *RestoStruct) closing() {
	color.Red("now is the resto closing time, stopping new customer")
	resto.IsOpen = false
	close(resto.CustomerChan)
	for x := 0; x < resto.NumberOfCustomer; x++ {
		<-resto.CustDoneChan
	}
	close(resto.CustDoneChan)
	close(resto.PesananChan)
	for i := 0; i < resto.NumberOfStaff; i++ {
		<-resto.StaffDoneChan
	}
	close(resto.WaiterChan)
	close(resto.ChefChan)
	close(resto.StaffDoneChan)
	color.Green("the resto is closed, all people is leave")
}
