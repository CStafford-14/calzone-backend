package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var devMode bool
var eventsLibrary = "/var/lib/calzone/"
var eventTypes []string = []string{"Meeting", "Cascade", "Computer Setup", "Appointment", "Other"}

func main() {
	// Capture the desired port. Defaults to 8090.
	port := flag.String("port", "8090", "Port to be used for the server")
	devModeRaw := flag.Bool("dev", false, "Development mode (Routes events to working directory)")
	flag.Parse()

	if *devModeRaw {
		devMode = true
		eventsLibrary = ""
	}

	// Start the Webserver
	staticFileServer := http.FileServer(http.Dir("./websource/static/"))

	http.Handle("/add/", staticFileServer)
	http.Handle("/manage/", staticFileServer)
	http.Handle("/today/", staticFileServer)
	http.Handle("/common/", staticFileServer)
	http.HandleFunc("/api", apiHandler)
	http.HandleFunc("/cal", calendarBuilder)
	http.HandleFunc("/mod", removalHandler)
	http.Handle("/lib/", http.StripPrefix("/lib/", http.FileServer(http.Dir(eventsLibrary))))
	http.Handle("/", http.RedirectHandler("/add", http.StatusSeeOther))

	if devMode {
		fmt.Print(strings.Join([]string{"Starting server at port ", *port, " in development mode\n"}, ""))
	} else {
		fmt.Print(strings.Join([]string{"Starting server at port ", *port, "\n"}, ""))
	}
	if err := http.ListenAndServe(strings.Join([]string{":", *port}, ""), nil); err != nil {
		log.Fatal(err)
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	errStatus := 0

	if r.URL.Path != "/api" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" && r.Method != "GET" {
		http.Error(w, "Method is not supported.\n\nInstead, please enjoy this ASCII Rick Astley.\n\ntttttfttttttttffffftfffftfftttttttttt111tttt1111111ttttttt111tttt111111tttttttttttttttttttt1111111tt\nttfttttttttttttttttffLLftfffffffftttt1t1111111t11tfffffftttt1111111111111111ttftttftttttttt111111ttt\nttttttttttttttttffffLftttfffffffLLfftt1tttttt111tfffffffftttttttt11111tttttffffftffLfttttttt11111ttt\nttttttttttttttttfffffttffffffLLLfffttt1ttttttt1tffffffffffttttttt11111tfffffffLffttfffffttt1111ttttt\nttttttttttttffffttfffttffffffffftttftttttt11ttt1ttfffffftttttt11tt111111ttffffLLLfttfLLft111111ttttt\nttttttttttffLLLLfttttffLLLfftttttfffttttt1tfffft1ttfffffttfffft11t111tt111ttfffLLftttftt1ttt11tttttt\nttttttttfffLLLLLLffttfLLLfttfffftfLftttttttfffffftttffttffffffftt1111tftt111ttffLLfttt11tffftttttttt\ntttttttfffLLLLLLLLLfttfttttfLLffftffttttt111i;iitfftttfffffffffftt111tft1tttttttfLLfttttfffffftttttt\ntttttttffLLLLLLLLLfttttffffffffftfftttt1;:,::,,,:;ittttffffffffftt1111t11ttttttttttttttffffffffffttt\ntttffftfLLLLLLLLLLfttfftfffttfffLLLftt1i:,,,,,,,,,,:1tttttfffffftt1111tft1ttt1tffttt11tffffLLLLLfttt\nttffffffffLLLLLLLfftfLLftttfLLffLLLftt1i;;;;;;;i;:,:tftt11tfffft11111tffft1tt11ttttt1ttfffLLLLLLLfff\ntfffffttttfLLLffttttfffffttLLLLfLLLft1;iiiii11111i::tffft11tttt11t111tfffttfftt1ttfftttffLLLLLLLLffL\ntttfftttttffLftttttttttttttfLLLffLLftt11;;;iiiiiii:;tffffft11tt1111111tfftfffftttttttttttfLLLLffffff\nttfLffffffftttfffffffffttffttLLfffLft11i;;;;i;;iii;tffffttt1tffft11111tfttfffttt1tttttttttfLLftffttt\nttffffffffftttffffffffftffffttfftfttt11i;;;i1iii1iitffft11tttttfft1111tt1tffttffttfffffLLfffftfLffLL\nttffffffffftttffffffffttttttttttffftt11i;;;iii111ii1tt111tfft1ttft111111t1ttttfftttffffLLLfttfLLLLLL\ntttffffffffttfffffffftt1ttt111ttffftttti;;;;iiiiii11111ttfffft11t11111tfftttttttt1tfffffLLftttfLLLLf\nttttttttfffttfffttttt1ttfft1111ttffttt1i;;;iiii;;itt11tffffffft1111111ttttt1ttfft11ttffffLfftfLLLLLf\ntftttttttttttt1tttt11tfffft1tt11tfft111;;;;;iii;;it111ttfffffft11t1111tttt1tt1tffftt1111ttfftfLffftt\ntfttffffftt11ttftfftttffftttfft11tti;;1;;;:;iii::i111111ttfft111tt1111ttt11tf1tffffftttttttttttttfft\n111tfffffftttffftttt11tfftttft1i;:,..,1i;;;;iii;i;:;;1t111t11111tt1111tt111tft1fffffttfffftt11tffLLf\n1tttffffft11ttffftt1111tt1ii:,,......,i1i;;;ii;11;...,:;i111tt11111111111t11tf1tfffttfffffffttfLffLf\ntfttttftt11t11ttt11ttt11;,,..........,11t1iiii1t1:,,.....,:i111111t1111111111t1tffttttffffttt1tfffLf\ntffttttttttfttt111ttttti,.............:;;iii11t1i,,.........,ittt11111111tttt11tttfLffttffttttttffft\nttttttttttttttt111tttt1:................:;;;;;i;:,...........:ttt1111111tttffttttffLLffttttffftttttf\ntfffttttftttttt11ttttti.................,;;;ii;;:............,1ttt11111ttttffftttffffffttttttttffttf\ntffft1ttfffftf111ttttti..................:;;;i;;:,...........,1tt11111111ttffttttffLLfLfttfffffffttf\ntffft1ttfttttt111ttttt;..................,;;;;;::,...........,1t1111111111tttttttffLLfLfttfffffffttf\ntffft1tttttttt111ttttt;...........,,.....,;;:;:::............,111111111tt111tttttfffLLLfttfffftffttf\nttttt1tttttttt111tttti,........,;ii:......:;:::::............,ittttt111ttt11tft1tfffffLftttfftfffttt\nttttt1tttt111t111tttt;........:;;;;:,.....,;:::::.............ittttt1111t1111tt1ttfttttttttttttffttt\nttttt1tttttttt11tttt1,.......,:::;;;,......::::::.............;tt1t1111111ttt1tttfttttttt1tttttffttt\n1tttt1tttttttt11tttti.........,:::;:.......,:::::..........,,:;1t11111111ttff1111ttttttt111ttttffttt\n111111tttttt111111t1,..........,:::........,:::::..........,:;;i1111t111ttttt11111ttttt1111ttttttt1t\n1111111111111111111;............,,.........,::::,.........,:;;;i1111111111111111111111111111111tt111\n11111111111111111111:.......................,:::,.. ......,:;;ii111111111111111111111111111111111111\n111111111111111111111:.......,..............,:::,..........,::i11t1111111111111111111111111111111111\n1111111111111111111111:,,..,, ..............,,,::,............,1111111111111111111111111111111111111\n1111111111111111111111111i11:...............,..,:,,.....    ..:1t11111111111111111111111111111111111\n111111111111111111111111111;................,::::,,.....::::;i11tttt111111t1111111111111111111111111\n111111111111111111111111111,.................ii;;::,....:1tttt1ttttttttttttttttt11tt111111ttttt11111\n1111111111111111111111111ti..................;i;;;:,.....;t1ttttttttttttttttttttt1ttttttttttttt11111", http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		throw(w, fmt.Sprintf("ParseForm() err: %v", err), true)
		return
	}

	evTime := r.FormValue("evTime") // HH:MM in 24-hour time
	evTimeEnd := r.FormValue("evTimeEnd")
	evDate := r.FormValue("evDate") // YYYY-MM-DD
	evType := r.FormValue("evType") // int, see documentation
	evName := r.FormValue("evName") // string. This is just some text to describe the event.
	user := r.FormValue("user")     // Hidden form element with the username and password... which I'll somehow use.

	if len(evTime) != 5 || (len(evTimeEnd) != 5 && len(evTimeEnd) != 0) || len(evDate) != 10 || len(evType) != 1 {
		throw(w, "Invalid request", true)
		fmt.Println("An error has occured in length test.")
		return
	}

	// Don't call this yet... it will be used to validate time when it's done.
	// valid := validateTime(evTime, evTimeEnd)

	// if !valid {
	// 	throw(w, "Invalid request: Invalid time", true)
	// 	fmt.Println("Time not valid")
	// 	return
	// }

	if len(evName) > 100 {
		throw(w, "Invalid request: Name exceeds 100 characters", true)
		fmt.Println("An error has occured in length test. Name exceeds 100 characters")
		return
	}

	if evTime != "" && evDate != "" && evType != "" && evName != "" && user != "" {
		fmt.Println("Event request by " + user + ":")
		fmt.Println(evDate + " at " + evTime + ":")
		fmt.Println(evName + "(" + evType + ")")

		hour := evTime[:2]
		minute := evTime[3:]
		hourEnd := "-1"
		minuteEnd := "-1"
		if len(evTimeEnd) == 5 {
			hourEnd = evTimeEnd[:2]
			minuteEnd = evTimeEnd[3:]
		}

		day := evDate[8:10]
		month := evDate[5:7]
		year := evDate[:4]

		eventLine := hour + "," + minute + "," + hourEnd + "," + minuteEnd + "," + day + "," + evType + ",\" " + evName + "\"\n"
		monthFile := month + "_" + year + ".csv"

		fmt.Println("Adding line to " + monthFile)
		fmt.Println(eventLine)

		if devMode {
			eventsLibrary = ""
		}

		arc, err := os.OpenFile(eventsLibrary+monthFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}

		if _, err := arc.Write([]byte(eventLine)); err != nil {
			log.Fatal(err)
		}
		if err := arc.Close(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("")

		throw(w, "Calendar event added", errStatus != 0)
	} else {
		throw(w, "Form not sent correctly.", true)
	}
}

// func authHandler(w http.ResponseWriter, r *http.Request) {

// 	if err := r.ParseForm(); err != nil {
// 		fmt.Fprintf(w, "ParseForm() err: %v", err)
// 		return
// 	}

// 	username := r.FormValue("username")
// 	// password_raw := r.FormValue("password")

// 	fmt.Println("Started session for " + username)

// 	fmt.Fprintf(w, "<h1>Empty response body</h1>")
// }

func throw(w http.ResponseWriter, m string, isError bool) {
	if isError {
		fmt.Fprintf(w, "<h1 id=\"response\" class=\"showError\">There was an error:</h1><br><h1 id=\"response\" class=\"showError\">"+m+"</h1>")
	} else {
		fmt.Fprintf(w, "<h1 id=\"response\" class=\"showOK\">"+m+"</h1>")
	}
}

// WIP Code that DOES NOT work. Don't call this function until I say so
func validateTime(start string, end string) bool {
	result := false

	startInt, errStart := strconv.Atoi(start)
	endInt, errEnd := strconv.Atoi(end)

	if errStart != nil || errEnd != nil {
		fmt.Println("Error while converting string to interger:")
		result = true
	}

	if startInt > endInt && startInt != endInt {
		result = true
	} else if startInt == endInt {
		result = true
	}
	return result
}

func calendarBuilder(w http.ResponseWriter, r *http.Request) {
	year := fmt.Sprint(time.Now().Year())
	nextYear := fmt.Sprint(int(time.Now().Year()) + 1)

	monthInt := int(time.Now().Month())
	monthNum := fmt.Sprint(monthInt)
	monthStr := time.Now().Month().String()

	var nextMonthInt int
	if monthInt == 12 {
		nextMonthInt = 1
	} else {
		nextMonthInt = monthInt + 1
		nextYear = year
	}
	nextMonthNum := fmt.Sprint(nextMonthInt)
	nextMonthStr := time.Month(nextMonthInt).String()

	if len(monthNum) != 2 {
		monthNum = "0" + monthNum
	}

	if len(nextMonthNum) != 2 {
		nextMonthNum = "0" + nextMonthNum
	}

	thisMonthCal, err := os.Open(eventsLibrary + monthNum + "_" + year + ".csv")
	if err != nil {
		fmt.Println(err)
	}
	nextMonthCal, err := os.Open(eventsLibrary + nextMonthNum + "_" + nextYear + ".csv")
	if err != nil {
		fmt.Print(err)
	}

	thisIOR := csv.NewReader(thisMonthCal)
	nextIOR := csv.NewReader(nextMonthCal)

	thisCal, _ := thisIOR.ReadAll()
	nextCal, _ := nextIOR.ReadAll()

	var response string
	var thisResponse string
	var nextResponse string

	for i := 0; i < len(thisCal); i++ {
		timeOfDay := "AM"
		timeOfEnd := "AM"
		var start string

		title := strings.Replace(strings.Replace(thisCal[i][6], ">", "&gt;", -1), "<", "&lt;", -1)
		startInt, err := strconv.Atoi(thisCal[i][0])
		if err != nil {
			fmt.Println(err)
			return
		} else if startInt > 12 && startInt != 12 {
			start = fmt.Sprint(startInt-12) + ":" + thisCal[i][1]
			timeOfDay = "PM"
		} else if startInt == 12 {
			timeOfDay = "PM"
			start = thisCal[i][0] + ":" + thisCal[i][1]
		} else {
			start = thisCal[i][0] + ":" + thisCal[i][1]
		}

		end := "-1"
		if thisCal[i][2] != "-1" {
			endInt, err := strconv.Atoi(thisCal[i][2])
			if err != nil {
				fmt.Println(err)
				return
			} else if endInt > 12 && endInt != 12 {
				end = fmt.Sprint(endInt-12) + ":" + thisCal[i][3]
				timeOfEnd = "PM"
			} else if endInt == 12 {
				timeOfEnd = "PM"
				end = thisCal[i][2] + ":" + thisCal[i][3]
			} else {
				end = thisCal[i][2] + ":" + thisCal[i][3]
			}
		}
		evDay := thisCal[i][4]
		evTypeInt, err := strconv.Atoi(thisCal[i][5])
		if err != nil {
			fmt.Println(err)
			return
		}
		evTypeStr := eventTypes[evTypeInt]

		record := "No events for this month."
		if end != "-1" {
			record = "<p><span class=\"evTitle\">" + title + "</span><p>\n<p>" + start + " " + timeOfDay + " - " + end + " " + timeOfEnd + "</p>\n<p>" + evTypeStr + " - " + monthStr + " " + evDay + "</p>"
		} else {
			record = "<p><span class=\"evTitle\">" + title + "</span><p>\n<p>" + start + " " + timeOfDay + "</p>\n<p>" + evTypeStr + " - " + monthStr + " " + evDay + "</p>"
		}

		thisResponse = thisResponse + "\n<div class=\"eventWrapper\"><div class=\"eventContainer\">" + record + "</div><form hx-post=\"/mod\" hx-target=\"#dialog\" hx-swap=\"outerHTML\" hx-indicator=\"#throbber\" style=\"display: grid;\" method=\"post\">\n<input type=\"text\" name=\"month\" value=\"" + monthNum + "\" style=\"display: none;\">\n<input type=\"text\" name=\"year\" value=\"" + year + "\" style=\"display: none;\">\n<button type=\"submit\" name=\"del\" class=\"delButton material-symbols-rounded\" value=\"" + fmt.Sprint(i) + "\">delete</button>\n</form>\n</div>"
	}

	for i := 0; len(nextCal) > i; i++ {
		title := strings.Replace(strings.Replace(nextCal[i][6], ">", "&gt;", -1), "<", "&lt;", -1)
		startInt, err := strconv.Atoi(nextCal[i][0])
		var start string
		timeOfDay := "AM"
		timeOfEnd := "PM"

		if err != nil {
			fmt.Println(err)
			return
		} else if startInt > 12 && startInt != 12 {
			start = fmt.Sprint(startInt-12) + ":" + nextCal[i][1]
			timeOfDay = "PM"
		} else if startInt == 12 {
			timeOfDay = "PM"
			start = nextCal[i][0] + ":" + nextCal[i][1]
		} else {
			start = nextCal[i][0] + ":" + nextCal[i][1]
		}

		end := "-1"
		if thisCal[i][2] != "-1" {
			endInt, err := strconv.Atoi(nextCal[i][2])
			if err != nil {
				fmt.Println(err)
				return
			} else if endInt > 12 && endInt != 12 {
				end = fmt.Sprint(endInt-12) + ":" + nextCal[i][3]
				timeOfEnd = "PM"
			} else if endInt == 12 {
				timeOfEnd = "PM"
				end = nextCal[i][2] + ":" + nextCal[i][3]
			} else {
				end = nextCal[i][2] + ":" + nextCal[i][3]
			}
		}
		evDay := nextCal[i][4]
		evTypeInt, err := strconv.Atoi(nextCal[i][5])
		if err != nil {
			fmt.Println(err)
			return
		}
		evTypeStr := eventTypes[evTypeInt]

		record := "No events for this month."
		if end != "-1" {
			record = "<p><span class=\"evTitle\">" + title + "</span><p>\n<p>" + start + " " + timeOfDay + " - " + end + " " + timeOfEnd + "</p>\n<p>" + evTypeStr + " - " + nextMonthStr + " " + evDay + "</p>"
		} else {
			record = "<p><span class=\"evTitle\">" + title + "</span><p>\n<p>" + start + " " + timeOfDay + "</p>\n<p>" + evTypeStr + " - " + nextMonthStr + " " + evDay + "</p>"
		}

		nextResponse = nextResponse + "\n<div class=\"eventWrapper\"><div class=\"eventContainer\">" + record + "</div><form hx-post=\"/mod\" hx-target=\"#dialog\" hx-swap=\"outerHTML\" hx-indicator=\"#throbber\" style=\"display: grid;\" method=\"post\">\n<input type=\"text\" name=\"month\" value=\"" + nextMonthNum + "\" style=\"display: none;\">\n<input type=\"text\" name=\"year\" value=\"" + nextYear + "\" style=\"display: none;\">\n<button type=\"submit\" name=\"del\" class=\"delButton material-symbols-rounded\" value=\"" + fmt.Sprint(i) + "\">delete</button>\n</form>\n</div>"
	}

	if len(thisCal) == 0 {
		thisResponse = "<p>No events for this month.</p>"
	}
	if len(nextCal) == 0 {
		nextResponse = "<p>No events for this month.</p>"
	}

	response = "<h1>" + monthStr + "</h1>\n<div class=\"monthContainer\">" + thisResponse + "</div>\n<h1>" + nextMonthStr + "</h1>\n<div class=\"monthContainer\">" + nextResponse + "</div>"

	fmt.Fprint(w, response)

	thisMonthCal.Close()
	nextMonthCal.Close()
}

func removalHandler(w http.ResponseWriter, r *http.Request) {
	// Get which event the user wants to delete
	var deathRow int
	var whereFrom string
	var yearFrom string
	var err error
	var err2 error
	var file *os.File
	var writable *os.File

	deathRow, err = strconv.Atoi(r.FormValue("del"))
	whereFrom = r.FormValue("month")
	yearFrom = r.FormValue("year")
	if err != nil || len(yearFrom) != 4 || (len(whereFrom) != 1 && len(whereFrom) != 2) {
		fmt.Println(err)
		fmt.Fprint(w, "<div id=\"dialog\" class=\"dgOpen\"><p>Invalid request.</p></div>")
		return
	}

	file, err = os.Open(eventsLibrary + whereFrom + "_" + yearFrom + ".csv")
	writable, err2 = os.OpenFile(eventsLibrary+whereFrom+"_"+yearFrom+".csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	if err2 != nil {
		fmt.Println(err2)
	}

	read := csv.NewReader(file)
	write := csv.NewWriter(writable)

	dav, _ := read.ReadAll()
	if len(r.FormValue("response")) == 0 {
		fmt.Fprint(w, "<div id=\"dialog\" class=\"dgOpen\"><p>You are about to remove event \"<strong>"+dav[deathRow][6]+"</strong>\" from the calendar.<br/><br/>Are you sure you want to continue?</p><br/><br/>\n<form hx-post=\"/mod\" hx-target=\"#dialog\" hx-swap=\"outerHTML\" hx-indicator=\"#throbber\" method=\"post\">\n<input type=\"text\" name=\"month\" value=\""+whereFrom+"\" style=\"display: none;\">\n<input type=\"text\" name=\"del\" value=\""+fmt.Sprint(deathRow)+"\" style=\"display: none;\">\n<input type=\"text\" name=\"year\" value=\""+yearFrom+"\" style=\"display: none;\">\n<input type=\"submit\" name=\"response\" value=\"Yes\" />\n<input type=\"submit\" name=\"response\" value=\"No\" />\n</form></div>")
	} else {
		if r.FormValue("response") == "Yes" {
			tmp := append(dav[:deathRow], dav[deathRow+1:]...)

			if err := os.Truncate(eventsLibrary+whereFrom+"_"+yearFrom+".csv", 0); err != nil {
				log.Printf("Failed to truncate: %v", err)
			}
			// for i := 0; i < len(tmp); i++ {
			// 	tmp[i][6] = `.\` + tmp[i][6]
			// 	// fmt.Println(tmp[i])
			// }
			write.WriteAll(tmp)
			fmt.Fprint(w, "<div id=\"dialog\" style=\"animation: popup reverse 0.5s ease;\"></div><script>location.reload()</script>")
		} else {
			fmt.Fprint(w, "<div id=\"dialog\" style=\"animation: popup reverse 0.5s ease;\"></div>")
		}
	}
	file.Close()
}
