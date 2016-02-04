# gocurrency
Currency Converter utilizing http://fixer.io api

Simple currency conversion utilizing shopspring/decimal and exchange rates from http://fixer.io API

```
moneyUSD, _ := currency.Money("13.33", "USD")
moneyGBP, _ := moneyUSD.Convert("GBP")
fmt.Println(moneyGBP)
```
