# things you need to install before running this app:
- golang 1.22+
- mysql
- redis
- erlang
- rabbitmq
- rabbitmq delay exchange(depands on your rabbit mq version)
# and prepare for ./config:
- rabbitmq_config.yaml
    - dialstr
- redis_config.yaml
    - addr
    - password
    - db
- web_server_config.yaml
    - port
    - jwt_key
- mysql_config.yaml
    - driver_name
    - port
    - host
    - database
    - username
    - password
    - charset
# and ...
- go mod tidy
- go run .
