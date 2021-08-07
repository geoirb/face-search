FROM golang:1.16-alpine AS build
WORKDIR /go/src
COPY . ./

ENV CGO_ENABLED=0
RUN go build \
    -o /out/face-search \
    ./cmd/face-search/main.go

FROM ubuntu:trusty
ENV LANG="C.UTF-8"

# install utilities
RUN apt-get update
RUN apt-get -y install wget --fix-missing
RUN apt-get -y install xvfb xorg xvfb firefox dbus-x11 xfonts-100dpi xfonts-75dpi xfonts-cyrillic --fix-missing # chrome will use this to run headlessly
RUN apt-get -y install unzip --fix-missing

# install go
RUN wget -O - 'https://storage.googleapis.com/golang/go1.7.linux-amd64.tar.gz' | tar xz -C /usr/local/
ENV PATH="$PATH:/usr/local/go/bin"

# install dbus - chromedriver needs this to talk to google-chrome
RUN apt-get -y install dbus --fix-missing
RUN apt-get -y install dbus-x11 --fix-missing
RUN ln -s /bin/dbus-daemon /usr/bin/dbus-daemon     # /etc/init.d/dbus has the wrong location
RUN ln -s /bin/dbus-uuidgen /usr/bin/dbus-uuidgen   # /etc/init.d/dbus has the wrong location

# install chrome
RUN wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add -
RUN sh -c 'echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google-chrome.list'
RUN apt-get update
RUN apt-get -y install google-chrome-stable

# install chromedriver
# NOTE: this is a relatively old version.  Try a newer version if this does not work.
RUN wget -N http://chromedriver.storage.googleapis.com/2.25/chromedriver_linux64.zip
RUN unzip chromedriver_linux64.zip
RUN chmod +x chromedriver
RUN mv -f chromedriver /usr/local/bin/chromedriver

ENV DISPLAY=:99
ENV XVFB_WHD=1280x720x16

RUN apt-get install -y ca-certificates tzdata
COPY --from=build /out/face-search /bin/face-search
COPY scripts/entrypoint.sh /
RUN chmod +x /entrypoint.sh
CMD ["./entrypoint.sh", "/bin/face-search"]