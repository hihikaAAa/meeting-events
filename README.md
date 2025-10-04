# Meeting Events Service

Небольшой сервис для управления встречами и их планирования.  
Разработан с использованием **Go**, **Postgres**, **Docker**, методологий **DDD** и **TDD**.  

Сервис реализует CRUD-функционал для сущности `Meeting` и поднимается в контейнерах через `docker-compose`.

---

## Архитектура проекта

Проект разделён по слоям по подходу **DDD**:

- **domain/** — бизнес-сущности, инварианты и правила.
- **internal/app/usecase/** — прикладные сценарии (Create/Update/Cancel/Get).
- **internal/adapters/** — работа с инфраструктурой (Postgres, миграции, репозитории).
- **internal/httpserver/** — слой HTTP-обработчиков и роутинг.
- **config/** — конфигурационные файлы (`prod.yaml`).
- **tests/** — интеграционные и e2e тесты.

---

## Сущность Meeting

`Meeting` содержит поля:

- `id` (UUID)
- `title` — название встречи
- `starts_at` — время начала
- `duration` — продолжительность (в минутах)
- `status` — текущий статус (`planned`, `canceled`, …)
- `created_at`, `updated_at`
- `events` - события, примененные к конкретному meeting

---

## Требования

- [Go 1.24+]
- [Docker Desktop]
- `docker-compose`

---

## Быстрый старт

1. Склонируйте репозиторий:

   ```bash
   git clone https://github.com/<your-repo>/meeting-events.git
   cd meeting-events
   ```

2. Поднимите сервис:
  
   ```bash
   docker compose up -d --build
   ```

3. Проверьте, что все поднялось:

   ```bash
   docker ps
   ```

4. API доступно по адресу:

   ```bash
   http://localhost:8081
   ```
   

## Аутентификация
 
   ```bash
   username: user
   password: pass
   ```

## Примеры запросов к API(curl)

1. Создание встречи(Create)
   ```bash 
   curl -X POST http://localhost:8081/v1/meetings/ \
  -u user:pass \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Daily sync",
    "starts_at": "2030-01-01T09:00:00Z",
    "duration": 45
  }'
  ```

Ответ: 
    ```json
    {"id": "f6e4f9a9-3c58-4f83-83b8-b9ad2bcd9c24"}
    ```

2. Получение встречи по id(Get)

   ```bash
   curl -X GET http://localhost:8081/v1/meetings/<id> \
  -u user:pass
  ```

3. Обновление информации о встрече(Update)

   ```bash
   curl -X PATCH http://localhost:8081/v1/meetings/<id> \
  -u user:pass \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated meeting",
    "starts_at": "2030-01-01T10:00:00Z",
    "duration": 60
  }'
  ```

4. Удаление встречи (Delete)

   ```bash
   curl -X DELETE http://localhost:8081/v1/meetings/<id> \
  -u user:pass
  ```


## Тестирование
Unit-тесты 
На каждый use-case, handler,dto написаны unit-тесты. Они распологаются рядом с кодом их логики.

Функциональные/интеграционные тесты
Находятся в tests/

## Технологии

- Go — основной язык
- PostgreSQL — хранилище данных
- Docker / Docker Compose — контейнеризация
- testcontainers-go — интеграционные тесты
- slog — логирование
- chi — маршрутизация HTTP


## Разработчик

Автор: Апанасевич Сергей
Email: heheka800@gmail.com



