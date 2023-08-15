package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)


type PollQuestion struct {
  Id string `json:"id"`
  Question string `json:"question"`
  Options [2]string `json:"options"`
  submissions [2]int
}

type User struct {
  Id string `json:"id"`
  submittedPolls []string 
}

func init() {
  rand.Seed(time.Now().UnixNano())
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// taken from Stack overflow answer https://stackoverflow.com/a/22892986/8077711
func randSeq(n int) string {
  b := make([]rune, n)
  for i := range b {
      b[i] = letters[rand.Intn(len(letters))]
  }
  return string(b)
}


func idFactory() func () string {
  i :=0

  return func () string {
    i++;
    return fmt.Sprintf("%v",i)
  }
}

var mockdb = make(map[string]PollQuestion)
var mockUserdb = make(map[string]User)

func mockDataSetup() {
  idGenerator := idFactory()
  for idx:=0; idx<10;idx++ {
    options:= [2]string {"Option 1", "Option 2"}
    id:= idGenerator()
    mockdb[id] = PollQuestion{
      Id: id,
      Question: fmt.Sprintf("This is the quesiton number %v", idx),
      Options: options,
      submissions: [2]int {0,0},
    }
  }

  fmt.Printf("Created %v mock data", len(mockdb))
}

func VerifyCookie () gin.HandlerFunc {
  

  return func(ctx *gin.Context) {
    assingCookie:= func ()  {
      id := randSeq(10)
      ctx.SetCookie("x-user-id", id, 86400, "/","", false, false)
      user:=User{Id: id}
      mockUserdb[id] = user
      ctx.Set("user", user)
    }

    if cookie, err := ctx.Cookie("x-user-id"); err == nil {
      if user, ok:= mockUserdb[cookie]; ok {
        ctx.Next();
        ctx.Set("user", user)
        return
      }
      assingCookie()
      ctx.Next()
      return;
    }
    assingCookie()
    ctx.Next()
    
  }
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

  r.GET("/poll/:id", VerifyCookie() ,func(ctx *gin.Context) {
    id:= ctx.Params.ByName("id")
    poll, valid := mockdb[id]


    if !valid {
      ctx.JSON(http.StatusNotFound, gin.H{"error" : "poll not found", "message" : fmt.Sprintf("Poll with ID %v not found", id)})
    }else {
      ctx.JSON(http.StatusOK, poll)
    }
  })

  r.GET("/polls",VerifyCookie(), func(ctx *gin.Context) {
    
    // Whats the best way of doing this??
    type PollWithViewStatus struct {
      PollQuestion
      Viewed bool `json:"viewed"`
    }
    polls := make([]PollWithViewStatus, 0,len(mockdb))
    
    maybeUser, ok := ctx.Get("user")

    if !ok {
      // not possible
      ctx.JSON(http.StatusInternalServerError, gin.H{"error" : "something went wrong", "message" : "Unable to recognise the user"})
      return;
    }

    user, ok := maybeUser.(User);

    fmt.Println(maybeUser)

    if !ok {
       // not possible
       ctx.JSON(http.StatusInternalServerError, gin.H{"error" : "something went wrong", "message" : "Unable to recognise the user"})
       return
    }

    findInPollId:= func (slice []string, id string) bool  {
      for idx := range(slice){
        if slice[idx] == id{ 
          return true
        }
      }
      return false
    }

    for _, val :=range(mockdb) {
      polls = append(polls, PollWithViewStatus{PollQuestion: val, Viewed: findInPollId(user.submittedPolls, val.Id)})
    }
    ctx.JSON(http.StatusOK, polls)

  })

  r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}