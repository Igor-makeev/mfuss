/*
Multichecker staticlint нужен для статического анализа кода
и поиска возможных ошибок
# Как запускать

	./cmd/staticlint/mymultichecker -mymultichecker ./cmd/shortener/main.go
	./cmd/staticlint/mymultichecker ./...

# Перечень используемых анализаторов:
Стандартные статические анализаторы пакета golang.org/x/tools/go/analysis/passes;
Всех анализаторы класса SA пакета staticcheck.io;
  - анализаторы класса SA;
  - quickfix для рефакторинга кода;
  - simple для упрощения кода;
  - stylecheck для соблюдения правил стиля.

go-critic — расширяемый линтер исходного кода Go, предоставляющий проверки, отсутствующие в других линтерах;
errcheck — анализатор для проверки непроверенных ошибок в коде Go.;
Osexit,кастомный анализатор, который чекает функцию main пакета main на наличие прямого вызова функции os.Exit().
*/
package main
