---
tags: ["dockerfile"]
---
# ubuntu
---
## Example

# dockerfile
```dockerfile
FROM ubuntu:22.04

RUN apt-get update \
    && apt-get install -y net-tools \
    && apt-get install -y iputils-ping \
    && apt-get install -y vim wget curl

CMD [ "/bin/bash" ]
```
# iamge
```bash
docker build -t myubuntu .
docker run -it --rm myubuntu /bin/bash
```
