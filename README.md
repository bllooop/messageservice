# Микросервис обработки сообщений
Данный сервис реализован на языке Go с использованием библиотеки HTTP. Для работы с PostgreSQL использовался драйвер lib/pq. Для работы с файлом конфигураций библиотека Viper.
Сообщение принимается в JSON формате, далее сохраняется в базу данных PostgreSQL и отправляется в Kafka. Producer выводит в консоль строку о доставке сообщения. 
Consumer выводит полученное собщение в JSON формате, указывая Topic, текст сообщения и время отправки. Для пользователя выводится текст сообщения и id, с которым оно было сохранено в базу данных.
Посмотреть метрики и количество обработанных сообщений можно в UI запущенного Prometheus.
### Есть 3 способа запуска микросервиса:
1. Локально.

   Необходимо в файле конфигурации configs/config.yaml указать
   ```
   host = "localhost"
   ```
   Запустить контейнер postgres
   ```
   docker run --name=db -e POSTGRES_PASSWORD='54321' -p 5432:5432 -d postgres
   ```
   Контейнер kafka
   ```
   confluent local kafka start
   ```
   Контейнер prometheus
   ```
   docker run \
    -p 9090:9090 \
    -v ./prometheus.yml:/etc/prometheus/prometheus.yml \
    prom/prometheus
   ```
   Ввести в консоль команду
   ```
   go run cmd/main.go
   ```
   При необходимости, можно заменить названия контейнера базы данных, пароль и порты. Соответствующие параметры для такого же изменения находятся в configs/config.yaml.
3. Локально при использовании docker-compose.
   
   Необходимо в файле конфигурации configs/config.yaml указать (По умолчанию в проекте стоит такое значение)
   ```
   host = "db"
   ```
   Ввести в консоль команду
   ```
   docker-compose up --build
   ```
   При появлении ошибки о платформе использовать команду
   ```
   DOCKER_DEFAULT_PLATFORM=linux/amd64 docker-compose up --build
   ```

### Для сохранения сообщения необходимо сделать запрос:
```
curl --location 'http://{address}:8000/create' \
--header 'Content-Type: application/json' \
--data '{
    "text": {text}
}'
```
Вместо text ввести желаемое сообщение.
Вместо address необходимо ввести 
1. localhost при запуске 1 способом
2. 0.0.0.0 при запуске вторым способом
### При выполнении запроса в консоль выводятся сообщения от producer и consumer об отправке в Kafka:
```
messageservice-1  | Delivered message to messages[0]@0
messageservice-1  | Consumed event {
messageservice-1  |   "text": "new",
messageservice-1  |   "timestamp": "2024-07-26T16:46:51.959Z",
messageservice-1  |   "topic": "messages"
messageservice-1  | }
```
Для пользователя выводится 
```
{
    "id": 1,
    "text": "new"
}
```
## Prometheus
По адресу http://{address}:9090/graph находится график количества обработанных сообщений. Для отображения надо ввести в поле expression параметр kafka_proccessed_messages. 
По адресу http://{address}:8000/metrics находятся различные метрики для сервиса.
