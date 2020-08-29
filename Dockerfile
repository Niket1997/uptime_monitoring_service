FROM golang:1.15.0 AS build
RUN mkdir /app
ENV SRC_DIR=/app
ENV APP_NAME=UMS
ARG GIT_COMMIT_HASH
ENV GIT_COMMIT_HASH=${GIT_COMMIT_HASH}
COPY . /app
WORKDIR /app
RUN go get -d
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s -X main.AppVersion=$GIT_COMMIT_HASH" -o $APP_NAME .


FROM alpine:latest AS runtime
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2


FROM nginx
ARG GIT_COMMIT_HASH
ENV GIT_COMMIT_HASH=${GIT_COMMIT_HASH}
RUN echo "$GIT_COMMIT_HASH"
RUN mkdir -p /app/dockerconf  /app/public && \
    echo "$GIT_COMMIT_HASH" > /app/public/commit.txt

RUN set -eux \
    && apt-get update \
    && apt-get install -y ca-certificates \
    && apt-get install -y gosu

ENV SRC_DIR=/app
ENV APP_NAME=UMS

WORKDIR /app
COPY  --from=build /app/$APP_NAME /app/
COPY --from=build /app/dockerconf/entrypoint.sh /app/dockerconf/
COPY --from=build --chown=nginx:nginx /app/dockerconf/nginx.conf /etc/nginx/conf.d/default.conf
RUN chmod +x /app/dockerconf/entrypoint.sh
EXPOSE 80
ENTRYPOINT ["/app/dockerconf/entrypoint.sh"]