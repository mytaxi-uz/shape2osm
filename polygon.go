package main

import (
	"reflect"

	"github.com/mytaxi-uz/shape2osm/tags"
	"github.com/mytaxi-uz/shape2osm/util"
	"github.com/mytaxi-uz/shape2osm/util/osm"
	"github.com/mytaxi-uz/shape2osm/util/shp"
)

// building, landuse, water

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

		var osmTags osm.Tags

		switch shapeType {
		case "building":
			osmTags = tags.BuildingAttrToOSMTag(num, fields, shapeReader)
		case "landuse":
			osmTags = tags.LanduseAttrToOSMTag(num, fields, shapeReader)
		case "water":
			osmTags = tags.WaterAttrToOSMTag(num, fields, shapeReader)
		}

		osmOut.Ways = append(osmOut.Ways, ways...)

		if len(polygon.Parts) == 1 {
			ways[0].Tags = osmTags
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

			osmTags = append(osmTags, tag)

			relation := osm.Relation{
				ID:        osmID,
				Tags:      osmTags,
				Members:   members,
				Version:   1,
				Timestamp: nowTime,
			}

			osmOut.Relations = append(osmOut.Relations, &relation)
		}
	}
}
