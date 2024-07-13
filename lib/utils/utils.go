package utils

import (
	"log"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/params"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/rand"
)

func WeiToEther(wei *big.Int) *big.Float {
	f := new(big.Float)
	f.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	f.SetMode(big.ToNearestEven)
	fWei := new(big.Float)
	fWei.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	fWei.SetMode(big.ToNearestEven)
	return f.Quo(fWei.SetInt(wei), big.NewFloat(params.Ether))
}

func StringToBigInt(s string) (*big.Int, bool) {
	n := big.Int{}
	return n.SetString(s, 10)
}

func LogReuqest(c *fiber.Ctx) (time.Time, func(time.Time)) {
	tempId := rand.Intn(100_000_000_000)
	reqId := strconv.Itoa(tempId)
	log.Printf("%s %s request(%s)\n", c.Method(), c.OriginalURL(), string(reqId))
	start := time.Now()
	return start, func(t time.Time) {
		log.Printf("request(%s) took %d ms\n", reqId, time.Since(t).Milliseconds())
	}
}
