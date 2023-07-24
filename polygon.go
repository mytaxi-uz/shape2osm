package main

import (
	"reflect"
	"strings"

	"github.com/mytaxi-uz/shape2osm/util"
	"github.com/mytaxi-uz/shape2osm/util/osm"
	"github.com/mytaxi-uz/shape2osm/util/shp"
)

// house, forest, water

func convertPolygonToOSMWay(shapeReader *shp.Reader, shapeType string) {
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

		var ways []*osm.Way

		for i, part := range polygon.Parts {

			var wayNodes osm.WayNodes

			end := 0

			if i+1 < int(polygon.NumParts) {
				end = int(polygon.Parts[i+1])
			} else {
				end = int(polygon.NumPoints)
			}
			points := polygon.Points[part:end]

			for _, point := range points {
				lat := util.TruncateFloat64(point.Y)
				lon := util.TruncateFloat64(point.X)
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

			way := osm.Way{
				ID:        osmID,
				Nodes:     wayNodes,
				Version:   1,
				Timestamp: nowTime,
			}
			ways = append(ways, &way)
		}

		var tags osm.Tags

		switch shapeType {
		case "building":
			tags = convertBuildingAttrToOSMTag(num, fields, shapeReader)
		case "landuse":
			tags = convertLanduseAttrToOSMTag(num, fields, shapeReader)
		case "water":
			tags = convertWaterAttrToOSMTag(num, fields, shapeReader)
		case "place_a":
			tags = convertPolygonPlaceAttrToOSMTag(num, fields, shapeReader)
		}

		osmOut.Ways = append(osmOut.Ways, ways...)

		if len(polygon.Parts) == 1 {
			ways[0].Tags = tags
		} else {
			osmID++
			var members []osm.Member

			for _, way := range ways {
				member := osm.Member{
					Type: "way",
					Ref:  way.ID,
					Role: "inner",
				}
				members = append(members, member)
			}

			members[0].Role = "outer"

			tag := osm.Tag{
				Key:   "type",
				Value: "multipolygon",
			}

			tags = append(tags, tag)

			relation := osm.Relation{
				ID:        osmID,
				Tags:      tags,
				Members:   members,
				Version:   1,
				Timestamp: nowTime,
			}

			osmOut.Relations = append(osmOut.Relations, &relation)
		}
	}
}

func convertBuildingAttrToOSMTag(num int, fields []shp.Field, reader *shp.Reader) (tags osm.Tags) {

	tag := osm.Tag{
		Key:   "building",
		Value: "yes",
	}

	tags = append(tags, tag)

	for k, f := range fields {
		attr := reader.ReadAttribute(num, k)

		if attr == "" {
			continue
		}

		key := ""
		value := ""
		field := strings.ToUpper(f.String())

		switch field {
		case "STOREY":
			key = "building:levels"
			value = strings.Split(attr, ".")[0]
		case "HOUSE_NUM":
			key = "addr:housenumber"
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

func convertLanduseAttrToOSMTag(num int, fields []shp.Field, reader *shp.Reader) (tags osm.Tags) {

	for k, f := range fields {
		attr := reader.ReadAttribute(num, k)

		if attr == "" {
			continue
		}

		key := ""
		value := ""
		field := strings.ToUpper(f.String())

		switch field {
		case "TYP_COD":
			switch attr {
			case "103", "111":
				key = "leisure"
				value = "park"
			case "92":
				key = "landuse"
				value = "grass"
			case "340", "349", "350":
				key = "amenity"
				value = "hospital"
			case "207", "213", "216", "201":
				key = "landuse"
				value = "university"
			default:
				key = ""
			}
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

func convertWaterAttrToOSMTag(num int, fields []shp.Field, reader *shp.Reader) (tags osm.Tags) {
	var key, value string

	for k, f := range fields {
		attr := reader.ReadAttribute(num, k)

		if attr == "" {
			continue
		}

		key = ""
		field := strings.ToUpper(f.String())

		switch field {
		case "TYP_COD":
			key = "water"
			switch attr {
			case "65":
				value = "canal"
			case "74":
				value = "river"
			case "75":
				value = "stream"
			case "271":
				key = "leisure"
				value = "swimming_pool"
			default:
				key = ""
			}
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

func convertPolygonPlaceAttrToOSMTag(num int, fields []shp.Field, reader *shp.Reader) (tags osm.Tags) {
	for k, f := range fields {
		attr := reader.ReadAttribute(num, k)

		if attr == "" {
			continue
		}

		key := ""
		value := ""
		field := strings.ToUpper(f.String())

		switch field {
		case "TYP_COD":
			switch attr {
			case "13":
				tag := osm.Tag{
					Key:   "boundary",
					Value: "administrative",
				}

				tags = append(tags, tag)

				tag = osm.Tag{
					Key:   "admin_level",
					Value: "2",
				}

				tags = append(tags, tag)
			case "9":
				tag := osm.Tag{
					Key:   "boundary",
					Value: "administrative",
				}

				tags = append(tags, tag)

				tag = osm.Tag{
					Key:   "admin_level",
					Value: "4",
				}

				tags = append(tags, tag)
			}

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
