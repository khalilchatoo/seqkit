// Copyright © 2016 Wei Shen <shenwei356@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"math/rand"

	"github.com/brentp/xopen"
	"github.com/shenwei356/bio/seqio/fasta"
	"github.com/shenwei356/util/randutil"
	"github.com/spf13/cobra"
)

// shuffleCmd represents the seq command
var shuffleCmd = &cobra.Command{
	Use:   "shuffle",
	Short: "shuffle sequences",
	Long: `shuffle sequences.

`,
	Run: func(cmd *cobra.Command, args []string) {
		alphabet := getAlphabet(cmd, "seq-type")
		idRegexp := getFlagString(cmd, "id-regexp")
		chunkSize := getFlagInt(cmd, "chunk-size")
		threads := getFlagInt(cmd, "threads")
		lineWidth := getFlagInt(cmd, "line-width")
		outFile := getFlagString(cmd, "out-file")
		quiet := getFlagBool(cmd, "quiet")

		files := getFileList(args)

		seed := getFlagInt64(cmd, "rand-seed")

		outfh, err := xopen.Wopen(outFile)
		checkError(err)
		defer outfh.Close()

		sequences := make(map[string]*fasta.FastaRecord)
		index2name := make(map[int]string)

		if !quiet {
			log.Infof("read sequences ...")
		}
		i := 0
		for _, file := range files {
			fastaReader, err := fasta.NewFastaReader(alphabet, file, chunkSize, threads, idRegexp)
			checkError(err)
			for chunk := range fastaReader.Ch {
				checkError(chunk.Err)

				for _, record := range chunk.Data {
					sequences[string(record.Name)] = record
					index2name[i] = string(record.Name)
					i++
				}
			}
		}

		if !quiet {
			log.Infof("%d sequences loaded", len(sequences))
			log.Infof("shuffle ...")
		}
		rand.Seed(seed)
		indices := make([]int, len(index2name))
		for i := 0; i < len(index2name); i++ {
			indices[i] = i
		}
		randutil.Shuffle(indices)

		if !quiet {
			log.Infof("output ...")
		}
		var record *fasta.FastaRecord
		for _, i := range indices {
			record = sequences[index2name[i]]
			outfh.WriteString(fmt.Sprintf(">%s\n%s\n", record.Name, record.FormatSeq(lineWidth)))
		}
	},
}

func init() {
	RootCmd.AddCommand(shuffleCmd)
	shuffleCmd.Flags().Int64P("rand-seed", "s", 23, "rand seed for shuffle")
}
