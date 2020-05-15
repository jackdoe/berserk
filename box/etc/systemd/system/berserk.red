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
Environment="GEMINI_PORT=1965"
Environment="GEMINI_CRT=/etc/letsencrypt/live/berserk.red/cert.pem"
Environment="GEMINI_KEY=/etc/letsencrypt/live/berserk.red/privkey.pem"
[Install]
WantedBy=multi-user.target