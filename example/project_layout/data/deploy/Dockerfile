# 构建
FROM golang:1.18-alpine AS builder

RUN go env -w GOPROXY=https://goproxy.cn,https://goproxy.io,direct

ARG PROJECT_PATH=Project/project_layout

RUN mkdir -p /$PROJECT_PATH
WORKDIR /$PROJECT_PATH

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

WORKDIR /$PROJECT_PATH/cmd/project_layout

RUN CGO_ENABLED=0 go build -o /project_layout .

FROM centurylink/ca-certs

ARG PROJECT_PATH=Project/project_layout
WORKDIR /$PROJECT_PATH

COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip

ENV TZ Asia/Shanghai

COPY --from=builder /$PROJECT_PATH/data/conf/project_layout.toml /data/conf/project_layout.toml
COPY --from=builder /project_layout /usr/local/bin/project_layout
