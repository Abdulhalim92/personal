Данный проект представляет собой сервис для учета личных финансов. Сервис позволяет пользователям вносить данные о движении денежных средств и отслеживать свое финансовое состояние.

При внесении записи, пользователь указывает: 
⁃	тип операции (например, доход),
⁃	категорию (например, зарплата),
⁃	счет (например, банковская карта ***4321),
⁃	сумму (например, 10 000),
⁃	дату (например, 30.11.2022).

Приложение позволяет генерировать отчеты на основе заданных фильтров (по типам, категориям, счетам, за определенный период). Например:
⁃	общий баланс на текущую дату по всем счетам (т.е. сумма доступных средств), 
⁃	список операций по определенному счету за указанный период, 
⁃	и так далее на основе указанных фильтров.
Кроме того, в приложении реализована возможность генерации и экспорта отчетов в формате .xlsx.

Некоторые особенности логики приложения:
⁃	Все операции разделены на 3 типа: доход, расход и перевод средств между своими счетами (например, снятие наличных денежных средств с банковской карты).
⁃	Типы операций «доход» и «расход» имеют свои под-типы (категории): доход – зарплата, подработка и т.д.; расход – аренда жилья, питание, транспорт и т.д. Реализована возможность добавлять свои кастомные категории.
⁃	Счета: наличные, безналичные и т.д. Реализована возможность добавлять свои кастомные счета.
⁃	По умолчанию, приложение ведет учет в одной единственной валюте – таджикский сомони (TJS).
⁃	В приложении реализованы возможности регистрации новых пользователей, идентификация и аутентификация, а также работа с токенами.
⁃	Кроме того, реализовано логирование.