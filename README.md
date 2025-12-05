# URL Shortener

## О проекте

**url-shortener** — простой REST API сервис для сокращения длинных URL. Он позволяет:

- Создавать "короткие" ссылки из длинных.
- При переходе по короткой ссылке — перенаправлять пользователя на оригинальный URL.

## Технологии и библиотеки

- **Язык:** Go
- **Маршрутизация:** chi
- **Middleware:** RequestID, Logger, Recoverer, URLFormat, BasicAuth
- **Логирование:** slog
- **База данных:** SQLite
- **Конфигурация:** модуль config

## Основные возможности

- REST API для создания и получения коротких URL
- Редирект по короткой ссылке на оригинальный URL
- Поддержка авторизации через BasicAuth для защищённых маршрутов

## Эндпоинты / API Endpoints

### 1. Создание короткой ссылки

**POST** `/url`

**Описание:** Создаёт короткую ссылку из длинного URL.

**Запрос (JSON):**

```json
{
  "url": "https://example.com/very/long/url"
}
```

**Ответ (JSON):**

```json
{
  "short_url": "http://localhost:8080/abcd123"
}
```

### 2. Редирект по короткой ссылке

**GET** `/{alias}`

**Описание:** При переходе по короткой ссылке перенаправляет пользователя на исходный URL.

**Пример запроса:**

```
GET /abcd123
```

**Поведение:** Редирект на https://example.com/very/long/url

### 3. Защищённый эндпоинт создания ссылки с BasicAuth

**POST** `/url/url`

**Описание:** Создаёт короткую ссылку с использованием BasicAuth для дополнительной защиты.

**Пример запроса с BasicAuth:**

```
POST /url/url
Authorization: Basic base64(username:password)
Content-Type: application/json

{
  "url": "https://example.com/another/long/url"
}
```

**Ответ (JSON):**

```json
{
  "short_url": "http://localhost:8080/xyz789"
}
```
