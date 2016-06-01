# Kitty Monitor

#### Project Status: Beta

## Synopsis
Kitty Monitor is a system for monitoring the temperature of a remote location.
* For simplicity and cost-effectiveness, the internal board temperature of an Odroid C1+/C2 (buy at hardkernel.com or ameridroid.com) is the sensor used
* The client is designed to run on a Ubuntu based operating system
* The Odroid locally polls it's internal system temperature and saves it to a SQLite database (typically every 2mins)
* The Odroid sends any unsent readings to the server (typically every 4mins)
* The server collects and displays these readings in a webserver
* There is much more that can be done on the webserver side, but displaying raw readings for now: http://gonotes.net:9080/q/all/l/20
* The temperature provided is not an exact reading of ambient temperature, but an indication in case of say failure of an airconditioning unit

## Getting Setup

### Download
Get Go for your operating system: http://golang.org/dl/ and install. The installation should set the GOROOT environment variable to the Golang toplevel folder.
If you installed Go to the folder ```~/apps/``` then your GOROOT env. var should contain: ```~/apps/go```.

### Go Workspace Setup
The environment variable *GOPATH* must point to your Go project workspace. Google 'environment variable set Windows', for example, to do that for your operating system.
Your can check your current GOPATH and GOROOT environment variables by running the following command
```
go env
```

Your Go projects should live under a folder path of this format:
```
GOPATH/src/yourdomain.com/your_project
```

## Getting and Building KittyMonitor
Note that both the server and client is provided from the same binary.
(You may need Mercurial installed)

```
sudo apt-get install mercurial # may not be necessary
# Get the KittyMonitor archive
go get github.com/rohanthewiz/kitty_mon
cd $GOPATH/src/github.com/rohanthewiz/kitty_mon
go build # this will produce the executable 'kitty_mon' in the current directory
```
You will require the ego package to make changes to the template files
Install it with

```
go get github.com/benbjohnson/ego/...  # yes the three dots at the end are necessary
```
If you make changes to an ego template file, you will need to 'compile' it before the main build.
For example if you updated *query.ego* you will need to run

```
$GOPATH/bin/ego -package main  # compile the template before doing 'go build'
```

## Using KittyMonitor
KittyMonitor launches a goroutine that polls system board temperature of the Odroid.
The main thread continues in an infinite loop to:
 1. poll temperature every 2mins
 2. connect to the server every 4mins and send any unsent readings.
 
Usage info: `kitty_mon -h` 

### Get the server's secret token (you'll need to copy the token unto the client)

```
./kitty_mon -get_server_secret # => 9A7blahblah123...
```

### Starting the Gob/Web Server

```
$ ./kitty_mon
```

### Starting the Client in Development mode (Readings every 8s. The server_secret is from above)

```
$ ./kitty_mon -synch_client ServerNameOrIP -server_secret 9A7blahblah123...
```

### Starting the Client in Production mode (Readings every 2 minutes)

```
$ ./kitty_mon -synch_client ServerNameOrIP -server_secret 9A7blahblah123... -env prod
```
(The secret code is only needed the first time the client talks to the server)

### Monitor Raw data
In a browser go to http://ServerNameOrIP:9080
See my example at http://gonotes.net:9080

### Listing Readings

	/q/:query/l/:limit # replace :limit with a numeric limit

Examples:
```
    http://localhost:8080/q/all/l/60 # List 60 readings
```
You will be able to delete (soon) by clicking the title of the reading. This takes you to single reading view.
From there you can click the 'Delete' link

### Other Options
    
    -h -- List available options with defaults
    -db "" -- Sqlite DB path. KittyMon will try to create the database 'kitty_mon.sqlite' in your home directory by default
    -l "-1" -- Limit the number of readings returned - default: -1 (no limit)
    -admin="" -- Privileged actions like 'delete_tables' (drops the readings table)

### Local Testing
For testing you can run a client and server on the same machine, just use different databases
Example:

```
# Start the server with the default database
./kitty_mon
# Start the client with test.sqlite
./kitty_mon -db test.sqlite -synch_client localhost -node_name TestClient
# View at localhost:9080
```
If no readings are listed then your workstation may not support internal temperature reporting.
That's okay, the Odroid C1/C2 do!

### TODO
- A lot since we are in Beta
- Summarization, graphs, dashboard
- SMS Alerts
- Token based auth webserver mode

### TIPS
- There is a great article on ego at http://blog.gopheracademy.com/advent-2014/ego/
- gcc is required to compile SQLite. On Windows you can get 64bit MinGW here http://mingw-w64.sourceforge.net/download.php. Install it and make sure to add the bin directory to your path
  (Could not get this to work with Windows 8.1, Windows 7 did work though)
- For less typing, you might want to do 'go build -o km' (km.exe on Windows) to produce the executable 'km'
- This is a sweet way to learn a modern, highly performant language - The Go Programming Language using a database with an ORM (object relational manager) while building a useful tool!
I recommend using git to checkout an early version of kitty_mon so you can start out simple
- Firefox has a great addon called SQLite Manager which you can use to peek into the database file
- Feel free to create a pull request if you'd like to pitch in.

### Credits
- Go -- http://golang.org/  Thanks Google!!
- GORM -- https://github.com/jinzhu/gorm  - Who needs sluggish ActiveRecord, or other interpreted code interfacing to your database.
- SQLite -- http://www.sqlite.org/ - A great place to start. Actually GORM includes all the things needed for SQLite so SQLite gets compiled into KittyMon!
