
# CCFIW

A bridge between cloudcampus and Telefonica's IotAgentJSON. This software runs and exits after single update. One should combine this with a cronjob in order to update with the desired frequency



## Features

- Quality of signal based on user's RSSI
- OK/KO AP's 
- Nr of online users
- Sites and Devices metadata mapping
- Distinct Users per Hour


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

`iotagent.force_update | default: false`
`iotagent.ignore_sites | default: `

`redis.servers | default: redis:6379`
`redis.database | default: 0`
`redis.cluster | default: false`
`redis.master_name | default: `
`redis.pool_size | default: empty(auto)`
`redis.password | default: `
`redis.tls_enabled | default: false`

## Requirements

- docker
- docker-compose

## Setup
 First, make sure you have all the required dependencies installed.
 Second, edit your docker-compose.yml file in order to setup the credencials and hosts desired.
 Third, run 
 ```bash 
    docker-compose up --build ccfiw
 ```
and wait for completion.

Hint: You can run the command above as a cronjob in order to keep data syncd.

## Migration Guide
Every new release that changes Datamodel Structure, should be migrated by enabling the env_var `iotagent.force_update` for the 1st run. After that, set the force_update to false to continue normaly.

## Ignoring unwanted Sites/PoI
By setting `iotagent.ignore_sites` (comma splitted PoI ID), sites (and their related entities) calculated from CloudCampus will not be updated on iotagent, if present on the list.