package main

import (
	"reflect"

	"github.com/mytaxi-uz/shape2osm/tags"
	"github.com/mytaxi-uz/shape2osm/util"
	"github.com/mytaxi-uz/shape2osm/util/osm"
	"github.com/mytaxi-uz/shape2osm/util/shp"
)

// poi, place, entrance

func convertPointToOSMNode(shapeReader *shp.Reader, shapeType string) {
	// fields from the attribute table (DBF)
	fields := shapeReader.Fields()

	t := reflect.TypeOf(&shp.Point{}).Elem()

	// loop through all features in the shapefile
	for shapeReader.Next() {
		num, p := shapeReader.Shape()
		if reflect.TypeOf(p).Elem() != t {
			continue
		}
		point := p.(*shp.Point)
		osmID++
		var osmTags osm.Tags

		switch shapeType {
		case "poi":
			osmTags = tags.PoiAttrToOSMTag(num, fields, shapeReader)
		case "place":
			osmTags = tags.PlaceAttrToOSMTag(num, fields, shapeReader)
		case "entrance":
			osmTags = tags.EntranceAttrToOSMTag(num, fields, shapeReader)
		}
		node := osm.Node{
			ID:        osmID,
			Lat:       util.TruncateFloat64(point.Y),
			Lon:       util.TruncateFloat64(point.X),
			Tags:      osmTags,
			Version:   1,
			Timestamp: nowTime,
		}
		osmOut.Nodes = append(osmOut.Nodes, &node)
		if osmOut.Bounds.MinLat > point.Y {
			osmOut.Bounds.MinLat = point.Y
		}
		if osmOut.Bounds.MaxLat < point.Y {
			osmOut.Bounds.MaxLat = point.Y
		}
		if osmOut.Bounds.MinLon > point.X {
			osmOut.Bounds.MinLon = point.X
		}
		if osmOut.Bounds.MaxLon < point.X {
			osmOut.Bounds.MaxLon = point.X
		}
	}
}
