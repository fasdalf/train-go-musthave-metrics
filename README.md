# go-musthave-metrics-tpl

Шаблон репозитория для трека «Сервер сбора метрик и алертинга».
 
## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m main template https://github.com/Yandex-Practicum/go-musthave-metrics-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/main .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).


`go install github.com/go-delve/delve/cmd/dlv@latest; go mod tidy; go mod vendor`
`dlv debug ./cmd/server/main.go --headless=true --api-version=2  -- -k 1234 --filestoragepath new.json --restore`

### profiling

go tool pprof -http=":9090" -seconds=30 http://localhost:8093/debug/pprof/goroutine
можно использовать опции focus и ignore — 
go tool pprof -focus=github.com/fasdalf/train-go-musthave-metrics -ignore="^runtime|^net/http|^[github.com|^golang.org](http://github.com%7C%5Egolang.org)" base.pprof
Это позволит сфокусироваться на вызовах, связанных именно с нашим проектом.

curl -s http://localhost:8093/debug/pprof/heap?seconds=30 > server.heap.result.pprof
go tool pprof -top -diff_base=profiles/base.pprof profiles/result.pprof


### godoc

* go install -v golang.org/x/tools/cmd/godoc@latest
* godoc -http=:8080 -goroot=.
* in browser http://localhost:8080/pkg/github.com/fasdalf/train-go-musthave-metrics/?m=all
* inline example http://localhost:8075/pkg/github.com/fasdalf/train-go-musthave-metrics/internal/common/retryattempt/?m=all#Retryer.Try

