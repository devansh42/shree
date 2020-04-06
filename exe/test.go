package exe

import (
	"fmt"
	"log"
	"net/http"
)

//starts testing http server on given port
func StartTestHttpServer(port int) {

	http.HandleFunc("/"+fmt.Sprint(port), func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world\n"))
		w.Write([]byte("Here is the remote addr\n"))
		w.Write([]byte(r.RemoteAddr))
		w.WriteHeader(200)
	})
	log.Println("Testing server is listening at ", port)
	go http.ListenAndServe(JoinHost("", port), nil)

}
