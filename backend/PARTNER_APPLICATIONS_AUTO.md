# Автоматическая обработка заявок на партнерство

Этот документ описывает систему автоматической обработки просроченных заявок на партнерство.

## Как это работает

### Логика обработки

1. **Заявка живет 4 часа** - после создания заявки у партнера есть 4 часа на ее рассмотрение
2. **Автоматическое принятие** - если заявка не обработана в течение 4 часов, система автоматически пытается ее принять
3. **Проверка свободных слотов** - система проверяет, есть ли у партнера свободные слоты для новых партнеров
4. **Эскалация по спонсорской линии** - если у текущего партнера нет свободных слотов, заявка эскалируется вверх по спонсорской линии
5. **Автоматическое принятие спонсором** - первый спонсор в линии, у которого есть свободный слот, автоматически принимает заявку

### Компоненты системы

#### 1. Логика автоматической обработки
- **Файл**: `internal/service/team/partner_application_auto.go`
- **Основная функция**: `ProcessExpiredApplications()` - обрабатывает все просроченные заявки
- **Эскалация**: `escalateApplicationUpSponsorLine()` - поднимает заявку по спонсорской линии

#### 2. База данных
- **Файл**: `internal/repository/mongodb/partner_application.go`
- **Новый метод**: `PartnerApplicationGetExpired()` - получает заявки старше 4 часов со статусом "pending"

#### 3. Cron Job / Background Service
- **Файл**: `cmd/partner_applications_processor/main.go`
- **Частота**: каждые 30 минут
- **Функция**: запускает обработку просроченных заявок

## Запуск системы

### Локальная разработка

```bash
# Собрать processor
make build-partner-processor

# Запустить processor
make run-partner-processor
```

### Docker

```bash
# Запустить всю систему включая processor
docker-compose up -d

# Запустить только processor
docker-compose up -d partner-applications-processor

# Посмотреть логи processor
docker-compose logs -f partner-applications-processor
```

### Только processor сервис

```bash
# Собрать образ только для processor
docker build --build-arg BUILD_TARGET=partner-processor -t partner-processor .

# Запустить processor
docker run -d \
  --name partner-processor \
  --env MONGODB_URL="mongodb://user:pass@mongo:27017/aibetrade?authSource=admin&replicaSet=rs0" \
  --env MONGODB_DATABASE=aibetrade \
  --network aibetrade_default \
  partner-processor
```

## Конфигурация

### Переменные окружения

- `MONGODB_URL` - URL подключения к MongoDB
- `MONGODB_DATABASE` - название базы данных
- `LOG_LEVEL` - уровень логирования (debug, info, warn, error)

### Настройки в коде

- **Время жизни заявки**: 4 часа (константа в `partner_application_auto.go`)
- **Частота проверки**: 30 минут (настройка в `main.go`)
- **Максимальное количество уровней эскалации**: 10 (константа в `escalateApplicationUpSponsorLine()`)

## Мониторинг и логи

Система пишет подробные логи обо всех операциях:

```
Starting processing of expired partner applications...
Found 2 expired applications to process
Processing expired application abc123: user456 -> partner789
Partner partner789 has no available slots, escalating application abc123 up the sponsor line
Checking sponsor sponsor123 (level 1) for application abc123
Sponsor sponsor123 can accept application abc123, auto-accepting
Successfully auto-accepted application abc123 by partner sponsor123
```

## Алгоритм эскалации

1. Получить просроченные заявки (статус "pending", создана более 4 часов назад)
2. Для каждой заявки:
   - Проверить, может ли текущий партнер принять заявку (есть ли свободные слоты)
   - Если да - автоматически принять заявку
   - Если нет - найти спонсора партнера
   - Проверить, может ли спонсор принять заявку
   - Повторять до максимального количества уровней (10)
   - Если никто не может принять - отклонить заявку с причиной

## Безопасность

- Система использует те же проверки лимитов партнеров, что и ручная обработка
- Все операции логируются для аудита
- В случае ошибки обработка продолжается для остальных заявок
- Graceful shutdown при получении сигналов SIGINT/SIGTERM

## Тестирование

Для тестирования можно:

1. Создать заявку на партнерство
2. Изменить время создания заявки в базе данных на 5 часов назад
3. Запустить processor вручную или дождаться следующего цикла
4. Проверить, что заявка была обработана автоматически
