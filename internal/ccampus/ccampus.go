package ccampus

import (
	"ccfiw/internal/config"
	"net/http"
	"strings"
	"time"
)

type CCampus interface {
	RefreshData() error
	GetData() *SiteMeta
}

var cc CCampus

// Get returns the CCampus client.
func Get() CCampus {
	return cc
}

func Setup(c config.Config) error {
	cc = &ccampus{
		Username: c.CCAMPUS.Username,
		Password: c.CCAMPUS.Password,
		Data:     &SiteMeta{},
		Client: &Client{
			c: &http.Client{
				Timeout: 30 * time.Second,
			},
			basePath:    c.CCAMPUS.BasePath,
			WorkerSSIDS: strings.Split(c.CCAMPUS.Workerssid, ","),
		},
	}
	return nil
}

type ccampus struct {
	Data     *SiteMeta
	Client   *Client
	Username string
	Password string
}

func (c *ccampus) RefreshData() error {
	if err := c.Client.Login(c.Username, c.Password); err != nil {
		return err
	}
	return c.Data.updateData(c.Client)
}

func (c *ccampus) GetData() *SiteMeta {
	return c.Data
}
