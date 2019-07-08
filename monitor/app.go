package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

func Home(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("inicio.html"))
	tmpl.Execute(w, nil)
}

func main() {
	//b := 200
	mux := http.NewServeMux()

	mux.HandleFunc("/", Home)
	mux.HandleFunc("/receive", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "POST" {
			ajax_post_data := r.FormValue("ajax_post_data")
			fmt.Println("Receive ajax post data string ", ajax_post_data)
			//output := fmt.Sprintf("%d", b)
			w.Write([]byte(getdata()))
			//b += 1
		}
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println("Servidor funcionando en el puerto 5000..")
	http.ListenAndServe(":5000", mux)
}

func getdata() string {
	read, err := ReadMemInfo("/proc/meminfo")
	dealwithErr(err)
	total := (strconv.FormatUint((read.MemTotal / 1000), 10))
	consumida := (strconv.FormatUint(((read.MemTotal - read.MemAvailable) / 1000), 10))
	percent := ((read.MemTotal - read.MemAvailable) * 100) / read.MemTotal
	porcentaje := strconv.FormatUint(percent, 10)

	idle0, total0 := getCPUSample()
	time.Sleep(3 * time.Second)
	idle1, total1 := getCPUSample()

	idleTicks := float64(idle1 - idle0)
	totalTicks := float64(total1 - total0)
	porcentajeCPU := 100 * (totalTicks - idleTicks) / totalTicks

	//porcentajeCPU := 75.36
	return fmt.Sprintf("%s|%s|%s|%.2f", total, consumida, porcentaje, porcentajeCPU)
}

func dealwithErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
