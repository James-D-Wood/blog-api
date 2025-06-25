FROM golang:1.23.5-bullseye

WORKDIR /app
COPY . .
RUN make build
EXPOSE 8080
ENV NODB true

CMD [ "./app" ]