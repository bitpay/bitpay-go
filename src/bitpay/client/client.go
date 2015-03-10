package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Client struct {
	ApiUri string
}

func (client *Client) PairWithCode(str string) (token map[string]string, err error) {
	token = make(map[string]string)
	match, _ := regexp.MatchString("^[[:alnum:]]{7}$", str)
	if !match {
		token = nil
		err = errors.New("BitPayArgumentError: invalid pairing code")
		return token, err
	}
	url := client.ApiUri + "/tokens"
	request, _ := http.NewRequest("POST", url, strings.NewReader(""))
	response, _ := http.DefaultClient.Do(request)
	contents, _ := ioutil.ReadAll(response.Body)
	var jsonContents map[string]interface{}
	json.Unmarshal(contents, &jsonContents)
	if response.StatusCode/100 != 2 {
		responseStatus := strconv.Itoa(response.StatusCode)
		contentError := responseStatus + fmt.Sprintf(": %s", jsonContents["error"])
		err = errors.New(contentError)
	} else {
		token["this"] = fmt.Sprintf("%s", contents)
		err = nil
	}
	return token, err
}
