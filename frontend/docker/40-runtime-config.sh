#!/bin/sh
set -eu

: "${BACKEND_BASE_URL:=http://localhost:8080}"

# Strip trailing slash and escape characters that would break the JS string
url=$(printf '%s' "${BACKEND_BASE_URL%/}" | sed "s/[\\\\']/\\\\&/g")

printf "window.__env = { BACKEND_BASE_URL: '%s' };\n" "$url" > /usr/share/nginx/html/config.js
