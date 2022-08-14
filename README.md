**Система хранения и распространения приватной информации между даверенными клиентами**

Архитектура
- в качествет СХД использьется Postgresql 12+ версии
- система обмена данными реализована с помошью технологии webSocket

Принцеп работы
Пользователь регистрируется в системе, далее у пользователя появляется авторизоваться - получив токен которым в дальнейшем подписываются все запросы
для авторизированного пользователя система предоставляет функционал создания комнат, в которые пользователь может публиковать свои приватные данные, также пользователь - который создал комнату , может приглошать предоставляя права доступа других зарегистрированным пользователям системы. тем самым все кто является участником "комнаты" могут получать и отправлять приватные данные

Клиент 
Реализовано cli приложение позволяющее совершать операция 
- Регистрации
- Авторизации
- Получения публичной информации любого зарегистрированного в системе пользователя
- Получение публичной информации Лобби
- Приглошение пользователей в лобби
- Добавление данны в лобби


Также сервер работает с любым из клиентов поддерживающим генерацию Json запросов к серверу

**Запросы** 

```
Регистрация
{
   "action" : "register",
   "login" : "login",
   "password" : "password"
}
Авторизация
{
   "action" : "auth",
   "login" : "login",
   "password" : "password"
}

Создание Лобби
{
   "action" : "create_lobby",
   "token" : "927d8e77-ecf8-42eb-808d-01d9c6c3dd42",
   "lobby_name" : "first_lobby"
}

Получение публичного ключа пользователя
{
   "action" : "get_public_token",
   "login" : "login"
}

Получение Id Лобби
{
   "action" : "lobby_id",
   "lobby_name" : "first_lobby"
}

Добавление пользователя к Лобби
{
   "action" : "invite_lobby",
   "lobby_name" : "first_lobby",
   "token" : "927d8e77-ecf8-42eb-808d-01d9c6c3dd42",
   "public_token" : "0fbbe04c-87bb-4c33-9aae-f6eea9c19b4f"
}

Додписаться на сообщения из Лобби
{
   "action" : "subscribe",
   "lobby_name" : "invite_lobby",
   "token" : "927d8e77-ecf8-42eb-808d-01d9c6c3dd42"
}

Постинг сообщений в Лобби 
{
   "action" : "message",
   "lobby_name" : "223",
   "token" : "927d8e77-ecf8-42eb-808d-01d9c6c3dd42",
   "lobby_id" : "927d8e77-ecf8-42eb-808d-01d9c6c3dd42",
   "row_data" : "Новый пароль для учетки админа 'test123'"
}
```

Запуск 
Сервер и БД запускаются в Docker контейнераз, командами `make ecosysten_up`
Клиент компилируется из дериктории `./client`