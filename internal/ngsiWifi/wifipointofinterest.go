package ngsiwifi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mozillazg/go-slugify"

	log "github.com/sirupsen/logrus"
)

type WIFIPointOfInterest struct {
	ID          string   `json:"id"`
	Source      string   `json:"source"`      //auto-generated
	TimeInstant string   `json:"TimeInstant"` //auto-generated
	Name        string   `json:"name"`
	Latitude    string   `json:"-"`
	Longitude   string   `json:"-"`
	Location    []string `json:"location"` // auto-generated
	// Location    struct {
	// 	Type        string   `json:"type"`        // auto-generated
	// 	Coordinates []string `json:"coordinates"` // auto-generated
	// } `json:"location"` // auto-generated
	Type         string   `json:"type"` // auto-generated
	Category     []string `json:"category"`
	Service      []string `json:"service"`
	Address      string   `json:"address"`
	Email        string   `json:"contactPoint"`
	DataProvider string   `json:"dataProvider"` //auto-generated

	Description string `json:"description"`
	Zip         string `json:"zip"`
	// total AP's
	NrOfAPs       int16  `json:"rtNumberOfAPs"`
	NrOfAPsSource string `json:"rtNumberOfAPsSource"`    //auto-generated
	NrOfAPsUpdAT  string `json:"rtNumberOfAPsUpdatedAt"` //auto-generated
	// OK AP's
	NrOfAPsOK       int16  `json:"rtNumberOfAPsOK"`
	NrOfAPsOKSource string `json:"rtNumberOfAPsOKSource"`    //auto-generated
	NrOfAPsOKUpdAT  string `json:"rtNumberOfAPsOKUpdatedAt"` //auto-generated
	// KO AP's
	NrOfAPsKO       int16  `json:"rtNumberOfAPsKO"`
	NrOfAPsKOSource string `json:"rtNumberOfAPsKOSource"`    //auto-generated
	NrOfAPsKOUpdAT  string `json:"rtNumberOfAPsKOUpdatedAt"` //auto-generated
	// wifi status
	WifiStatus       string `json:"wifiStatus"`          //auto-generated
	WifiStatusSource string `json:"wifiStatusSource"`    //auto-generated
	WifiStatusUpdAt  string `json:"wifiStatusUpdatedAt"` //auto-generated
	// total Users online
	NrOfUsersConnected       int    `json:"rtNumberOfUsersConnected"`
	NrOfUsersConnectedSource string `json:"rtNumberOfUsersConnectedSource"`    //auto-generated
	NrOfUsersConnectedUpdAT  string `json:"rtNumberOfUsersConnectedUpdatedAt"` //auto-generated
	// total Citizens online
	NrOfCitizensConnected       int    `json:"rtNumberOfCitizensConnected"`
	NrOfCitizensConnectedSource string `json:"rtNumberOfCitizensConnectedSource"`    //auto-generated
	NrOfCitizensConnectedUpdAT  string `json:"rtNumberOfCitizensConnectedUpdatedAt"` //auto-generated
	// total Workers online
	NrOfWorkersConnected       int    `json:"rtNumberOfWorkersConnected"`
	NrOfWorkersConnectedSource string `json:"rtNumberOfWorkersConnectedSource"`    //auto-generated
	NrOfWorkersConnectedUpdAT  string `json:"rtNumberOfWorkersConnectedUpdatedAt"` //auto-generated

	// total Users Good connection
	NrOfUsersGoodQuality       int    `json:"rtNumberOfUsersGoodQuality"`
	NrOfUsersGoodQualitySource string `json:"rtNumberOfUsersGoodQualitySource"`    //auto-generated
	NrOfUsersGoodQualityUpdAT  string `json:"rtNumberOfUsersGoodQualityUpdatedAt"` //auto-generated
	// total Users Medium connection
	NrOfUsersMediumQuality       int    `json:"rtNumberOfUsersMediumQuality"`
	NrOfUsersMediumQualitySource string `json:"rtNumberOfUsersMediumQualitySource"`    //auto-generated
	NrOfUsersMediumQualityUpdAT  string `json:"rtNumberOfUsersMediumQualityUpdatedAt"` //auto-generated
	// total Users Poor connection
	NrOfUsersPoorQuality       int    `json:"rtNumberOfUsersPoorQuality"`
	NrOfUsersPoorQualitySource string `json:"rtNumberOfUsersPoorQualitySource"`    //auto-generated
	NrOfUsersPoorQualityUpdAT  string `json:"rtNumberOfUsersPoorQualityUpdatedAt"` //auto-generated
	// Distinct Users connected in the current hour
	HrNumberOfUsersConnected                      int64  `json:"hrNumberOfUsersConnected"`
	HrNumberOfUsersConnectedSource                string `json:"hrNumberOfUsersConnectedSource"`                //auto-generated
	HrNumberOfUsersConnectedUpdatedAt             string `json:"hrNumberOfUsersConnectedUpdatedAt"`             //auto-generated
	HrNumberOfUsersConnectedCalculationPeriodFrom string `json:"hrNumberOfUsersConnectedCalculationPeriodFrom"` // auto-generated and truncated by hour
	// Distinct Citizens connected in the current hour
	HrNumberOfCitizensConnected                      int64  `json:"hrNumberOfCitizensConnected"`
	HrNumberOfCitizensConnectedSource                string `json:"hrNumberOfCitizensConnectedSource"`                //auto-generated
	HrNumberOfCitizensConnectedUpdatedAt             string `json:"hrNumberOfCitizensConnectedUpdatedAt"`             //auto-generated
	HrNumberOfCitizensConnectedCalculationPeriodFrom string `json:"hrNumberOfCitizensConnectedCalculationPeriodFrom"` // auto-generated and truncated by hour
	// Distinct Workers connected in the current hour
	HrNumberOfWorkersConnected                      int64  `json:"hrNumberOfWorkersConnected"`
	HrNumberOfWorkersConnectedSource                string `json:"hrNumberOfWorkersConnectedSource"`                //auto-generated
	HrNumberOfWorkersConnectedUpdatedAt             string `json:"hrNumberOfWorkersConnectedUpdatedAt"`             //auto-generated
	HrNumberOfWorkersConnectedCalculationPeriodFrom string `json:"hrNumberOfWorkersConnectedCalculationPeriodFrom"` // auto-generated and truncated by hour
}

func (d *WIFIPointOfInterest) Validate(source string) error {
	d.Type = "WifiPointOfInterest"
	d.Name = slugify.Slugify(d.Name)
	// fill source
	d.Source = source
	d.NrOfAPsKOSource = source
	d.NrOfAPsOKSource = source
	d.NrOfAPsSource = source
	d.WifiStatusSource = source
	d.DataProvider = source
	d.NrOfUsersConnectedSource = source
	d.NrOfWorkersConnectedSource = source
	d.NrOfCitizensConnectedSource = source
	d.NrOfUsersGoodQualitySource = source
	d.NrOfUsersMediumQualitySource = source
	d.NrOfUsersPoorQualitySource = source
	d.HrNumberOfUsersConnectedSource = source
	d.HrNumberOfCitizensConnectedSource = source
	d.HrNumberOfWorkersConnectedSource = source

	// resolve wifi status
	if d.NrOfAPs == 0 {
		d.WifiStatus = "noService"
	} else if d.NrOfAPsKO == 0 {
		d.WifiStatus = "working"
	} else if d.NrOfAPsOK == 0 {
		d.WifiStatus = "totalFailure"
	} else {
		d.WifiStatus = "workingPartially"
	}

	// resolve location
	d.Location = []string{d.Longitude, d.Latitude}

	// updated at fields
	ntime := time.Now()
	ftime := ntime.Format(time.RFC3339)
	d.TimeInstant = ftime
	d.WifiStatusUpdAt = ftime
	d.NrOfAPsKOUpdAT = ftime
	d.NrOfAPsOKUpdAT = ftime
	d.NrOfAPsUpdAT = ftime
	d.NrOfUsersConnectedUpdAT = ftime
	d.NrOfWorkersConnectedUpdAT = ftime
	d.NrOfCitizensConnectedUpdAT = ftime
	d.NrOfUsersGoodQualityUpdAT = ftime
	d.NrOfUsersMediumQualityUpdAT = ftime
	d.NrOfUsersPoorQualityUpdAT = ftime
	d.HrNumberOfUsersConnectedUpdatedAt = ftime
	d.HrNumberOfCitizensConnectedUpdatedAt = ftime
	d.HrNumberOfWorkersConnectedUpdatedAt = ftime

	// update truncation time
	ttime := ntime.Truncate(time.Hour).Format(time.RFC3339)
	d.HrNumberOfUsersConnectedCalculationPeriodFrom = ttime
	d.HrNumberOfWorkersConnectedCalculationPeriodFrom = ttime
	d.HrNumberOfCitizensConnectedCalculationPeriodFrom = ttime

	return nil
}

func (c *ngsiwifi) UpdatePointOfInterest(source string, data WIFIPointOfInterest, force_update bool) error {
	if err := data.Validate(source); err != nil {
		return err
	}
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	log.Debugf("----POI----(force update = %v)", force_update)
	// fmt.Println(string(body))
	if !c.Client.DeviceExists(data.ID) {
		log.Infof("POI device %s does not exists", data.ID)
		if err := c.RegisterPOI(data.ID, data.Name, false); err != nil {
			log.Errorf("failed to register POI %s", data.ID)
			return err
		}
	} else if force_update {
		if err := c.RegisterPOI(data.ID, data.Name, true); err != nil {
			log.Errorf("failed to force_update POI fields for %s", data.ID)
			return err
		}
	}
	log.Debugf("sending measurements for POI %s", data.ID)
	// fmt.Printf("POST /iot/json | BODY: %v", string(body))

	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("http://%s:%d/iot/json?k=%s&i=%s&getCmd=0", c.Client.hostname, c.Client.json_port, c.Client.apikey, data.ID),
		bytes.NewReader(body))
	if err != nil {
		log.Error("ngsi: send measurement for POI %s error: %w", data.ID, err)
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Fiware-Service", c.Client.service)
	req.Header.Add("Fiware-ServicePath", c.Client.servicePath)
	req.Header.Add("X-Auth-Token", c.Client.token)
	req.Header.Add("Accept", "*/*")
	// log.Debug(req.Header)
	resp, err := c.Client.c.Do(req)
	if err != nil {
		log.Error("ngsi: sending measurements for POI %s request: %w", data.ID, err)
		return err
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("ngsi: failed to send measurements for POI %s | Status: %s", data.ID, resp.Status)
	}
	log.Infof("sent measurements for POI %s | Status %s", data.ID, resp.Status)
	return nil
}

func (c *ngsiwifi) ProvisionPOIGroup() error {
	body := map[string]interface{}{
		"services": []map[string]interface{}{
			{
				"attributes": []map[string]interface{}{
					{
						"object_id": "TimeInstant",
						"name":      "TimeInstant",
						"type":      "DateTime",
					},
					{
						"object_id": "Service",
						"name":      "Service",
						"type":      "Text",
					},
					{
						"object_id": "name",
						"name":      "name",
						"type":      "Text",
					},
					{
						"object_id": "location",
						"name":      "location",
						"type":      "geo:json",
					},
					{
						"object_id": "category",
						"name":      "category",
						"type":      "Array",
					},
					{
						"object_id": "address",
						"name":      "address",
						"type":      "Text",
					},
					{
						"object_id": "email",
						"name":      "email",
						"type":      "Text",
					},
					{
						"object_id": "dataProvider",
						"name":      "dataProvider",
						"type":      "Text",
					},
					{
						"object_id": "description",
						"name":      "description",
						"type":      "Text",
					},
					{
						"object_id": "zip",
						"name":      "zip",
						"type":      "Text",
					},
					{
						"object_id": "zip",
						"name":      "zip",
						"type":      "Text",
					},
					{
						"object_id": "rtNumberOfAPs",
						"name":      "rtNumberOfAPs",
						"type":      "Number",
					},
					{
						"object_id": "rtNumberOfAPsSource",
						"name":      "rtNumberOfAPsSource",
						"type":      "Text",
					},
					{
						"object_id": "rtNumberOfAPsUpdatedAt",
						"name":      "rtNumberOfAPsUpdatedAt",
						"type":      "DateTime",
					},
					{
						"object_id": "rtNumberOfAPsOK",
						"name":      "rtNumberOfAPsOK",
						"type":      "Number",
					},
					{
						"object_id": "rtNumberOfAPsOKSource",
						"name":      "rtNumberOfAPsOKSource",
						"type":      "Text",
					},
					{
						"object_id": "rtNumberOfAPsOKUpdatedAt",
						"name":      "rtNumberOfAPsOKUpdatedAt",
						"type":      "DateTime",
					},
					{
						"object_id": "rtNumberOfAPsKO",
						"name":      "rtNumberOfAPsKO",
						"type":      "Number",
					},
					{
						"object_id": "rtNumberOfAPsKOSource",
						"name":      "rtNumberOfAPsKOSource",
						"type":      "Text",
					},
					{
						"object_id": "rtNumberOfAPsKOUpdatedAt",
						"name":      "rtNumberOfAPsKOUpdatedAt",
						"type":      "DateTime",
					},
					{
						"object_id": "wifiStatus",
						"name":      "wifiStatus",
						"type":      "Text",
					},
					{
						"object_id": "wifiStatusSource",
						"name":      "wifiStatusSource",
						"type":      "Text",
					},
					{
						"object_id": "wifiStatusUpdatedAt",
						"name":      "wifiStatusUpdatedAt",
						"type":      "DateTime",
					},

					{
						"object_id": "rtNumberOfUsersConnected",
						"name":      "rtNumberOfUsersConnected",
						"type":      "Number",
					},
					{
						"object_id": "rtNumberOfUsersConnectedSource",
						"name":      "rtNumberOfUsersConnectedSource",
						"type":      "Text",
					},
					{
						"object_id": "rtNumberOfUsersConnectedUpdatedAt",
						"name":      "rtNumberOfUsersConnectedUpdatedAt",
						"type":      "DateTime",
					},

					{
						"object_id": "rtNumberOfCitizensConnected",
						"name":      "rtNumberOfCitizensConnected",
						"type":      "Number",
					},
					{
						"object_id": "rtNumberOfCitizensConnectedSource",
						"name":      "rtNumberOfCitizensConnectedSource",
						"type":      "Text",
					},
					{
						"object_id": "rtNumberOfCitizensConnectedUpdatedAt",
						"name":      "rtNumberOfCitizensConnectedUpdatedAt",
						"type":      "DateTime",
					},

					{
						"object_id": "rtNumberOfWorkersConnected",
						"name":      "rtNumberOfWorkersConnected",
						"type":      "Number",
					},
					{
						"object_id": "rtNumberOfWorkersConnectedSource",
						"name":      "rtNumberOfWorkersConnectedSource",
						"type":      "Text",
					},
					{
						"object_id": "rtNumberOfWorkersConnectedUpdatedAt",
						"name":      "rtNumberOfWorkersConnectedUpdatedAt",
						"type":      "DateTime",
					},

					{
						"object_id": "hrNumberOfUsersConnected",
						"name":      "hrNumberOfUsersConnected",
						"type":      "Number",
					},
					{
						"object_id": "hrNumberOfWorkersConnected",
						"name":      "hrNumberOfWorkersConnected",
						"type":      "Number",
					},
					{
						"object_id": "hrNumberOfCitizensConnected",
						"name":      "hrNumberOfCitizensConnected",
						"type":      "Number",
					},
					{
						"object_id": "hrNumberOfUsersConnectedUpdatedAt",
						"name":      "hrNumberOfUsersConnectedUpdatedAt",
						"type":      "DateTime",
					},
					{
						"object_id": "hrNumberOfWorkersConnectedUpdatedAt",
						"name":      "hrNumberOfWorkersConnectedUpdatedAt",
						"type":      "DateTime",
					}, {
						"object_id": "hrNumberOfCitizensConnectedUpdatedAt",
						"name":      "hrNumberOfCitizensConnectedUpdatedAt",
						"type":      "DateTime",
					},
					{
						"object_id": "hrNumberOfUsersConnectedCalculationPeriodFrom",
						"name":      "hrNumberOfUsersConnectedCalculationPeriodFrom",
						"type":      "DateTime",
					},
					{
						"object_id": "hrNumberOfWorkersConnectedCalculationPeriodFrom",
						"name":      "hrNumberOfWorkersConnectedCalculationPeriodFrom",
						"type":      "DateTime",
					}, {
						"object_id": "hrNumberOfCitizensConnectedCalculationPeriodFrom",
						"name":      "hrNumberOfCitizensConnectedCalculationPeriodFrom",
						"type":      "DateTime",
					},
					{
						"object_id": "hrNumberOfUsersConnectedSource",
						"name":      "hrNumberOfUsersConnectedSource",
						"type":      "Text",
					}, {
						"object_id": "hrNumberOfWorkersConnectedSource",
						"name":      "hrNumberOfWorkersConnectedSource",
						"type":      "Text",
					}, {
						"object_id": "hrNumberOfCitizensConnectedSource",
						"name":      "hrNumberOfCitizensConnectedSource",
						"type":      "Text",
					},
				},
				"static_attributes": []map[string]interface{}{},
				"apikey":            c.Client.apikey,
				"description":       "Provision Group WIFIPointOfInterest",
				"protocol": []string{
					"IoTA-JSON",
				},
				"entity_type": "WifiPointOfInterest",
			},
		},
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"body": fmt.Sprintf("%+v", body),
	}).Debugln("POST /iot/services")
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://%s:%d/iot/services", c.Client.hostname, c.Client.iota_port), bytes.NewReader(bodyBytes))
	if err != nil {
		log.Error("ngsi: provision poi error: %w", err)
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Fiware-Service", c.Client.service)
	req.Header.Add("Fiware-ServicePath", c.Client.servicePath)
	req.Header.Add("X-Auth-Token", c.Client.token)
	req.Header.Add("Accept", "*/*")
	fmt.Println(req.Header)

	resp, err := c.Client.c.Do(req)
	if err != nil {
		log.Error("ngsi: provision poi request: %w", err)
		return err
	}

	if resp.StatusCode > 299 && resp.StatusCode != 409 {
		return fmt.Errorf("ngsi: failed to provision poi. Status: %s", resp.Status)
	}

	log.Debugf("provision poi success. Status %s", resp.Status)
	return nil
}

func (c *ngsiwifi) RegisterPOI(device_id, entity_name string, force_update bool) error {
	attributes := []map[string]interface{}{
		{
			"object_id": "TimeInstant",
			"name":      "TimeInstant",
			"type":      "DateTime",
		},
		{
			"object_id": "Service",
			"name":      "Service",
			"type":      "Text",
		},
		{
			"object_id": "name",
			"name":      "name",
			"type":      "Text",
		},
		{
			"object_id": "location",
			"name":      "location",
			"type":      "geo:json",
		},
		{
			"object_id": "category",
			"name":      "category",
			"type":      "Array",
		},
		{
			"object_id": "address",
			"name":      "address",
			"type":      "Text",
		},
		{
			"object_id": "email",
			"name":      "email",
			"type":      "Text",
		},
		{
			"object_id": "dataProvider",
			"name":      "dataProvider",
			"type":      "Text",
		},
		{
			"object_id": "description",
			"name":      "description",
			"type":      "Text",
		},
		{
			"object_id": "zip",
			"name":      "zip",
			"type":      "Text",
		},
		{
			"object_id": "rtNumberOfAPs",
			"name":      "rtNumberOfAPs",
			"type":      "Number",
		},
		{
			"object_id": "rtNumberOfAPsSource",
			"name":      "rtNumberOfAPsSource",
			"type":      "Text",
		},
		{
			"object_id": "rtNumberOfAPsUpdatedAt",
			"name":      "rtNumberOfAPsUpdatedAt",
			"type":      "DateTime",
		},
		{
			"object_id": "rtNumberOfAPsOK",
			"name":      "rtNumberOfAPsOK",
			"type":      "Number",
		},
		{
			"object_id": "rtNumberOfAPsOKSource",
			"name":      "rtNumberOfAPsOKSource",
			"type":      "Text",
		},
		{
			"object_id": "rtNumberOfAPsOKUpdatedAt",
			"name":      "rtNumberOfAPsOKUpdatedAt",
			"type":      "DateTime",
		},
		{
			"object_id": "rtNumberOfAPsKO",
			"name":      "rtNumberOfAPsKO",
			"type":      "Number",
		},
		{
			"object_id": "rtNumberOfAPsKOSource",
			"name":      "rtNumberOfAPsKOSource",
			"type":      "Text",
		},
		{
			"object_id": "rtNumberOfAPsKOUpdatedAt",
			"name":      "rtNumberOfAPsKOUpdatedAt",
			"type":      "DateTime",
		},
		{
			"object_id": "wifiStatus",
			"name":      "wifiStatus",
			"type":      "Text",
		},
		{
			"object_id": "wifiStatusSource",
			"name":      "wifiStatusSource",
			"type":      "Text",
		},
		{
			"object_id": "wifiStatusUpdatedAt",
			"name":      "wifiStatusUpdatedAt",
			"type":      "DateTime",
		},
		{
			"object_id": "rtNumberOfUsersConnected",
			"name":      "rtNumberOfUsersConnected",
			"type":      "Number",
		},
		{
			"object_id": "rtNumberOfUsersConnectedSource",
			"name":      "rtNumberOfUsersConnectedSource",
			"type":      "Text",
		},
		{
			"object_id": "rtNumberOfUsersConnectedUpdatedAt",
			"name":      "rtNumberOfUsersConnectedUpdatedAt",
			"type":      "DateTime",
		},

		{
			"object_id": "rtNumberOfCitizensConnected",
			"name":      "rtNumberOfCitizensConnected",
			"type":      "Number",
		},
		{
			"object_id": "rtNumberOfCitizensConnectedSource",
			"name":      "rtNumberOfCitizensConnectedSource",
			"type":      "Text",
		},
		{
			"object_id": "rtNumberOfCitizensConnectedUpdatedAt",
			"name":      "rtNumberOfCitizensConnectedUpdatedAt",
			"type":      "DateTime",
		},

		{
			"object_id": "rtNumberOfWorkersConnected",
			"name":      "rtNumberOfWorkersConnected",
			"type":      "Number",
		},
		{
			"object_id": "rtNumberOfWorkersConnectedSource",
			"name":      "rtNumberOfWorkersConnectedSource",
			"type":      "Text",
		},
		{
			"object_id": "rtNumberOfWorkersConnectedUpdatedAt",
			"name":      "rtNumberOfWorkersConnectedUpdatedAt",
			"type":      "DateTime",
		},

		{
			"object_id": "hrNumberOfUsersConnected",
			"name":      "hrNumberOfUsersConnected",
			"type":      "Number",
		},
		{
			"object_id": "hrNumberOfWorkersConnected",
			"name":      "hrNumberOfWorkersConnected",
			"type":      "Number",
		},
		{
			"object_id": "hrNumberOfCitizensConnected",
			"name":      "hrNumberOfCitizensConnected",
			"type":      "Number",
		},
		{
			"object_id": "hrNumberOfUsersConnectedUpdatedAt",
			"name":      "hrNumberOfUsersConnectedUpdatedAt",
			"type":      "DateTime",
		},
		{
			"object_id": "hrNumberOfWorkersConnectedUpdatedAt",
			"name":      "hrNumberOfWorkersConnectedUpdatedAt",
			"type":      "DateTime",
		}, {
			"object_id": "hrNumberOfCitizensConnectedUpdatedAt",
			"name":      "hrNumberOfCitizensConnectedUpdatedAt",
			"type":      "DateTime",
		},
		{
			"object_id": "hrNumberOfUsersConnectedCalculationPeriodFrom",
			"name":      "hrNumberOfUsersConnectedCalculationPeriodFrom",
			"type":      "DateTime",
		},
		{
			"object_id": "hrNumberOfWorkersConnectedCalculationPeriodFrom",
			"name":      "hrNumberOfWorkersConnectedCalculationPeriodFrom",
			"type":      "DateTime",
		}, {
			"object_id": "hrNumberOfCitizensConnectedCalculationPeriodFrom",
			"name":      "hrNumberOfCitizensConnectedCalculationPeriodFrom",
			"type":      "DateTime",
		},
		{
			"object_id": "hrNumberOfUsersConnectedSource",
			"name":      "hrNumberOfUsersConnectedSource",
			"type":      "Text",
		}, {
			"object_id": "hrNumberOfWorkersConnectedSource",
			"name":      "hrNumberOfWorkersConnectedSource",
			"type":      "Text",
		}, {
			"object_id": "hrNumberOfCitizensConnectedSource",
			"name":      "hrNumberOfCitizensConnectedSource",
			"type":      "Text",
		},
	}

	if force_update {
		path := fmt.Sprintf("iot/devices/%s?protocol=IoTA-JSON", device_id)
		log.Warnf("force_updating device %s with name %s", device_id, entity_name)
		body := map[string]interface{}{
			"attributes": attributes,
		}
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return err
		}

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s://%s:%d/%s", c.Client.method, c.Client.hostname, c.Client.iota_port, path), bytes.NewReader(bodyBytes))
		if err != nil {
			log.Error("ngsi: register poi %s error: %w", device_id, err)
			return err
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Fiware-Service", c.Client.service)
		req.Header.Add("Fiware-ServicePath", c.Client.servicePath)
		req.Header.Add("X-Auth-Token", c.Client.token)
		req.Header.Add("Accept", "*/*")
		log.WithFields(log.Fields{
			"body":    fmt.Sprintf("%+v", string(bodyBytes)),
			"headers": req.Header,
		}).Debugf("PUT /%s", path)

		resp, err := c.Client.c.Do(req)
		if err != nil {
			log.Error("ngsi: force_update poi %s request: %w", device_id, err)
			return err
		}
		if resp.StatusCode > 299 {
			log.Error(resp)
		}
		// force continue to add if not found
		if resp.StatusCode != 404 {
			return nil
		}
	}
	log.Infof("registering device %s with name %s", device_id, entity_name)

	body := map[string]interface{}{
		"devices": []map[string]interface{}{
			{
				"device_id":         device_id,
				"entity_name":       entity_name,
				"entity_type":       "WifiPointOfInterest",
				"attributes":        attributes,
				"static_attributes": []map[string]interface{}{},
				"apikey":            c.Client.apikey,
				"protocol":          "IoTA-JSON",
				"transport":         "HTTP",
			},
		},
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}
	// fmt.Printf("POST /iot/devices | BODY: %v", string(bodyBytes))

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s://%s:%d/iot/devices", c.Client.method, c.Client.hostname, c.Client.iota_port), bytes.NewReader(bodyBytes))
	if err != nil {
		log.Error("ngsi: register poi %s error: %w", device_id, err)
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Fiware-Service", c.Client.service)
	req.Header.Add("Fiware-ServicePath", c.Client.servicePath)
	req.Header.Add("X-Auth-Token", c.Client.token)
	req.Header.Add("Accept", "*/*")
	log.Debug(req.Header)

	resp, err := c.Client.c.Do(req)
	if err != nil {
		log.Error("ngsi: provision poi %s request: %w", device_id, err)
		return err
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("ngsi: failed to provision poi %s. Status: %s", device_id, resp.Status)
	}
	log.Infof("provision poi %s success. Status %s", device_id, resp.Status)
	return nil
}
