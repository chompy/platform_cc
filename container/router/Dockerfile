FROM nginx:1.19-alpine

RUN apk add go git && \
    go get github.com/jsha/minica && \
    apk del go git
