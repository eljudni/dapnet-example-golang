package main

import "flag"
import "fmt"
import "os"
import "encoding/json"
import "net/http"
import "bytes"

var owncall = os.Getenv("DAP_OWNCALL")
var api_user = os.Getenv("DAP_APIUSER")
var api_pass = os.Getenv("DAP_APIPASS")

type outmessage struct {
	Text	string	`json:"text"`
	CallSignNames	[]string	`json:"callSignNames"`
	TransmitterGroupNames	[]string	`json:"transmitterGroupNames"`
}
func main() {
	if owncall == "" || api_user == "" || api_pass == "" {
		fmt.Println("You must set environment variables accordingly: DAP_OWNCALL, DAP_APIUSER, DAP_APIPASS")
		os.Exit(1)
	}
	callPtr := flag.String("call","","Call to send to")	
	msgPtr := flag.String("message","","Message to send")
	txPtr := flag.String("tx","dl-all","TX group to use")

	flag.Parse()

	if *callPtr == "" || *msgPtr == "" {
		fmt.Println("Call and message must not be empty.")
		os.Exit(1) 	
	}

	send(*callPtr,*msgPtr,*txPtr)
	os.Exit(0)

}

func send (c,m,tx string) {
	url := "http://www.hampager.de:8080/calls"
	callout := &outmessage{
		Text: owncall + ": " + m,
		CallSignNames: []string{c},
		TransmitterGroupNames: []string{tx}}
	fmt.Printf("Sending to %s\n",url)
	marshalled,_ := json.Marshal(callout)
	req,err := http.NewRequest("POST", url, bytes.NewBuffer(marshalled))
	req.SetBasicAuth(api_user,api_pass)
	req.Header.Set("Content-Type","application/json")
	client := &http.Client{}
	resp,err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 201 {
		fmt.Println("Message(s) sent.")
	} else {
		fmt.Println("Messages failed!")
		fmt.Println(resp.Status)
	}
	return
}