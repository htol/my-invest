package cmd

import (
    "fmt"
    "log"
    "math/rand"
    "time"
	"context"
    "strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
)

// operationsCmd represents the operations command
var operationsCmd = &cobra.Command{
	Use:   "operations",
	Short: "Get operations from broker",
	Long: ``,
	Run: opsRun,
}

func init() {
	rootCmd.AddCommand(operationsCmd)
}

func opsRun(cmd *cobra.Command, args []string) {
    var token string = viper.GetString("token")
    TZ, err := time.LoadLocation("Europe/Moscow")
    cobra.CheckErr(err)
    from := time.Date(2020, 01, 01, 00, 00, 00, 00, TZ)
    to := time.Date(2020, 12, 31, 23, 59, 59, 9999, TZ)

    if token == ""  {
        log.Fatalln("Please provide TINKOFF_OPENAPI_TOKEN in .env file")
    }

    client := sdk.NewRestClient(token)
	rand.Seed(time.Now().UnixNano())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	operations, err := client.Operations(ctx, sdk.DefaultAccount, from, to, "")
	if err != nil {
		log.Fatalln(err)
	}

    // onlyCurrencies(operations)
    allOperations(operations)
}

func allOperations(data []sdk.Operation) {
    for _, entry := range data {
        fmt.Printf("%+v\n", entry)
    }
}

func onlyCurrencies(data []sdk.Operation) {
    fmt.Println("ID,Status,Commission,Currency,Payment,Price,Quantity,QuantityExecuted,FIGI,Time,Operation")
    for _, entry := range data {
        if (entry.InstrumentType == sdk.InstrumentTypeCurrency) &&
        (entry.Status == sdk.OperationStatusDone) &&
        (entry.OperationType != sdk.OperationTypeBrokerCommission) {
            var b strings.Builder
            //fmt.Printf("%d: %+v\n", idx, entry)

            // ID:589588285
            fmt.Fprintf(&b, "%s,", entry.ID)
            // Status:Done
            fmt.Fprintf(&b, "%s,", entry.Status)
            // Trades:[]
            // Commission:{Currency: Value:0}
            fmt.Fprintf(&b, "%+v,", entry.Commission.Value)
            // Currency:RUB
            fmt.Fprintf(&b, "%s,", entry.Currency)
            // Payment:-183.65
            fmt.Fprintf(&b, "%f,", entry.Payment)
            // Price:0
            fmt.Fprintf(&b, "%f,", entry.Price)
            // Quantity:0
            fmt.Fprintf(&b, "%d,", entry.Quantity)
            // QuantityExecuted:0
            fmt.Fprintf(&b, "%d,", entry.QuantityExecuted)
            // FIGI:BBG0013HGFT4
            fmt.Fprintf(&b, "%s,", entry.FIGI)
            // InstrumentType:Currency
            // fmt.Fprintf(&b, "%s,", entry.InstrumentType)
            // IsMarginCall:false
            // fmt.Fprintf(&b, "%t,", entry.IsMarginCall)
            // DateTime:2020-12-10 12:41:09 +0300 MSK
            fmt.Fprintf(&b, "%s,", entry.DateTime)
            // OperationType:BrokerCommission
            fmt.Fprintf(&b, "%s", entry.OperationType)

            fmt.Println(b.String())
        }
    }
}
