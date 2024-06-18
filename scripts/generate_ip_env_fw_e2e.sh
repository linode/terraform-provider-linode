#!/bin/bash

# Fetch public IPv4 address
PUBLIC_IPV4=$(curl -s https://api.ipify.org)

# Fetch public IPv6 address
PUBLIC_IPV6=$(curl -s https://ifconfig.co/ip)

# Create .env file in /tmp directory
ENV_FILE="/tmp/ip_vars.env"

cat << EOF > "$ENV_FILE"
PUBLIC_IPV4="$PUBLIC_IPV4"
PUBLIC_IPV6="$PUBLIC_IPV6"
EOF

# Display the path to the created .env file
echo "Generated .env file: $ENV_FILE"
