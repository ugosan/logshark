FROM progrium/busybox
RUN opkg-install bash

ENV TERM "xterm-256color"

ADD dist/logshark-linux-amd64 /usr/local/bin/logshark

ENTRYPOINT ["logshark"]

