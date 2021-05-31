FROM alpine

COPY bin/config.yaml bin/etserver bin/server.pem bin/server.key bin/client.pem bin/client.key /

EXPOSE 443

RUN chown 777 etserver

CMD ["./etserver"]