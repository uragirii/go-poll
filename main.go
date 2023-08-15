package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PollQuestion struct {
  Id string `json:"id"`
  Question string `json:"question"`
  Options [2]string `json:"options"`
}

func idFactory() func () string {
  i :=0

  return func () string {
    i++;
    return fmt.Sprintf("%v",i)
  }
}

var mockdb = make(map[string]PollQuestion)

func mockDataSetup() {
  idGenerator := idFactory()
  for idx:=0; idx<10;idx++ {
    options:= [2]string {"Option 1", "Option 2"}
    id:= idGenerator()
    mockdb[id] = PollQuestion{
      Id: id,
      Question: fmt.Sprintf("This is the quesiton number %v", idx),
      Options: options,
    }
  }

  fmt.Printf("Created %v mock data", len(mockdb))
}

func main() {
  mockDataSetup()
  r := gin.Default()
  
  

  // just here for testing purposes
  r.GET("/ping", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "message": "pong",
    })
  })

  r.GET("/poll/:id", func(ctx *gin.Context) {
    id:= ctx.Params.ByName("id")
    poll, valid := mockdb[id]


    if !valid {
      ctx.JSON(http.StatusNotFound, gin.H{"error" : "poll not found", "message" : fmt.Sprintf("Poll with ID %v not found", id)})
    }else {
      ctx.JSON(http.StatusOK, poll)
    }
  })

  r.GET("/polls", func(ctx *gin.Context) {
    var polls []PollQuestion

    for _, val :=range(mockdb) {
      polls = append(polls, val)
    }

    ctx.JSON(http.StatusOK, polls)

  })


  r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}