FROM golang
ENV PORT 8080
EXPOSE 8080

WORKDIR /go/src/app
COPY . .

RUN go mod vendor
RUN CGO_ENABLED=0 go build -installsuffix 'static' -o app cmd/app/*
RUN mv ./app /go/bin/



FROM scratch

COPY --from=0 /go/bin/app /

CMD ["/app"]