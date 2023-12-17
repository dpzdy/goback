FROM docker.io/library/alpine:latest
RUN apk add --no-cache sqlite
CMD ["sqlite3"]

#FROM ubuntu:trusty
#RUN sudo apt-get -y update
#RUN sudo apt-get -y upgrade
#RUN sudo apt-get install -y sqlite3 libsqlite3-dev
#RUN mkdir /db
#RUN /usr/bin/sqlite3 /db/dome.db
#CMD /bin/bash