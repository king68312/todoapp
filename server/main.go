package main

import (
  "database/sql"
  "fmt"
  "net/http"
  "github.com/gin-gonic/gin"
  _ "github.com/mattn/go-sqlite3"
  "log"
)

var A = "I am a student"
var B = "I recommend this course"
var C = "I use Arch linux"
var taskItems = []string{A, B, C}

func main(){
  db,err := sql.Open("sqlite3","./todo.db")
  if err != nil{
    log.Fatal(err)
  } 
  defer db.Close()

  tableName := "todolist"
  createTable := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s ( id INTEGER PRIMARY KEY AUTOINCREMENT,
  rank INTEGER NOT NULL,                                              task TEXT NOT NULL,
  done BOOLEAN NOT NULL DEFAULT 0
  );`, tableName)
  _, err = db.Exec(createTable)
  if err != nil{
    log.Fatal(err)
  }
  _, err = db.Exec(`INSERT INTO todolist (task,rank) VALUES(?, ?)`,"算数の宿題", 2)
  if err != nil{
    log.Fatal(err)
  }

  rows, err := db.Query("SELECT id, rank,task,done FROM todolist")
  if err != nil{
    log.Fatal(err)
  }
  defer rows.Close()
 
for rows.Next() {
  var id int
  var rank int
  var task string
  var done bool
  err = rows.Scan(&id, &rank, &task,&done)
  if err != nil{
    log.Fatal(err)
  }
  fmt.Printf("%d %d %s %t\n",id, rank , task, done)
}
  router := gin.Default() 
  router.GET("/", IndexGET)
  router.Run(":8080")
}

func IndexGET(c *gin.Context) {
	c.String(http.StatusOK, "Hello world!")
}

func helloUser(writer http.ResponseWriter, request *http.Request){
  var greet = "Hello!"
  fmt.Fprintln(writer, greet)
}
