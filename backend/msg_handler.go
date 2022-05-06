package main

import ( // {{{
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
) // }}}

// conversations {{{
//! MOSTLY OK, httpGetBody IS NOT WORKING
func getConversations(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: get conversations")

	type Convs struct {
		Id          string `json:"convId"db:"id"`
		Name        string `json:"name"db:"name"`
		Description string `json:"description"db:description`
	}

	b, err := ioutil.ReadAll(r.Body)

	type ReqData struct {
		UserId int `json:"userid"`
	}

	var resp ReqData

	// err = httpGetBody(r, &resp)
	err = json.Unmarshal(b, &resp)

	if err != nil {
		httpError(&w, 500, "error getting body: "+err.Error())
	}

	db, err := sql.Open("mysql", databaseString)

	if err != nil {
		httpError(&w, 500, err)
		return
	}

	defer db.Close()

	rows, err := db.Query(`
		SELECT p.id AS id, u.username AS name, u.state AS description
		FROM PrivateMessages p INNER JOIN Users u ON p.user = u.id
		WHERE u.id <> ? AND p.id IN (
			SELECT id
			FROM PrivateMessages
			WHERE user = ?
		)
	`, resp.UserId, resp.UserId)

	if err != nil {
		httpError(&w, 500, "error doing query: "+err.Error())
		return
	}

	var Conversations []Convs

	for rows.Next() {
		var conv Convs

		var tmpId int

		if err := rows.Scan(&tmpId, &conv.Name, &conv.Description); err != nil {
			httpError(&w, 500, "error getting query data: "+err.Error())
			return
		}

		conv.Id = "P" + strconv.Itoa(tmpId)

		Conversations = append(Conversations, conv)
	}

	rows, err = db.Query(`
		SELECT *
		FROM GroupNames
		WHERE id IN (
			SELECT id
			FROM GroupMembers
			WHERE user = ?
		)
	`, resp.UserId)

	if err != nil {
		httpError(&w, 500, "error doing query: "+err.Error())
		return
	}

	for rows.Next() {
		var conv Convs

		var tmpId int

		if err := rows.Scan(&tmpId, &conv.Name, &conv.Description); err != nil {
			httpError(&w, 500, "error getting query data: "+err.Error())
			return
		}

		conv.Id = "G" + strconv.Itoa(tmpId)

		Conversations = append(Conversations, conv)
	}

	retStr, _ := json.Marshal(Conversations)

	httpSuccessf(&w, 200, "%v", string(retStr))
}

func addToGroup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: get conversations")

	b, err := ioutil.ReadAll(r.Body)

	type ReqData struct {
		Sender  int `json:"s"`
		GroupId int `json:"groupid"`
		UserId  int `json:"userid"`
	}

	var resp ReqData

	// err = httpGetBody(r, &resp)
	err = json.Unmarshal(b, &resp)

	if err != nil {
		httpError(&w, 500, "error getting body: "+err.Error())
	}

	db, err := sql.Open("mysql", databaseString)

	if err != nil {
		httpError(&w, 500, err)
		return
	}

	defer db.Close()

	Debugln(resp)

	var i int

	err = db.QueryRow(`SELECT id FROM GroupMembers WHERE id = ? AND user = ?;`, resp.GroupId, resp.Sender).Scan(&i)

	if err != nil && err == sql.ErrNoRows {
		httpError(&w, 300, "you are not in the selected group")
		return
	}

	if err != nil && err != sql.ErrNoRows {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	err = db.QueryRow(`SELECT id FROM GroupMembers WHERE id = (?) AND user = (?);`, resp.GroupId, resp.UserId).Scan(&i)

	Debugln(i)

	if err != nil {
		if err != sql.ErrNoRows {
			httpError(&w, 500, "error doing control query: "+err.Error())
			return
		}
	} else {
		httpError(&w, 300, "user is in the selected group")
		return
	}

	_, err = db.Exec(`INSERT INTO GroupMembers (id, user) VALUES (?, ?);`, resp.GroupId, resp.UserId) // .Scan(convId)

	if err != nil && err != sql.ErrNoRows {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	httpSuccess(&w, 200, "success")
}

// }}}

// firends requests {{{

func makeFriendRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: get friend request")

	b, err := ioutil.ReadAll(r.Body)

	type ReqData struct {
		S      int `json:"s"`
		UserId int `json:"userid"`
	}

	var resp ReqData

	// err = httpGetBody(r, &resp)
	err = json.Unmarshal(b, &resp)

	if err != nil {
		httpError(&w, 500, "backend error - "+err.Error())
		return
	}

	requesterId := resp.S

	requesteeId := resp.UserId

	// requesteeId, _ := strconv.Atoi(resp.UserId)

	if requesterId == requesteeId {
		httpError(&w, 300, "you can't ask a friend request to yourself")
		return
	}

	db, err := sql.Open("mysql", databaseString)

	if err != nil {
		httpError(&w, 500, "backend error")
		return
	}

	defer db.Close()

	var ret int

	err = db.QueryRow(`SELECT id FROM FriendRequests WHERE senderId = ? AND reciverId = ?;`, requesteeId, requesterId).Scan(&ret) // .Scan(convId)

	if err != sql.ErrNoRows {
		httpError(&w, 500, "request already exists")
		return
	}

	err = db.QueryRow(`SELECT id FROM FriendRequests WHERE senderId = ? AND reciverId = ?;`, requesterId, requesteeId).Scan(&ret) // .Scan(convId)

	if err != sql.ErrNoRows {
		httpError(&w, 500, "firend request already exists")
		return
	}

	if err != nil && err != sql.ErrNoRows {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	_, err = db.Exec(`INSERT INTO FriendRequests (senderId, reciverId) VALUES (?, ?);`, requesterId, requesteeId) // .Scan(convId)

	if err != nil && err != sql.ErrNoRows {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	httpSuccess(&w, 200, "request successfully sent")
}

func getFriendRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: get friend request")

	type Res struct {
		S int `json:"s"`
	}

	b, err := ioutil.ReadAll(r.Body)

	var resp Res

	// err = httpGetBody(r, &resp)
	err = json.Unmarshal(b, &resp)

	userId := resp.S

	if err != nil {
		httpError(&w, 500, "error getting body: "+err.Error())
	}

	type RespData struct {
		Id       int `json:"id"`
		SenderId int `json:"senderId"`
		Username int `json:"usr"`
	}

	db, err := sql.Open("mysql", databaseString)

	defer db.Close()

	if err != nil {
		httpError(&w, 500, err)
		return
	}

	res, err := db.Query(`
	SELECT f.id as id, f.senderId as sender, u.username as usr
	FROM FriendRequests f INNER JOIN Users u ON u.id = f.senderId
	WHERE f.reciverId = ? ;
	`, userId)

	defer res.Close()

	if err == sql.ErrNoRows {
		var r []string
		httpSuccessf(&w, 200, "data:%v", r)
		return
	}

	if err != nil {
		httpError(&w, 500, err)
		return
	}

	Debugln(res)

	var requests []RespData

	for res.Next() {
		var req RespData

		err := res.Scan(&req.Id, &req.SenderId, &req.Username)

		if err != nil {
			httpError(&w, 500, err)
			return
		}

		requests = append(requests, req)
	}

	a, _ := json.Marshal(requests)

	Debugln(a)

	httpSuccessf(&w, 200, "data: %v", string(a))
}

//}}}
