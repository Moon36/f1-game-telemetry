#!/bin/sh

# Set default port if DASHBOARD_PORT is not provided
PORT=${DASHBOARD_PORT:-80}

# Custom nginx config
cat > /etc/nginx/conf.d/default.conf << EOF
server {
    listen ${PORT};
    server_name localhost;
    location / {
        root /usr/share/nginx/html;
        index index.html index.htm;
        try_files \$uri \$uri/ /index.html;
    }
}
EOF

echo "Dashboard will be available on port ${PORT}"

# Start nginx
exec nginx -g "daemon off;"
