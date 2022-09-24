package main

import (
	"reflect"
	"strings"

	"github.com/mytaxi-uz/shape2osm/utils/osm"
	"github.com/mytaxi-uz/shape2osm/utils/shp"
)

func convertPolygonToOSMWay(shapeReader *shp.Reader) {
	// fields from the attribute table (DBF)
	fields := shapeReader.Fields()

	t := reflect.TypeOf(&shp.Polygon{}).Elem()

	// loop through all features in the shapefile
	for shapeReader.Next() {
		num, p := shapeReader.Shape()
		if reflect.TypeOf(p).Elem() != t {
			continue
		}
		polygon := p.(*shp.Polygon)

		var wayNodes osm.WayNodes

		for _, point := range polygon.Points {
			lat := truncateFloat64(point.Y)
			lon := truncateFloat64(point.X)
			nodeID, ok := nodesIDMap[[2]float64{lat, lon}]
			if !ok {
				osmID++
				nodeID = osmID
				node := osm.Node{
					ID:        osmID,
					Lat:       lat,
					Lon:       lon,
					Version:   1,
					Timestamp: nowTime,
				}
				nodesIDMap[[2]float64{lat, lon}] = osmID
				osmOut.Nodes = append(osmOut.Nodes, &node)
				if osmOut.Bounds.MinLat > lat {
					osmOut.Bounds.MinLat = lat
				}
				if osmOut.Bounds.MaxLat < lat {
					osmOut.Bounds.MaxLat = lat
				}
				if osmOut.Bounds.MinLon > lon {
					osmOut.Bounds.MinLon = lon
				}
				if osmOut.Bounds.MaxLon < lon {
					osmOut.Bounds.MaxLon = lon
				}
			}
			wayNodes = append(wayNodes, osm.WayNode{ID: nodeID})
		}

		osmID++

		tags := convertPolygonAttrToOSMTag(num, fields, shapeReader)

		way := osm.Way{
			ID:        osmID,
			Nodes:     wayNodes,
			Tags:      tags,
			Version:   1,
			Timestamp: nowTime,
		}
		osmOut.Ways = append(osmOut.Ways, &way)
	}
}

func convertPolygonAttrToOSMTag(num int, fields []shp.Field, reader *shp.Reader) (tags osm.Tags) {
	var key, value string

	for k, f := range fields {
		attr := reader.ReadAttribute(num, k)

		if attr == "" {
			continue
		}

		key = ""
		field := strings.ToUpper(f.String())

		switch field {
		/*
			case "ID", "id":
				key = "id"
				for i, c := range attr {
					if c == '.' {
						attr = attr[:i]
						break
					}
				}
				value = attr
		*/
		case "TYP_COD":
			switch attr {
			case "103", "111":
				key = "leisure"
				value = "park"
			}
		case "STOREY":
			tag := osm.Tag{
				Key:   "building",
				Value: "yes",
			}
			tags = append(tags, tag)
			key = "building:levels"
			value = attr
		case "NAME_UZ":
			key = "name:uz"
			value = attr
		case "NAME", "NAME_RU":
			key = "name:ru"
			value = attr
		case "NAME_EN":
			key = "name:en"
			value = attr
		}

		if key != "" {
			tag := osm.Tag{
				Key:   key,
				Value: value,
			}
			tags = append(tags, tag)
		}
	}

	return
}
