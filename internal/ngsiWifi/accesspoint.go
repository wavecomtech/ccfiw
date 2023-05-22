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

type WIFIAccessPoint struct {
	ID                    string   `json:"id"`
	Type                  string   `json:"type"`             //auto-generated
	TimeInstant           string   `json:"TimeInstant"`      //auto-generated
	Source                string   `json:"source"`           //auto-generated
	APState               string   `json:"apState"`          //auto-generated
	APStateSource         string   `json:"apStateSource"`    //auto-generated
	APStateSourceUpdtAt   string   `json:"apStateUpdatedAt"` //auto-generated
	IsOK                  bool     `json:"-"`
	AreaServed            string   `json:"areaServed"`
	ConnectionType        string   `json:"connectionType"` //auto-generated
	DataProvider          string   `json:"dataProvider"`   //auto-generated
	DateInstalled         string   `json:"dateInstalled,omitempty"`
	DateLastReboot        string   `json:"dateLastReboot,omitempty"`
	DateLastValueReported string   `json:"dateLastValueReported,omitempty"`
	Description           string   `json:"description"`
	FirmwareVersion       string   `json:"firmwareVersion"`
	IPAddress             []string `json:"ipAddress"`
	MACAddress            []string `json:"macAddress"`
	Manufacturer          string   `json:"manufacturer"`
	ModelName             string   `json:"modelName"`
	Name                  string   `json:"name"`
	PoiId                 string   `json:"poiId"`
	Provider              string   `json:"provider"` //auto-generated
	Municipality          string   `json:"municipality"`

	// total Users online
	ClientsConnected       int    `json:"clientsConnected"`
	ClientsConnectedSource string `json:"clientsConnectedSource"`    //auto-generated
	ClientsConnectedUpdAT  string `json:"clientsConnectedUpdatedAt"` //auto-generated

}

func (wa *WIFIAccessPoint) Validate(source string) error {
	wa.Type = "AccessPoint"
	wa.ConnectionType = "wireless"
	wa.Name = slugify.Slugify(wa.Name)
	wa.PoiId = slugify.Slugify(wa.PoiId)
	// fill source
	wa.Source = source
	wa.Provider = source
	wa.APStateSource = source
	wa.DataProvider = source
	wa.ClientsConnectedSource = source
	// updated at fields
	time := time.Now().Format(time.RFC3339)
	wa.TimeInstant = time
	wa.APStateSourceUpdtAt = time
	wa.ClientsConnectedUpdAT = time
	wa.DateLastValueReported = time
	if wa.Municipality == "" {
		wa.Municipality = "NA"
	}

	// fill state
	if wa.IsOK {
		wa.APState = "up"
	} else {
		wa.APState = "down"
	}
	return nil
}

func (c *ngsiwifi) UpdateAccessPoint(source string, data WIFIAccessPoint) error {
	if err := data.Validate(source); err != nil {
		return err
	}
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"data": fmt.Sprintf("%+v", data),
	}).Debugln("----AP----")
	if !c.Client.DeviceExists(data.ID) {
		log.Infof("AP device %s does not exists", data.ID)
		if err := c.RegisterAP(data.ID, data.Name); err != nil {
			log.Errorf("failed to register AP %s", data.ID)
			return err
		}
	}
	log.Debugf("sending measurements for AP %s", data.ID)
	url := fmt.Sprintf("http://%s:%d/iot/json?k=%s&i=%s&getCmd=0", c.Client.hostname, c.Client.json_port, c.Client.apikey, data.ID)
	req, err := http.NewRequest(http.MethodPost,
		url, bytes.NewReader(body))
	if err != nil {
		log.Error("ngsi: send measurement for AP %s error: %w", data.ID, err)
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
		log.Error("ngsi: sending measurements for AP %s request: %w", data.ID, err)
		return err
	}

	if resp.StatusCode > 299 {
		fmt.Println(url)
		return fmt.Errorf("ngsi: failed to send measurements for AP %s | Status: %s", data.ID, resp.Status)
	}
	log.Infof("sent measurements for AP %s | Status %s", data.ID, resp.Status)
	return nil
}

func (c *ngsiwifi) RegisterAP(device_id, entity_name string) error {
	log.Infof("registering ap %s with name %s", device_id, entity_name)
	body := map[string]interface{}{
		"devices": []map[string]interface{}{
			{
				"device_id":   device_id,
				"entity_name": entity_name,
				"entity_type": "AccessPoint",
				"attributes": []map[string]interface{}{
					{
						"object_id": "TimeInstant",
						"name":      "TimeInstant",
						"type":      "DateTime",
					},
					{
						"object_id": "address",
						"name":      "address",
						"type":      "Text",
					},
					{
						"object_id": "apState",
						"name":      "apState",
						"type":      "Text",
					},
					{
						"object_id": "apStateSource",
						"name":      "apStateSource",
						"type":      "Text",
					},
					{
						"object_id": "apStateUpdatedAt",
						"name":      "apStateUpdatedAt",
						"type":      "DateTime",
					},
					{
						"object_id": "areaServed",
						"name":      "areaServed",
						"type":      "Text",
					},
					{
						"object_id": "clientsConnected",
						"name":      "clientsConnected",
						"type":      "Number",
					},
					{
						"object_id": "clientsConnectedSource",
						"name":      "clientsConnectedSource",
						"type":      "Text",
					},
					{
						"object_id": "clientsConnectedUpdatedAt",
						"name":      "clientsConnectedUpdatedAt",
						"type":      "DateTime",
					},
					{
						"object_id": "dataProvider",
						"name":      "dataProvider",
						"type":      "Text",
					},
					{
						"object_id": "contactPoint",
						"name":      "contactPoint",
						"type":      "Text",
					},
					{
						"object_id": "connectionType",
						"name":      "connectionType",
						"type":      "Text",
					},
					{
						"object_id": "dateInstalled",
						"name":      "dateInstalled",
						"type":      "DateTime",
					},
					{
						"object_id": "dateLastReboot",
						"name":      "dateLastReboot",
						"type":      "DateTime",
					},
					{
						"object_id": "dateLastValueReported",
						"name":      "dateLastValueReported",
						"type":      "DateTime",
					},
					{
						"object_id": "description",
						"name":      "description",
						"type":      "Text",
					},
					{
						"object_id": "firmwareVersion",
						"name":      "firmwareVersion",
						"type":      "Text",
					},
					{
						"object_id": "hardwareVersion",
						"name":      "hardwareVersion",
						"type":      "Text",
					},
					{
						"object_id": "ipAddress",
						"name":      "ipAddress",
						"type":      "Array",
					},
					{
						"object_id": "location",
						"name":      "location",
						"type":      "geo:json",
					},
					{
						"object_id": "macAddress",
						"name":      "macAddress",
						"type":      "Array",
					},
					{
						"object_id": "manufacturer",
						"name":      "manufacturer",
						"type":      "Text",
					},
					{
						"object_id": "modelID",
						"name":      "modelID",
						"type":      "Number",
					},
					{
						"object_id": "modelName",
						"name":      "modelName",
						"type":      "Text",
					},
					{
						"object_id": "name",
						"name":      "name",
						"type":      "Text",
					},
					{
						"object_id": "osVersion",
						"name":      "osVersion",
						"type":      "Text",
					},
					{
						"object_id": "owner",
						"name":      "owner",
						"type":      "Array",
					},
					{
						"object_id": "poiId",
						"name":      "poiId",
						"type":      "Relationship",
					},
					{
						"object_id": "provider",
						"name":      "provider",
						"type":      "Text",
					},
					{
						"object_id": "refSwitch",
						"name":      "refSwitch",
						"type":      "Text",
					},
					{
						"object_id": "serialNumber",
						"name":      "serialNumber",
						"type":      "Text",
					},
					{
						"object_id": "service",
						"name":      "service",
						"type":      "Text",
					},
					{
						"object_id": "softwareVersion",
						"name":      "softwareVersion",
						"type":      "Text",
					},
					{
						"object_id": "source",
						"name":      "source",
						"type":      "Text",
					},
					{
						"object_id": "ssid",
						"name":      "ssid",
						"type":      "Array",
					},
					{
						"object_id": "zip",
						"name":      "zip",
						"type":      "Text",
					},
					{
						"object_id": "zone",
						"name":      "zone",
						"type":      "Text",
					},
					{
						"object_id": "district",
						"name":      "district",
						"type":      "Text",
					},
					{
						"object_id": "municipality",
						"name":      "municipality",
						"type":      "Text",
					},
					{
						"object_id": "province",
						"name":      "province",
						"type":      "Text",
					},
					{
						"object_id": "region",
						"name":      "region",
						"type":      "Text",
					},
					{
						"object_id": "community",
						"name":      "community",
						"type":      "Text",
					},
					{
						"object_id": "country",
						"name":      "country",
						"type":      "Text",
					},
					{
						"object_id": "streetAddress",
						"name":      "streetAddress",
						"type":      "Text",
					},
					{
						"object_id": "postalCode",
						"name":      "postalCode",
						"type":      "Text",
					},
					{
						"object_id": "addressLocality",
						"name":      "addressLocality",
						"type":      "Text",
					},
					{
						"object_id": "addressRegion",
						"name":      "addressRegion",
						"type":      "Text",
					},
					{
						"object_id": "addressCommunity",
						"name":      "addressCommunity",
						"type":      "Text",
					},
					{
						"object_id": "addressCountry",
						"name":      "addressCountry",
						"type":      "Text",
					},
				},
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
	fmt.Println(string(bodyBytes))

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s://%s:%d/iot/devices", c.Client.method, c.Client.hostname, c.Client.iota_port), bytes.NewReader(bodyBytes))
	if err != nil {
		log.Error("ngsi: register ap %s error: %w", device_id, err)
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
		log.Error("ngsi: provision ap %s request: %w", device_id, err)
		return err
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("ngsi: failed to provision ap %s. Status: %s", device_id, resp.Status)
	}
	log.Infof("provision ap %s success. Status %s", device_id, resp.Status)
	return nil
}
