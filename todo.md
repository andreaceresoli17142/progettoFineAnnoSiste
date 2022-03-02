 _            _       
| |_ ___   __| | ___  
| __/ _ \ / _` |/ _ \ 
| || (_) | (_| | (_) |
 \__\___/ \__,_|\___/ 

---             
         
 - migrate all access token to headers and all other data to / separated url
 - clean up comments
 - comment code
 - change all q:= "" db.Exec() into dc.Exec("... (?) ... ", x)
 - frontend interface
 - oauth
 - change user: send confirmation email when changing email

 - implementing: get conversations
	 - returns only first conversation ( i think ) 

 - implement:
	 - send friend request
	 - accept firend request
	 - get friend requests
	 - send messages
	 - get messages ( with api pages )
