package twocaptcha

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var base = "http://2captcha.com"

type captcha struct {
	key    string // api key
	client http.Client
}

func New(key string) *captcha {
	return &captcha{
		key:    key,
		client: *http.DefaultClient,
	}
}

func (c *captcha) Solve(recaptcha, origin string, invisible bool) (string, error) {
	id, err := c.submit(recaptcha, origin, invisible)
	if err != nil {
		return "", err
	}

	resp, err := c.fetch(id)
	if err != nil {
		return "", err
	}

	return resp, nil
}

func (c *captcha) submit(recaptcha, origin string, invisible bool) (string, error) {
	var invisibleBit uint8
	if invisible {
		invisibleBit = 1
	}

	data := url.Values{}
	data.Set("key", c.key)
	data.Set("method", "userrecaptcha")
	data.Set("googlekey", recaptcha)
	data.Set("pageurl", origin)
	data.Set("invisible", strconv.Itoa(int(invisibleBit)))

	req, err := http.NewRequest("POST", base+"/in.php", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", errors.New("network error")
	}

	bin, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	body := string(bin)

	if body == "ERROR_NO_SLOT_AVAILABLE" {
		time.Sleep(2 * time.Second)
		return c.submit(recaptcha, origin, invisible)
	}

	if body[:2] == "OK" {
		return body[3:], nil
	}

	return "", errors.New(body)
}

func (c *captcha) fetch(id string) (string, error) {
	time.Sleep(5 * time.Second)

	data := url.Values{}
	data.Set("key", c.key)
	data.Set("action", "get")
	data.Set("id", id)

	req, err := http.NewRequest("POST", base+"/res.php", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", errors.New("network error")
	}

	bin, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	body := string(bin)

	if body == "CAPCHA_NOT_READY" { // not a typo
		return c.fetch(id)
	}

	if body[:2] == "OK" {
		return body[3:], nil
	}

	return "", errors.New(body)
}
