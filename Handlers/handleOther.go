package handlers

import (
	"net/http"
	"fmt"
)

func HandleOther(w http.ResponseWriter, r *http.Request) {
	fmt.Println( "Unknown endpoint requested" )
	fmt.Println( r.URL.String() )
	fmt.Println( r )

	return
}