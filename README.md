# Currency Conversion Program

This program allows you to convert amounts between different currencies using real-time exchange rates from a third-party API.

## Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/yourusername/currency-conversion.git
   cd currency-conversion
   ```

2. Install dependencies:

   ```sh
   go mod tidy
   ```

3. Create a `.env` file with your API key and base URL for the exchange rate provider:

   ```sh
   API_KEY=your_api_key_here
   BASE_URL=https://api.exchangerate-api.com/v4/latest
   ```

4. Create a `config.yaml` file with the supported currency codes and their symbols:

   ```yaml
   supported_currencies:
     USD: "$"
     EUR: "€"
     GBP: "£"
     JPY: "¥"
     AUD: "A$"
     CAD: "C$"
     CHF: "CHF"
     CNY: "¥"
     SEK: "kr"
     NZD: "NZ$"
     TWD: "NT$"
     // ...existing code...
   ```

## Usage

To convert an amount from one currency to another, run the program:

```sh
go run main.go
```

### Command Line Options

You can also use the following command line options:

- `-from`: The currency code to convert from (default: USD)
- `-to`: The currency code to convert to (required)
- `-amount`: The amount to convert (required)

Example:

```sh
go run main.go -from USD -to JPY -amount 100
```

## Supported Currencies

The following currency codes and their symbols are supported:

- USD: $
- EUR: €
- GBP: £
- JPY: ¥
- AUD: A$
- CAD: C$
- CHF: CHF
- CNY: ¥
- SEK: kr
- NZD: NZ$
- TWD: NT$

## Configuration

The program uses environment variables and a configuration file for settings:

- `.env` file:

  ```sh
  API_KEY=your_api_key_here
  BASE_URL=https://api.exchangerate-api.com/v4/latest
  ```

- `config.yaml` file:

  ```yaml
  supported_currencies:
    USD: "$"
    EUR: "€"
    GBP: "£"
    JPY: "¥"
    AUD: "A$"
    CAD: "C$"
    CHF: "CHF"
    CNY: "¥"
    SEK: "kr"
    NZD: "NZ$"
    TWD: "NT$"
    // ...existing code...
  ```

## Testing

To run the tests, use the following command:

```sh
go test -v
```

## Deployment

To deploy the program, follow these steps:

1. Build the program:

   ```sh
   go build -o currency-conversion
   ```

2. Run the executable:

   ```sh
   ./currency-conversion
   ```

3. Ensure that the `.env` and `config.yaml` files are properly configured on the deployment environment.
