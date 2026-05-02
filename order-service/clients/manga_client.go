package clients

import (
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

var mangaClient = resty.New().
	OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		log.Printf("[RESTY] Запрос к Manga-Service: %s %s", req.Method, req.URL)
		return nil
	}).
	OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		log.Printf("[RESTY] Ответ от Manga-Service: %d", resp.StatusCode())
		return nil
	})

func CheckMangaExists(mangaId string) (bool, error) {
	url := fmt.Sprintf("http://manga-service:8082/manga/%s", mangaId)
	resp, err := mangaClient.R().Get(url)
	if err != nil {
		return false, err
	}
	return resp.StatusCode() == 200, nil
}
