FROM alpine:3.18.4

RUN apk add --no-cache bash
RUN apk add --no-cache tzdata
RUN apk add --no-cache dumb-init

WORKDIR /opt/app

# Copy files to docker image
COPY configs/.env configs/.env
COPY backend .

ENTRYPOINT [ "/usr/bin/dumb-init", "--" ]

CMD [ "./backend" ]