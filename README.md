# Веб-сервис Калькулятора на Go

## Описание проекта

Этот проект представляет собой веб-сервис арифметического калькулятора, разработанный на языке программирования Go. Сервис позволяет пользователям отправлять арифметические выражения через HTTP-запросы, а в ответ получать вычисленные результаты.

## Быстрый старт

### Предварительные требования

- **Go:** Убедитесь, что у вас установлен Go версии 1.23 или выше.

### Клонирование репозитория

```bash
git clone https://github.com/mpkelevra23/arithmetic-web-service.git
cd arithmetic-web-service
```

### Установка зависимостей

Проект использует модули Go для управления зависимостями. Выполните следующую команду для установки необходимых пакетов:

```bash
go mod tidy
```

### Настройка конфигурации

Создайте файл `.env` в корне проекта и добавьте в него следующие переменные:

```env
PORT=8080
LOG_LEVEL=info
```

- **PORT:** Порт, на котором будет запущен сервер (по умолчанию `8080`).
- **LOG_LEVEL:** Уровень логирования (`debug`, `info`, `warn`, `error`).

### Запуск сервера

Вы можете запустить сервер с помощью следующей команды:

```bash
go run ./cmd/server/main.go
```

Сервер запустится на указанном в `.env` порту (по умолчанию `8080`).

## Использование

### Пример успешного запроса

**Запрос:**

```bash
curl -i --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*2"
}'
```

**Ответ:**

```json
{
    "result": "6"
}
```

### Примеры ошибок

#### 1. Недопустимый ввод (422)

Возникает при наличии недопустимых символов во вводе.

**Запрос:**

```bash
curl -i --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2a"
}'
```

**Ответ:**

```json
{
    "error": "Expression is not valid"
}
```

#### 2. Деление на ноль (422)

Попытка деления на ноль.

**Запрос:**

```bash
curl -i --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "10/0"
}'
```

**Ответ:**

```json
{
    "error": "Division by zero"
}
```

#### 3. Отсутствие поля `expression` (422)

**Запрос:**

```bash
curl -i --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expr": "2+2"
}'
```

**Ответ:**

```json
{
    "error": "Missing field: expression"
}
```

#### 4. Искусственно вызванная ошибка 500

Добавлена возможность принудительного вызова ошибки 500 через заголовок `X-Trigger-500` или определенные условия:

**Запрос с заголовком:**

```bash
curl -i --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--header 'X-Trigger-500: true' \
--data '{
  "expression": "1+1"
}'
```

**Ответ:**

```json
{
    "error": "Internal Server Error"
}
```

**Запрос с длинным выражением:**

```bash
curl -i --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "1+1+...(очень длинная строка)"
}'
```

Если длина выражения находится в диапазоне от 500 до 1000 символов, сервер возвращает:

**Ответ:**

```json
{
    "error": "Expression length triggered server error"
}
```

## Тестирование

### Запуск тестов

Проект содержит модульные и интеграционные тесты, расположенные в директории `tests/`. Для запуска всех тестов используйте следующую команду:

```bash
go test ./tests/...
```

### Описание тестов

- **calculator_test.go:** Тестирует функцию `Calc` для различных арифметических выражений и проверяет корректность вычислений.
- **handler_test.go:** Интеграционные тесты для обработчика HTTP-запросов, проверяющие правильность обработки различных сценариев запросов.
- **middleware_test.go:** Тестирует middleware для логирования запросов, убеждаясь в корректности логирования и передачи запросов.

## Логирование

Проект использует библиотеку Zap для структурированного логирования. Логи записываются как в консоль, так и в файлы:

- **Основные логи:** `logs/app.log`
- **Ошибки:** `logs/error.log`

Уровень логирования настраивается через переменную окружения `LOG_LEVEL` в файле `.env`.

## Примеры использования с дополнительными сценариями

### Пример с использованием скобок и приоритетов операций

**Запрос:**

```bash
curl -i --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "(2+3)*4 - 5/5"
}'
```

**Ответ:**

```json
{
    "result": "19"
}
```

### Пример с отрицательными числами

**Запрос:**

```bash
curl -i --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "-5 + 3"
}'
```

**Ответ:**

```json
{
    "result": "-2"
}
```

### Пример с десятичными числами

**Запрос:**

```bash
curl -i --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "3.5 * 2"
}'
```

**Ответ:**

```json
{
    "result": "7"
}
```
