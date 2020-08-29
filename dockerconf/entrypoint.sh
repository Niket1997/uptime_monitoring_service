#!/bin/sh
#set -euo pipefail
initialize()
{
  echo "$GIT_COMMIT_HASH" > /app/public/commit.txt
}

start_webservers()
{
    start_nginx
    start_app
}

start_app()
{
    echo "Starting app"
#    su-exec appuser $SRC_DIR/$APP_NAME
    gosu root $SRC_DIR/$APP_NAME
}

start_nginx()
{
    echo "Starting Nginx"
    /usr/sbin/nginx
}
initialize
start_webservers