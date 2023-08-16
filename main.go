package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"simple-server/db"
	"simple-server/poll"

	"github.com/gin-gonic/gin"
)

func VerifyCookie (pollDb *db.Poll) gin.HandlerFunc {

  return func(ctx *gin.Context) {

    if cookie, err := ctx.Cookie("x-user-id"); err == nil {

      if user, err := pollDb.GetUser(cookie); err == nil {
        ctx.Set("user", user);
        ctx.Next();
        return;
      }

      fmt.Println("Cannot find user", err);
      ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error" :"internal server error", "message" : "couldn't find the user"});
      return;
    }
    // if no cookie is set we set the cookie
    if user, err:= pollDb.CreateUser(); err == nil {
      ctx.SetCookie("x-user-id", user.Id, 86400, "/","", false, false)
      ctx.Set("user", user);
      ctx.Next();
      return;
    }
    fmt.Println("Cannot create user");
    ctx.JSON(http.StatusInternalServerError, gin.H{"error" :"internal server error", "message" : "cannot create new user"});

  }
}

func getUserFromContext(ctx  *gin.Context) (poll.User, bool) {
  maybeUser, ok := ctx.Get("user")

    if !ok {
      // not possible
      ctx.JSON(http.StatusInternalServerError, gin.H{"error" : "something went wrong", "message" : "Unable to recognise the user"})
      return poll.User{}, false;
    }

    user, ok := maybeUser.(poll.User);

    if !ok {
       // not possible
       ctx.JSON(http.StatusInternalServerError, gin.H{"error" : "something went wrong", "message" : "Unable to recognise the user"})
       return poll.User{}, false;
    }

    return user, true
}

func main() {
  isFlyEnv := os.Getenv("IS_FLY_ENV")

  dbString := "sqldb/go.db"

  if isFlyEnv == "true" {
    dbString = "/data/go.db"
  }

  pollDb, err := db.NewPoll(dbString);

  if err != nil {
    fmt.Println("Error creating connection to db");
    fmt.Println(err);
    os.Exit(1);
    return;
  }

  var polls []poll.PollQuestion

  if polls,err= pollDb.GetAll(); err!=nil {
    fmt.Println("Error reading polls from db");
    fmt.Println(err);
    os.Exit(1);
    return;
  }

  if len(polls) == 0 {
    fmt.Println("Creating new dummy polls")
    // create dummy polls
    for i := 0; i < 10; i++ {
      _, err := pollDb.Create(fmt.Sprintf("This is the question number %v", i+1), [2]string{"Option1", "Option2"})
      if err != nil {
        fmt.Println("Error creating poll ");
        fmt.Println(err);
        os.Exit(1);
        return;
      }
    }
  }

  if polls,err= pollDb.GetAll(); err!=nil {
    fmt.Println("Error reading polls from db");
    fmt.Println(err);
    os.Exit(1);
    return;
  }

  fmt.Printf("Created %v polls\n", len(polls));

  r := gin.Default()

  // just here for testing purposes
  r.GET("/ping", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "message": "pong",
    })
  })

  r.GET("/poll/:id", VerifyCookie(pollDb) ,func(ctx *gin.Context) {
    id:= ctx.Params.ByName("id")

    pollData, err := pollDb.Get(id);

    if err == sql.ErrNoRows {
      ctx.JSON(http.StatusNotFound, gin.H{"error" : "poll not found", "message" : fmt.Sprintf("Poll with ID %v not found", id)})
      return;
    }

    if err !=nil {
      fmt.Println("Error while getting poll.id", id, err)
      ctx.JSON(http.StatusInternalServerError, gin.H{"error" :"internal server error", "message" : "cannot get poll"});
      return;
    }

    user, ok := getUserFromContext(ctx);

    if !ok {
      return;
    }

    submitted:= user.HasSubmitted(id)

    if !submitted {
      ctx.JSON(http.StatusOK, pollData)
      return;
    }

    type PollWithSubmissions struct {
      poll.PollQuestion
      Submissions [2] int `json:"submissions"` 
    }

    ctx.JSON(http.StatusOK, PollWithSubmissions{PollQuestion: pollData, Submissions: pollData.GetSubmissions()})

  })

  r.GET("/polls",VerifyCookie(pollDb), func(ctx *gin.Context) {
    
    // Whats the best way of doing this??
    type PollWithViewStatus struct {
      poll.PollQuestion
      Submitted bool `json:"submitted"`
    }
    pollsWithViewStatus := make([]PollWithViewStatus, 0,len(polls))
    
    user, ok := getUserFromContext(ctx)

    if !ok {
      return;
    }

    for _, val :=range(polls) {
      pollsWithViewStatus = append(pollsWithViewStatus, PollWithViewStatus{ Submitted: user.HasSubmitted(val.Id), PollQuestion: val})
    }
    ctx.JSON(http.StatusOK, pollsWithViewStatus)

  })

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

  r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}