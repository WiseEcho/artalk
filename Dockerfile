FROM d-x.cmstop.net/base:3.13

WORKDIR /app

COPY ./bin/artalk /app/artalk
COPY ./artalk.yml /app/artalk.yml

#CMD ["server", "--host", "0.0.0.0", "--port", "23366"]
#ENTRYPOINT ["/app/artalk"]
