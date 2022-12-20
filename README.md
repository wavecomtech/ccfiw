
# CCFIW

A bridge between cloudcampus and Telefonica's IotAgentJSON. This software runs and exits after single update. One should combine this with a cronjob in order to update with the desired frequency



## Features

- Quality of signal based on user's RSSI
- OK/KO AP's 
- Nr of online users
- Sites and Devices metadata mapping


## Support

**CloudCampus version:** iMaster NCE-Campus V300R021C00SPC111

**IoTAgent version:**
```javascript
{
    "libVersion": "2.24.0",
    "port": 8082,
    "baseRoot": "/",
    "version": "1.22.0"
}
```


## Environment Variables

To run this project, you will need to add the following environment variables to your .env file or directly to the docker-compose.yml file.

`ccampus.basepath | default: https://dashboard.lanaas.eu:18002`

`ccampus.username | default: `

`ccampus.password | default: `

`ccampus.workerssid | default: Corp,test1`

`idm.basepath | default: https://auth.iotplatform.telefonica.com:15001`

`idm.username | default: `

`idm.password | default: `

`idm.service | default: `

`idm.servicepath | default: `

`iotagent.apikey | default: <get one from iotagent POST iot/services> `

`iotagent.hostname | default: iota.iotplatform.telefonica.com`

`iotagent.iota_port | default: 8088`

`iotagent.json_port | default: 8185`
