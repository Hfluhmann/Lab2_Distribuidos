# Pozo
En esta máquina se ejecuta el pozo

## Ejecucion
```bash
make proto_pozo
make pozo
```

## Consideraciones
1. Debe estar corriendo RabbitMQ
```bash
service rabbitmq-server start
```

2. Debe estar autorizado el usuario cliente
```bash
sudo rabbitmqctl add_user 'client' '1234'
sudo rabbitmqctl set_user_tags 'client' administrator
sudo rabbitmqctl set_permissions -p / "client" ".*" ".*" ".*"
```
3. Debe existir el archivo `pozo.txt` y para iniciar el juego solo debe contener la linea:
```
Jugador_0 Ronda_0 0
```

# Data Node
En esta máquina se debe ejecutar un data node.

## Ejecución
```bash
make proto_data
make data
```

## Consideraciones
Los archivos de las jugadas se guardarán en la carpeta data.
