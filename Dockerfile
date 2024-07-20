FROM alpine:3.20.1

ENV LD_LIBRARY_PATH="."

WORKDIR /app

COPY . .

RUN ./build.sh

CMD ["/app/libcsv"]
