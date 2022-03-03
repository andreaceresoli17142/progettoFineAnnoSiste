package main

import ( // {{{
	// "crypto/sha256"
	"database/sql"
	"fmt"
	// "log"
	"net/http"
	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
) // }}}

type Conversation struct {
	Id int `db:"id"`
	Name string `db:"name"`
	Description string `db:"description"`
}

// get conversations {{{
// QUERY: SELECT c.id, cn.name, cn.description FROM Conversations c INNER JOIN ConversationName cn WHERE c.participantId = 0 GROUP BY c.id;
func getConversations(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: get conversations")
	// fmt.Println( r.Header.Authorization )
	 access_token := BearerAuthHeader(r.Header.Get("Authorization"))

	// err := r.ParseForm()

	// if err != nil {
	// 	fmt.Fprintf(w, "{ \"resp_code\":500, error: \"%v\" }", err)
	// 	return
	// }

	// access_token := validate(r.PostForm.Get("access_token"))
	
	fmt.Printf( "debug - act: %s\n", access_token )

	usr_id, err := accessToken_get_usrid( access_token )

	fmt.Printf( "debug - userid: %d\n", usr_id )

	if err != nil {
		fmt.Fprintf(w, "{ \"resp_code\":500, error: \"%v\" }", err)
		return
	}

	if usr_id == -1 {
		fmt.Fprint(w, "{ \"resp_code\":300, error: \"invalid access token\" }")
		return
	}

	db, err := sql.Open("mysql", "root:root@tcp("+sqlServerIp+")/"+dbname)

	defer db.Close()

	if err != nil {
		fmt.Fprintf(w, "{ \"resp_code\":300, error: \"%v\" }", err)
		return 
	}

	// q := fmt.Sprintf("SELECT salt, pHash FROM Users WHERE id = (?);", usr_id)
	res, err := db.Query("SELECT c.id as id, cn.name as name, cn.description as description FROM Conversations c INNER JOIN ConversationName cn WHERE c.participantId = (?) GROUP BY c.id;", usr_id)

	defer res.Close()

	if err == sql.ErrNoRows {
		fmt.Fprint(w, "{ \"resp_code\":200, data:[] }")
		return
	}	

	if err != nil {
		fmt.Fprintf(w, "{ \"resp_code\":500, error: \"%v\" }", err)
		return
	}

	var convs []Conversation

	for res.Next() {
		var conv Conversation

		err := res.Scan(&conv.Id, &conv.Name, &conv.Description)
		
		if err != nil {
			// fmt.Printf( "debug - wqe: %d\n", usr_id )
			fmt.Fprintf(w, "{ \"resp_code\":500, error: \"%v\" }", err)
			return
		}

		convs = append(convs, conv)
	}

	a, _ := json.Marshal(convs)

	fmt.Fprintf(w, "{ \"resp_code\":200, data:\"%s\" }", string(a) )

} // }}}
