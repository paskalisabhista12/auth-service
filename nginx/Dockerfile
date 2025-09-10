FROM openresty/openresty:alpine

# Install dependencies + luarocks
RUN apk add --no-cache \
    build-base \
    curl \
    git \
    unzip \
    lua5.1-dev \
    lua5.1 \
    luarocks5.1

# Install lua-resty-http and lua-cjson via luarocks
RUN luarocks-5.1 install lua-resty-http && \
    luarocks-5.1 install lua-cjson

# Copy nginx.conf
COPY nginx.conf /usr/local/openresty/nginx/conf/nginx.conf

