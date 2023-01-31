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

const totalConferenceTicket int = 2

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
	http.HandleFunc("/allTicketDetails", allTicketDetailsHandler)
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
		_, err := db.Exec("insert into conference(firstName,lastName,address,email,noOfTickets)values(?,?,?,?,?)", d.FirstName, d.LastName, d.Address, d.Email, d.NoOfTickets)
		if err != nil {
			fmt.Println(err.Error())
			tpl.ExecuteTemplate(w, "failed.html", "Already Booked using This email")
		} else {
			purchasedTickets = append(purchasedTickets, d)
			fmt.Println(purchasedTickets)
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
		fmt.Println("No Tickets available")
		return false, "No Tickets available"
	}
	return true, "OK"
}

func ticketDetailsHandler(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "ticketDetails.html", nil)
}

func fetchDetailHandler(w http.ResponseWriter, r *http.Request) {
	db = getMySQLDB()
	var d ticketDetails
	if r.FormValue("submit") == "fetch" {
		d.Email = r.FormValue("email")
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
				for i, value := range purchasedTickets {
					if value.Email == d.Email {
						purchasedTickets[i].FirstName = d.FirstName
						purchasedTickets[i].LastName = d.LastName
						purchasedTickets[i].Address = d.Address
					}
				}
				fmt.Println("After Correction: ", purchasedTickets)
				tpl.ExecuteTemplate(w, "thanks.html", d)
			}
		}
	} else if r.FormValue("submit") == "delete" {
		d.Email = r.FormValue("email3")
		result, _ := db.Exec("delete from conference where email = ?", d.Email)
		output, err := result.RowsAffected()
		if err != nil {
			tpl.ExecuteTemplate(w, "failed.html", "Ticket Cancelling failed. Kindly Try again later")
		} else {
			if output == 0 {
				tpl.ExecuteTemplate(w, "failed.html", "This Email Id doesnot exist")
			} else {
				for i, value := range purchasedTickets {
					if value.Email == d.Email {
						purchasedTickets = append(purchasedTickets[:i], purchasedTickets[i+1:]...)
					}
				}
				d.NoOfTickets = 0
				fmt.Println("After Deletion: ", purchasedTickets)
				tpl.ExecuteTemplate(w, "thanks.html", d)
			}
		}
	}
}

func allTicketDetailsHandler(w http.ResponseWriter, r *http.Request) {
	db = getMySQLDB()
	var d []ticketDetails
	res, err := db.Query("select * from conference")
	if err != nil {
		log.Fatal(err)
		tpl.ExecuteTemplate(w, "failed.html", "Error while fetching Data from Database")
	} else {
		var e ticketDetails
		for res.Next() {
			err := res.Scan(&e.FirstName, &e.LastName, &e.Address, &e.Email, &e.NoOfTickets)
			if err != nil {
				log.Fatal(err)
			}
			d = append(d, e)
		}
		tpl.ExecuteTemplate(w, "allTicketDetails.html", d)
	}

}
