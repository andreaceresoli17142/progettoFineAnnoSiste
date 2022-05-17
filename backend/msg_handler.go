package main

import ( // {{{
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
) // }}}

//TODO: modifiable profile pics
//TODO: modifiable statuses
//? add read message?

// conversations {{{
func getConversations(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: get conversations")

	type Convs struct {
		Id          string `json:"convId"db:"id"`
		Name        string `json:"name"db:"name"`
		Description string `json:"description"db:"description"`
		ProfilePic  bool   `json:"pfp"`
	}

	// b, err := ioutil.ReadAll(r.Body)

	v := r.URL.Query()

	userId, _ := strconv.Atoi(v.Get("S"))

	if userId == 0 {
		httpError(&w, 300, "missing parameters")
		return
	}

	// type ReqData struct {
	// 	UserId int `json:"userid"`
	// }

	// var resp ReqData

	// err := httpGetBody(r, &resp)
	// // err = json.Unmarshal(b, &resp)

	// if err != nil {
	// 	httpError(&w, 500, "error getting body: "+err.Error())
	// 	return
	// }

	db, err := sql.Open("mysql", databaseString)

	if err != nil {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	defer db.Close()

	rows, err := db.Query(`
		SELECT p.id AS id, u.username AS name, u.state AS description, u.profilePic
		FROM PrivateMessages p INNER JOIN Users u ON p.user = u.id
		WHERE u.id <> ? AND p.id IN (
			SELECT id
			FROM PrivateMessages
			WHERE user = ?
		)
	`, userId, userId)

	if err != nil {
		httpError(&w, 500, "error doing query: "+err.Error())
		return
	}

	var Conversations []Convs

	for rows.Next() {
		var conv Convs

		var (
			tmpId int
			pfp   int
		)

		if err := rows.Scan(&tmpId, &conv.Name, &conv.Description, &pfp); err != nil {
			httpError(&w, 500, "error getting query data: "+err.Error())
			return
		}

		conv.ProfilePic = intToBool(pfp)

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
	`, userId)

	if err != nil {
		httpError(&w, 500, "error doing query: "+err.Error())
		return
	}

	for rows.Next() {
		var conv Convs

		var (
			tmpId int
			pfp   int
		)

		if err := rows.Scan(&tmpId, &conv.Name, &conv.Description, &pfp); err != nil {
			httpError(&w, 500, "error getting query data: "+err.Error())
			return
		}

		conv.ProfilePic = intToBool(pfp)

		conv.Id = "G" + strconv.Itoa(tmpId)

		Conversations = append(Conversations, conv)
	}

	retStr, _ := json.Marshal(Conversations)

	httpSuccessf(&w, 200, `"data":%v`, string(retStr))
}

func getSingleConversation(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: get single conversations")

	type Convs struct {
		Id          string `json:"convId"db:"id"`
		Name        string `json:"name"db:"name"`
		Description string `json:"description"db:"description"`
		ProfilePic  bool   `json:"pfp"`
	}

	var convData Convs

	// b, err := ioutil.ReadAll(r.Body)

	v := r.URL.Query()

	userId, _ := strconv.Atoi(v.Get("S"))
	tmpId := v.Get("convid")
	convId, _ := strconv.Atoi(tmpId[1:])
	conType := string(tmpId[0])

	if userId == 0 {
		httpError(&w, 300, "missing parameters")
		return
	}

	db, err := sql.Open("mysql", databaseString)

	if err != nil {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	defer db.Close()

	var usrD *sql.Row

	// Debugln(tmpId)
	// Debugln(convId)
	// Debugln(userId)

	if conType == "P" {
		usrD = db.QueryRow(`
			SELECT id, username AS name, state, profilePic
			FROM Users
			WHERE id = (
				SELECT user
				FROM PrivateMessages
				WHERE id = ? AND user <> ?
			)
		`, convId, userId)
	} else if conType == "G" {
		usrD = db.QueryRow(`
			SELECT id, name AS name, description, profilePic
			FROM GroupNames
			WHERE id = (
				SELECT id
				FROM GroupMembers
				WHERE id = ? AND user = ?
			)
		`, convId, userId)
	} else {
		httpError(&w, 300, "coversation type does not exist")
		return
	}

	err = usrD.Scan(&convData.Id, &convData.Name, &convData.Description, &convData.ProfilePic)

	if err != nil {
		if err == sql.ErrNoRows {
			httpError(&w, 300, "conv does not exists")
			return
		}
		httpError(&w, 500, "error doing query: "+err.Error())
		return
	}

	retStr, _ := json.Marshal(convData)

	httpSuccessf(&w, 200, `"data":%v`, string(retStr))
}

func addToGroup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: add user to group")

	// b, err := ioutil.ReadAll(r.Body)

	type ReqData struct {
		Sender  int `json:"s"`
		GroupId int `json:"groupid"`
		UserId  int `json:"userid"`
	}

	var resp ReqData

	err := httpGetBody(r, &resp)
	// err = json.Unmarshal(b, &resp)

	if err != nil {
		httpError(&w, 500, "error getting body: "+err.Error())
		return
	}

	db, err := sql.Open("mysql", databaseString)

	if err != nil {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	defer db.Close()

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

	if err != nil {
		if err != sql.ErrNoRows {
			httpError(&w, 500, "error doing control query: "+err.Error())
			return
		}
	} else {
		httpError(&w, 300, "user is not in the selected group")
		return
	}

	err = db.QueryRow(`SELECT isAdmin FROM GroupMembers WHERE id = ? AND user = ?;`, resp.GroupId, resp.Sender).Scan(&i)

	if i == 0 {
		httpError(&w, 300, "you are not an administrator of this group")
		return
	}

	if err != nil {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	_, err = db.Exec(`INSERT INTO GroupMembers (id, user) VALUES (?, ?);`, resp.GroupId, resp.UserId)

	if err != nil && err != sql.ErrNoRows {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	httpSuccess(&w, 200, "success")
}

func quitGroup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: add user to group")

	type ReqData struct {
		Sender  int `json:"s"`
		GroupId int `json:"groupid"`
	}

	var resp ReqData

	err := httpGetBody(r, &resp)

	if err != nil {
		httpError(&w, 500, "error getting body: "+err.Error())
		return
	}

	db, err := sql.Open("mysql", databaseString)

	if err != nil {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	defer db.Close()

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

	err = db.QueryRow(`SELECT isAdmin FROM GroupMembers WHERE id = ? AND user = ?;`, resp.GroupId, resp.Sender).Scan(&i)

	if i == 0 {
		httpError(&w, 300, "you are not an administrator of this group")
		return
	}

	if err != nil {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	err = db.QueryRow(`SELECT COUNT(id) FROM GroupMembers WHERE id = ? AND isAdmin = 1 GROUP BY(id);`, resp.GroupId).Scan(&i)

	if i <= 1 {
		httpError(&w, 300, "there has to be at least one administator in the group")
		return
	}

	if err != nil {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	_, err = db.Exec(`DELETE FROM GroupMembers WHERE id = ? AND user = ?`, resp.GroupId, resp.Sender)

	if err != nil && err != sql.ErrNoRows {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	httpSuccess(&w, 200, "success")
}

func changeGroupData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: change group data")

	type Gd struct {
		S       int    `json:"S"`
		GroupId int    `json:"groupid"`
		Name    string `json:"name"`
		Desc    string `json:"desc"`
		Pfp     string `json:"pfp"`
	}

	var resp Gd

	err := httpGetBody(r, &resp)
	// err = json.Unmarshal(b, &resp)
	userId := resp.S

	var Pfp = false

	if resp.Pfp != "" {
		Pfp = true
	}

	if err != nil {
		httpError(&w, 500, "error getting body: "+err.Error())
		return
	}

	db, err := sql.Open("mysql", databaseString)

	if err != nil {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	defer db.Close()

	var i int

	err = db.QueryRow(`SELECT id FROM GroupMembers WHERE id = ? AND user = ?;`, resp.GroupId, userId).Scan(&i)

	if err != nil && err == sql.ErrNoRows {
		httpError(&w, 300, "you are not in the selected group")
		return
	}

	if err != nil && err != sql.ErrNoRows {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	err = db.QueryRow(`SELECT isAdmin FROM GroupMembers WHERE id = ? AND user = ?;`, resp.GroupId, userId).Scan(&i)

	if i == 0 {
		httpError(&w, 300, "you are not an administrator of this group")
		return
	}

	if err != nil {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	_, err = db.Exec(`UPDATE GroupNames SET name = ?, description = ?, profilePic = ? WHERE id = ?;`, resp.Name, resp.Desc, Pfp, resp.GroupId)

	if err != nil {
		httpError(&w, 500, "backend error")
		return
	}

	// Debugln(resp)

	if Pfp {
		img, err := base64.StdEncoding.DecodeString(strings.Split(resp.Pfp, ",")[1])
		if err != nil {
			httpError(&w, 500, "backend error:"+err.Error())
			return
		}
		permissions := 0644 // or whatever you need
		err = ioutil.WriteFile("../frontend/fe/assets/G"+fmt.Sprint(resp.GroupId), img, fs.FileMode(permissions))
		if err != nil {
			httpError(&w, 500, "backend error:"+err.Error())
			return
		}
	}

	httpSuccess(&w, 200, "success")

}

func createGroup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: add user to group")

	type ReqData struct {
		S         int    `json:"s"`
		GroupName string `json:"name"`
		GroupDesc string `json:"desc"`
		Pfp       string `json:"pfp"`
	}

	var resp ReqData

	err := httpGetBody(r, &resp)
	// err = json.Unmarshal(b, &resp)
	userId := resp.S

	var Pfp = false

	if resp.Pfp != "" {
		Pfp = true
	}

	if err != nil {
		httpError(&w, 500, "error getting body: "+err.Error())
		return
	}

	db, err := sql.Open("mysql", databaseString)

	if err != nil {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	defer db.Close()

	var groupId int64

	res, err := db.Exec(`INSERT INTO GroupNames (name, description, profilePic) VALUES (?, ?, ?);`, resp.GroupName, resp.GroupDesc, Pfp)

	if err != nil {
		httpError(&w, 500, "backend error")
		return
	}

	groupId, err = res.LastInsertId()

	if err != nil {
		httpError(&w, 500, "backend error")
		return
	}

	_, err = db.Exec(`INSERT INTO GroupMembers (id, user, isAdmin) VALUES (?, ?, true);`, groupId, userId)

	if err != nil {
		httpError(&w, 500, "backend error")
		return
	}

	// Debugln(resp)

	if Pfp {
		img, err := base64.StdEncoding.DecodeString(strings.Split(resp.Pfp, ",")[1])
		if err != nil {
			httpError(&w, 500, "backend error:"+err.Error())
			return
		}
		permissions := 0644 // or whatever you need
		err = ioutil.WriteFile("../frontend/fe/assets/G"+fmt.Sprint(groupId), img, fs.FileMode(permissions))
		if err != nil {
			httpError(&w, 500, "backend error:"+err.Error())
			return
		}
	}

	httpSuccess(&w, 200, "success")
}

func adminManage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: admin manage")

	// b, err := ioutil.ReadAll(r.Body)

	type ReqData struct {
		Sender  int  `json:"s"`
		GroupId int  `json:"groupid"`
		UserId  int  `json:"userid"`
		Fvalue  bool `json:"isadmin"`
	}

	var resp ReqData

	err := httpGetBody(r, &resp)
	// err = json.Unmarshal(b, &resp)

	db, err := sql.Open("mysql", databaseString)

	if err != nil {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	defer db.Close()

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

	if err != nil {
		if err != sql.ErrNoRows {
			httpError(&w, 500, "error doing control query: "+err.Error())
			return
		}
		httpError(&w, 300, "user is not in the selected group")
		return
	}

	err = db.QueryRow(`SELECT isAdmin FROM GroupMembers WHERE id = ? AND user = ?;`, resp.GroupId, resp.Sender).Scan(&i)

	if i == 0 {
		httpError(&w, 300, "you are not an administrator of this group")
		return
	}

	if err != nil {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	err = db.QueryRow(`SELECT COUNT(id) FROM GroupMembers WHERE id = ? AND isAdmin = 1 GROUP BY(id);`, resp.GroupId).Scan(&i)

	if i <= 1 && resp.Fvalue == false {
		httpError(&w, 300, "there has to be at least one administator in the group")
		return
	}

	if err != nil {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	_, err = db.Exec(`UPDATE GroupMembers set isAdmin = ? WHERE id = ? AND user = ?;`, boolToInt(resp.Fvalue), resp.GroupId, resp.UserId)

	if err != nil && err != sql.ErrNoRows {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	httpSuccess(&w, 200, "success")
}

func adminKickUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: admin remove user")

	// b, err := ioutil.ReadAll(r.Body)

	type ReqData struct {
		Sender  int `json:"s"`
		GroupId int `json:"groupid"`
		UserId  int `json:"userid"`
	}

	var resp ReqData

	err := httpGetBody(r, &resp)
	// err = json.Unmarshal(b, &resp)

	db, err := sql.Open("mysql", databaseString)

	if err != nil {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	defer db.Close()

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

	if err != nil {
		if err != sql.ErrNoRows {
			httpError(&w, 500, "error doing control query: "+err.Error())
			return
		}
		httpError(&w, 300, "user is not in the selected group")
		return
	}

	err = db.QueryRow(`SELECT isAdmin FROM GroupMembers WHERE id = ? AND user = ?;`, resp.GroupId, resp.Sender).Scan(&i)

	if i == 0 {
		httpError(&w, 300, "you are not an administrator of this group")
		return
	}

	if err != nil {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	err = db.QueryRow(`SELECT isAdmin FROM GroupMembers WHERE id = ? AND user = ?;`, resp.GroupId, resp.UserId).Scan(&i)

	if i == 1 {
		err = db.QueryRow(`SELECT COUNT(id) FROM GroupMembers WHERE id = ? AND isAdmin = 1 GROUP BY(id);`, resp.GroupId).Scan(&i)

		if i <= 1 {
			httpError(&w, 300, "there has to be at least one administator in the group")
			return
		}
		if err != nil {
			httpError(&w, 500, "backend error: "+err.Error())
			return
		}
	}

	if err != nil {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	_, err = db.Exec(`DELETE FROM GroupMembers WHERE id = ? AND user = ?;`, resp.GroupId, resp.UserId)

	if err != nil && err != sql.ErrNoRows {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	httpSuccess(&w, 200, "success")
}

func getGroupData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: admin remove user")
	v := r.URL.Query()
	// b, err := ioutil.ReadAll(r.Body)

	type ReqData struct {
		Sender  int `json:"s"`
		GroupId int `json:"groupid"`
	}

	type User struct {
		Id       int    `db:"id"`
		Username string `db:"username"`
		Pfp      bool   `db:"pfp"`
		IsAdmin  bool   `db:"isAdmin"`
	}

	type GroupData struct {
		Id    int    `db:"id"`
		Name  string `db:"name"`
		Desc  string `db:"description"`
		Pfp   bool   `db:"profilePic"`
		Users []User
	}

	var resp ReqData

	var ret GroupData

	resp.Sender, _ = strconv.Atoi(v.Get("S"))
	resp.GroupId, _ = strconv.Atoi(v.Get("groupid"))

	if resp.Sender == 0 || resp.GroupId == 0 {
		httpError(&w, 300, "missing parameters")
		return
	} // err = json.Unmarshal(b, &resp)

	db, err := sql.Open("mysql", databaseString)

	if err != nil {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	defer db.Close()

	var i int

	err = db.QueryRow(`SELECT id FROM GroupMembers WHERE id = ? AND user = ?;`, resp.GroupId, resp.Sender).Scan(&i)

	if err != nil {
		if err != nil {
			httpError(&w, 300, "you are not in the selected group")
			return
		}
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	err = db.QueryRow(`
		SELECT id, name, description, profilePic
		FROM GroupNames
		WHERE id = ?
	`, resp.GroupId).Scan(&ret.Id, &ret.Name, &ret.Desc, &ret.Pfp)

	if err != nil {
		if err == sql.ErrNoRows {
			httpError(&w, 300, "group does not exist")
			return
		}
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	rows, err := db.Query(`
		SELECT u.id AS id, u.username AS username, u.profilePic AS pfp, g.isAdmin AS isAdmin
		FROM Users u INNER JOIN GroupMembers g ON u.id = g.user
		WHERE g.id = ?
	`, resp.GroupId)

	if err != nil {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	var tmp []User

	for rows.Next() {
		var gu User
		err = rows.Scan(&gu.Id, &gu.Username, &gu.Pfp, &gu.IsAdmin)

		if err != nil {
			httpError(&w, 500, "backend error: "+err.Error())
			return
		}
		tmp = append(tmp, gu)
	}

	ret.Users = tmp

	a, _ := json.Marshal(ret)

	httpSuccessf(&w, 200, `"data":%v`, string(a))
}

// }}}

// firends requests {{{
func makeFriendRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: make friend request")

	// b, err := ioutil.ReadAll(r.Body)

	type ReqData struct {
		S      int `json:"s"`
		UserId int `json:"userid"`
	}

	var resp ReqData

	err := httpGetBody(r, &resp)
	// err = json.Unmarshal(b, &resp)

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

	err = db.QueryRow(`SELECT * FROM PrivateMessages p1, PrivateMessages p2 WHERE p1.id = p2.id AND p1.user <> p2.user AND p1.user = ? AND p2.user = ?;`, requesteeId, requesterId).Scan(&ret)

	if err != sql.ErrNoRows {
		httpError(&w, 500, "you are already friends with this user")
		return
	}

	err = db.QueryRow(`SELECT id FROM FriendRequests WHERE senderId = ? AND reciverId = ?;`, requesteeId, requesterId).Scan(&ret)

	if err != sql.ErrNoRows {
		httpError(&w, 500, "request already exists")
		return
	}

	err = db.QueryRow(`SELECT id FROM FriendRequests WHERE senderId = ? AND reciverId = ?;`, requesterId, requesteeId).Scan(&ret)

	if err != sql.ErrNoRows {
		httpError(&w, 500, "firend request already exists")
		return
	}

	if err != nil && err != sql.ErrNoRows {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	_, err = db.Exec(`INSERT INTO FriendRequests (senderId, reciverId) VALUES (?, ?);`, requesterId, requesteeId)

	if err != nil && err != sql.ErrNoRows {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	httpSuccess(&w, 200, "request successfully sent")
}

func getFriendRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: get friend request")

	// type Res struct {
	// 	S int `json:"s"`
	// }

	// // b, err := ioutil.ReadAll(r.Body)

	// var resp Res

	// err := httpGetBody(r, &resp)
	// // err = json.Unmarshal(b, &resp)

	// userId := resp.S

	v := r.URL.Query()

	userId, err := strconv.Atoi(v.Get("S"))

	if userId == 0 {
		httpError(&w, 300, "missing parameters")
		return
	}

	if err != nil {
		httpError(&w, 500, "error getting urldata: "+err.Error())
		return
	}

	type RespData struct {
		Id       int    `json:"id"`
		SenderId int    `json:"senderId"`
		Username string `json:"usr"`
	}

	db, err := sql.Open("mysql", databaseString)

	defer db.Close()

	if err != nil {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	res, err := db.Query(`
		SELECT f.id as id, f.senderId as sender, u.username as usr
		FROM FriendRequests f INNER JOIN Users u ON u.id = f.senderId
		WHERE f.reciverId = ? ;
	`, userId)

	//! is this even userful?
	defer res.Close()

	if err == sql.ErrNoRows {
		var r []string
		httpSuccessf(&w, 200, `"data":%v`, r)
		return
	}

	if err != nil {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	var requests []RespData

	for res.Next() {
		var req RespData

		err := res.Scan(&req.Id, &req.SenderId, &req.Username)

		if err != nil {
			httpError(&w, 500, "backend error: "+err.Error())
			return
		}

		requests = append(requests, req)
	}

	a, _ := json.Marshal(requests)

	httpSuccessf(&w, 200, `"data":%v`, string(a))
}

func acceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: accept friend request")

	type Res struct {
		S   int  `json:"s"`
		Id  int  `json:"id"`
		Acc bool `json:"accept"`
	}

	// b, err := ioutil.ReadAll(r.Body)

	var resp Res

	err := httpGetBody(r, &resp)
	// err = json.Unmarshal(b, &resp)

	userId := resp.S
	reqId := resp.Id

	Debugln(resp)

	if err != nil {
		httpError(&w, 500, "error getting body: "+err.Error())
		return
	}

	db, err := sql.Open("mysql", databaseString)

	defer db.Close()

	if err != nil {
		httpError(&w, 500, err)
		return
	}

	var sender int

	err = db.QueryRow(`
		SELECT senderId
		FROM FriendRequests
		WHERE id = ? AND reciverId = ? ;
	`, reqId, userId).Scan(&sender)

	if err == sql.ErrNoRows {
		httpError(&w, 300, "request does not exists")
		return
	}

	if err != nil {
		httpError(&w, 500, "error executing query: "+err.Error())
		return
	}

	_, err = db.Exec(`
		DELETE FROM FriendRequests
		WHERE id = ? AND reciverId = ?;
	`, reqId, userId)

	if err != nil {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	if resp.Acc {

		var maxPm int

		err = db.QueryRow(`SELECT MAX(id) FROM PrivateMessages;`).Scan(&maxPm)

		if err == sql.ErrNoRows {
			maxPm = 0
		}

		if err != nil {
			httpError(&w, 500, "error executing query: "+err.Error())
			return
		}

		maxPm++

		_, err = db.Exec(`
			INSERT INTO PrivateMessages (id, user)
			VALUES
			( ?, ? ),
			( ?, ? );
		`, maxPm, sender, maxPm, userId)

		if err != nil {
			httpError(&w, 500, "backend error: "+err.Error())
			return
		}
	}

	httpSuccess(&w, 200, "request accepted succesfully")
}

//}}}

// messages {{{
func sendMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: send messages")

	type Res struct {
		S    int    `json:"s"`
		Id   string `json:"convid"`
		Text string `json:"text"`
	}

	var resp Res

	err := httpGetBody(r, &resp)

	if err != nil {
		httpError(&w, 500, "error getting body: "+err.Error())
		return
	}

	user := resp.S

	conType := string(resp.Id[0])

	id, _ := strconv.Atoi(resp.Id[len(resp.Id)-1:])

	db, err := sql.Open("mysql", databaseString)

	defer db.Close()

	if err != nil {
		httpError(&w, 500, err)
		return
	}

	var sender int

	table := ""

	if conType == "P" {
		table = "PrivateMessages"
	} else if conType == "G" {
		table = "GroupMembers"
	} else {
		httpError(&w, 300, "coversation type does not exist")
		return
	}

	err = db.QueryRow(fmt.Sprintf(`
		SELECT user
		FROM %s
		WHERE id = ? AND user = ?
	`, table), id, user).Scan(&sender)

	if err == sql.ErrNoRows {
		httpError(&w, 300, "conversation does not exists")
		return
	}

	if err != nil {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	ret, err := db.Exec(`
		INSERT INTO Messages ( conv, senderId, content ) VALUES ( ?, ?, ? ) ;
	`, resp.Id, user, resp.Text)

	if err != nil {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	type Message struct {
		ConvId   string
		Id       int    `db:"id"`
		UserId   int    `db:"senderId"`
		Username string `db:"username"`
		Text     string `db:"content"`
		Time     string `db:"time"`
	}

	var mess Message

	mess.ConvId = resp.Id

	messId, _ := ret.LastInsertId()

	err = db.QueryRow(`
		SELECT m.id AS id, m.senderId AS senderId, u.username AS username, m.content AS content, m.time AS time
		FROM Messages m INNER JOIN Users u ON m.senderId = u.id
		WHERE m.id = ?
	`, messId).Scan(&mess.Id, &mess.UserId, &mess.Username, &mess.Text, &mess.Time)

	if err != nil {
		httpError(&w, 500, "error doing query: "+err.Error())
		return
	}

	retStr, _ := json.Marshal(mess)

	rows, _ := db.Query(fmt.Sprintf(`
		SELECT user
		FROM %s
		WHERE id = ?
	`, table), id)

	for rows.Next() {
		var Uid int
		rows.Scan(&Uid)
		socketSendNotification(Uid, string(retStr))
	}

	httpSuccess(&w, 200, "message sent succesfully")
}

//? implementare paginazione?
func getMessages(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: get messages")
	v := r.URL.Query()

	user, _ := strconv.Atoi(v.Get("S"))
	convid := v.Get("convid")

	if user == 0 || convid == "" {
		httpError(&w, 300, "missing parameters")
		return
	}

	// err := httpGetBody(r, &resp)

	// if err != nil {
	// 	httpError(&w, 500, "error getting body: "+err.Error())
	// 	return
	// }

	// user := resp.S

	// convid := resp.Id

	conType := string(convid[0])

	id, _ := strconv.Atoi(convid[len(convid)-1:])

	db, err := sql.Open("mysql", databaseString)

	defer db.Close()

	if err != nil {
		httpError(&w, 500, err)
		return
	}

	var sender int

	table := ""

	if conType == "P" {
		table = "PrivateMessages"
	} else if conType == "G" {
		table = "GroupMembers"
	} else {
		httpError(&w, 300, "coversation type does not exist")
		return
	}

	err = db.QueryRow(fmt.Sprintf(`
		SELECT user
		FROM %s
		WHERE id = ? AND user = ?
	`, table), id, user).Scan(&sender)

	if err == sql.ErrNoRows {
		httpError(&w, 300, "conversation does not exists")
		return
	}

	if err != nil {
		httpError(&w, 500, "backend error: "+err.Error())
		return
	}

	rows, err := db.Query(`
		SELECT m.id AS id, m.senderId AS senderId, u.username AS username, m.content AS content, m.time AS time
		FROM Messages m INNER JOIN Users u ON m.senderId = u.id
		WHERE m.conv = ?
	`, convid)

	if err != nil {
		httpError(&w, 500, "error doing query: "+err.Error())
		return
	}

	type Message struct {
		Id       int    `db:"id"`
		UserId   int    `db:"senderId"`
		Username string `db:"username"`
		Text     string `db:"content"`
		Time     string `db:"time"`
	}

	var messages []Message

	for rows.Next() {
		var mess Message

		if err := rows.Scan(&mess.Id, &mess.UserId, &mess.Username, &mess.Text, &mess.Time); err != nil {
			httpError(&w, 500, "error getting query data: "+err.Error())
			return
		}

		messages = append(messages, mess)
	}

	retStr, _ := json.Marshal(messages)

	httpSuccessf(&w, 200, `"data":%v`, string(retStr))
}

// }}}
