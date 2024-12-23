# Оглавление
- [Web Service for Arithmetic Expression Calculation](#web-service-for-arithmetic-expression-calculation)
- [Структура проекта](#структура-проекта)
- [Требования](#требования)
- [Установка и запуск проекта](#установка-и-запуск-проекта)
- [Инструкция по отправке POST-запроса на endpoint с использованием curl](#инструкция-по-отправке-post-запроса-на-endpoint-с-использованием-curl)
- [Обработка ошибок](#обработка-ошибок)
- [Инструкция по отправке POST-запроса на endpoint с использованием Postman](#инструкция-по-отправке-post-запроса-на-endpoint-с-использованием-postman)
- [Тестирование](#тестирование)

# Web Service for Arithmetic Expression Calculation

Этот проект представляет собой веб-сервис, который принимает арифметические выражения через HTTP-запросы, выполняет их вычисление и возвращает результат. Проект написан на языке Go.

## Структура проекта
```markdown
[GO_Yandex](https://github.com/syirnik/GO_Yandex)
├── README.md                  # Описание проекта, инструкция по запуску и примеры использования
├── cmd                         # Папка с основной логикой запуска сервиса
│   └── main.go                 # Главный файл, запускающий HTTP-сервис
├── go.mod                      # Модуль Go, описывает зависимости проекта
├── internal                    # Папка для внутренней логики приложения (не для публичного использования)
│   └── application             # Логика приложения, работающая с запросами и бизнес-логикой
│       ├── application.go      # Реализация приложения: обработка запросов, вычисления и ошибки
│       └── application_test.go # Тесты для логики приложения
└── pkg                         # Публичная часть приложения, содержащая основные вычисления
    └── calculation             # Логика вычислений
        ├── calc.go             # Основной файл с реализацией логики вычисления выражений
        ├── calc_test.go        # Тесты для calc.go
        └── errors.go           # Обработка ошибок 
```

## Требования

Перед тем как запустить проект, убедитесь, что на вашем компьютере установлены:

- [Go 1.18+](https://golang.org/dl/) (для разработки и запуска проекта)
- [Git](https://git-scm.com/) (для клонирования репозитория)
- [curl](https://curl.se/) (для тестирования API запросов)


## Установка и запуск проекта

### 1. Клонируйте репозиторий

Клонируйте репозиторий на свой компьютер с GitHub с помощью команды:

```bash
git clone https://github.com/syirnik/GO_Yandex
```
### 2. Перейдите в каталог проекта

```bash
cd GO_Yandex
```
### 3. Установите все необходимые зависимости
Если вы используете Go Modules, выполните следующую команду для установки зависимостей:

```bash
go mod tidy
```
### 4. Перейдите в папку cmd
Для запуска сервера необходимо перейти в папку cmd:

```bash
cd cmd
```
### 5. Запуск проекта
Выполните команду из директории cmd:
```bash
go run main.go
```
### 6. Проверка работы сервера
После этого сервер будет слушать на порту 8080. В консоли появится сообщение:

```console
Server is listening on port 8080...
```

После того как сервер успешно запущен, он готов к приему запросов.

Для вычисления арифметических выражений с помощью сервиса нужно открыть **новое окно терминала** и выполнить команду **curl** или использовать **Postman** для отправки запросов к API.
Чтобы проверить работу сервиса, отправьте POST-запрос на следующий endpoint: http://localhost:8080/api/v1/calculate

Этот endpoint принимает POST-запрос с арифметическим выражением в формате JSON и возвращает результат вычисления.


## Инструкция по отправке POST-запроса на endpoint с использованием `curl`

Для того чтобы отправить POST-запрос на сервер, используйте команду `curl`. В зависимости от операционной системы синтаксис запроса будет различаться.

## 1. Для пользователей **Windows**

### 1.1 **Command Prompt (CMD)**

В **CMD** требуется экранировать кавычки с помощью символа `\"`.

#### Пример запроса:
```bash
curl -s -X POST http://localhost:8080/api/v1/calculate -d "{\"expression\": \"3 + 5 * (2 - 8)\"}"
```

### 1.2 **PowerShell**, **Git Bash** и **Windows Subsystem for Linux (WSL)**

В этих оболочках можно использовать стандартный синтаксис с **двойными кавычками** для строк и **одиночными** для JSON-формата. Экранировать кавычки не требуется.

#### Пример запроса:

```bash
curl -s -X POST http://localhost:8080/api/v1/calculate -d '{"expression": "3 + 5 * (2 - 8)"}'
```
## 2. Для пользователей **Linux**
Откройте терминал (Ctrl + Alt + T) и введите команду:

```bash
curl -s -X POST http://localhost:8080/api/v1/calculate -d '{"expression": "3 + 5 * (2 - 8)"}'
```
## 3. Для пользователей **macOS**
Откройте терминал (Cmd + Space → "Terminal") и введите команду:
```bash
curl -s -X POST http://localhost:8080/api/v1/calculate -d '{"expression": "3 + 5 * (2 - 8)"}'
```
После выполнения команды вы получите результат вычислений выражения в формате JSON.

Пример ответа:

```json
{"result":"-27.000000"}
```
Чтобы увидеть коды ответа сервера (например, 200, 422, 500), используйте ключи -i или -v с командой curl.
## Пример успешного запроса (код 200)

```bash
curl -i -s -X POST http://localhost:8080/api/v1/calculate -d "{\"expression\": \"(7 + 3 * (2 + 5)) / (4 - 2)\"}"
```
Ответ:
```json
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sun, 22 Dec 2024 12:14:42 GMT
Content-Length: 23

{"result":"14.000000"}
```
## Обработка ошибок

Сервис обрабатывает различные виды ошибок и возвращает соответствующие HTTP-коды:
- **400 Bad Request** – Ошибка в запросе (например, неправильный формат JSON).
- **422 Unprocessable Entity** – Некорректное арифметическое выражение.
- **500 Internal Server Error** – Ошибка на сервере.

### Пример ошибки 400 (Bad Request):

```bash
curl -i -s -X POST http://localhost:8080/api/v1/calculate -d "{\"expression\": \"\"}"
```
Ответ:
```json
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Sun, 22 Dec 2024 12:10:04 GMT
Content-Length: 40

{"error": "Expression cannot be empty"}
```
### Пример ошибки 422 (Unprocessable Entity)

```bash
curl -i -s -X POST http://localhost:8080/api/v1/calculate -d "{\"expression\": \"(3 + 5 * (2 - 8)\"}"
```
Ответ:
```json
HTTP/1.1 422 Unprocessable Entity
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Sun, 22 Dec 2024 12:15:53 GMT
Content-Length: 65

{"error": "Mismatched parentheses. Please check your brackets."}
```
## Инструкция по отправке POST-запроса на endpoint с использованием **Postman**

Если вы предпочитаете работать с графическим интерфейсом, то для отправки запросов вместо curl можно использовать **Postman**.

Откройте приложение **Postman**  и выполните следующие шаги:
- Создайте новый запрос
- Настройте тип запроса **POST**
- В поле URL введите адрес сервера:
http://localhost:8080/api/v1/calculate
- Добавьте тела запроса: перейдите на вкладку Body, выберите опцию raw, выберите формат JSON,
введите тело запроса в формате JSON
, например:
```json
{
  "expression": "3 + 5 * (2 - 8)"
}
```
- Нажмите кнопку Send, чтобы отправить запрос на сервер.

После того как вы отправите запрос, Postman покажет ответ сервера внизу экрана:

**Status:** Показывает статусный код HTTP (200 для успешного запроса или 400/422 для ошибок).

**Response Body:** Здесь будет отображаться результат выполнения арифметического выражения или ошибка.
## Тестирование

Для написания тестов используется стандартная библиотека Go testing, которая предоставляет базовые функции для написания, выполнения и проверки результатов тестов.
Для проверки правильности работы кода можно использовать стандартный инструмент Go для тестирования — `go test`.

## Запуск тестов

### 1. Перейдите в корневую папку проекта
```bash
cd путь/к/проекту
```
### 2. Запуск тестов
Чтобы запустить тесты, выполните следующую команду:

```bash
go test ./...
```

Эта команда выполнит все тесты, находящиеся в проекте, включая те, которые расположены в папках internal и pkg.

Если вам нужно запустить тесты только для конкретного пакета, например, для пакета с логикой вычислений, используйте команду:

```bash
go test ./pkg/calculation
```
Эта команда выполнит тесты только для пакета pkg/calculation.

## Что проверяют тесты
Тесты в проекте проверяют правильность работы функций и алгоритмов, таких как:
- Корректность вычислений арифметических выражений.
- Обработку ошибок, например, деление на ноль или некорректные выражения.
- Соответствие возвращаемых значений ожидаемым результатам для различных тестовых случаев.

