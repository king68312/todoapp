package main

import (
  "database/sql"
  "fmt"
  "net/http"
  "strconv"
  "github.com/gin-gonic/gin"
  _ "github.com/mattn/go-sqlite3"
  "log"
)

var db *sql.DB

func createTask(task string, rank int) error {//C
  _, err := db.Exec("INSERT INTO todolist (rank, task) VALUES (?, ?)", rank, task)//データベースにデータを追加  
  if err != nil{
    log.Fatal(err)
  }
  return nil
}

func getTasks() ([]map[string]interface{}, error) {//R
	rows, err := db.Query("SELECT id, rank,task,done FROM todolist")//データベースからデータの取得
  if err != nil{
    log.Fatal(err)
  }
  defer rows.Close()//処理を遅延させ、Closeを確実に呼び出される

  var tasks []map[string]interface{}//map型のtasksを作成
  for rows.Next() {//rows.Next()でデータがある限り繰り返す
    var id int
    var rank int
    var task string
    var done bool
    err = rows.Scan(&id, &rank, &task, &done)
    if err != nil{ 
      log.Fatal(err)
    }
    tasks = append(tasks, map[string]interface{}{//appendの第一引数はスライス、第二引数は追加する要素
      "id":    id,
      "task":  task,
      "rank":  rank,
      "done":  done,
    })//json形式でデータを追加
  }
  return tasks, nil
}

func updateTask(id int, done bool) error {//U
  _, err := db.Exec("UPDATE todolist SET done = ? WHERE id = ?", done, id)//データベースのデータを更新
  if err != nil{
    log.Fatal(err)
  }
  return nil
}

func deleteTask(id int) error {//D
  _, err := db.Exec("DELETE FROM todolist WHERE id = ?", id)//データベースからデータを削除
  if err != nil{
    log.Fatal(err)
  }
  return nil
}

func main(){
  //データベースの作成
  var err error 
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

  router := gin.Default() 
  
  router.GET("/getTask",func(c *gin.Context) {  //タスクの一覧取得
    c.String(http.StatusOK, "Hello World") 
    fmt.Println("getTask") 
    tasks, err := getTasks()
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})//エラーが発生した場合に500エラーを返し、gin.HでJSONを返す
      return
    }
    c.JSON(200, tasks)//正常に処理が終了した場合に200を返し、tasksを返す
  })

  router.POST("/addTask", func(c *gin.Context) {//タスクの追加
    var json struct {
      Task string `json:"task"`
      Rank int    `json:"rank"`
    }
    if err := c.ShouldBindJSON(&json); err != nil {
      fmt.Println("health") 
      c.JSON(400, gin.H{"aaaa": err.Error()})//http.StatusInternalServerErrorは500エラー
    }
    err := createTask(json.Task, json.Rank)
    if err != nil {
    c.JSON(500, gin.H{"error": "Invalid input"})
    return
    }
    c.JSON(201, gin.H{"message": "Task created"})//http.StatusCreatedは201
  })

  router.PUT("/updateTask", func(c *gin.Context) {//タスクの更新
    var json struct {
      Task string `json:"task"`
      Rank int    `json:"rank"`
      Done bool   `json:"done"`
    }
    if err := c.ShouldBindJSON(&json); err != nil {
      c.JSON(400, gin.H{"error": err.Error()})
      return
    }
    err := updateTask(json.Rank, json.Done)
    if err != nil {
      c.JSON(500, gin.H{"error": "Invalid input"})
      return
    }
    c.JSON(200, gin.H{"message": "Task updated"})
  })

  router.DELETE("/deleteTask/:id", func(c *gin.Context) {//タスクの削除
    id := c.Param("id")
    idInt, err := strconv.Atoi(id)
    err = deleteTask(idInt)
    if err != nil { 
      c.JSON(500, gin.H{"error": "Invalid input"})
      return
    }
    c.JSON(200, gin.H{"message": "Task deleted"})
  })
  router.Run(":8080")
}
