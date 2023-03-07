package ngsiwifi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Client struct {
	c           *http.Client
	token       string
	method      string
	apikey      string
	hostname    string
	iota_port   int16
	json_port   int16
	IDMbasePath string
	service     string
	servicePath string
}

func (t *Client) Login(username, password string) error {
	body := map[string]interface{}{
		"auth": map[string]interface{}{
			"identity": map[string]interface{}{
				"methods": []string{"password"},
				"password": map[string]interface{}{
					"user": map[string]interface{}{
						"domain": map[string]interface{}{
							"name": t.service,
						},
						"name":     username,
						"password": password,
					},
				},
			},
			"scope": map[string]interface{}{
				"project": map[string]interface{}{
					"domain": map[string]interface{}{
						"name": t.service,
					},
					"name": t.servicePath,
				},
			},
		},
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}
	// fmt.Println(string(bodyBytes))
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v3/auth/tokens", t.IDMbasePath), bytes.NewReader(bodyBytes))
	if err != nil {
		log.Error("ngsi: get token error: %w", err)
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	resp, err := t.c.Do(req)
	if err != nil {
		log.Error("ngsi: get token request: %w", err)
		return err
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("ngsi: failed to get token. Status: %s", resp.Status)
	}
	t.token = resp.Header.Get("X-Subject-Token")
	return nil
}

func (t *Client) DeviceExists(id string) bool {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s://%s:%d/iot/devices/%s?protocol=IoTA-JSON", t.method, t.hostname, t.iota_port, id), nil)
	if err != nil {
		log.Error("ngsi: get device error: %w", err)
		return false
	}

	req.Header.Add("Accept", "*/*")
	req.Header.Add("Fiware-Service", t.service)
	req.Header.Add("Fiware-ServicePath", t.servicePath)
	req.Header.Add("X-Auth-Token", t.token)
	resp, err := t.c.Do(req)
	if err != nil {
		log.Error("ngsi: get device request: %w", err)
		return false
	}

	if resp.StatusCode > 299 {
		log.Errorf("ngsi: get device %s. Does not exist. Status %s", id, resp.Status)
		return false
	}
	log.Debugf("ngsi: get device %s. Exists", id)
	return true
}
