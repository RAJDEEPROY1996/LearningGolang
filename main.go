package main

import (
	"fmt"
	teacher "schoolAddmissionSystem/TeacherDetails"
)

func main() {
	var choice uint
	var remainingSeat uint = 50
	for {
		fmt.Printf("\tPress -> 1 to see School Details \n \tPress -> 2 to see Teacher details\n \tPress -> 3 Addmission details\n \tPress -> 4 to see Remaining Seats\n \tPress -> 5 to see all Student's Detail List\n")
		fmt.Println("\t XXXXXXXXXXXXXXXXXX--- Press any no greater than 5 to exit ---XXXXXXXXXXXXXXXXXX")
		fmt.Println("Enter Your Choice:")
		fmt.Scan(&choice)
		switch choice {
		case 1:
			schoolDetails()
		case 2:
			teacher.TeacherDetails()
		case 3:
			remainingSeat = addmissionSection(remainingSeat)
		case 4:
			fmt.Printf("Remaining seat: %d\n", remainingSeat)
		case 5:
			fmt.Println("List of Addmitted Student in the School", studentDetailsList)
		default:
			fmt.Println("Wrong Choice")
		}

		if choice > 5 || remainingSeat == 0 {
			break
		}
	}
}
