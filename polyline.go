package main

import (
	"reflect"

	"github.com/mytaxi-uz/shape2osm/tags"
	"github.com/mytaxi-uz/shape2osm/util"
	"github.com/mytaxi-uz/shape2osm/util/osm"
	"github.com/mytaxi-uz/shape2osm/util/shp"
)

// road, river

func convertPolylineToOSMWay(shapeReader *shp.Reader, shapeType string) {
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

		var osmTags osm.Tags

		switch shapeType {
		case "road":
			osmTags = tags.RoadAttrToOSMTag(num, fields, shapeReader)
		case "river":
			osmTags = tags.RiverAttrToOSMTag(num, fields, shapeReader)
		case "railway":
			osmTags = tags.RailwayAttrToOSMTag(num, fields, shapeReader)
		}

		way := osm.Way{
			ID:        osmID,
			Nodes:     wayNodes,
			Tags:      osmTags,
			Version:   1,
			Timestamp: nowTime,
		}
		osmOut.Ways = append(osmOut.Ways, &way)
	}
}
