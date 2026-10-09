# Deploy on Railway

The repo includes a `Dockerfile` and `railway.json`. Railway builds the image, terminates HTTPS, and routes traffic to the app over HTTP on `$PORT`. One MySQL service can be shared by irvis, scvis, and eevis; each app creates its own database on first start.

1. In the Railway project, use the shared **MySQL** database service (add one if the project has none).
2. Add a service from this GitHub repo (`vaccs/scvis-go`), using the `railway` branch until it is merged.
3. In the service's **Variables**, set:

```text
DB_HOST=${{MySQL.MYSQLHOST}}
DB_PORT=${{MySQL.MYSQLPORT}}
DB_USER=${{MySQL.MYSQLUSER}}
DB_PASSWD=${{MySQL.MYSQLPASSWORD}}
DB_NAME=scvis
DB_CREATE_DATABASE=true
VACCS_ADDR=<host:port of the VACCS analysis server>
ADMIN_EMAIL=<your email>
ADMIN_PASSWORD=<initial admin password>
GITHUB_TOKEN=<token, if the issue feature is used>
```

4. Under **Settings → Networking**, click **Generate Domain**. The app is at `https://<domain>/scvis/`.

On startup the app creates any missing tables, and creates the admin user from `ADMIN_*` only when the users table is empty. After the first login you can remove `ADMIN_PASSWORD`.

## VACCS analysis server

Analyses are sent over TCP to the VACCS analysis server (default `localhost:3580`). The server is not part of this repo. Set `VACCS_ADDR` to wherever it runs, for example a second Railway service reached over private networking (`<service>.railway.internal:3580`). Without it, pages and logins work but analyses fail.

## Settings used by the image

- `SERVER_PROTOCOL=http`: serve plain HTTP (Railway provides HTTPS)
- `LOG_DIR=stdout`: write logs to the console, shown in Railway's log view
- `PORT`: set by Railway
