# Connect multiple ```kubectl port-forward svc```

### Usage
* Add/Remove Service list under main func
* Run
```
$ go run main.go
[redis] Forwarding from 127.0.0.1:6379 -> 6379
[mongo] Forwarding from 127.0.0.1:27017 -> 27017
[rabbitmq] Forwarding from 127.0.0.1:5672 -> 5672
[postgres] Forwarding from 127.0.0.1:5432 -> 5432
```
