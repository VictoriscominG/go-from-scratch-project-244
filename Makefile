# Компилируем Go‑пакет в исполняемый файл
build:
	go build -o bin/gendiff ./cmd/gendiff

# Запускаем линтер для проверки кода
lint:
	golangci-lint run

# Автоматически исправляем ошибки, которые может исправить линтер
lint-fix:
	golangci-lint run --fix

#Тестируем функцию DiffFile на сравнении двух json файлах
test:
	go test -v ./internal/diff/