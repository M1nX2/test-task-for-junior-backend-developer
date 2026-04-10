# Task Service

Сервис для управления задачами с HTTP API на Go.

## Требования

- Go `1.23+`
- Docker и Docker Compose

## Быстрый запуск через Docker Compose

```bash
docker compose up --build
```

После запуска сервис будет доступен по адресу `http://localhost:8080`.

Если `postgres` уже запускался ранее со старой схемой, пересоздай volume:

```bash
docker compose down -v
docker compose up --build
```

Причина в том, что SQL-файл из `migrations/0001_create_tasks.up.sql` монтируется в `docker-entrypoint-initdb.d` и применяется только при инициализации пустого data volume.

## Swagger

Swagger UI:

```text
http://localhost:8080/swagger/
```

OpenAPI JSON:

```text
http://localhost:8080/swagger/openapi.json
```

## API

Базовый префикс API:

```text
/api/v1
```

Основные маршруты:

- `POST /api/v1/tasks`
- `GET /api/v1/tasks`
- `GET /api/v1/tasks/schedule?date=YYYY-MM-DD`
- `GET /api/v1/tasks/{id}`
- `PUT /api/v1/tasks/{id}`
- `DELETE /api/v1/tasks/{id}`

## Периодичность задач

Для `POST /api/v1/tasks` и `PUT /api/v1/tasks/{id}` добавлен объект `recurrence`:

- `type: "none"` — задача только на дату создания (по умолчанию).
- `type: "daily"` + `every_n_days` — каждый `n`-й день, начиная с даты создания.
- `type: "monthly"` + `day_of_month` — каждый месяц в указанное число (`1..30`).
- `type: "specific_dates"` + `dates` — только на перечисленные даты (`YYYY-MM-DD`).
- `type: "even_odd"` + `parity` — только на четные или нечетные дни месяца.

Пример:

```json
{
  "title": "Ежедневный обзвон",
  "description": "Утренний обход",
  "status": "new",
  "recurrence": {
    "type": "daily",
    "every_n_days": 2
  }
}
```

## Принятые допущения

- Сервис хранит задачу как шаблон с правилами периодичности.
- Фактический "план задач на дату" строится через `GET /api/v1/tasks/schedule?date=...`.
- Для `type=none` задача показывается только в день создания.
- Для `monthly` число ограничено диапазоном `1..30`, чтобы избежать невалидных дней месяца.
- Для существующих инсталляций после изменения схемы рекомендуется пересоздать volume: `docker compose down -v`.
