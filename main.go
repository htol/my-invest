package main

import (
    "encoding/json"
	"context"
	"log"
    "fmt"
	"math/rand"
	"time"
    "os"
    "github.com/joho/godotenv"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var token string

func init() {
    // loads values from .env into the system
    if err := godotenv.Load(); err != nil {
        log.Print("No .env file found")
    }
}

func main() {
    token, exists := os.LookupEnv("TINKOFF_OPENAPI_TOKEN")

    if !exists {
        log.Fatalln("Please provide TINKOFF_OPENAPI_TOKEN in .env file")
    }

    if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

    client := sdk.NewRestClient(token)
	rand.Seed(time.Now().UnixNano()) // инициируем Seed рандома для функции requestID

    cashTable := makeCashTable(getCash(client))
    positionsTable := makePositionsTable(getPositions(client))

    ui.Render(cashTable)
    ui.Render(positionsTable)

	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			break
		}
	}
}


var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Генерируем уникальный ID для запроса
func requestID() string {
	b := make([]rune, 12)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}


func errorHandle(err error) error {
	if err == nil {
		return nil
	}

	if tradingErr, ok := err.(sdk.TradingError); ok {
		if tradingErr.InvalidTokenSpace() {
			tradingErr.Hint = "Do you use sandbox token in production environment or vise verse?"
			return tradingErr
		}
	}

	return err
}


type Currency struct {
    Symbol string
    Ticker string
}


var Currencies = map[sdk.Currency]Currency{
    "RUB": { "₽", "" },
    "USD": { "$", "USD000UTSTOM" },
    "EUR": { "€", "EUR_RUB__TOM" },
}


func makeCashTable(data [][]string) *widgets.Table {
	table := widgets.NewTable()
    table.Rows = data
    table.TextStyle = ui.NewStyle(ui.ColorMagenta)
    table.TextAlignment = ui.AlignRight
    table.RowSeparator = false
	table.SetRect(61, 0, 61+18, len(table.Rows)+2)

    return table
}


func getCash(client *sdk.RestClient) [][]string {
    var table  [][]string
    table = append(table, []string{"Currency balance"})
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	currencies, err := client.CurrenciesPortfolio(ctx, sdk.DefaultAccount)

	if err != nil {
		log.Fatalln(err)
	}
    //log.Printf("%+v\n", currencies)

    for _, cur := range currencies {
        if cur.Balance == 0 {
            continue
        }
        var row []string
        row = append(row, fmt.Sprintf("%.2f %s", cur.Balance, Currencies[cur.Currency].Symbol))
        table = append(table, row)
    }
    return table
}


func makePositionsTable(data [][]string) *widgets.Table {
	table := widgets.NewTable()
    table.Rows = data
    table.TextStyle = ui.NewStyle(ui.ColorMagenta)
    table.TextAlignment = ui.AlignRight
    table.RowSeparator = false
    table.RowStyles[0] = ui.NewStyle(ui.ColorBlue, ui.ColorBlack, ui.ModifierBold)
	table.SetRect(0, 0, 60, len(table.Rows)+2)

    return table
}


func getPositions(client *sdk.RestClient) [][]string {
    var table [][]string
    table = append(table, []string{"Type", "Name", "Amount", "Avg. price"})

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    stocks, err := client.PositionsPortfolio(ctx, sdk.DefaultAccount)
    if err != nil {
        log.Fatalln(err)
    }

    for _, p := range stocks {
        var row []string
        row = append(row, string(p.InstrumentType))
        row = append(row, string(p.Name))
        row = append(row, fmt.Sprintf("%.2f", p.Balance))
        row = append(row, fmt.Sprintf("%.2f %s", p.AveragePositionPrice.Value, Currencies[p.AveragePositionPrice.Currency].Symbol))

        // fmt.Printf("%+v\n", row)
        table = append(table, row)
    }

    return table
}


func toJSON(data interface{}) string{
     dataJSON, err := json.MarshalIndent(data, "", "  ")
     if err != nil {
         log.Fatalln(err)
     }
    return string(dataJSON)
}
