package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/bamstats"
	"github.com/spf13/cobra"
)

var (
	bam, annotation, loglevel, output string
	cpu, maxBuf, reads                int
)

func run(cmd *cobra.Command, args []string) (err error) {
	err = nil

	if version, e := cmd.Flags().GetBool("version"); e == nil && version {
		fmt.Printf("bamstats version %s\n", bamstats.Version())
		return
	}

	// Set loglevel
	level, err := log.ParseLevel(loglevel)
	if err != nil {
		return
	}
	log.SetLevel(level)

	// Get stats
	log.Infof("Running %s", cmd.Use)
	stats, err := bamstats.Process(bam, annotation, cpu, maxBuf, reads)
	if err != nil {
		return
	}

	out := bamstats.NewOutput(output)
	bamstats.OutputJson(out, stats)

	return
}

func setBamstatsFlags(c *cobra.Command) {
	c.PersistentFlags().StringVarP(&bam, "bam", "b", "", "input file")
	c.PersistentFlags().StringVarP(&annotation, "annotaion", "a", "", "element annotation file")
	c.PersistentFlags().StringVarP(&loglevel, "loglevel", "", "warn", "logging level")
	c.PersistentFlags().StringVarP(&output, "output", "o", "-", "output file")
	c.PersistentFlags().IntVarP(&cpu, "cpu", "c", 1, "number of cpus to be used")
	c.PersistentFlags().IntVarP(&maxBuf, "max-buf", "", 1000000, "maximum number of buffered records")
	c.PersistentFlags().IntVarP(&reads, "reads", "n", -1, "number of records to process")
	c.PersistentFlags().Bool("version", false, "show version and exit")
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "bamstats",
		Short: "Compute mapping statistics",
		RunE:  run,
	}

	setBamstatsFlags(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Debug(err)
	}
}