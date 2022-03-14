package main

import ( // {{{
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
) // }}}

// get conversations {{{
func getConversations(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: get conversations")
	access_token := BearerAuthHeader(r.Header.Get("Authorization"))

	usr_id, err := getAccessToken_usrid(access_token)

	if err != nil {
		httpError(&w, 500, err)
		return
	}

	if usr_id == -1 {
		httpError(&w, 300, "invalid access token")
		return
	}

	db, err := sql.Open("mysql", databaseString)

	defer db.Close()

	if err != nil {
		httpError(&w, 500, err)
		return
	}

	res, err := db.Query("SELECT c.id as id, cn.name as name, cn.description as description FROM Conversations c INNER JOIN ConversationName cn WHERE c.participantId = (?) GROUP BY c.id;", usr_id)

	defer res.Close()

	if err == sql.ErrNoRows {
		var r []string
		httpSuccessf(&w, 200, "data", r)
		return
	}

	if err != nil {
		httpError(&w, 500, err)
		return
	}

	var convs []Conversation

	for res.Next() {
		var conv Conversation

		err := res.Scan(&conv.Id, &conv.Name, &conv.Description)

		if err != nil {
			httpError(&w, 500, err)
			return
		}

		convs = append(convs, conv)
	}

	a, _ := json.Marshal(convs)

	httpSuccessf(&w, 200, "data", string(a))
}

// takes userid, name and description, if description is empty we are creaing a conversation otherwise we are creating a group
// then create conversation, add user and data, return conversation id
func addConversations(userId int, name string, description string) (int, error) {

	// generating conversation id
	db, err := sql.Open("mysql", databaseString)

	if err != nil {
		return -1, AppendError("error: addConversation", err)
	}

	defer db.Close()

	var convId int
	err = db.QueryRow(`SELECT COUNT(*) FROM (SELECT * FROM Conversations GROUP BY id) as a;`).Scan(convId)

	if err == sql.ErrNoRows {
		return -1, AppendError("error: addConversation", err)
	}

	if err != nil {
		return -1, AppendError("error: addConversation", err)
	}

	convId = convId + 1
	// finished conversation id

	// adding conversation
	db, err = sql.Open("mysql", databaseString)

	if err != nil {
		return -1, AppendError("error: addConversation", err)
	}

	defer db.Close()

	err = db.QueryRow(`INSERT INTO Conversations (id, participantId)(?, ?);`, convId, userId).Scan() // .Scan(convId)

	if err == sql.ErrNoRows {
		return -1, AppendError("error: addConversation", err)
	}

	if err != nil {
		return -1, AppendError("error: addConversation", err)
	}
	// finised adding conversation

	// adding conversation name and description
	// db, err = sql.Open("mysql", databaseString)

	// if err != nil {
	// 	return -1, AppendError("error: addConversation", err)
	// }

	// defer db.Close()

	// q := `INSERT INTO GroupName (id, participantId)(?, ?);`

	// err = db.QueryRow(q, convId, userId).Scan() // .Scan(convId)

	// if err == sql.ErrNoRows {
	// 	return -1, AppendError("error: addConversation", err)
	// }

	// if err != nil {
	// 	return -1, AppendError("error: addConversation", err)
	// }

	return convId, nil
}

// }}}
