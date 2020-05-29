# VotingApp - a docker tutorial

## Start

```bash
git clone git@github.com:mehiX/VotingApp.git
cd VotingApp
mv .env.tmpl .env
# Fill in the missing values in .env

# Start the log collector
docker-compose up -d --build fluentd

# Start the databases first
docker-compose up -d --build mysql redis

# Start the rest of the containers
docker-compose up -d --build voting proxy results worker

# Check installation (Linux systems)
if [ 200 = $(curl -sI --url http://localhost:8080/ | grep HTTP | awk '{ print $2 }') ]; \
  then echo Installation successful!; \
  else echo Installation failed!; \
fi
```

## Use the application

In a web browser navigate to `http://localhost:8080`. On Windows systems replace `localhost` with the host IP.

## Votes generator

Build the generator image:

```bash
cd generator
docker build -t votingapp/generator:1.0 .
```

Start generating votes:

```bash
export CONTAINER_NAME=$(docker container ls --filter ancestor=voting/proxy:1.0 --format "{{.Names}}")

docker run --rm \
  --name voting-generator \
  --network container:${CONTAINER_NAME} \
  --log-driver fluentd \
  --log-opt tag=generator \
  votingapp/generator:1.0 \
  -url http://proxy/voting -workers 5
```

If the last command returns an error it may be that the proxy is not running. In that case the output of `echo ${CONTAINER_NAME}` is empty.

## NGINX

Check the default server configuration:

```bash
docker run \
  -ti --rm \
  nginx \
  cat /etc/nginx/conf.d/default.conf
```

## MYSQL

Start the service:

```bash
docker-compose up -d mysql
```

Connect to the running instance:

```bash
export CONTAINER_NAME=$(docker container ls --filter name=voting --filter ancestor=mysql --format "{{.Names}}")

docker run \
  -it --rm \
  --env-file .env \
  --network container:${CONTAINER_NAME} \
  mysql \
  sh -c 'mysql -h '${CONTAINER_NAME}' -u${MYSQL_USER} -p${MYSQL_PASSWORD}'
```

## REDIS

Connect to the Redis container:

```bash
export CONTAINER_NAME=$(docker container ls --filter ancestor=redis --filter name=voting --format "{{.Names}}")

docker run -it --rm \
  --network container:${CONTAINER_NAME} \
  redis \
  redis-cli -h ${CONTAINER_NAME}
```

## Docker topics

### Docker daemon

- Use a different logging driver. Can be setup in `daemon.json` or using `docker run --log-driver`. Official [documentation](https://docs.docker.com/config/containers/logging/configure/). For a working example see the fluentd service and the votes generating section.

### Docker run

- Use `--filter` to limit the results (see connecting to the mysql container). Mulptiple filters can be used to apply multiple criteria
- Use `--format` to help automation: show results in a way that can be further used in automation scripts (for an example see getting only the name of a container)
- Use `--rm` to run one-off containers and keep the system clean. Otherwise containers will still exist as stopped.
- Use `--name` to give the container a meaningful name. Otherwise a random name is given by Docker and that is not really helpful
- Use `--network` to connect the container to a desired network. This can be a network where another specific container is running.

### Dockerfile

- Review [best practices](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/). `redis/Dockerfile` shows installing additional packages based on the best practices.
- Multi-stage builds. All the GO projects show examples of multi-stage builds. The Docker plugin for VS Code generates a Dockerfile for GO projects that is not performant: any code change invalidates the docker build cache usage for installing `git`.
- When the command needed to run in a container becomes too big, it is a good practice to save it in a script and then call the script. Example in `redis/Dockerfile` where using `envsubst` makes the command too long.

## Docker-compose topics

### env_file vs. environment

`environment` gives precendence to host environment variables over the ones in `.env`. So if the same variable is declared in the environment and also in the `.env`, then the value in the file will be ignored. Always check the compiled docker-compose.yml with `docker-compose config`

The disadvantage of the `env_file` is that all the variables end up in all environments and that might not be desirable.
