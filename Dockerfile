FROM golang:1.23.5-bullseye

WORKDIR /app
COPY . .
RUN make build
EXPOSE 8080
# remove env var when DB implementation is added
ENV NODB true

CMD [ "./app" ]