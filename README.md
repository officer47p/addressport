# Addressport

Addressport is an investigation tool for analysing transactions on different blockchains.

## How to start the project?

1. Create a `.env` file and put your etherscan api key there(you can use the `.env.example` file to find all possible environment variables)
2. Cd into cloned repository and run: `go mod tidy`
3. Run: `make run_api` to start the server(default port is 3000)
4. Open your browser and enter the following url(change address and depth parameters to your desire):
   `http://localhost:3000/api/v1/investigation/tools/transaction-association/0xC1D8E8f14b6AA1cf2F2321348Cbb51d94dc73152?depth=2&`
