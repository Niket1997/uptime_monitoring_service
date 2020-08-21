package platform

import (
	"fmt"
	"net/http"
	"time"
)

// GetRequest function to perform requests of type GET
func GetRequest(url string, timeout int) string {
	client := http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}
	// context read about it : create deadlines

	// Will throw error as it's not quick enough
	_, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		return "Inactive"
	}
	return "Active"
}
