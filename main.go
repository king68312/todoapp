package main

import (
  "fmt"
  "net/http"
)
var A = "I am a student"
var B = "I recommend this course"
var C = "I use Arch linux"
var taskItems = []string{A, B, C}

func main(){
  fmt.Print("Hello world")

  http.HandleFunc("/hello-go", helloUser)
  http.HandleFunc("/show-tasks", showTasks)
  http.ListenAndServe(":8080", nil)
}

func showTasks(writer http.ResponseWriter, request *http.Request){
  for _, task := range taskItems {
    fmt.Fprintln(writer, task)
  }
}

func helloUser(writer http.ResponseWriter, request *http.Request){
  var greet = "Hello!"
  fmt.Fprintln(writer, greet)
}
