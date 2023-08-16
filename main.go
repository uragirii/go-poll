package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"simple-server/db"

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
        ctx.Set("user", user)
        ctx.Next();
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

func getUserFromContext(ctx  *gin.Context) (User, bool) {
  maybeUser, ok := ctx.Get("user")

    if !ok {
      // not possible
      ctx.JSON(http.StatusInternalServerError, gin.H{"error" : "something went wrong", "message" : "Unable to recognise the user"})
      return User{}, false;
    }

    user, ok := maybeUser.(User);

    if !ok {
       // not possible
       ctx.JSON(http.StatusInternalServerError, gin.H{"error" : "something went wrong", "message" : "Unable to recognise the user"})
       return User{}, false;
    }

    return user, true
}

func hasUserSubmitted(user User, pollId string) bool {
  findInPollId:= func (slice []string, id string) bool  {
    for idx := range(slice){
      if slice[idx] == id{ 
        return true
      }
    }
    return false
  }

  return findInPollId(user.submittedPolls, pollId)

  
}

func main() {
  // mockDataSetup()
  // r := gin.Default()

  db.NewPoll("data/go.db");

  return;

  // just here for testing purposes
  // r.GET("/ping", func(c *gin.Context) {
  //   c.JSON(http.StatusOK, gin.H{
  //     "message": "pong",
  //   })
  // })

  // r.GET("/poll/:id", VerifyCookie() ,func(ctx *gin.Context) {
  //   id:= ctx.Params.ByName("id")
  //   poll, valid := mockdb[id]


  //   if !valid {
  //     ctx.JSON(http.StatusNotFound, gin.H{"error" : "poll not found", "message" : fmt.Sprintf("Poll with ID %v not found", id)})
  //   }else {

  //     user, ok := getUserFromContext(ctx);

  //     if !ok {
  //       return;
  //     }

  //     submitted:= hasUserSubmitted(user, id)

  //     if !submitted {
  //       ctx.JSON(http.StatusOK, poll)
  //       return;
  //     }

  //     type PollWithSubmissions struct {
  //       PollQuestion
  //       Submissions [2] int `json:"submissions"` 
  //     }

  //     ctx.JSON(http.StatusOK, PollWithSubmissions{PollQuestion: poll, Submissions: poll.submissions})

  //   }
  // })

  // r.GET("/polls",VerifyCookie(), func(ctx *gin.Context) {
    
  //   // Whats the best way of doing this??
  //   type PollWithViewStatus struct {
  //     PollQuestion
  //     Viewed bool `json:"viewed"`
  //   }
  //   polls := make([]PollWithViewStatus, 0,len(mockdb))
    
  //   user, ok := getUserFromContext(ctx)

  //   if !ok {
  //     return;
  //   }

  //   for _, val :=range(mockdb) {
  //     polls = append(polls, PollWithViewStatus{PollQuestion: val, Viewed: hasUserSubmitted(user, val.Id)})
  //   }
  //   ctx.JSON(http.StatusOK, polls)

  // })

  // r.POST("/poll/:id", VerifyCookie(), func(ctx *gin.Context) {
  //   // verify poll
  //   pollId:= ctx.Params.ByName("id");

  //   poll, valid := mockdb[pollId]

  //   user,ok := getUserFromContext(ctx);

  //   if !ok {
  //     return;
  //   }

  //   if hasUserSubmitted(user, pollId) {
  //     ctx.JSON(http.StatusBadRequest, gin.H{"error" : "invalid data", "message": "user has already submitted poll"})
  //     return;
  //   }

  //   type Submission struct {
  //     SelectedOption int `json:"selectedOption" binding:"required" uri:"selectedOption"`
  //   }

  //   if !valid {
  //     ctx.JSON(http.StatusNotFound, gin.H{"error" : "poll not found", "message" : fmt.Sprintf("Poll with ID %v not found", pollId)})
  //   }

  //   var submission Submission

  //   err:=ctx.BindJSON(&submission)
  //   if err != nil{
  //     fmt.Println(err);
  //     ctx.JSON(http.StatusBadRequest, gin.H{"error":"invalid data", "message": "cannot parse the data"})
  //     return;
  //   }

  //   optionIdx := submission.SelectedOption

  //   if optionIdx > len(poll.submissions) {
  //     ctx.JSON(http.StatusBadRequest, gin.H{"error" :"invalid data", "message" : "selectedOption cannot be more than available options"})
  //   }

  //   poll.submissions[optionIdx] = poll.submissions[optionIdx]+1

  //   mockdb[pollId] = poll

  //   user.submittedPolls = append(user.submittedPolls, pollId)

  //   mockUserdb[user.Id] = user;

  //   fmt.Println(mockUserdb)

  //   ctx.JSON(http.StatusOK, gin.H{"data": poll.submissions})

  // })

  // r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}