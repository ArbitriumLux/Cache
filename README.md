# Cache
Приложение имплементация im memory Redis кеша.

Инструкция по запуску:
  1.Перейти в директорию с репозиторием.
  2.Вызвать "make".
  3.Перейти в директорию "artifcacts".
  4.Запустить "main".
  5.Перейти в браузере на "localhost:8080/map"
  
По умолчанию время жизни ключа равно пяти минутам, а интервал, через который запускается механизм очистки кеша десяти минутам.
Сохранение кеша на диск производится в директорию из которой производился запуск.