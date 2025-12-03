# cmd/shortener

## Сборка проекта

При сборке приложения можно установить значения переменных `buildVersion`, `buildDate` и `buildCommit`, используя флаг `-ldflags`.

Пример команды сборки:
```
go build -ldflags "-X main.buildVersion=v1.0.1 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')' -X main.buildCommit=$(git rev-parse HEAD)" -o shortener cmd/shortener/main.go
```

После запуска собранного бинарного файла, в консоль будет выведена информация о версии:
```
Build version: v1.0.1
Build date: 2024/05/20 12:34:56
Build commit: a1b2c3d4...
```

Если флаги не переданы, будут использоваться значения по умолчанию (`N/A`).

Инкремент 25.