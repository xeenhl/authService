1) to build docker image:
    docker build -t <image_name> .
    example: docker build -t my-golang-app .

2) to run docker container:
    docker run -i -t -d -p <host_port>:<container_port> <image_name>
    example:  docker run -i -t -d -p 8081:8081 my-golang-app
