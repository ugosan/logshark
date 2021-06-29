FROM scratch

ADD logshark /logshark

CMD ["/logshark"]