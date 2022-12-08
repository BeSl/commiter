## Docker Compose
### Start
`docker-compose --project-name="committer-pg-14" up -d`

### Stop
`docker-compose --project-name="committer-pg-14" down`

### TODO
* утвердить состав таблиц
* перенести общие структуры в отдельный модуль
* доделать коммиты
* перенести прекоммит сюда
* подумать как оповещать по ошибкам

* добавление через ui пользователей


- проверка настроек среды
- ручка для инструкций

## краткое описание

есть доступные http-ручки:
- ping - GET - без параметров. Для проверки доступности сервиса. Ответ Код 200. Содержание "pong"
- uploadtoquery - POST - в загаловках параметров нет. В теле запроса строка json. Принимается тело запроса и пишется в очередь обработки

### uploadtoquery
Тело запроса:
строка json. 
Состав:
```
"author":
{
    extId: "id36",
    name:"@BeSl"
},
"DataProccessor":
{
    "base64data": "",
    "name":"",
    "extID" "",
    "type":"обработка|отчет|заполнение тч"

},
"textCommit":"any text",
"dataevent":"YYYY.MM.DD.HH.mm.ss"


```
