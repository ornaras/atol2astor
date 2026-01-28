# ATOL-TO-ASTOR

Автоматический конвертер экспорта Frontol 6 из формата Atol в формат ASTOR

![Затраченное время](https://waka.ornaras.ru/api/badge/ornaras/interval:week/project:atol2astor)
![Поддерживаемые ОС](https://img.shields.io/badge/%D0%9F%D0%BE%D0%B4%D0%B4%D0%B5%D1%80%D0%B6%D0%BA%D0%B0-Windows_7%2B-blue?logo=windows)
![GitHub License](https://img.shields.io/github/license/ornaras/atol2astor?label=%D0%9B%D0%B8%D1%86%D0%B5%D0%BD%D0%B7%D0%B8%D1%8F)
![GitHub Downloads](https://img.shields.io/github/downloads/ornaras/atol2astor/total?label=%D0%A1%D0%BA%D0%B0%D1%87%D0%B0%D0%BD%D0%BE)
![GitHub repo size](https://img.shields.io/github/repo-size/ornaras/atol2astor?label=%D0%A0%D0%B0%D0%B7%D0%BC%D0%B5%D1%80%20%D1%80%D0%B5%D0%BF%D0%BE%D0%B7%D0%B8%D1%82%D0%BE%D1%80%D0%B8%D1%8F)


## Параметры запуска

```
C:\ProgramData\atol2astor>./atol2astor.exe -h
Usage of atol2astor.exe:
  -d    Запуск в режиме отладки
  -s    Запуск в режиме сервиса
```

## Пример конфигурации

> [!TIP]
> Место расположения каталога с файлами приложения: `%PROGRAMDATA%\atol2astor`

```xml
<!--C:\ProgramData\atol2astor\config.xml-->

<configuration>
    <!--Тег 'interval' устанавливает интервал между проверками файлов в минутах-->
    <interval>5</interval>
    <!--В теге 'imports' хранятся пути к файлам, которые будут конвертироваться в формат ASTOR-->
    <!--Конвертированный файл сохраняется в той же директории, в которой находится оригинальный файл, с имененем 'export.txt'-->
    <!--Внимание! Не рекомендуется использовать несколько раз одну и ту же директорию.-->
    <!--Внимание! Оригинальный файл должен иметь имя отличное от 'export.txt'.-->
    <imports>
        <import path="C:\atol1\import.txt"/>
        <import path="C:\atol2\import.txt"/>
    </imports>
</configuration>
```

## Порядок установки
1) Скачать [последнюю версию](https://github.com/ornaras/atol2astor/releases/latest)
2) Запустить от имени админстратора
3) В меню выбрать первый пункт:
   ```batch
   Возможные действия:
   1) Установка службы
   2) Удаление службы
   3) Открыть конфигурацию
   
   Номер действия: 1
   ```

## Порядок сборки

> [!TIP]
> Go 1.20.x является последней версией поддерживающей сборку приложений для Windows 7 и выше

1) Установить [Go 1.20+](https://go.dev/dl/)
2) Клонировать репозиторий: `git clone https://github.com/ornaras/atol2astor.git`
3) Настроить следующие переменные окружения:
   - **GOARCH**: `386`
   - **GOOS**: `windows`
4) Запустить сборку приложения: `go build -ldflags="-s -w" .`
