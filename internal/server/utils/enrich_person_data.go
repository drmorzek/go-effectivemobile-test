package utils

import (
	"encoding/json"
	"go-test/internal/db"
	"net/http"
)

type AgifyResponse struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type GenderizeResponse struct {
	Name   string `json:"name"`
	Gender string `json:"gender"`
}

type NationalizeResponse struct {
	Name    string `json:"name"`
	Country []struct {
		CountryID   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	}
}

func EnrichPersonData(person *db.Person) error {
	client := &http.Client{}

	resp, err := client.Get("https://api.agify.io/?name=" + person.Name)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var agifyResponse AgifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&agifyResponse); err != nil {
		return err
	}
	person.Age = agifyResponse.Age

	resp, err = client.Get("https://api.genderize.io/?name=" + person.Name)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var genderizeResponse GenderizeResponse
	if err := json.NewDecoder(resp.Body).Decode(&genderizeResponse); err != nil {
		return err
	}
	person.Gender = genderizeResponse.Gender

	resp, err = client.Get("https://api.nationalize.io/?name=" + person.Name)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var nationalizeReponse NationalizeResponse
	if err := json.NewDecoder(resp.Body).Decode(&nationalizeReponse); err != nil {
		return err
	}

	if len(nationalizeReponse.Country) > 0 { // Используем страну с наибольшей вероятностью
		person.Nationality = nationalizeReponse.Country[0].CountryID
	}

	return nil
}
