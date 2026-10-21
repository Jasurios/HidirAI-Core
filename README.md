# HidirAi

HidirAi — прототип агентного ИИ-ассистента на Go. Одно ядро с HTTP API и два клиента — Telegram и WhatsApp. Ядро общается с OpenAI-совместимой моделью, хранит историю пользователей и вызывает инструменты: встроенные и внешние MCP-серверы. Отвечает на русском, английском и таджикском.

## Архитектура

```
Telegram Bot ─┐
WhatsApp Bot ─┼── HTTP :2010 ── Core (Gin) ── LLM API (OpenAI-compatible)
              │                    ├── SQLite: keys, blacklist
              │                    ├── users/<id>/*.jsonl: история, MCP
              │                    └── tools + MCP clients
              └── сервисный ключ
```

Клиенты ходят в ядро на `http://127.0.0.1:2010`. Пользовательские запросы к `/chat` авторизуются пользовательским ключом, клиенты — сервисным.

## Компоненты

| Путь | Что это |
|---|---|
| `Core/` | Gin-сервер, агентный цикл, авторизация, SQLite, история, инструменты, vision/STT |
| `Telegram/` | Бот на telebot: текст, голос, фото; команды `/start`, `/delmydata`, `/ban`, `/newmcp` |
| `WhatsApp/` | Бот на whatsmeow: личные чаты, QR-сессия, текст/голос/фото; те же команды текстом |
| `MCP/` | Опциональный MCP-сервер на `127.0.0.1:9090`: `make_qrcode`, `image_gen` |
| `HidirAi_v3/` | HTML/Three.js-презентация проекта и сценарий выступления |
| `Site/` | Статический сайт проекта |
| `HidirAi_server/` | Готовая серверная сборка: бинарники, конфиги, `start.sh`/`stop.sh` |
| `OldCore/`, `v1.0.0/`, `v2.0.0/`, `v3.0.0/`, `test*/` | Архивные версии и эксперименты; активный код — в `Core/`, `Telegram/`, `WhatsApp/`, `MCP/` |

## Стек

- Go 1.27.1, Gin, `sashabaranov/go-openai`
- Telegram: `telebot`; WhatsApp: `whatsmeow` + `mattn/go-sqlite3` (CGO)
- SQLite: `modernc.org/sqlite` в ядре, `mattn/go-sqlite3` в WhatsApp
- MCP-клиент и MCP-сервер, инструменты QR и генерации изображений (Cloudflare)
- Groq для vision/STT, OpenAI-совместимый endpoint для основной модели
- Презентация: HTML/CSS + Three.js

## Быстрый старт

Требования: Go 1.27.1; для WhatsApp — `gcc` (CGO); сетевой доступ к Telegram, WhatsApp и API модели.

```bash
# ядро — первым
cd Core && go run .

# Telegram
cd Telegram && go run .

# WhatsApp (CGO)
cd WhatsApp && CGO_ENABLED=1 go run .

# MCP-сервер, опционально
cd MCP && go run .
```

Каждый сервис читает `config.env` из своей папки. Клиенты ожидают ядро на `127.0.0.1:2010`.

Для сервера:

```bash
./build_server.sh          # собрать Core, Telegram, WhatsApp в HidirAi_server/
# затем на сервере — HidirAi_server/README.md, ./start.sh, ./stop.sh
```

## Конфигурация

| Сервис | Переменные |
|---|---|
| Core | `ADMIN_PASSWORD`, `DATABASE`, `API`, `URL`, `MODEL`, `GROQ_API`, `GROQ_URL`, `MODEL_VISION`, `MODEL_SPEECH`, `GHAPI`, `TAVILY_API`, `CF_ACCOUNT_ID`, `CF_API_TOKEN`, `MESSAGE`, `BLACKLISTFILE` |
| Telegram | `TOKEN`, `ADMIN`, `ADMINPASSWORD`, `API` |
| WhatsApp | `ADMINWHATSAPP`, `ADMINPASSWORD`, `API` |
| MCP | `CF_ACCOUNT_ID`, `CF_API_TOKEN` |

Важно:

- В **Core** `API` — ключ модели; в **Telegram/WhatsApp** `API` — сервисный ключ. Это разные значения.
- `ADMIN_PASSWORD` в Core и `ADMINPASSWORD` в клиентах должны совпадать.
- Сервисный ключ хранится в `Core/database/maindb.db` (таблица `keys`, тип `Service`).
- `DATABASE=maindb.db` открывается как `./database/maindb.db` относительно рабочей папки Core.
- `config.env` содержат секреты и не должны попадать в публичный репозиторий.

## HTTP API ядра

| Метод и путь | Авторизация | Назначение |
|---|---|---|
| `POST /chat` | пользовательский ключ | диалог с ассистентом |
| `POST /service` | сервисный ключ | диалог от Telegram/WhatsApp |
| `POST /newusermcp` | пользовательский ключ | подключить MCP пользователю |
| `POST /newservicemcp` | сервисный ключ | подключить MCP сервису |
| `POST /newkey` | `ADMIN_PASSWORD` | создать ключ |
| `POST /ban/:user` | `ADMIN_PASSWORD` | бан пользователя |
| `POST /del/:user` | `ADMIN_PASSWORD` | удалить данные пользователя |

Ключи и бан-лист лежат в SQLite; таблицы создаются при старте (`database.Init`).

## Агентный цикл

- До 20 раундов вызова инструментов на один запрос (`Core/ai/ai.go`).
- Встроенные инструменты: `create_qrcode`, `image_gen`, `new_mcp`.
- MCP по умолчанию: Tavily и GitHub; пользователь может добавить свой MCP (URL + ключ).
- История: `Core/users/<id>/<id>.jsonl`; MCP: `Core/users/<id>/MCP.jsonl`.
- Поддержка изображений и голоса: Groq vision/STT.

## Данные

```
Core/database/maindb.db       # ключи, blacklist
Core/users/<id>/<id>.jsonl    # история диалога
Core/users/<id>/MCP.jsonl     # MCP-серверы пользователя
WhatsApp/session.db           # сессия WhatsApp
```

Папки `database/`, `users/`, `session.db` и все `config.env` — локальные данные, не для коммита.

## Презентация

- `HidirAi_v3/HidirAi_motion.html` — исходная колода;
- `HidirAi_v3/HidirAi_motion_standalone.html` — самодостаточная версия для показа;
- `HidirAi_v3/HidirAi_script.md` — сценарий на ~7 минут.

## Лицензия

Apache License 2.0 — см. `LICENSE`.
