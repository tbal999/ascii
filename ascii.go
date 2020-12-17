package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strings"

	"github.com/go-vgo/robotgo"
)

var mapage = `<title>ASCII OUTPUT</title>
</head>
<body>
    <style>
        input {
            display: inline-block;
            float: left;
            margin-right: 20px;
            background-color: #00FFFF
        }
    </style>
    <b>¯\_(ツ)_/¯</b>
    <p>Now you can easily paste ASCII's on IRC!</p>
    <form action="/submit" method="POST">
		<textarea name="asciiQuery" cols="90" rows="20" placeholder="Paste the ASCII art here..."></textarea>
		<br>
		(There is a 2 second delay between submitting something here and the keyboard output - so you have a chance to click on IRC chat input)<p>
		<button type="submit" style="background-color: #00FFFF;" name="submit" value="IRC">IRC etc</button>
		<button type="submit" style="background-color: #00FFFF;" name="submit" value="Telegram">Telegram etc</button>
		</form>
		
	 <br/>
    <br>
</body>`

//PageVariables - GUI variables that change on webpages.
type PageVariables struct {
	Input string
}

//Main page handler - the front page.
func mainpage(w http.ResponseWriter, r *http.Request) {
	HomePageVars := PageVariables{}
	t, _ := template.New("mainpage").Parse(mapage)
	err2 := t.Execute(w, HomePageVars)
	if err2 != nil { // if there is an error
		log.Print("template executing error: ", err2) //log it
	}
}

func start(w http.ResponseWriter, r *http.Request) {
	HomePageVars := PageVariables{}
	t, _ := template.New("mainpage").Parse(mapage)
	r.ParseForm()
	submit := r.FormValue("submit")
	fmt.Println(submit)
	switch submit {
	case "IRC":
		if len(r.Form["asciiQuery"]) != 0 && r.Form["asciiQuery"][0] != "" { //tablequery
			output := strings.Split(r.Form["asciiQuery"][0], "\n")
			robotgo.Sleep(2)
			for index := range output {
				robotgo.TypeStr(output[index] + "\n")
				robotgo.KeyTap("enter")
			}
		} else {
			fmt.Println("No input.")
		}
	case "Telegram":
		if len(r.Form["asciiQuery"]) != 0 && r.Form["asciiQuery"][0] != "" { //tablequery
			output := strings.Split(r.Form["asciiQuery"][0], "\n")
			robotgo.Sleep(2)
			for index := range output {
				robotgo.TypeStr(output[index])
				robotgo.KeyTap("enter", "shift")
			}
		} else {
			fmt.Println("No input.")
		}
		robotgo.KeyTap("enter")
		robotgo.KeyTap("enter")
	default:
		fmt.Println(r.FormValue("submit"))
	}
	err2 := t.Execute(w, HomePageVars)
	if err2 != nil { // if there is an error
		log.Print("template executing error: ", err2) //log it
	}
}

//OpenBrowser - Opens your default browser, depending on the OS you are on.
func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println("Started...")
	http.HandleFunc("/", mainpage)
	http.HandleFunc("/submit", start)
	openBrowser("http://127.0.0.1:8080/")
	err := http.ListenAndServe(":8080", nil) // setting listening port
	if err != nil {
		fmt.Println(err)
	}

}
