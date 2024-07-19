# Addressport

## Introduction

**Yo, what the heck is Addressport?** Addressport is an investigation tool for analysing transactions on different blockchains.
Basically it maps blockchain transactions into graphs, using addresses as nodes and transactions as links(edges).

**Yo cool idea, but what am I going to use it for?** That's a great question, which I don't know the exact answer :)
Basically there was a guy who's trying to scam me that I wanted to track down, I wanted to find out at least which local exchange was he going to use to exchange his/her crypto so I could contact them and at least let them know, so I build this tool.
I hope you never gonna need to use this project, but as grandma always says, it's gonna rain 💦.

## Features and components

Addressport currently has one feature, which is "Address Info Graph", you can access it by opening the url below:

```
/api/v1/investigation/tools/transaction-association/<TARGET-ADDRESS>?depth=<DEPTH-OF-SEARCH>&format=<RESPONSE-FORMAT>
```

- `<TARGET-ADDRESS>` can be any EVM account address(addressport currently searches ethereum network, but adding any EVM based blockchains should be a piece of 🍰)
- `<DEPTH-OF-SEARCH>` is the level of transaction search, basically when you pass 1 as depth, it gives you all the direct transactions **from** or **to** that address. When you pass 2, beside the direct transactions, it also returns all the direct transactions which any of the previous level addresses where associated. Imagine this scenario: 1. Scammer creates a fresh address, 2. Uses it for some purposes, 3. Sends the amount to another fresh address, 4. Then sends the amount to an exchange user address, 5. Exchange sends the amount to their hot wallet which is mostly labled in the internet. Here, using depth of 4 whould give you some really good insights about the transactions flow.
- `<RESPONSE-FORMAT>` can be either `nodesandlinks` or `graph`. Use `graph` which is the default if you wanna get the result in graph format. Use `nodesandlinks` if you wanna get the result in json format and plot it yourself.

## How to start the project?

1. `git clone https://github.com/officer47p/addressport.git`
2. `cd addressport`
3. Create a `.env` file and put your etherscan api key there(you can use the `.env.example` file to find all possible environment variables)
4. Run: `make start` to start the server(default port is 3000)
5. Open your browser and open http://localhost:3000
