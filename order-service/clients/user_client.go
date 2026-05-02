package clients

import (
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

var client = resty.New().
	OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		log.Printf("[RESTY] Запрос к User-Service: %s %s", req.Method, req.URL)
		return nil
	}).
	OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		log.Printf("[RESTY] Ответ от User-Service: %d", resp.StatusCode())
		return nil
	})

func CheckUserExists(userId string) (bool, error) {
	url := fmt.Sprintf("http://user-service:8081/users/%s", userId)
	resp, err := client.R().Get(url)
	if err != nil {
		return false, err
	}
	return resp.StatusCode() == 200, nil
}
