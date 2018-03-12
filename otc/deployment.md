# OTC Deployment instructions

## Environement
Following software must be installed on the deployment environment:
1. Go
2. Nginx
3. NodeJS (>8)
4. Yarn

The 'otc' user is created on the deployment environment.

## Backend
1. Copy built OTC to the `/usr/share/otc`.
2. Install and configure Btcwallet (`go get github.com/btcsuite/btcwallet`) according to its Getting started guide
3. Create **Btcwallet service** config file (`/etc/systemd/system/multi-user.target.wants/btcw.service`):
```
[Unit]
Description=Btcwallet
Documentation=
After=network-online.target remote-fs.target nss-lookup.target
Wants=network-online.target

[Service]
User=otc
Type=simple
PIDFile=/tmp/btcwallet.service.pid
ExecStart=path_to_btcwallet -C path_to_btcwallet_config_file
ExecReload=/bin/pkill -9 -F /tmp/btcwallet.service.pid
ExecStop=/bin/pkill -9 -F /tmp/btcwallet.service.pid
Restart=always

[Install]
WantedBy=multi-user.target
```
4. Install and configure **Skycoin node**
5. Create **Skycoin node** config file (`/etc/systemd/system/multi-user.target.wants/skycoin.service`) with following content:

```
[Unit]
Description=SkyCoin node
Documentation=
After=network-online.target remote-fs.target nss-lookup.target
Wants=network-online.target

[Service]
User=otc
Type=simple
PIDFile=/tmp/skycoin.service.pid
ExecStart=path_to_skycoin_node
ExecReload=/bin/pkill -9 -F /tmp/skycoin.service.pid
ExecStop=/bin/pkill -9 -F /tmp/skycoin.service.pid
Restart=always

[Install]
WantedBy=multi-user.target
```

6. Create OTC service config file (`/etc/systemd/system/multi-user.target.wants/otc.service`) with following content:

```
[Unit]
Description=SkyCoin OTC
Documentation=
After=network-online.target remote-fs.target nss-lookup.target
Wants=network-online.target btcw.service skycoin.service

[Service]
User=otc
Type=simple
WorkingDirectory=/usr/share/otc/
PIDFile=/tmp/otc.service.pid
ExecStart=/usr/share/otc/otc
ExecReload=/bin/pkill -9 -F /tmp/otc.service.pid
ExecStop=/bin/pkill -9 -F /tmp/otc.service.pid
Restart=always

[Install]
WantedBy=multi-user.target
```
7. Run `sudo systemctl daemon-reload`

## Frontend

### Nginx configuraion
1. Remove the `/etc/nginx/sites-enabled/default` file.
2. Create two folders: `/var/www/otc-admin-ui` and `/var/www/otc-ui`.
3. Add following files to the `/etc/nginx/sites-enabled/` folder:
#### otc-admin-ui
```
server {
        listen 8080 default_server;
        listen [::]:8080 default_server;

        root /var/www/otc-admin-ui;

        index index.html index.htm 

        server_name _;

       location /api {
                proxy_pass      http://127.0.0.1:8000;
                proxy_redirect  off;
        }

        location / {
                try_files $uri $uri/ =404;
        }
}
```
#### otc-ui
```
server {
        listen 80 default_server;
        listen [::]:80 default_server;

        root /var/www/otc-ui;

        index index.html index.htm 

        server_name _;

        location /api {
                proxy_pass      http://127.0.0.1:8081;
                proxy_redirect  off;
        }

        location / {
                try_files $uri $uri/ =404;
        }
}
```
4. Restore packages in the `otc-web` and `oct-web-admin` folders of the **services** repository using `yarn` command.
5. Build both client apps using `npm run build` command.
6. Copy `build` folders on the `otc-web` and `oct-web-admin` folders to the `/var/www/otc-admin-ui` and `/var/www/otc-ui`.
7. Run `systemctl restart nginx` command.