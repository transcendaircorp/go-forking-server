# go-forking-server
Simple go server which forwards serial stream to all connected tcp clients, and vice versa

# Run as service with systemd
Put the following in /etc/systemd/system/serial-tcp.service:
```
[Unit]
Description=Serial Forwarding TCP Server
After=network.target

[Service]
Type=simple
Restart=always
RestartSec=1
User=root
ExecStart=/path/to/executable /path/to/device port

[Install]
WantedBy=multi-user.target
```
Then run `systemctl enable serial-tcp && systemctl start serial-tcp` and authenticate. If that fails run the previous command as root.
