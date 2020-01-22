[Unit]
Description=berserk.red
After=sshd.service
Requires=

[Service]
ExecStart=/root/berserk/red/red
Restart=always
RestartSec=10s
Type=simple
LimitNOFILE=30002
Environment="PORT=9000" 
[Install]
WantedBy=multi-user.target