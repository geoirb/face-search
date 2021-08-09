# find-face

## Формат ответа

```json
{
    "is_ok": bool,
    "payload": {},
}
```

> is_ok - признак успешности запроса
> 
> payload - данные, если is_ok=true, либо сообщение об ошибке, если is_ok=false

## Получение настройки поиска лиц

URI: `/v1/api/config`

Пример запроса:

```shell
curl --request GET \
--url http://127.0.0.1:8080/v1/api/config
```

Пример ответа `payload`:
```json
{
    "timeout": 20000000000,
    "actions": [
      {
        "timeout": "navigate",
        "params": [
          "https://search4.com/vk01/index.html"
        ]
      },
      {
        "timeout": "click",
        "params": [
          "upload-button"
        ]
      },
      {
        "timeout": "set_upload_files",
        "params": [
          "input[type='file']",
          "/home/geoirb/Pictures/file.png"
        ]
      },
      {
        "timeout": "click",
        "params": [
          "effects-continue--upload"
        ]
      },
      {
        "timeout": "wait_not_visible",
        "params": [
          "uppload-container"
        ]
      },
      {
        "timeout": "click",
        "params": [
          "search-button"
        ]
      },
      {
        "timeout": "wait_visible",
        "params": [
          "div.row.no-gutters"
        ]
      }
    ]

}
```

# Изменения настройки поиска лиц

URI: `/v1/api/config`

Изначальные параметры поиска считываются из файла `/config/config.yml`

Пример запроса:

```shell
curl --request PUT \
  --url http://127.0.0.1:8080/v1/api/config \
  --header 'Content-Type: application/json' \
  --data '{
    "timeout": 20000000000,
    "actions": [
      {
        "timeout": "navigate",
        "params": [
          "https://search4.com/vk01/index.html"
        ]
      },
      {
        "timeout": "click",
        "params": [
          "upload-button"
        ]
      },
      {
        "timeout": "set_upload_files",
        "params": [
          "input[type='\''file'\'']",
          "/home/geoirb/Pictures/file.png"
        ]
      },
      {
        "timeout": "click",
        "params": [
          "effects-continue--upload"
        ]
      },
      {
        "timeout": "wait_not_visible",
        "params": [
          "uppload-container"
        ]
      },
      {
        "timeout": "click",
        "params": [
          "search-button"
        ]
      },
      {
        "timeout": "wait_visible",
        "params": [
          "div.row.no-gutters"
        ]
      }
		]
}'
```

# Запуск поиска по лицам

URI: `/v1/api/face_search`

Запускается поиск по изображению лежащего по `url` из запроса

Пример запроса:
```shell
curl --request POST \
  --url http://127.0.0.1:8081/v1/api/face_search \
  --header 'Content-Type: application/json' \
  --data '{
   "url":"https://sun9-85.userapi.com/impf/c850736/v850736489/16c130/rMcdDHy0HJ4.jpg?size=864x1080&quality=96&sign=86112d4bfa33cd27dd321d13a96bfc45&type=album"
}'
```

Пример `payload`:
```json
{
  "status": false,
  "uuid": "f81c4ff2-dd3a-43b9-a8b9-708faf527d20",
  "photo_hash": "yM6wX7Eikr69XaV3jDoIg9NKtPa+9+1h7tPga6QITnQ="
}
```

* Если для данного изображения выполнялся поиск, то в ответе будет результат предыдущего поиска

Пример `payload`:
```json
{
  "status": true,
  "uuid": "f81c4ff2-dd3a-43b9-a8b9-708faf527d20",
  "photo_hash": "yM6wX7Eikr69XaV3jDoIg9NKtPa+9+1h7tPga6QITnQ=",
  "profiles": [
    {
      "full_name": "Полина Дорошенко ",
      "link_profile": "https://vk.com/id38169513",
      "link_photo": "https://vk.com/id38169513?z=photo38169513_457241192%2Fphotos38169513",
      "confidence": "98.32%"
    },
    {
      "full_name": "Полина Дорошенко ",
      "link_profile": "https://vk.com/id38169513",
      "link_photo": "https://vk.com/id38169513?z=photo38169513_286737917%2Fphotos38169513",
      "confidence": "61.81%"
    },
    {
      "full_name": "Полина Дорошенко ",
      "link_profile": "https://vk.com/id38169513",
      "link_photo": "https://vk.com/id38169513?z=photo38169513_201925188%2Fphotos38169513",
      "confidence": "61.58%"
    }
  ],
  "create_at": 1628371204
}
```

* Если предыдущий поиск завершился с ошибкой, то поиск будет перезапущен

Пример `payload`:
```json
{
  "status": false,
  "error": "текст ошибки",
  "uuid": "f81c4ff2-dd3a-43b9-a8b9-708faf527d20",
  "photo_hash": "yM6wX7Eikr69XaV3jDoIg9NKtPa+9+1h7tPga6QITnQ=",
  "create_at": 1628371204
}
```

# Запрос результатов поиска

URI: `/v1/api/face_search/{uuid}`

Получение результатов по UUID поиска

```shell
curl --request GET \
  --url http://127.0.0.1:8081/v1/api/face_search/c98fa475-a3a3-4e1f-9bd2-cc45cbeef058
```

Пример `payload`:
```json
{
  "status": true,
  "uuid": "f81c4ff2-dd3a-43b9-a8b9-708faf527d20",
  "photo_hash": "yM6wX7Eikr69XaV3jDoIg9NKtPa+9+1h7tPga6QITnQ=",
  "profiles": [
    {
      "full_name": "Полина Дорошенко ",
      "link_profile": "https://vk.com/id38169513",
      "link_photo": "https://vk.com/id38169513?z=photo38169513_457241192%2Fphotos38169513",
      "confidence": "98.32%"
    },
    {
      "full_name": "Полина Дорошенко ",
      "link_profile": "https://vk.com/id38169513",
      "link_photo": "https://vk.com/id38169513?z=photo38169513_286737917%2Fphotos38169513",
      "confidence": "61.81%"
    },
    {
      "full_name": "Полина Дорошенко ",
      "link_profile": "https://vk.com/id38169513",
      "link_photo": "https://vk.com/id38169513?z=photo38169513_201925188%2Fphotos38169513",
      "confidence": "61.58%"
    }
  ],
  "create_at": 1628371204
}
```

* Если предыдущий поиск завершился с ошибкой, то поиск будет перезапущен

Пример `payload`:
```json
{
  "status": false,
  "error": "текст ошибки",
  "uuid": "f81c4ff2-dd3a-43b9-a8b9-708faf527d20",
  "photo_hash": "yM6wX7Eikr69XaV3jDoIg9NKtPa+9+1h7tPga6QITnQ=",
  "create_at": 1628371204
}

expected: &result.Facade{Result:service.Result{Status:"success", Error:"", UUID:"test-uuid", PhotoHash:"test-hash", Profiles:[]service.Profile{service.Profile{FullName:"test-name-1", LinkProfile:"test-link-profile-1", LinkPhoto:"test-link-photo-1", Confidence:"test-confidence-1"}, service.Profile{FullName:"test-name-2", LinkProfile:"test-link-profile-2", LinkPhoto:"test-link-photo-2", Confidence:"test-confidence-2"}}, UpdateAt:1, CreateAt:1}, timeFunc:(func() int64)(0x89c6a0), uuidFunc:(func() string)(0x89c6c0), storage:(*mongo.Mock)(0xc0001805a0)}

actual  : &result.Facade{Result:service.Result{Status:"success", Error:"", UUID:"test-uuid", PhotoHash:"test-hash", Profiles:[]service.Profile{service.Profile{FullName:"test-name-1", LinkProfile:"test-link-profile-1", LinkPhoto:"test-link-photo-1", Confidence:"test-confidence-1"}, service.Profile{FullName:"test-name-2", LinkProfile:"test-link-profile-2", LinkPhoto:"test-link-photo-2", Confidence:"test-confidence-2"}}, UpdateAt:1, CreateAt:1}, timeFunc:(func() int64)(0x89c6a0), uuidFunc:(func() string)(0x89c6c0), storage:(*mongo.Mock)(0xc0001805a0)}