FROM golang:1.20.4-alpine as builder
LABEL maintainer = "Chenghao Yu <yuchenghao0624@qq.com>"

WORKDIR /usr/src/xxl_job_alert
COPY ./ /usr/src/xxl_job_alert/
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct
RUN go mod tidy && go build

FROM alpine:3
RUN mkdir -p /etc/xxl_job_alert
COPY --from=builder /usr/src/xxl_job_alert/xxl_job_alert /usr/bin/
COPY ./etc/xxl_job_alert/conf /etc/xxl_job_alert/conf
RUN chmod u+rw /etc/xxl_job_alert/conf/conf.toml
EXPOSE 30000
CMD ["sh", "-c", "xxl_job_alert"]