# Nio with docker

1. Build docker image

```sh
docker build -t dkrnio .
```

2. Run container

```sh
 docker run -it --rm -p 1323:1323 dkrnio:latest
 ```

3. Test running container

```sh
curl localhost:1323
```
