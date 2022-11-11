FROM progrium/busybox
RUN opkg-install bash

ENV TERM "xterm-256color"

ADD dist/linux-amd64/logshark /usr/local/bin/logshark

ENTRYPOINT ["logshark"]