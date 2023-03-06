package ccampus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

type DStatsData struct {
	SSID string `json:"ssid"`
	RSSI int    `json:"rssi"`
}

type DeviceStats struct {
	Total         int `json:"total"`
	TotalCitizens int
	TotalWorkers  int
	Data          []DStatsData `json:"data"`
}

type Device struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	SiteName     string `json:"siteName"`
	RegisterTime string `json:"registerTime"`
	StartupTime  string `json:"startupTime"`
	Description  string `json:"description"`
	Version      string `json:"version"`
	IP           string `json:"ip"`
	MAC          string `json:"mac"`
	Vendor       string `json:"vendor"`
	DeviceModel  string `json:"deviceModel"`
	Name         string `json:"name"`
	SiteID       string `json:"siteId"`
	IsOK         bool   `json:"is_ok"`
	Stats        DeviceStats
}

// DeviceMeta ...
type DeviceMeta struct {
	Total              int16    `json:"totalRecords"`
	Devices            []Device `json:"data"`
	OK                 int16
	KO                 int16
	TotalUsers         int
	TotalCitizens      int
	TotalWorkers       int
	TotalGoodQuality   int
	TotalMediumQuality int
	TotalPoorQuality   int
}

// Site is a CloudCampus Site
type Site struct {
	ID string `json:"id"`
	// TenantID    string   `json:"tenantId"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        []string `json:"type"`
	Tag         []string `json:"tag"`
	Latitude    string   `json:"latitude"`
	Longitude   string   `json:"longitude"`
	Address     string   `json:"address"`
	Email       string   `json:"email"`
	PostCode    string   `json:"postcode"`
	DeviceMeta  DeviceMeta
}

// SiteMeta to store all Sites and SiteMetadata
type SiteMeta struct {
	Total int16  `json:"totalRecords"`
	Sites []Site `json:"data"`
}

func (sm *SiteMeta) updateData(client *Client) error {
	if err := sm.getSiteData(client); err != nil {
		return err
	}
	for i := range sm.Sites {
		if err := sm.Sites[i].getDeviceData(client); err != nil {
			return err
		}
		fmt.Println("-----")

		for j, dev := range sm.Sites[i].DeviceMeta.Devices {
			status, err := strconv.Atoi(dev.Status)
			fmt.Println("status: ", status)
			if err != nil {
				return err
			}
			if status > 1 {
				sm.Sites[i].DeviceMeta.KO += 1
			} else {
				sm.Sites[i].DeviceMeta.OK += 1
				sm.Sites[i].DeviceMeta.Devices[j].IsOK = true
			}

			if dev.RegisterTime != "" {

				var regTime time.Time
				regTime, err = time.Parse("2006-01-02 15:04:05 DST", dev.RegisterTime)
				if err != nil {
					regTime, err = time.Parse("2006-01-02 15:04:05", dev.RegisterTime)
					if err != nil {
						return err
					}
				}
				sm.Sites[i].DeviceMeta.Devices[j].RegisterTime = regTime.Format(time.RFC3339)

			}

			if dev.StartupTime != "" {
				var startTime time.Time
				startTime, err = time.Parse("2006-01-02 15:04:05 DST", dev.StartupTime)
				if err != nil {
					startTime, err = time.Parse("2006-01-02 15:04:05", dev.StartupTime)
					if err != nil {
						return err
					}
				}
				sm.Sites[i].DeviceMeta.Devices[j].StartupTime = startTime.Format(time.RFC3339)
			}

			// get device stats
			if err := sm.Sites[i].getDeviceStats(client, j); err != nil {
				return err
			}
		}
		fmt.Println("OK: ", sm.Sites[i].DeviceMeta.OK)
		fmt.Println("KO:", sm.Sites[i].DeviceMeta.KO)
		fmt.Println("TotalUsers:", sm.Sites[i].DeviceMeta.TotalUsers)
		fmt.Println("TotalCitizens:", sm.Sites[i].DeviceMeta.TotalCitizens)
		fmt.Println("TotalWorkers:", sm.Sites[i].DeviceMeta.TotalWorkers)
	}

	fmt.Println(sm.Sites[0])

	return nil
}

func (sm *SiteMeta) getSiteData(client *Client) error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/controller/campus/v3/sites", client.basePath), nil)
	if err != nil {
		log.Error("get sites new request error: %w", err)
		return err
	}

	req.Header.Add("x-access-token", client.token)
	req.Header.Add("Accept", "*/*")

	resp, err := client.c.Do(req)
	if err != nil {
		log.Error("get sites request: %w", err)
		return err
	}
	fmt.Println(resp.StatusCode)

	if resp.StatusCode > 299 {
		return fmt.Errorf("failed to get sites. Status: %s", resp.Status)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("get sites body: %w", err)
		return err
	}
	if err := json.Unmarshal(b, sm); err != nil {
		log.Error("get sites body: %w", err)
		return err
	}
	return nil
}

func (s *Site) getDeviceData(client *Client) error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/controller/campus/v3/devices?siteId=%s", client.basePath, s.ID), nil)
	if err != nil {
		log.Error("get devices new request error: %w", err)
		return err
	}

	req.Header.Add("x-access-token", client.token)
	req.Header.Add("Accept", "*/*")
	resp, err := client.c.Do(req)
	if err != nil {
		log.Error("get devices request: %w", err)
		return err
	}
	// fmt.Println(resp.StatusCode)

	if resp.StatusCode > 299 {
		return fmt.Errorf("failed to get devices. Status: %s", resp.Status)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("get devices body: %w", err)
		return err
	}

	if err := json.Unmarshal(b, &s.DeviceMeta); err != nil {
		log.Error("get devices body: %w", err)
		return err
	}

	return nil
}

func (s *Site) getDeviceStats(client *Client, devidx int) error {
	devid := s.DeviceMeta.Devices[devidx].ID
	pageSize := 100
	pageIdx := 1
	for {
		stats, err := requestStationStats(client, devid, pageIdx, pageSize)
		if err != nil {
			return err
		}
		//s.DeviceMeta.Devices[devidx].Stats.Data = append(s.DeviceMeta.Devices[devidx].Stats.Data, stats.Data...)
		for _, st := range stats.Data {
			isWorkerSSID := false
			for _, wssid := range client.WorkerSSIDS {
				if st.SSID == wssid {
					isWorkerSSID = true
				}
			}

			if isWorkerSSID {
				s.DeviceMeta.Devices[devidx].Stats.TotalWorkers++
				s.DeviceMeta.TotalWorkers++
			} else {
				s.DeviceMeta.Devices[devidx].Stats.TotalCitizens++
				s.DeviceMeta.TotalCitizens++

			}

			// analyse RSSI
			if st.RSSI > -65 {
				//good
				s.DeviceMeta.TotalGoodQuality++
			} else if st.RSSI <= -80 {
				//poor
				s.DeviceMeta.TotalPoorQuality++
			} else {
				//medium
				s.DeviceMeta.TotalMediumQuality++
			}

			s.DeviceMeta.Devices[devidx].Stats.Total++
			s.DeviceMeta.TotalUsers++

		}
		if stats.Total <= pageIdx*pageSize {
			break
		}
		pageIdx++
	}
	// fmt.Printf("device %s | total users: %d citizens: %d workers: %d\n", devid,
	// 	s.DeviceMeta.Devices[devidx].Stats.Total,
	// 	s.DeviceMeta.Devices[devidx].Stats.TotalCitizens,
	// 	s.DeviceMeta.Devices[devidx].Stats.TotalWorkers)
	return nil
}

func requestStationStats(client *Client, devid string, pageIdx, pageSize int) (DeviceStats, error) {
	boDD := DeviceStats{}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/controller/campus/v1/performanceservice/station/client/device/%s?pageIndex=%d&pageSize=%d&status=online", client.basePath, devid, pageIdx, pageSize), nil)
	if err != nil {
		log.Error("get site stats new request error for page %d: %w", pageIdx, err)
		return boDD, err
	}
	req.Header.Add("x-access-token", client.token)
	req.Header.Add("Accept", "*/*")
	resp, err := client.c.Do(req)
	if err != nil {
		log.Error("get site stats request: %w", err)
		return boDD, err
	}

	if resp.StatusCode > 299 {
		return boDD, fmt.Errorf("failed to get site stats for page %d. Status: %s", pageIdx, resp.Status)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("get site stats body for page %d: %w", pageIdx, err)
		return boDD, err
	}
	if err := json.Unmarshal(b, &boDD); err != nil {
		log.Error("get site stats body for page %d: %w", pageIdx, err)
		return boDD, err
	}

	return boDD, nil
}
