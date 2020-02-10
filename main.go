package main

import ( 
	"fmt"
    "net/http"
)

func Handler(w http.ResponseWriter, req *http.Request){
	w.Write([]byte("hello"))
}

func main(){
	fmt.Println("start")
	http.HandleFunc("/hello", Handler)
	http.ListenAndServe(":12345", nil)
}