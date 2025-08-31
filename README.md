# WB Tech School — Order Service
Проект разработан в рамках обучения в **WB Техношколе**.  
Цель — реализовать сервис для приёма заказов из Kafka, сохранения их в PostgreSQL и предоставления API для получения информации о заказах.  

## Задание
Текст задания вынесен в отдельный файл: [TASK.md](TASK.md)

## Запуск проекта
1. Клонировать репозиторий:
   ```bash
   git clone https://github.com/RozmiDan/wb_tech_testtask.git
   cd wb_tech_testtask
   ```
2. Создать файл .env в корне проекта (пример доступен в .env.example):
    ```bash
    cp .env.example .env
    ```
3. Запустить проект с помощью Docker Compose:
    ```bash
    docker compose up --build
    ```
4. После запуска, сервисы будут доступны по адресам:
    - API: http://localhost:8080
    - Swagger UI: http://localhost:8080/swagger/index.html

### Примеры запросов
- Добавление заказа
    ```bash
    curl -X POST http://localhost:8080/order/b563feb7b2b84b6test \
    -H "Content-Type: application/json" \
    -d @producer_samples/order1.json
    ```
- Получение заказа
    ```bash
    curl -X GET http://localhost:8080/order/b563feb7b2b84b6test
    ```
- Принудительный перезапуск сервиса (для тестирования восстановления кеша из БД)
    ```bash
    curl -X GET http://localhost:8080/drop/service
    ```

## Архитектура проекта
Проект построен на основе шаблона [go-clean-template](https://github.com/evrone/go-clean-template), что позволяет разделить ответственность между слоями приложения и упростить тестирование.

Основные слои:
- **Controller (Delivery)** — обработка HTTP-запросов, middleware, внешние запросы.
- **UseCase (Business logic)** — бизнес-правила, работа с кэшем и т.д.
- **Repository (Data access)** — работа с PostgreSQL.
- **Entity** — описание сущностей (Order, Delivery, Payment, Item и пр.).

---

## Структура проекта
```
├── build/                # Dockerfile и сборочные скрипты
│   └── producer-dir/     # Dockerfile для Kafka producer
├── cmd/                  # main пакеты
│   ├── main.go           # запуск основного сервиса
│   └── producer-dir/     # запуск утилиты-продьюсера (имитация потока заказов)
├── db/                   # миграции для PostgreSQL (goose)
│   ├── 001_create_orders_table.sql
│   ├── 002_create_deliveries_table.sql
│   ├── 003_create_payments_table.sql
│   └── 004_create_items_table.sql
├── internal/             # внутренняя бизнес-логика (чистая архитектура)
│   ├── app/              # точка входа, запуск приложения, инициализация зависимостей
│   │   └── app.go
│   ├── config/           # конфигурация (env → структура)
│   ├── controller/       # контроллеры, middleware, HTTP-сервер, Kafka consumer
│   ├── entity/           # сущности (Order, Delivery, Payment, Item и пр.)
│   ├── repo/             # слой репозитория (PostgreSQL)
│   └── usecase/          # бизнес-логика (валидация, кэш, сценарии)
├── pkg/                  # утилиты и переиспользуемые пакеты
│   ├── cache/            # потокобезопасный LRU-кэш + двусвязный список
│   ├── logger/           # инициализация zap-логгера
│   └── postgres/         # подключение к PostgreSQL
├── producer_samples/     # примеры JSON-заказов для тестирования producer'а
├── .env.example          # пример конфигурации (шаблон .env)
├── .gitignore
├── docker-compose.yaml   # docker-compose для запуска всех сервисов
├── Dockerfile            # Dockerfile для основного приложения
├── go.mod
└── README.md
```

## Что реализовано
- **Конфигурация**  
  Чтение переменных окружения из `.env` (директория `internal/config`).  
- **Миграции БД**  
  Автоматический запуск миграций через `goose` при старте сервиса (директория `db/`).  
- **Архитектура**  
  Реализация по принципам [Go Clean Template](https://github.com/evrone/go-clean-template).  
- **Swagger-документация**  
  Автогенерация документации для API.  
- **HTTP API**  
  - `POST /order/{order_uid}` — добавление заказа  
  - `GET /order/{order_uid}` — получение заказа (сначала из кэша, если нет — из БД)  
  - `GET /drop/service` — тестовая ручка для завершения приложения (используется для проверки перезапуска, работы кэша и Kafka)  
- **LRU-кэш**  
  Собственная потокобезопасная реализация на основе двусвязного списка и мапы (директория `pkg/cache`).  
- **Автовосстановление кеша при перезапуске сервиса**  
  При старте приложения в кэш загружается N последних заказов из БД (лимит задается в `.env`).
- **Kafka consumer**  
  Получение сообщений из топика `orders`, валидация, сохранение в PostgreSQL, добавление в кэш.  
- **Kafka producer**  
  Отдельный сервис для эмуляции потока заказов: читает JSON-файлы из каталога `producer_samples/` и публикует их в Kafka с задержками.  
- **Фронтенд**  
  Простая HTML/JS-страница для поиска заказа по `order_uid` и отображения информации (обращается к API).  

## Что не успел реализовать (в процессе)
- тесты (интеграционные, юнит)
- метрики (Prometheus + Grafana)