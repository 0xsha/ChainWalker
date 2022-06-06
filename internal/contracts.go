package internal

import (
	"context"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
)

func DisasmContractsEVM(outDir string) {

	evmPath := "/usr/local/bin/evm"

	files, err := ioutil.ReadDir(outDir)
	if err != nil {
		log.Fatal().Err(err).Msg("Can not open dir")
	}

	for _, file := range files {
		//fmt.Println(file.Name(), file.IsDir())

		if !file.IsDir() {

			getwd, err := os.Getwd()
			if err != nil {
				log.Err(err).Msg("Failed to get dir")

			}

			// careful /
			args := []string{"disasm", getwd + "/" + outDir + file.Name()}

			getDisasm, err := ExecuteCommand(evmPath, 600, args...)

			if err != nil {
				log.Err(err).Msg("Failed to disasm")
			}

			// sanity
			if strings.HasPrefix(file.Name(), "0x") {
				//fmt.Println(getDisasm)
				err = os.WriteFile("output/"+file.Name()+"_opcode", []byte(getDisasm), 0666)
				if err != nil {
					log.Err(err).Msg("Failed to write disasm")

				}
			}

		}
	}

}

func DownloadContractsEVM(server string, start int64, end int64, balance float64, concurrency int) {

	client, err := ethclient.Dial(server)
	if err != nil {
		log.Fatal().Err(err).Msg("Can not connect")
	}

	// current block and current amount of transactions
	currentBlock, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatal().Msg("Can not get the current block")

	}

	log.Info().Msg("Current block : " + strconv.FormatUint(currentBlock, 10))

	blockNumber := big.NewInt(int64(currentBlock))
	blockData, err := client.BlockByNumber(context.Background(), blockNumber)

	if err != nil {
		log.Fatal().Msg("Can not get transactions on the current block")

	}

	log.Info().Msg("Total number of transactions on the current block : " + strconv.Itoa(len(blockData.Transactions())))

	wg := sync.WaitGroup{}

	maxGoRoutines := make(chan struct{}, concurrency)

	for i := start; i < end; i++ {

		wg.Add(1)

		//i := i
		go func(i int64) {

			defer wg.Done()
			maxGoRoutines <- struct{}{}
			blockData, err := client.BlockByNumber(context.Background(), big.NewInt(i))

			if err != nil {

				log.Err(err).Msg("Can not get block data")
				return
			}

			transactions := blockData.Transactions()
			// sanity
			if len(transactions) > 0 {

				for _, tx := range transactions {

					// to is null, so we assume we got contract creation transaction
					if tx.To() == nil {

						transaction, err := client.TransactionReceipt(context.Background(), tx.Hash())

						if err != nil {
							log.Err(err).Msg("can not get bytecode")
							continue
						}

						// trying to grab bytecode
						bytecode, err := client.CodeAt(context.Background(), transaction.ContractAddress, nil) // nil is the latest block
						if err != nil {
							log.Err(err).Msg("can not get bytecode")
						}

						// at the moment balance of the contract itself
						accountBalance, err := client.BalanceAt(context.Background(), transaction.ContractAddress, nil)
						if err != nil {
							log.Err(err).Msg("Can not get balance")
						}

						y := big.NewFloat(balance)

						if balance > 0 {
							if WeiToEther(accountBalance).Cmp(y) == 1 {
								log.Info().Msg("Current block : " + strconv.FormatInt(i, 10))
								log.Info().Msg("Contract address : " + transaction.ContractAddress.String())
								log.Info().Msg("Contract balance : " + WeiToEther(accountBalance).String())

								// debug
								log.Debug().Msg(hex.EncodeToString(bytecode))
								// write evm hex string to file
								WriteHexToFile(transaction.ContractAddress.String(), hex.EncodeToString(bytecode))
								log.Info().Msg("----------------------------------------------")
							}
						} else {

							log.Info().Msg("Current block : " + strconv.FormatInt(i, 10))
							log.Info().Msg("Contract address : " + transaction.ContractAddress.String())
							log.Info().Msg("Contract balance : " + WeiToEther(accountBalance).String())

							// debug
							log.Debug().Msg(hex.EncodeToString(bytecode))
							// write evm hex string to file
							WriteHexToFile(transaction.ContractAddress.String(), hex.EncodeToString(bytecode))
							log.Info().Msg("----------------------------------------------")

						}

					}

				}

			}

			<-maxGoRoutines

		}(i)

	}

	wg.Wait()

}

// WeiToEther thanks @kimxilxyong
func WeiToEther(wei *big.Int) *big.Float {
	f := new(big.Float)
	f.SetPrec(236)
	f.SetMode(big.ToNearestEven)
	fWei := new(big.Float)
	fWei.SetPrec(236)
	fWei.SetMode(big.ToNearestEven)
	return f.Quo(fWei.SetInt(wei), big.NewFloat(params.Ether))
}
