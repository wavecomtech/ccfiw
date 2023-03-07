package cmd

import (
	"ccfiw/internal/cache"
	"ccfiw/internal/ccampus"
	"ccfiw/internal/config"
	ngsiwifi "ccfiw/internal/ngsiWifi"
	"context"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func run(cmd *cobra.Command, args []string) error {

	ctx := context.Background()
	_, cancel := context.WithCancel(ctx)
	defer cancel()

	tasks := []func() error{
		setLogLevel,
		printStartMessage,
		setupCache,
		initCloudCampus,
		initNgsiWifi,
		readData,
		upsertData,
	}

	for _, t := range tasks {
		if err := t(); err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func setLogLevel() error {
	log.SetLevel(log.Level(uint8(config.C.General.LogLevel)))
	return nil
}

func printStartMessage() error {
	log.WithFields(log.Fields{
		"version":  version,
		"logLevel": log.GetLevel(),
	}).Info("starting CCFIW")
	return nil
}

func setupCache() error {
	if err := cache.Setup(config.C); err != nil {
		return errors.Wrap(err, "setup cache error")
	}

	return nil
}

func initCloudCampus() error {
	return ccampus.Setup(config.C)
}

func initNgsiWifi() error {
	return ngsiwifi.Setup(config.C)
}

func readData() error {
	return ccampus.Get().RefreshData()
}

func upsertData() error {
	data := ccampus.Get().GetData()
	source := "Huawei iMaster NCE"
	for _, site := range data.Sites {
		if strings.Contains(config.C.IoTAgent.IgnoreSites, site.ID) {
			log.Debugf("ignoring %s because is set on IgnoreSites", site.ID)
			continue
		}
		pl := ngsiwifi.WIFIPointOfInterest{
			ID:                          site.ID,
			Address:                     site.Address,
			Email:                       site.Email,
			Name:                        site.Name,
			Latitude:                    site.Latitude,
			Longitude:                   site.Longitude,
			Category:                    site.Type,
			Service:                     site.Tag,
			Description:                 site.Description,
			NrOfAPs:                     site.DeviceMeta.Total,
			NrOfAPsOK:                   site.DeviceMeta.OK,
			NrOfAPsKO:                   site.DeviceMeta.KO,
			Zip:                         site.PostCode,
			NrOfUsersConnected:          site.DeviceMeta.TotalUsers,
			NrOfWorkersConnected:        site.DeviceMeta.TotalWorkers,
			NrOfCitizensConnected:       site.DeviceMeta.TotalCitizens,
			NrOfUsersGoodQuality:        site.DeviceMeta.TotalGoodQuality,
			NrOfUsersMediumQuality:      site.DeviceMeta.TotalMediumQuality,
			NrOfUsersPoorQuality:        site.DeviceMeta.TotalPoorQuality,
			HrNumberOfUsersConnected:    site.DeviceMeta.CurHourUsers,
			HrNumberOfCitizensConnected: site.DeviceMeta.CurHourCitizens,
			HrNumberOfWorkersConnected:  site.DeviceMeta.CurHourWorkers,
		}
		if err := ngsiwifi.Get().UpdatePointOfInterest(source, pl, config.C.IoTAgent.ForceUpdate); err != nil {
			return err
		}
		for _, device := range site.DeviceMeta.Devices {
			pl := ngsiwifi.WIFIAccessPoint{
				ID:               device.ID,
				IsOK:             device.IsOK,
				AreaServed:       device.SiteName,
				DateInstalled:    device.RegisterTime,
				DateLastReboot:   device.StartupTime,
				Description:      device.Description,
				FirmwareVersion:  device.Version,
				IPAddress:        []string{device.IP},
				MACAddress:       []string{device.MAC},
				Manufacturer:     device.Vendor,
				ModelName:        device.DeviceModel,
				Name:             device.Name,
				PoiId:            site.Name,
				ClientsConnected: device.Stats.Total,
			}
			if err := ngsiwifi.Get().UpdateAccessPoint(source, pl); err != nil {
				return err
			}
		}
	}
	return nil
}
