#FROM golang
#
#RUN mkdir /go/src/demo
#COPY . /go/src/demo
#
#RUN export GO111MODULE="on"
#RUN export GOPROXY="https://mirrors.aliyun.com/goproxy/"
#
#RUN go get github.com/spf13/cast
#RUN go get github.com/gin-gonic/gin
#RUN go get github.com/jinzhu/gorm
#RUN cd /go/src/demo
#RUN go build -o greyRelease demo1

FROM centos:7

COPY grayReleaseProject /root/server
EXPOSE 8000
CMD /root/server