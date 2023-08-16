package db

import (
	"database/sql"
	"fmt"
	"simple-server/poll"

	_ "github.com/mattn/go-sqlite3"
)

const create string = `
CREATE TABLE IF NOT EXISTS PollQuestion (
  id INTEGER PRIMARY KEY NOT NULL,
  question TEXT NOT NULL,
  option1 TEXT NOT NULL,
  option2 TEXT NOT NULL,
  option1Count INTEGER DEFAULT 0,
  option2Count INTEGER DEFAULT 0
);
CREATE TABLE IF NOT EXISTS User (
  id TEXT PRIMARY KEY NOT NULL,
  selectedOptions TEXT DEFAULT ''
);
`

type Poll struct {
	db *sql.DB
}

func NewPoll(file string) (*Poll, error) {
	db, err := sql.Open("sqlite3", file)

	if err != nil {
		fmt.Println("Error opening connection to db ", file)
		fmt.Println(err)
		return nil, err
	}

	fmt.Printf("Database connection to '%v' completed\n", file)

	if _, err := db.Exec(create); err != nil {
		fmt.Println("Error creating PollQuestion/User table")
		return nil, err
	}

	fmt.Println("Created PollQuestion & User table")

	return &Poll{
		db: db,
	}, nil
}

func (p *Poll) Create(question string, options [2]string) (int, error) {

	res, err := p.db.Exec("INSERT INTO PollQuestion (question, option1, option2) VALUES (?, ?, ?);", question, options[0], options[1])

	if err != nil {
		fmt.Println("Error inserting into PollQuestion")
		return 0, nil
	}

	var id int64

	if id, err = res.LastInsertId(); err != nil {
		return 0, nil
	}

	return int(id), nil
}

func (p *Poll) GetAll() ([]poll.PollQuestion, error) {
	res, err := p.db.Query("SELECT * FROM PollQuestion;")

	if err != nil {
		fmt.Println("Error getting all the PollQuestions")
		return nil, err
	}

	defer res.Close()

	var rows []poll.PollQuestion

	for res.Next() {
		data := poll.PollQuestion{}

		var option1 string
		var option2 string
		var option1Count int
		var option2Count int

		err = res.Scan(&data.Id, &data.Question, &option1, &option2, &option1Count, &option2Count)

		if err != nil {
			fmt.Println("Error while scanning row of PollQuestion")
			return nil, err
		}
		data.Options = [2]string{option1, option2}
		data.AddSubmissions(option1Count, option2Count)
		rows = append(rows, data)
	}
	return rows, nil
}

func (p *Poll) CreateUser() (poll.User, error) {
	id := poll.RandUserId()
	_, err := p.db.Exec("INSERT INTO User (id) VALUES (?)", id)

	if err != nil {
		fmt.Println("Error inserting a new user")
		return poll.User{}, err
	}

	user, err := p.GetUser(id)

	if err != nil {
		fmt.Println("Error getting the created user")
		return poll.User{}, err
	}

	return user, nil
}

func (p *Poll) GetUser(id string) (poll.User, error) {
	row := p.db.QueryRow("SELECT * FROM User WHERE id=?", id)

	var userId string
	var submittedPolls string

	err := row.Scan(&userId, &submittedPolls)

	if err == sql.ErrNoRows {
		fmt.Printf("User with Id %v not found", id)
		return poll.User{}, err
	}

	return poll.User{Id: userId, SubmittedPolls: poll.ParseIds(submittedPolls)}, err
}

func (p *Poll) Get(id string) (poll.PollQuestion, error) {
	row := p.db.QueryRow("SELECT * FROM PollQuestion WHERE id=?", id)

	data := poll.PollQuestion{}

	var option1 string
	var option2 string
	var option1Count int
	var option2Count int

	err := row.Scan(&data.Id, &data.Question, &option1, &option2, &option1Count, &option2Count)

	if err == sql.ErrNoRows {
		fmt.Printf("Poll with Id %v not found", id)
		return data, err
	}
	data.Options = [2]string{option1, option2}
	data.AddSubmissions(option1Count, option2Count)

	return data, err

}
