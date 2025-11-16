# Сервис назначения ревьюеров для Pull Request’ов

Микросервис, который автоматически назначает ревьюеров на Pull Request’ы (PR), а также позволяет управлять командами и участниками.

### Используемые технологии
- Go (Реализация бизнес-логики и всех операций сервиса)
- PostgreSQL (в качестве хранилища данных)
- Docker (для запуска сервиса)
- Gin (веб фреймворк)
- pgx (драйвер для работы с PostgreSQL)
- Logrus (для логирования)
- k6 (Нагрузочное тестирование)

Сервис написан с использованием подхода Clean Architecture, что позволяет легко расширять функционал сервиса и тестировать его. Также был реализован Graceful Shutdown для корректного завершения работы сервиса.

## Запуск проекта
Для запуска проекта необходимо заполнить .env файл, по примеру .env.example.

Запустить сервис можно командой make up.

## Примеры запросов
- [Добавление команды с участниками](#team_add)
- [Получение команды с участниками](#team_get)
- [Установка флага активности пользователя](#users_set_is_active)
- [Создание Pull Request](#create_pull_request)
- [Merge Pull Request](#merge_pull_request)
- [Переназначение ревьювера](#reassign_pull_request)
- [Получение PR'ов, где пользователь назначен ревьювером](#get_review)
- [Статистика](#statistics)


### Добавление команды с участниками <a name="team_add"></a>

Создать команду с участниками (создаёт/обновляет пользователей):
```curl
curl -X POST http://localhost:8080/team/add \
  -H "Content-Type: application/json" \
  -d '{
    "team_name": "backend-team",
    "members": [
      {
        "user_id": "user_001",
        "username": "Ivan_Ivanov",
        "is_active": true
      },
      {
        "user_id": "user_002", 
        "username": "Alex_Petrov",
        "is_active": true
      }
    ]
  }'
```
Пример ответа:
```json
{
    "team": {
        "team_name": "backend-team",
        "members": [
            {
                "user_id": "user_001",
                "username": "Ivan_Ivanov",
                "is_active": true
            },
            {
                "user_id": "user_002",
                "username": "Alex_Petrov",
                "is_active": true
            }
        ]
    }
}
```

### Получить команду с участниками <a name="team_get"></a>

Получение команды с участниками:
```curl
curl -X GET "http://localhost:8080/team/get?team_name=backend-team" \
  -H "Content-Type: application/json" 
```

Пример ответа:
```json
{
    "team_name": "backend-team",
    "members": [
        {
            "user_id": "user_001",
            "username": "Ivan_Ivanov",
            "is_active": true
        },
        {
            "user_id": "user_002",
            "username": "Alex_Petrov",
            "is_active": true
        }
    ]
}
```

### Установить флаг активности пользователя <a name="users_set_is_active"></a>

Устанавливает флаг активности пользователя (так как сервис не предусматривает авторизации и аутентификации метод выполняется только с токеном-заглушкой):
```curl
curl -X POST http://localhost:8080/users/setIsActive \
  -H "X-Admin-Token: adminToken" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user_001",
    "is_active": false
  }'
```
Пример ответа:
```json
{
    "user": {
        "user_id": "user_001",
        "username": "Ivan_Ivanov",
        "team_name": "backend-team",
        "is_active": false
    }
}
```

### Создать Pull Request <a name="create_pull_request"></a>

Создать PR и автоматически назначить до 2 ревьюверов из команды автора:
```curl
curl -X POST http://localhost:8080/pullRequest/create \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-001",
    "pull_request_name": "PR1",
    "author_id": "user_001"
}'
```
Пример ответа:
```json
{
    "pr": {
        "pull_request_id": "pr-001",
        "pull_request_name": "PR1",
        "author_id": "user_001",
        "status": "OPEN",
        "assigned_reviewers": [
            "user_002"
        ]
    }
}
```

### Merge Pull Request <a name="merge_pull_request"></a>

Помечает PR как MERGED (идемпотентная операция):
```curl
curl -X POST http://localhost:8080/pullRequest/merge \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-001"
    }'
```
Пример ответа:
```json
{
    "pr": {
        "pull_request_id": "pr-001",
        "pull_request_name": "PR1",
        "author_id": "user_001",
        "status": "MERGED",
        "assigned_reviewers": [
            "user_002"
        ],
        "mergedAt": "2025-11-16T18:20:59Z"
    }
}
```


### Переназначить конкретного ревьювера <a name="reassign_pull_request"></a>

Переназначает конкретного ревьювера на другого из его команды:
```curl
curl -X POST http://localhost:8080/pullRequest/reassign \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-005",
    "old_user_id": "user_002"
    }'
```
Пример ответа:
```json
{
    "pr": {
        "pull_request_id": "pr-005",
        "pull_request_name": "PR5",
        "author_id": "user_001",
        "status": "OPEN",
        "assigned_reviewers": [
            "user_003",
            "user_004"
        ]
    },
    "replaced_by": "user_004"
}
```

### Получение PR'ов, где пользователь назначен ревьювером <a name="get_review"></a>

Получить PR'ы, где пользователь назначен ревьювером:
```curl
curl -X GET "http://localhost:8080/users/getReview?user_id=user_004" \
  -H "Content-Type: application/json" 
```
Пример ответа:
```json
{
    "user_id": "user_004",
    "pull_requests": [
        {
            "pull_request_id": "pr-005",
            "pull_request_name": "PR5",
            "author_id": "user_001",
            "status": "OPEN"
        }
    ]
}
```

### Статистика <a name="statistics"></a>

Получение статистики:
```curl
curl -X GET "http://localhost:8080/stats/reviews" \
  -H "Content-Type: application/json" 
```
Пример ответа:
```json
{
    "user_stats": [
        {
            "user_id": "user_002",
            "username": "Alex_Petrov",
            "review_count": 5
        },
        {
            "user_id": "user_003",
            "username": "Bob_Johnson",
            "review_count": 1
        },
        {
            "user_id": "user_004",
            "username": "Victor_Kozlov",
            "review_count": 1
        },
        {
            "user_id": "user_001",
            "username": "Ivan_Ivanov",
            "review_count": 0
        }
    ],
    "total_prs": 5
}
```

# Решения

В ходе разработки было сомнение по одному из вопросов, которое было решено следующим образом:

Спецификация OpenAPI предусматривает 401 ошибку для эндпоинта `/users/setIsActive`, но в задании явно не требуется реализация полноценной авторизации.
>Решил реализовать токен-заглушку `adminToken` в заголовке `X-Admin-Token`, таким образом при отсутствии откена возвращается 401 код.