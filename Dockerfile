FROM golang:1.12
# docker build -t vanessa/containerflow .
RUN apt-get update && apt-get install -y nginx git python build-essential
WORKDIR /opt
RUN git clone https://github.com/emscripten-core/emsdk.git && \
    cd emsdk && \
    git pull && \
    ./emsdk install latest && \
    ./emsdk activate latest

ENV PATH /opt/emsdk:/opt/emsdk/fastcomp/emscripten:/opt/emsdk/node/12.9.1_64bit/bin:$PATH

WORKDIR /var/www/html
COPY . /var/www/html
RUN make && \
    mv html/* /var/www/html && \
    echo "application/wasm                                                wasm" >> /etc/mime.types
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
