package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"go-service/internal/user/model"
)

type UserClient struct {
	Client *http.Client
	Url    string
}

func NewUserClient(client *http.Client, url string) *UserClient {
	return &UserClient{Client: client, Url: url}
}

func (c *UserClient) Load(ctx context.Context, id string) (*model.User, error) {
	requestURL := fmt.Sprintf("%s/%s", c.Url, id)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		panic(err)
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var res model.User
	err = json.Unmarshal(body, &res)
	if err != nil {
		panic(err)
	}

	return &res, err
}

func (c *UserClient) Create(ctx context.Context, user *model.User) (int64, error) {
	requestURL := c.Url

	data, err := json.Marshal(user)
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 201 {
		return 1, err
	} else {
		return 0, err
	}
}

func (c *UserClient) Update(ctx context.Context, user *model.User) (int64, error) {
	requestURL := fmt.Sprintf("%s/%s", c.Url, user.Id)

	data, err := json.Marshal(user)
	req, err := http.NewRequest("PUT", requestURL, bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 200 {
		return 1, err
	} else {
		return 0, err
	}
}

func (c *UserClient) Patch(ctx context.Context, user map[string]interface{}) (int64, error) {
	id := user["id"]
	requestURL := fmt.Sprintf("%s/%s", c.Url, id)

	data, err := json.Marshal(user)
	req, err := http.NewRequest("PATCH", requestURL, bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 200 {
		return 1, err
	} else {
		return 0, err
	}
}

func (c *UserClient) Delete(ctx context.Context, id string) (int64, error) {
	requestURL := fmt.Sprintf("%s/%s", c.Url, id)

	req, err := http.NewRequest("DELETE", requestURL, nil)
	if err != nil {
		panic(err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 200 {
		return 1, err
	} else {
		return 0, err
	}
}
