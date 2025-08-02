# Events Service

RESTful API-сервис для управления событиями на Go с PostgreSQL

## Описание

Сервис предоставляет API для:
- Создания и управления пользователями
- Работы с событиями (CRUD операции)
- Получения событий за разные периоды (день, неделя, месяц)

## Особенности

- Полный CRUD для событий
- Гибкая система выборки событий по периодам
- Конфигурация через YAML/переменные окружения
- Роутинг через `go-chi`
- Логирование с помощью `slog`
- Поддержка PostgreSQL

## Установка

### Требования
- Go 1.21+
- PostgreSQL 12+

### Запуск
```bash
git clone https://github.com/your-username/events-service.git
cd events-service

# Настройка конфига
cp config/config.yaml.example config/local.yaml
# Редактируем local.yaml

# Запуск
go run ./cmd/events-service
```

Сервер будет доступен на `http://localhost:8036`

## Структура базы данных

### Таблица пользователей
```sql
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY
);
```

### Таблица событий
```sql
CREATE TABLE event (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    date DATE NOT NULL,
    text TEXT NOT NULL DEFAULT ''
);
```

## API Endpoints

| Метод | Путь               | Описание                              |
|-------|--------------------|---------------------------------------|
| POST  | /create_user       | Создание пользователя                 |
| POST  | /create_event      | Создание события                      |
| POST  | /update_event      | Обновление события                    |
| POST  | /delete_event      | Удаление события                      |
| GET   | /events_for_day    | События за день (YYYY-MM-DD)          |
| GET   | /events_for_week   | События за неделю (от переданной даты)|
| GET   | /events_for_month  | События за месяц (YYYY-MM-DD)         |

## Конфигурация

Основной файл конфигурации `config/local.yaml`:

```yaml
env: "local"
storage_path: "host=${DB_HOST} port=${DB_PORT} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=disable"
```

Или через переменные окружения:
```bash
export DB_HOST=localhost
export DB_PORT=5432
# и т.д.
```

## Тестирование

Юнит-тесты:
```bash
go test ./...
```

Интеграционные тесты (требует запущенную БД):
```bash
go test -tags integration ./...
```

## Структура проекта

```
├── cmd/              # Основной пакет приложения
├── config/           # Конфигурационные файлы
├── internal/         # Внутренние пакеты
│   ├── config/       # Парсинг конфига
│   ├── http-server/  # HTTP-handlers и middleware-логгер
│   ├── lib/          # api и loggers
│   ├── models/       # Модели данных
│   └── storage/      # Работа с БД
└── tests/            # Интеграционные тесты
```

## Примеры запросов

Создание пользователя:
```bash
curl -X POST http://localhost:8080/create_user \
  -H "Content-Type: application/json" \
  -d '{}'
```

Получение событий за день:
```bash
curl -X GET "http://localhost:8080/events_for_day?user_id=1&date=2025-01-01"
```
