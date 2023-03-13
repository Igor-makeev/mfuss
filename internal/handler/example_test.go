package handler

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

func ExampleHandler_PostHandler() {
	client := resty.New()

	request := client.SetBaseURL("http://localhost:8080").
		R().
		SetCookie(
			&http.Cookie{
				Name:  generateCook(),
				Value: "[UUID]",
				Path:  "/",
			},
		)

	response, err := request.
		SetBody("https://kanobu.ru/").
		Post("/")

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode())
	}
}
