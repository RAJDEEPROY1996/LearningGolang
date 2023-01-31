package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

var tpl *template.Template
var db *sql.DB

func getMySQLDB() *sql.DB {
	db, err := sql.Open("mysql", "root:@(127.0.0.1:3306)/tickets?parseTime=true")
	if err != nil {
		log.Fatal((err))
	}
	return db
}

const totalConferenceTicket int = 50

type ticketDetails struct {
	FirstName   string
	LastName    string
	Address     string
	Email       string
	NoOfTickets int
}

var purchasedTickets []ticketDetails

func main() {
	tpl, _ = template.ParseGlob("../static/*.html")
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/postForm", postFormHandler)
	http.HandleFunc("/processPost", processPostHandler)
	http.HandleFunc("/ticketDetails", ticketDetailsHandler)
	http.HandleFunc("/fetchDetail", fetchDetailHandler)
	fmt.Println("Connecting to Server ...")
	http.ListenAndServe(":8000", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", remainingConferenceTickets())
}

func remainingConferenceTickets() int {
	db = getMySQLDB()
	var total int
	rows, err := db.Query("select count(*) from conference")
	if err != nil {
		return -1
	} else {
		for rows.Next() {
			rows.Scan(&total)
		}
		return totalConferenceTicket - total
	}
}

func postFormHandler(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "postForm.html", nil)
}
func processPostHandler(w http.ResponseWriter, r *http.Request) {
	db = getMySQLDB()
	var d ticketDetails
	d.FirstName = r.FormValue("firstName")
	d.LastName = r.FormValue("lastName")
	d.Address = r.FormValue("address")
	d.Email = r.FormValue("email")
	e := r.FormValue("noOfTickets")
	d.NoOfTickets, _ = strconv.Atoi(e)
	t, s := checking(d)
	if !t {
		tpl.ExecuteTemplate(w, "failed.html", s)
	} else {
		purchasedTickets = append(purchasedTickets, d)
		fmt.Println(purchasedTickets)
		_, err := db.Exec("insert into conference(firstName,lastName,address,email,noOfTickets)values(?,?,?,?,?)", d.FirstName, d.LastName, d.Address, d.Email, d.NoOfTickets)
		if err != nil {
			fmt.Println(err.Error())
			tpl.ExecuteTemplate(w, "failed.html", err.Error())
		} else {
			tpl.ExecuteTemplate(w, "thanks.html", d)
		}
	}
}
func checking(d ticketDetails) (bool, string) {
	if len(d.FirstName) < 2 || len(d.LastName) < 2 || len(d.Address) < 3 {
		return false, "less character passed"
	}
	if d.NoOfTickets != 1 {
		return false, "Sorry !!! You can Buy only 1 ticket"
	}
	for _, item := range purchasedTickets {
		if item.Email == d.Email {
			fmt.Println("Only Once")
			return false, "Already Booked"
		}
	}

	if remainingConferenceTickets()-d.NoOfTickets < 0 {
		fmt.Println("Low Tickets available")
		return false, "Low Tickets available"
	}
	return true, "OK"
}

func ticketDetailsHandler(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "ticketDetails.html", nil)
}

func fetchDetailHandler(w http.ResponseWriter, r *http.Request) {
	db = getMySQLDB()
	var d ticketDetails
	d.Email = r.FormValue("email")

	fmt.Println(r.FormValue("submit"))
	if r.FormValue("submit") == "fetch" {
		err := db.QueryRow("select * from conference where email = ?", d.Email).Scan(&d.FirstName, &d.LastName, &d.Address, &d.Email, &d.NoOfTickets) //for query
		if err != nil {
			fmt.Println(err.Error())
			tpl.ExecuteTemplate(w, "failed.html", "You have no booking in the conference")
		} else {
			tpl.ExecuteTemplate(w, "thanks.html", d)
		}
	} else if r.FormValue("submit") == "correction" {
		d.NoOfTickets = 1
		d.Email = r.FormValue("email2")
		fmt.Println("email", d.Email)
		d.FirstName = r.FormValue("firstName")
		d.LastName = r.FormValue("lastName")
		d.Address = r.FormValue("address")
		result, _ := db.Exec("update conference set firstName = ?, lastName = ?, address =? where email = ?", d.FirstName, d.LastName, d.Address, d.Email)
		output, err := result.RowsAffected()
		if err != nil {
			tpl.ExecuteTemplate(w, "failed.html", "Modification failed. Kindly Try again later")
		} else {
			if output == 0 {
				tpl.ExecuteTemplate(w, "failed.html", "This Email Id doesnot exist")
			} else {
				tpl.ExecuteTemplate(w, "thanks.html", d)
			}
		}
	}
}
