# upd от 05.12.2021
В ветке feature/Docker, все 4 уровная + рефакторинг

# Hackathon по GO

1. Level 1: полученный отчёт отобразить в консоли.   Минимальный функционал чтобы засчитать задание.
2. Level 2: полученный отчёт сохранить в JSON файле.
3. Level 3: полученный отчёт сохранить в CSV файле.
4. Level 404: можно пропустить сохранение в файле, а сохранить данные отчета в Database и через SQL запрос в консоли.
Читаем JSON -> формируем структуру данных отчёта -> сохраняем в бд

## О выполненом задании
Первый уровень - не стал выводить, дабы не засорять консоль. Реализуется методом `fmt.Println(reports)`

Второй уровень - генерируется файл в дериктории `/app` с названием `report-transactions.json`

Третий уровень - генерируется файл в дериктории `/app` с названием `report-transactions.csv`

## Для запуска
```
cd app/
go run main.go 
```

### Результатом будет выведенно в консоли:
1. Время парсинга Json-файла `Parse Json`
2. Время генерации Json-файла `Generate Json`
3. Время генерации Csv-файла  `Generate Csv`
4. Общее время выполнение программы`General time`

### В дирректории /app будут сгенерированны 2 файлы
1. report-transactions.json
2. report-transactions.csv
