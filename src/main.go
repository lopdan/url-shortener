package main

import (
	"fmt"
	"os"
)

func main() {

}

/** Port in 8000*/
func HttpPort() string {
	port := "8000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
}
