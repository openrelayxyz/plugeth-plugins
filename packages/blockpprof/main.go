package main

import (
	"bufio"
	"bytes"
	"math/big"
	"os"
	"time"
	"github.com/openrelayxyz/plugeth-utils/core"
	"runtime/pprof"
	"path/filepath"
	"flag"
)

var (
	log core.Logger
	buff *bytes.Buffer
	writer *bufio.Writer
	lastNewHead time.Time
	Flags = *flag.NewFlagSet("blockpprof", flag.ContinueOnError)
	blockGapFlag = Flags.Int("pprof.block.gap", 14, "The number of seconds between blocks that should trigger a pprof flush")
	blockGapFolder = Flags.String("pprof.folder", "", "Directory where pprof flushes should be written")
	blockGapDuration = 12 * time.Second

)

func Initialize(ctx core.Context, loader core.PluginLoader, logger core.Logger) {
	log = logger
	blockGapDuration = time.Duration(*blockGapFlag) * time.Second
}

func NewHead(blockBytes []byte, hash core.Hash, logsBytes [][]byte, td *big.Int) {
	pprof.StopCPUProfile()
	if buff != nil && time.Since(lastNewHead) > blockGapDuration {
		if fd, err := os.Create(filepath.Join(*blockGapFolder, hash.String())); err == nil {
			writer.Flush()
			fd.Write(buff.Bytes())
			fd.Close()
		} else {
			log.Warn("Could not create file to flush pprof")
		}
	}
	buff.Truncate(0)
	writer.Reset(buff)
	if err := pprof.StartCPUProfile(writer); err != nil {
		log.Warn("Error starting cpu profile", "err", err)
	}
}