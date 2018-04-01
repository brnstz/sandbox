package main

import (
	"image/png"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	influxdb "github.com/influxdata/influxdb/client/v2"
)

func main() {
	frame := 0
	iflx, err := influxdb.NewHTTPClient(
		influxdb.HTTPConfig{
			Addr: "http://localhost:8086",
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	t := time.Now()

	for true {

		img, err := png.Decode(os.Stdin)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		bounds := img.Bounds()
		for x := 0; x < bounds.Dx(); x++ {

			bp, err := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
				Database: "movies",
			})
			if err != nil {
				log.Fatal(err)
			}

			for y := 0; y < bounds.Dy(); y++ {
				tags := map[string]string{
					"x":     strconv.Itoa(x),
					"y":     strconv.Itoa(y),
					"frame": strconv.Itoa(frame),
				}

				r, g, b, a := img.At(x, y).RGBA()
				fields := map[string]interface{}{
					"r": r,
					"g": g,
					"b": b,
					"a": a,
				}

				log.Println(x, y, t)
				pt, err := influxdb.NewPoint(
					"movie_color", tags, fields, t,
				)
				if err != nil {
					log.Fatal(err)
				}
				bp.AddPoint(pt)
			}

			err = iflx.Write(bp)
			if err != nil {
				log.Fatal(err)
			}
		}

		frame++
		t = t.Add(time.Second * 1)
	}

}
