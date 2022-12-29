package ccampus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Client struct {
	c           *http.Client
	token       string
	basePath    string
	WorkerSSIDS []string
}

func (t *Client) Login(username, password string) error {
	body := struct {
		Username string `json:"userName"`
		Password string `json:"password"`
	}{
		Username: username,
		Password: password,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/controller/v2/tokens", t.basePath), bytes.NewReader(bodyBytes))
	if err != nil {
		log.Error("get token error: %w", err)
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	resp, err := t.c.Do(req)
	if err != nil {
		log.Error("get token request: %w", err)
		return err
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("failed to get token. Status: %s", resp.Status)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("get token body: %w", err)
		return err
	}
	sm := struct {
		Data struct {
			TokenID string `json:"token_id"`
			Expire  string `json:"expiredDate"`
		} `json:"data"`
	}{}
	if err := json.Unmarshal(b, &sm); err != nil {
		log.Error("get token body: %w", err)
		return err
	}

	t.token = sm.Data.TokenID
	return nil
}
