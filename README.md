    Запуск инфарструктуры
    make up
  
    запуск сервера:
    go run cmd/server/main.go
  
    отправка файла: 
    go run cmd/client/main.go ~/desktop/file.png
    возвращает uuid
    (куда будут сохраняться файл, задаётся через переменную окуржения в .env)/ По mime заголовкам файлы не валидируются
  
    получение файла 
    go run cmd/client/main.go uuid  ~/folder
  
    получить список(работает без пагинации, отдаёт все)
    go run cmd/client/main.go list
