## Инструкция к запуску
Для запуска выполнить косольную команду docker-compose up --build, миграции к БД выполнятся автоматически

## Инструкция к интеграционному тестированию
1) cd tests - переходим в директорию с тестами
2) docker-compose up - поднимаем тестовую бд
3) goose -dir ./migrations postgres "user=test password=test dbname=test host=localhost port=5434 sslmode=disable" up - применяем миграции
4) go test -v - запускаем тесты
5) docker-compose down - останавливаем и удаляем контейнер