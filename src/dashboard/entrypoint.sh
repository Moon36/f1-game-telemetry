#!/bin/sh

# Set default port if DASHBOARD_PORT is not provided
PORT=${DASHBOARD_PORT:-80}

# Set default backend port if BACKEND_PORT is not provided
BACKEND_PORT=${BACKEND_PORT:-8282}

# Inject configuration into the HTML
CONFIG_SCRIPT="<script>window.APP_CONFIG = { BACKEND_PORT: '${BACKEND_PORT}' };</script>"

# Insert the script before the closing head tag in index.html
sed -i "s|</head>|${CONFIG_SCRIPT}</head>|g" /usr/share/nginx/html/index.html

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
echo "Backend port configured as: ${BACKEND_PORT}"

# Start nginx
exec nginx -g "daemon off;"
