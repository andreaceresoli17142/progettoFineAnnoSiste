package main


var (
    Creset 	string = "\033[0m"
    Cred 		string = "\033[31m"
    Cgreen 	string = "\033[32m"
    Cyellow 	string = "\033[33m"
    Cblue 	string = "\033[34m"
    Cpurple 	string = "\033[35m"
    Ccyan 	string = "\033[36m"
    Cwhite 	string = "\033[37m"
)

func green( input_string string ) string {
	return Cgreen + input_string + Creset
} 
