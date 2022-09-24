package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/mytaxi-uz/shape2osm/utils/osm"
	"github.com/mytaxi-uz/shape2osm/utils/shp"
)

var (
	osmOut     osm.OSM
	osmID      int64 = 0
	nodesIDMap map[[2]float64]int64
	nowTime    string
)

func main() {

	osmOut = osm.OSM{
		Version:   0.6,
		Generator: "shape2osm by mytaxi.uz",
		Bounds:    &osm.Bounds{MinLat: 91, MinLon: 181},
	}

	nodesIDMap = make(map[[2]float64]int64)

	dirname := "./"

	fmt.Println("Starting shapefile to osm converter...")
	fmt.Println("Open shapefiles from directory:", dirname)

	startTime := time.Now().UTC().Round(time.Second)
	nowTime = startTime.Format("2006-01-02T15:04:05Z")

	f, err := os.Open(dirname)
	if err != nil {
		log.Fatal(err)
	}

	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}

	fileCount := 0

	for _, file := range files {
		name := file.Name()
		if !strings.HasSuffix(name, ".shp") {
			continue
		}
		fileCount++
		fmt.Println("Converting shapefile:", name)
		shapeReader, err := shp.Open(dirname + name)
		if err != nil {
			log.Fatal(err)
		}
		defer shapeReader.Close()

		switch shapeReader.GeometryType {
		case shp.POINT:
			convertPointToOSMNode(shapeReader)
		case shp.POLYLINE:
			convertPolylineToOSMWay(shapeReader)
		case shp.POLYGON:
			convertPolygonToOSMWay(shapeReader)
		}
	}

	if fileCount == 0 {
		fmt.Println("No shapefiles found from directory:", dirname)
		return
	}

	if len(osmOut.Nodes) == 0 {
		fmt.Println("Empty shapefiles reads from directory:", dirname)
		return
	}

	fmt.Println("Encode and write to file...")

	f, err = os.Create(dirname + "uzbekistan.osm")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err = f.Write([]byte(xml.Header))
	if err != nil {
		log.Fatal(err)
	}
	enc := xml.NewEncoder(f)
	enc.Indent("", "  ")
	start := xml.StartElement{}
	err = osmOut.MarshalXML(enc, start)
	if err != nil {
		log.Fatal(err)
	}

	enc.Flush()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Total OSM nodes writed:     ", len(osmOut.Nodes))
	fmt.Println("Total OSM ways writed:      ", len(osmOut.Ways))
	fmt.Println("Total OSM relations writed: ", len(osmOut.Relations))
	fmt.Println("Estimated time: ", time.Since(startTime))
}

// truncate float64 value with precision 7 (~1cm)
func truncateFloat64(num float64) float64 {
	const precesion = 10000000
	return float64(int(num*precesion+math.Copysign(0.5, num))) / precesion
}
