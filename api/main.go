package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sdslabs/katana/api/services/sensei"
)

var port uint32 = 3000

func RunKatanaAPIServer() {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      sensei.ServiceInit(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	server.ListenAndServe()
}
