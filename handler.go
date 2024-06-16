package main

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"

	"github.com/nuclio/nuclio-sdk-go"
	"github.com/otiai10/gosseract/v2"
)

type UserData struct {
	client *gosseract.Client
}

type request struct {
	ImageName   string `json:"name"`
	Base64Image string `json:"data"`
	ImageUrl    string `json:"url"`
}

type response struct {
	ImageName string `json:"name"`
	Content   string `json:"content"`
}

func Handler(context *nuclio.Context, event nuclio.Event) (interface{}, error) {
	var response response
	userData := context.UserData.(*UserData)

	imageRequst, err := parseRequest(event)
	if err != nil {
		return faasResponse(http.StatusBadRequest, nil, err)
	}

	// get image data from body or downlaod from url
	imageData, err := getImageData(imageRequst)
	if err != nil {
		return faasResponse(http.StatusBadRequest, nil, err)
	}

	err = userData.client.SetImageFromBytes(imageData)
	if err != nil {
		return faasResponse(http.StatusBadRequest, nil, err)
	}

	text, err := userData.client.Text()
	if err != nil {
		return faasResponse(http.StatusBadRequest, nil, err)
	}

	response.ImageName = imageRequst.ImageName
	response.Content = text

	body, err := json.Marshal(response)
	if err != nil {
		return faasResponse(http.StatusBadRequest, nil, err)
	}

	return faasResponse(http.StatusOK, body, err)
}

func InitContext(context *nuclio.Context) error {
	var userData UserData

	client := gosseract.NewClient()

	userData.client = client

	context.UserData = &userData

	return nil
}

func parseRequest(event nuclio.Event) (request, error) {
	var ocrRequest request

	err := json.Unmarshal(event.GetBody(), &ocrRequest)
	if err != nil {
		return ocrRequest, err
	}

	return ocrRequest, nil
}

func getImageData(imageRequest request) ([]byte, error) {
	if imageRequest.ImageUrl != "" {
		res, err := http.Get(imageRequest.ImageUrl)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		imageData, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		return imageData, nil
	}

	imageData, err := base64.StdEncoding.DecodeString(imageRequest.Base64Image)
	if err != nil {
		return nil, err
	}

	return imageData, nil
}

func faasResponse(statusCode int, body []byte, err error) (nuclio.Response, error) {
	return nuclio.Response{
		StatusCode:  statusCode,
		ContentType: "application/json",
		Body:        body,
	}, err
}
