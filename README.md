# VotingApp - a docker tutorial

## Start

```bash
git clone git@github.com:mehiX/VotingApp.git
cd VotingApp
mv .env.tmpl .env
# Fill in the missing values in .env

# Start the databases first
docker-compose up -d --build mysql redis

# Start the rest of the containers
docker-compose up -d --build

# Check installation (Linux systems)
if [ 200 = $(curl -sI --url http://localhost:8080/ | grep HTTP | awk '{ print $2 }') ]; \
  then echo Installation successful!; \
  else echo Installation failed!; \
fi
```

## Use the application

In a web browser navigate to `http://localhost:8080`. On Windows systems replace `localhost` with the host IP.

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
docker run \
  -it --rm \
  --env-file .env \
  --network container:votingapp_mysql_1 \
  mysql \
  sh -c 'mysql -h votingapp_mysql_1 -u${MYSQL_USER} -p${MYSQL_PASSWORD}'
```

## REDIS

Connect to the Redis container:

```bash
docker run -it --rm \
  --network container:votingapp_redis_1 \
  redis \
  redis-cli -h votingapp_redis_1
```

## Docker topics


## Docker-compose topics

### env_file vs. environment

`environment` gives precendence to host environment variables over the ones in `.env`. So if the same variable is declared in the environment and also in the `.env`, then the value in the file will be ignored. Always check the compiled docker-compose.yml with `docker-compose config`

The disadvantage of the `env_file` is that all the variables end up in all environments and that might not be desirable.
