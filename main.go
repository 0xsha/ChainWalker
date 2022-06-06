package main

import (
	"flag"
	"os"

	"github.com/0xsha/ChainWalker/internal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	// step I : get user settings
	startPtr := flag.Int64("s", 14000000, "start block (int)")
	endPtr := flag.Int64("e", 14010000, "end block (int)")
	urlPtr := flag.String("u", "https://cloudflare-eth.com", "RCP/IPC endpoint")
	evmPath := flag.String("ev", "/usr/local/bin/evm", "EVM path")
	conPtr := flag.Int("c", 1, "concurrency")
	debugPtr := flag.Bool("d", false, "sets log level to debug")
	balancePtr := flag.Float64("b", 0, "minimum balance (default 0)")
	outPtr := flag.String("o", "output/", "output directory")
	helpPtr := flag.Bool("h", false, "shows usage")
	printPtr := flag.Bool("p", false, "print on console only and do not download or disassemble contracts")

	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	banner := `ChainWalker 1.0.3-alpha - Usage`

	if *helpPtr {

		log.Info().Msg(banner)
		flag.PrintDefaults()
		return
	}

	// setup zero-log
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if *debugPtr {
		// print-outs the bytecode
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Info().Msg(banner)

	// step II : download the contracts from requested blocks
	internal.DownloadContractsEVM(*urlPtr, *startPtr, *endPtr, *balancePtr, *conPtr, *printPtr)

	if !(*printPtr) {
		// step III : disassemble EVM to opcode
		internal.DisasmContractsEVM(*outPtr, *evmPath)
	}

}
