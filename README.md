# Cloudstatus
Cloudstatus is a small Golang REST API which shows the status of public cloud services.

A Go routine (can be considered as a lightweight thread) polls the health status of Azure and GitHub periodicly and stores them into a PostgreSQL database.

By accessing the API endpoint the latest health status is returned.


## Running the app with docker:

The provided docker-compose.yml is the fastest way to spawn the app and the required PostgreSQL database.


```
# download latest docker-compose.yml or (alternatively you can clone the whole repository):
$ wget https://raw.githubusercontent.com/hofmann-works/cloudstatus/main/docker-compose.yml

# run:
$ docker-compose up
# run as deamon:
$ docker-compose up -d
```

Once the PostgreSQL and Cloudstatus instances are launched, you can test the application by polling the API endpoint.

Tested with:  
* Docker version 17.05.0-ce
* docker-compose version 1.27.4
* CentOS 7 and Fedora 32/Podman

## Get Status of last checks
`GET /v1/status` - returnes the latest check result stored in PostgreSQL database.

`$ curl -i -H 'Accept: application/json' http://localhost:8080/v1/status`

Initial Response when the cloud services status have not been pulled yet (wait for pull interval):

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Wed, 28 Oct 2020 08:09:13 GMT
Content-Length: 15

{"Clouds":null}
```

Response when there is status data in the PostgreSQL database:

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Wed, 28 Oct 2020 08:09:28 GMT
Content-Length: 190

{
   "Clouds":[
      {
         "Name":"Azure",
         "LastUpdated":"2020-10-27T17:30:05.573Z",
         "UnhealthyServices":null
      },
      {
         "Name":"GitHub",
         "LastUpdated":"2020-10-28T08:03:37.41Z",
         "UnhealthyServices":[
            "API Requests"
         ]
      }
   ]
}
```

In this example output Azure does currently not have any unhealthy services.
It is currently a incident at GitHub's service 'API Requests'.

`LastUpdated` describes the time when the respective cloud service has updated its status page.

## Usage

### CI/CD Pipelines
This application will be mainly used in CI/CD Pipelines which depends on Azure (e.g. as hosting environment) or GitHub (e.g. to fetch Golang Modules).  
A CI/CD Task can run before a build/deployment and return a warning in case there is an incident at the respective cloud service.

### Monitoring
Monitoring applications (like icinga/nagios) can check the status of multiple cloud services by polling this API, and send notifications if necessary.

## Folder Structure

`main.go` contains the entrypoint of the application.

`config/` reads environment variables and uses default values if none are set.

`checks/` fetches status APIs of different cloud services and parses them.

`db/` initializes the PostgreSQL database connection and inserts/selects data.

`handlers/` contains RestAPI Endpoints.

## Environment Variables

* `CLOUDSTATUS_PollInterval` - (Seconds) Interval in which Azure and GitHub Status APIs will be polled.  
  * Default: 100
* `CLOUDSTATUS_AzureStatusURL` - (URL) Which will be polled to retrieve the Status of Azure.
  * Default: https://status.dev.azure.com/_apis/status/health?geographies=EU,US&api-version=6.0-preview.1"
  * By modifying "geographies=" different Regions can be checked.
* `CLOUDSTATUS_GitHubStatusURL` - (URL) Which will be polled to retrieve the Status of GitHub.
  * Default: https://kctbh9vrtdwd.statuspage.io/api/v2/summary.json
* `CLOUDSTATUS_PGHost` - (hostname) PostgreSQL hostname
  * Default: localhost
* `CLOUDSTATUS_PGPort` - (port) PostgreSQL port
  * Default: 5432
* `CLOUDSTATUS_PGDatabase` - (Database Name) PostgreSQL Database
  * Defualt: cloudstatus
* `CLOUDSTATUS_PGUser` - (Username) PostgreSQL User with access to Database
  * Default: cloudstatus
* `CLOUDSTATUS_PGPassword` - (Password) PostgreSQL User's password
  * Default: mypw


## Build/Run locally without Docker:


### Prerequisites:
* Go Version 1.14.9 or higher
* PostgreSQL 13.0 or higher

### Run
1. Clone Repository:  
`$ git clone https://github.com/hofmann-works/cloudstatus.git`

2. Set environment variables or update `config/config.go` for PostgreSQL settings.

3. Build: `go build -o cloudstatus .`

4. run: `./cloudstatus`

## Future Improvements:

* Add API Endpoints (History, Single Cloud Status)
* Add Cloud Services and provide a possiblity to check only necessary services (AWS, GCP, Heroku, GitLab, NPM,...)
* Unify naming (Since this is my first Golang project and the requirements changed during development naming conventions does not always match best practices)