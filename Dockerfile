FROM busybox:1.33

ENV TERM "xterm-256color"

ADD logshark /usr/local/bin/logshark


ENTRYPOINT ["logshark"]

