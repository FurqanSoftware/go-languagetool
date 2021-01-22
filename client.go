package languagetool

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	DefaultBaseURL = "https://api.languagetoolplus.com/v2"
)

type Client struct {
	BaseURL string

	Username string
	APIKey   string
}

func (c Client) endpoint(path string) string {
	baseurl := c.BaseURL
	if baseurl == "" {
		baseurl = DefaultBaseURL
	}
	return c.BaseURL + path
}

func (c Client) authorize(data url.Values) {
	if c.Username == "" && c.APIKey == "" {
		return
	}
	data.Set("username", c.Username)
	data.Set("apiKey", c.APIKey)
}

func (c Client) Check(input CheckInput) (*CheckResult, error) {
	data := url.Values{}
	c.authorize(data)
	if input.Text != "" {
		data.Set("text", input.Text)
	}
	if input.Data != "" {
		data.Set("data", input.Data)
	}
	if input.Language != "" {
		data.Set("language", input.Language)
	}
	if input.Dicts != "" {
		data.Set("dicts", input.Dicts)
	}
	if input.MotherTongue != "" {
		data.Set("motherTongue", input.MotherTongue)
	}
	if input.PreferredVariants != "" {
		data.Set("preferredVariants", input.PreferredVariants)
	}
	if input.EnabledRules != "" {
		data.Set("enabledRules", input.EnabledRules)
	}
	if input.DisabledRules != "" {
		data.Set("disabledRules", input.DisabledRules)
	}
	if input.EnabledCategories != "" {
		data.Set("enabledCategories", input.EnabledCategories)
	}
	if input.DisabledCategories != "" {
		data.Set("disabledCategories", input.DisabledCategories)
	}
	data.Set("enabledOnly", strconv.FormatBool(input.EnabledOnly))
	if input.Level != "" {
		data.Set("level", string(input.Level))
	}
	req, err := http.NewRequest("POST", c.endpoint("/check"), bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, Error{StatusCode: resp.StatusCode}
	}

	body, _ := ioutil.ReadAll(resp.Body)
	result := CheckResult{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
