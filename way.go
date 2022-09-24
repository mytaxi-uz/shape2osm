package main

import (
	"reflect"
	"strings"

	"shape2osm/utils/osm"
	"shape2osm/utils/shp"
)

func convertPolylineToOSMWay(shapeReader *shp.Reader) {
	// fields from the attribute table (DBF)
	fields := shapeReader.Fields()

	t := reflect.TypeOf(&shp.PolyLine{}).Elem()

	// loop through all features in the shapefile
	for shapeReader.Next() {
		num, p := shapeReader.Shape()
		if reflect.TypeOf(p).Elem() != t {
			continue
		}
		polyline := p.(*shp.PolyLine)

		var wayNodes osm.WayNodes

		for _, point := range polyline.Points {
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

		tags := convertPolylineAttrToOSMTag(num, fields, shapeReader)

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

func convertPolylineAttrToOSMTag(num int, fields []shp.Field, reader *shp.Reader) (tags osm.Tags) {
	var key, value string
	typCod := false
	for _, f := range fields {
		field := strings.ToUpper(f.String())
		if field == "TYP_COD" {
			typCod = true
			break
		}
	}

	if !typCod {
		tag := osm.Tag{
			Key:   "highway",
			Value: "road",
		}
		tags = append(tags, tag)
	}

	for k, f := range fields {
		attr := reader.ReadAttribute(num, k)

		if attr == "" {
			continue
		}

		key = ""
		field := strings.ToUpper(f.String())

		switch field {
		/*
			case "ID":
				key = "id"
				for i, c := range attr {
					if c == '.' {
						attr = attr[:i]
						break
					}
				}
				value = attr
		*/
		case "NAME_UZ":
			key = "name"
			value = attr
		case "NAME", "NAME_RU":
			key = "name:ru"
			value = attr
		case "NAME_EN":
			key = "name:en"
			value = attr
		case "SPEED":
			key = "additional:maxspeed"
			value = attr
		case "DIRECTION":
			if attr == "FT" {
				key = "oneway"
				value = "yes"
			}
		case "TUNNEL":
			if attr == "1" {
				key = "tunnel"
				value = "yes"
			}
		case "TYP_COD":
			key = "highway"
			switch attr {
			case "34", "30":
				value = "trunk"
			case "31", "56", "453", "42":
				value = "trunk_link"
			case "49", "54":
				value = "primary"
			case "55", "53":
				value = "secondary"
			case "37", "51":
				value = "tertiary"
			case "46", "45", "52", "50":
				value = "residential"
			case "32", "33":
				value = "pedestrian"
			case "35", "38", "36":
				value = "track"
			case "39", "59", "146":
				value = "footway"
			case "43", "40":
				value = "service"
			default:
				value = "road"
			}
			/*
				case "CLASS":
					key = "highway"
					switch attr {
					case "1":
						value = "motorway"
					case "2":
						value = "trunk"
					case "3":
						value = "primary"
					case "4":
						value = "secondary"
					case "5":
						value = "tertiary"
					case "6":
						value = "unclassified"
					case "7":
						value = "residential"
					case "8":
						value = "living_street"
					}
			*/
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
