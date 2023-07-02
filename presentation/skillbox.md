[//]: # (https://revealjs.com/markdown/)
Сетевой многопоточный сервис для "StatusPage"

<img src="./assets/dancing-gopher.gif" alt="gopher" width="160px"/>
<div style="text-align: left; padding: 6px; font-size: 28px">
    <p>Дипломник: Марченков Виктор курс "Профессия GO-разработчик"</p>
    <p>Руководитель диплома: Сидоров Данил</p>
</div> 

---

Личная информация
- образование: высшее 
- курсы: Coursera, Udacity, Udemy ...
- стаж: в разработке ~ 7 лет
- email: victormv03@gmail.com 
- имеется опыт применения:

<div>
    <img src=".\assets\js.png" height="100px"/>
    <img src=".\assets\svelte.png" height="100px"/>
    <img src=".\assets\perl.png" height="100px"/>
    <img src=".\assets\mysql.png" height="100px"/>
</div>
<div>
    <img src=".\assets\python.png" height="100px"/>
    <img src=".\assets\github.png" height="100px"/>
    <img src=".\assets\gopher.png" height="120px"/>
</div>

----

## Подробнее о текущих разработках
Сервис построения графиков (временные ряды, пространственные профили ...) с использованием
<div>
    <img src=".\assets\svelte.png" height="100px"/>
    <img src=".\assets\js.png" height="100px"/>
</div> 

<a href="http://smisdev.vega.smislab.ru/geosmis_charts_v3/demo.html" target="_blank">Примеры работ</a>

---

## Цель

1) Систематизировать знания о GO
2) Получить навык разработки под управлением профессионалов
3) Начать применять на практике

---

### Постановка задачи

- Создать страницу информирования клиентов о текущем состоянии системы.
- Результат разработки разместить на сервере Heroku

---

#### Информационные системы:
- SMS (csv)
- Голосовые звонки (csv)
- E-mail (csv)
- Служба биллинга (csv)
- MMS (http)
- История инцидентов (http)
- Служба поддержки клинетов (http)

---

### Что вызвало затруднение
- ТЗ описывает не тот генератор данных который предложен для разработки
- Необходимость работать с файлами и с сервером одновременно
- Для обновления данных приходится перезапускать генератор руками
- Цель объявленная в ТЗ (многопоточный сервис) теряет смысл при размещении на Heroku

---

При реализации проекта использовались
- IDE - Goland
- Среда разработки Windows 10
- GO версия 1.16
<div>
    <img src=".\assets\goland.png" height="120px"/>
    <img padding="0 50px 0 50px" src=".\assets\windows10.png" height="120px"/>
    <img src=".\assets\gopher116.png" height="120px"/>
</div>

----

Структура проекта
```shell
│   build_app.cmd
│   build_generator.cmd
│   Dockerfile
│   go.mod
│   go.sum
│   README.md
│   start.cmd
│   start.sh
│   start_docker.sh
│   start_generator.cmd
│   stop_docker.sh
│   test.http
│
├───bind
│   └───app
│           sbdiplom.exe
│
├───cmd
│   └───app
│           main.go
│
├───config
│       config.go
│       config.json
│
├───generator
│       main.go
│       README.md
│
├───info
│       chart.min.js
│       false.png
│       index.html
│       main.css
│       main.js
│       status_page.html
│       true.png
│
├───internal
│   │   entities.go
│   │   handlers.go
│   │   handlers_test.go
│   │
│   ├───app
│   │       app.go
│   │       
│   └───test-data
│           billing.data
│           email.data
│           sms.data
│           voice.data
│
└───pkg
        csv_reader.go
        mylog.go
        valid_data.go
```

----

## main.go
```go
package main

import (
	"flag"
	"sbdaemon/internal/app"
)

var port = flag.Int("p", 8282, "server port")

func main() {
	flag.Parse()
	app.Run(*port)
}
```

----

## app.go
```go
package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	cfg "sb-diplom-v2/internal"
	"sb-diplom-v2/pkg"
	"syscall"
	"time"
)

func Run(port int) {

	var cfg_ cfg.Config
	var resultT cfg.StatusResult
	var service = "skillbox diploma"
	start := fmt.Sprintf(":%d", port)
	data, err := ioutil.ReadFile("./config/config.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	json.Unmarshal(data, &cfg_)

	resultT.HandlerFiles(cfg_)

	mux := http.NewServeMux()
	mux.HandleFunc("/api", resultT.HandlerHTTP)
	fs := http.FileServer(http.Dir("generator"))
	http.Handle("/", fs)
	
	info := http.FileServer(http.Dir("."))
	mux.Handle("/info/", info)
	mux.Handle("/presentation/", info)
	siteHandler := pkg.AccessLogMiddleware(mux)

	srv := &http.Server{
		Addr: fmt.Sprintf("%s", start),
		//Handler: mux,
		Handler: siteHandler,
	}

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error %v\n", err)
		}
	}()
	log.Printf("%s starting on %d\n", service, port)

	<-stop

	log.Printf("%s shutting down ...\n", service)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
}
```

---

Тестирование
- Для тестов связанных с обработкой файлов применялись unit тесты
- Для тестов связанных с запрсами по http контролировалось отсутствие ошибок соединения。

----

Пример теста
```go
func TestSmsHandler(t *testing.T) {

	var want [][]SMSData
	t.Run("is error string deleted and result arranged", func(t *testing.T) {
		got := SmsHandler("./test-data/sms.data")

		sms1 := SMSData{"BL", "68", "1594", "Kildy"}
		sms2 := SMSData{"US", "36", "1576", "Rond"}
		tmp1 := []SMSData{sms1, sms2}
		tmp2 := []SMSData{sms1, sms2}
		want = append(want, tmp1, tmp2)

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})
	t.Run("is some string deleted", func(t *testing.T) {
		got := SmsHandler("./test-data/sms.data")

		sms1 := SMSData{"BL", "68", "1594", "Kildy"}
		sms2 := SMSData{"US", "36", "1576", "Rond"}
		tmp1 := []SMSData{sms1, sms2}
		tmp2 := []SMSData{sms1, sms2}
		want = append(want, tmp1, tmp2)

		if len(got) == len(want) {
			t.Errorf("len(got) '%d' len(want) '%d'", len(got), len(want))
		}
	})

}
```

----

Результаты unut тестирования
```shell
=== RUN   TestSmsHandler
=== RUN   TestSmsHandler/is_error_string_deleted_and_result_arranged
=== RUN   TestSmsHandler/is_really_corrupted_string_deleted
--- PASS: TestSmsHandler (0.00s)
    --- PASS: TestSmsHandler/is_error_string_deleted_and_result_arranged (0.00s)
    --- PASS: TestSmsHandler/is_really_corrupted_string_deleted (0.00s)
=== RUN   TestVoiceHandler
=== RUN   TestVoiceHandler/is_error_string_deleted_and_result_arranged
=== RUN   TestVoiceHandler/is_some_string_deleted
--- PASS: TestVoiceHandler (0.00s)
    --- PASS: TestVoiceHandler/is_error_string_deleted_and_result_arranged (0.00s)
    --- PASS: TestVoiceHandler/is_some_string_deleted (0.00s)
=== RUN   TestEmailHandler
=== RUN   TestEmailHandler/is_error_string_deleted
len(result) 1
--- PASS: TestEmailHandler (0.00s)
    --- PASS: TestEmailHandler/is_error_string_deleted (0.00s)
=== RUN   TestBillingHandler
=== RUN   TestBillingHandler/is_string_decoded_correctly
BillingHandler billing {true true false false true false}
got internal.BillingData {true true false false true false}
--- PASS: TestBillingHandler (0.00s)
    --- PASS: TestBillingHandler/is_string_decoded_correctly (0.00s)
PASS
```

----

HTTP тесты (файл test.http)
```shell
GET http://127.0.0.1:8282/api
Accept: application/json

###

GET http://127.0.0.1:8383/mms
Accept: application/json

###

GET http://127.0.0.1:8383/support
Accept: application/json

###

GET http://127.0.0.1:8383/accendent
Accept: application/json

###
```

---

### Результат на Heroku
Для размещения на Heroku структура проекта была изменена и дополненна Dockerfile и heroku.yml

- <a href="https://github.com/VictorMarchenkov/sbdiplom" target="_blank">репозиторий на Github</a>
- <a href="https://pure-peak-92586.herokuapp.com/api" target="_blank">API на Heroku</a>
- <a href="https://pure-peak-92586.herokuapp.com/status_page.html" target="_blank">StatusPage на Heroku</a>

----

Структура проекта для размещения на Heroku
```shell
│   Dockerfile
│   go.mod
│   go.sum
│   heroku.yml
│   main.go
│   README.md
│   start.cmd
│   start.sh
│   stop.cmd
│   test.cmd
│   USAGE.md
│
└───internal
    │   app.go
    │   entities.go
    │
    ├───config
    │       config.go
    │       config.json
    │
    ├───generator
    │       billing.data
    │       chart.min.js
    │       email.data
    │       false.png
    │       voice.data
    │
    └───pkg
            csv_reader.go
            utils.go
            valid_data.go
```

----

Dockerfile

```dockerfile
FROM golang:1.16 AS api

WORKDIR /buildapp
COPY . .
ADD go.mod go.sum /buildapp/

RUN CGO_ENABLED=0 GOOS=linux go build -o diploma ./main.go

RUN CGO_ENABLED=0 GOOS=linux go build -o ./internal/generator/generator ./internal/generator/main.go

EXPOSE 8282

CMD ["./start.sh"]
```

----

heroku.yml
```yml
build:
  docker:
    web: Dockerfile
```

---

## Итог 
- Получил опыт разработки на GO
- Внедрил сервис в котором использован GO сервер 
- Получил представление о:
   <div>
        <img height="50px" src=".\assets\nginx.png"/>
        <img padding= "0 50px 0 50px" height="50px" src=".\assets\docker.png"/>
        <img height="50px" src=".\assets\heroku.png"/>
   </div>

 
---


Спасибо за внимание!

Вопросы