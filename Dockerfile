ARG VERSION=0.0.0
FROM --platform=$BUILDPLATFORM node:lts-alpine as uibuilder
WORKDIR /src
COPY ui .
RUN yarn config set registry https://registry.npm.taobao.org
RUN yarn config set network-timeout 600000
RUN npm config set registry https://registry.npm.taobao.org
RUN npm install
RUN yarn && yarn build 

FROM golang:1-alpine as gobuilder
ARG VERSION
WORKDIR /src
COPY . .
COPY --from=uibuilder /src/build ./ui/build
RUN apk add git
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go generate ./... && CGO_ENABLED=0 go build -ldflags "-s -w -X main.version=${VERSION}" -o rmfakecloud-docker ./cmd/rmfakecloud/

FROM scratch
EXPOSE 3000
RUN --mount=from=busybox:latest,src=/bin/,dst=/bin/ mkdir -m 1755 /tmp
COPY --from=gobuilder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=gobuilder /src/rmfakecloud-docker /
ENTRYPOINT ["/rmfakecloud-docker"]
