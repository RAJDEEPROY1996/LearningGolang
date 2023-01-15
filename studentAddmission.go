package main

import (
	"fmt"
)

type studentInfos struct {
	firstNameS string
	lastNameS  string
	emailS     string
	dobS       int
	mobS       string
}

var studentDetailsList = make([]studentInfos, 0)

func addmissionSection(remainingSeat uint) uint {
	var seat uint
	fmt.Println("Welcome to Addmission Section")
	seat = studentDetails(remainingSeat)
	return seat

}

func studentDetails(remainingSeat uint) uint {
	var firstName, lastName, email string
	var dob int
	var mob string
	fmt.Println("Enter First Name of Student :")
	fmt.Scan(&firstName)
	fmt.Println("Enter Last Name of Student :")
	fmt.Scan(&lastName)
	fmt.Println("Enter Email of Student :")
	fmt.Scan(&email)
	fmt.Println("Enter Mobile No of Student :")
	fmt.Scan(&mob)
	fmt.Println("Enter Date of Birth (ddMMyyyy) No of Student :")
	fmt.Scan(&dob)
	studentInfo := studentInfos{
		firstNameS: firstName,
		lastNameS:  lastName,
		emailS:     email,
		dobS:       dob,
		mobS:       mob,
	}
	studentDetailsList = append(studentDetailsList, studentInfo)
	remainingSeat = remainingSeat - 1
	return remainingSeat
}
