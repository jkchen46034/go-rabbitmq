# Go-RabbitMQ

## Producer
```
~/dev/go-rabbitmq/producer$ go run main.go
2024/12/20 14:06:13 RabbitMQ producer running
2024/12/20 14:06:13 listening on 3030
2024/12/20 14:06:34 temperature is tored in RabbitMQ channel temperature, &{102 -25}
2024/12/20 14:07:16 temperature is tored in RabbitMQ channel temperature, &{103 -8}
```

## Consumer
```
~/dev/go-rabbitmq/consumer$ go run main.go
2024/12/20 14:06:34 Message: {"timestamp":102,"degree":-25}
2024/12/20 14:07:16 Message: {"timestamp":103,"degree":-8}
```


## HTTP POST
```
curl -X POST http://localhost:3030/temperature \
  -H "Content-Type: application/json" \
  -d '{"timestamp": 103, "degree": -8}'
```
