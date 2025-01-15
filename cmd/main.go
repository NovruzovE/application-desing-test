// Ниже реализован сервис бронирования номеров в отеле. В предметной области
// выделены два понятия: Order — заказ, который включает в себя даты бронирования
// и контакты пользователя, и RoomAvailability — количество свободных номеров на
// конкретный день.
//
// Задание:
// - провести рефакторинг кода с выделением слоев и абстракций
// - применить best-practices там где это имеет смысл
// - исправить имеющиеся в реализации логические и технические ошибки и неточности
package main

import (
	"fmt"
	"github.com/NovruzovE/application-design-test/internal/app"
)

/*
// logger interface?? slog??
// transactions in use case layer???
// в каком таймзоне передается время, а в каком хранить?

[] выделены слои
[] продумано дальнейшее расширение
[x] graceful shutdown
[x] router
[x] middleware
[x] validation
[] tests
[x] config
[] documentation

*/

var (
	version = "N/A"
)

func main() {
	fmt.Printf("\nBooking app\nBuild version: %s\n", version)

	bookingApp := app.NewApp()
	bookingApp.MustRun()
}
