package ngsiwifi

import (
	"ccfiw/internal/config"
	"net/http"
	"time"
)

type NGSIWifi interface {
	UpdateAccessPoint(string, WIFIAccessPoint) error
	UpdatePointOfInterest(string, WIFIPointOfInterest) error
	Login() error
	ProvisionPOIGroup() error
}

var nw NGSIWifi

// Get returns the NGSIWifi client.
func Get() NGSIWifi {
	return nw
}

func Setup(c config.Config) error {
	nw = &ngsiwifi{
		Username: c.IDM.Username,
		Password: c.IDM.Password,
		Client: &Client{
			c: &http.Client{
				Timeout: 30 * time.Second,
			},
			method:      "https",
			apikey:      c.IoTAgent.APIKey,
			hostname:    c.IoTAgent.HostName,
			iota_port:   c.IoTAgent.IoTAPort,
			json_port:   c.IoTAgent.JSONPort,
			IDMbasePath: c.IDM.BasePath,
			service:     c.IDM.Service,
			servicePath: c.IDM.ServicePath,
		},
	}
	if err := nw.Login(); err != nil {
		return err
	}
	if err := nw.ProvisionPOIGroup(); err != nil {
		return err
	}
	return nil
}

type ngsiwifi struct {
	Client   *Client
	Username string
	Password string
}

func (c *ngsiwifi) Login() error {
	if err := c.Client.Login(c.Username, c.Password); err != nil {
		return err
	}
	return nil
}
