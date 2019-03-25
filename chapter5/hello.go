package main

import (
    "io"
    "log"
    "net/http"
    "fmt"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "Hello world!")
}

type Integer int
    
func (a Integer) Less(b Integer) bool{
    return a < b
}
func (a *Integer) Add (b Integer) {
    *a += b
}
func main() {
    var a Integer = 1 
    if a.Less(2){
        fmt.Println(a, " less 2 ")
    }
    a.Add(5);
    fmt.Println(a);
    var b1 = [3] int{1,2,3}
    var b = &b1
    b[1]++
    fmt.Println(b1, *b)
    http.HandleFunc("/hello", helloHandler)
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServer:", err.Error())
    }
}
