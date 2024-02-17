## API сервис по созданию сокращенных ссылок(gRPC реализация)

- Принимает на вход URL и возвращает ссылку

#### Условия для ссылки:

- Уникальная: на один оригинальный URL одна сокращенная ссылка
- Длина 10 символов
- Из символов латинского алфавита в нижнем и верхнем регистре, цифр и символа \_ (подчеркивание)

#### Методы:

- Метод Post, сохраняет оригинальный URL в базе и возвращает сокращённый
- Метод Get, принимает сокращённую ссылку и возвращает оригинальный URL

#### Хранилище:

- PostgreSQL или inmemory (выбирается параметром при запуске сервиса)

#### Запуск:

- Для Unix-систем запуск производится вызовом скрипта с параметром
  ('./app.sh inmemory' или './app.sh database')
- Для Windows запуск для PostgreSQL команда docker-compose up
  для inmemory команда docker build -t inmemory . && docker run --rm --name inmemory -p 8080:8080 inmemory

## Использование:

#### GET

- curl -X GET http://localhost:8080/short-length-link/\<short link\>
  -H 'Accept: application/json'

#### POST

- curl -X POST http://localhost:8080/full-length-link/
  -H 'Content-Type: application/json'
  -d '{"longLink: \<URL\>"}'
