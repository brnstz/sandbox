package main

import (
	"flag"
	"log"
	"os"

	"github.com/go-audio/wav"
	"github.com/mattetti/audio/riff"
)

func processChunk(d *wav.Decoder, chunk *riff.Chunk) error {
	var (
		n   int
		err error
	)

	p := make([]byte, d.BitDepth/8)

	for !chunk.IsFullyRead() {
		n, err = chunk.Read(p)
		if err != nil {
			return err
		}
		log.Printf("%v %v %v", d.BitDepth, n, p)
	}

	return nil
}

func checkBeats(input string) (int, error) {
	var (
		chunk *riff.Chunk
		err   error
		i     int
	)

	fh, err := os.Open(input)
	if err != nil {
		return 0, err
	}

	d := wav.NewDecoder(fh)
	err = d.FwdToPCM()
	if err != nil {
		return 0, err
	}

	for !d.EOF() {

		chunk, err = d.NextChunk()
		if err != nil {
			if d.EOF() {
				break
			} else {
				return 0, err
			}
		}

		processChunk(d, chunk)

		log.Printf("%v %v", i, chunk)
		chunk.Done()
		i++
	}

	return i, nil
}

func main() {
	var input = flag.String("input", "", "path to input file")
	flag.Parse()

	if len(*input) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	bpm, err := checkBeats(*input)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("file %v is %v bpm", input, bpm)
}
