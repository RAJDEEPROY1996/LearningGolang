package main

import (
	"fmt"
	"strings"
	"time"
)

var fullName [50]string //other format   var fullName = [50]string{}    or fullName := [50]string{} 	-----> fix size array

func main() {
	conferenceName := "Blockchain Seminar"
	const totalConferenceTicket uint = 50
	var remainingConferenceTicket uint = 50
	var userFirstName, userLastName string
	var userTicketPurchase uint
	var timeOfPurchase []string //slice

	fmt.Printf("Welcome to the %s\n", conferenceName)
	fmt.Printf("Total Conference Tickets: %d\n", totalConferenceTicket)
	i := 0
	for {

		fmt.Print("Enter Your First Name ")
		fmt.Scan(&userFirstName)
		fmt.Print("Enter Your Last Name ")
		fmt.Scan(&userLastName)
		fmt.Print("Enter No of Tickets to buy ")
		fmt.Scan(&userTicketPurchase)

		if userTicketPurchase == 0 {
			println("At least you need one ticket to buy")
			continue
		} else if remainingConferenceTicket < userTicketPurchase {
			println("Less tickets available..")
			continue
		} else {
			t := time.Now().Format("2006 01 02")
			fmt.Println("time is ", t)
			// fmt.Println("ds", t.Format("2006 01 02 MST "))
			fmt.Printf("%T \n", t)
			fullName[i] = userFirstName + " " + userLastName
			timeOfPurchase = append(timeOfPurchase, t)
			remainingConferenceTicket = remainingConferenceTicket - userTicketPurchase
			buyerDetails(i, timeOfPurchase)
			fmt.Println("Remaining tickets are", remainingConferenceTicket)
			if remainingConferenceTicket == 0 {
				break
			}
			i++
		}
	}
}
func buyerDetails(i int, PurchaseTime []string) {
	fmt.Println("Full Name is: ", fullName[i])
	fmt.Println("Total list of Buyer", fullName)
	//range in loop testing
	var firstName string
	for _, value := range fullName {
		nameSlice := strings.Fields(value)
		if len(nameSlice) != 0 { // this checking is required as I have used array instead of slice
			firstName = nameSlice[0]
		}
	}
	fmt.Printf("\tSee You ----- %s ----- in the Conference\n", firstName)
	fmt.Println("Date of Purchase ", PurchaseTime[i])
	fmt.Println("Total Number of Buyer", len(PurchaseTime))
	fmt.Println("Total list of Time of Purchase", PurchaseTime)
}
