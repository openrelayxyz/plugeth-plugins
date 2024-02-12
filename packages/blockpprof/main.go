package main

import (
	"bufio"
	"bytes"
	"math/big"
	"os"
	"time"
	"github.com/openrelayxyz/plugeth-utils/core"
	"runtime"
	"runtime/pprof"
	"path/filepath"
	"flag"
	"fmt"
)

var (
	log core.Logger
	buff *bytes.Buffer
	writer *bufio.Writer
	lastNewHead time.Time
	Flags = *flag.NewFlagSet("blockpprof", flag.ContinueOnError)
	blockGapFlag = Flags.Int("pprof.block.gap", 6, "The number of seconds between blocks that should trigger a pprof flush")
	blockProfile = Flags.Bool("pprof.block.profile", false, "Enable blocking profiling")
	blockGapFolder = Flags.String("pprof.folder", "/var/lib/ethereum/pprofs", "Directory where pprof flushes should be written")
	blockGapDuration = 12 * time.Second

)

func Initialize(ctx core.Context, loader core.PluginLoader, logger core.Logger) {
	log = logger
	blockGapDuration = time.Duration(*blockGapFlag) * time.Second
	buff = new(bytes.Buffer)
	writer = bufio.NewWriter(buff)
}

func PreProcessBlock(hash core.Hash, num uint64, blockrlp []byte) {
	if *blockProfile {
		runtime.SetBlockProfileRate(1)
	}
}

func NewHead(blockBytes []byte, hash core.Hash, logsBytes [][]byte, td *big.Int) {
	pprof.StopCPUProfile()
	if buff != nil && time.Since(lastNewHead) > blockGapDuration {
		pprofFilePath := filepath.Join(*blockGapFolder, hash.String())
		if fd, err := os.Create(pprofFilePath); err == nil {
			writer.Flush()
			fd.Write(buff.Bytes())
			fd.Close()
		} else {
			log.Warn("Could not create file to flush pprof", "err", err, "path", pprofFilePath)
		}
		if *blockProfile {
			blockProfileFilePath := filepath.Join(*blockGapFolder, fmt.Sprintf("%v.blockprof", hash.String()))
			blockProfile, err := os.Create(blockProfileFilePath)
			if err != nil {
				log.Warn("could not create block profile", "err", err, "path", blockProfileFilePath)
			}
			defer blockProfile.Close()
			if err := pprof.Lookup("block").WriteTo(blockProfile, 0); err != nil {
				log.Warn("could not write block profile", "err", err, "block", hash)
			}
		}
			
	}
	buff.Truncate(0)
	writer.Reset(buff)
	if err := pprof.StartCPUProfile(writer); err != nil {
		log.Warn("Error starting cpu profile", "err", err)
	}
	runtime.SetBlockProfileRate(0) // Reset block profile rate after the function returns.
	lastNewHead = time.Now()
}