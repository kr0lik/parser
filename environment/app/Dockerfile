FROM golang:alpine

RUN apk --update --no-cache add  \
    bash  git  curl wget ca-certificates unzip libexif xvfb udev

RUN apk --update --no-cache add \
    chromium=117.0.5938.62-r0 chromium-chromedriver=117.0.5938.62-r0

RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /app
COPY ../../../app  .

RUN go mod download && \
    CGO_ENABLED=0 go build -o bin/ozon -gcflags "all=-N -l" ./cmd/ozon/main.go && \
    CGO_ENABLED=0 go build -o bin/wildberries  -gcflags "all=-N -l" ./cmd/wildberries/main.go


#RUN apk --update --no-cache add openjdk11
#RUN wget https://github.com/SeleniumHQ/selenium/releases/download/selenium-4.11.0/selenium-server-4.11.0.jar -P /src/

#CMD bash -c "chromedriver & java -jar /src/selenium-server-4.11.0.jar standalone"
