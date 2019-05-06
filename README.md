# incrementor


https://github.com/square/certstrap
Tools to bootstrap CAs, certificate requests, and signed certificates.

TODO:
* docker-compose
* Version & Build
* Unit tests
* Integration tests
* Docs


Deal within ENV & Config

Фибдек

Не было тестов, нужно починить


- В `client_database.go` откатывание транзакции `Rollback` нужно делать под `defer` на случай аварийного завершения.
- Связь между таблицами пользователей и счетчиков лучше сделать черкз `user_id`. Сейчас используется `username`, поэтому если в будущем позволить пользователю редактировать свое имя, то все развалится.
- Поле `mutex` в структуре `Increment` не используется.
- В методе `SetSettings` и `IncrementNumber` нет блокировки данных, поэтому при параллельном использовании можно перетереть данные.
- Миграцию нужно сделать отдельным скриптом. Это не сервис, чтобы каждый раз его запускать в `docker-compose.yaml`.
- Миграцию нужно выполнять, при остановленном сервисе. Сейчас все делается одновременно, т.е. сервис доступен и может сохранить данные в неготовое хранилище.
- Нет тестов.
- При запуске сервисов через `docker-compose.yaml` произошла ошибка `psql:/docker-entrypoint-initdb.d/init_db.sql:5: ERROR:  role "incrementor" already exists`.
- При запуске произошла ошибка `.... cannot parse 'Incrementor.MaxValue' as int: strconv.ParseInt: parsing \"12312312312313\": value out of range","time":"2019-04-19T11:05:45Z"}`. Параметр `MaxValue` имеет тип `int32`, а тестовое значение не помещается в него.
- Добавлен endpoint для prometheus метрик и есть логи ошибок.
- Есть авторизация с jwt токеном.
- Используется DI и флаги для запуска.
- Нет документации на proto API.

https://habr.com/ru/company/funcorp/blog/418081/
https://habr.com/ru/company/funcorp/blog/418329/

https://habr.com/ru/post/425025/
https://habr.com/ru/post/271239/
https://segment.com/blog/5-advanced-testing-techniques-in-go/
https://medium.com/@thejasbabu/testing-in-golang-c378b351002d

Unit test
https://medium.com/@matryer/5-simple-tips-and-tricks-for-writing-unit-tests-in-golang-619653f90742

Integration test
https://hackernoon.com/integration-test-with-database-in-golang-355dc123fdc9