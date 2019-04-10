# Atlas Notification Visualization

## Setup
Use a reverse proxy (nginx, haproxy) with tls certificates installed (letsencrypt acme) and proxy traffic to this service on `localhost:<port>`.

### Web Service
The following environment variables can be set to reconfigure the service:

`PORT` webservice port. default 8880

`STATICDIR` location of static server files. default ./www

`REDIS_ADDRESS`  Redis Address. default is localhost:6379. NOT YOUR ATLAS REDIS!

`REDIS_PASSWORD`  Redis Password. default is blank. NOT YOUR ATLAS REDIS!

`GLOBAL_ADMIN_STEAMID` The owners SteamID

`WEBHOOK_KEY` Random private string to prevent malicious notifications

### Web App
The client file "www/config.js" holds some cluster specific information like the grid size.
```
const config = {
    //Number of columns in the grid
    ServersX: 15,
	
    // Number of rows in the grid
    ServersY: 15,
}
```