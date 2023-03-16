package main

import (
	"reflect"
	"strings"

	"github.com/mytaxi-uz/shape2osm/util"
	"github.com/mytaxi-uz/shape2osm/util/osm"
	"github.com/mytaxi-uz/shape2osm/util/shp"
)

// poi, place

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
		var tags osm.Tags

		switch shapeType {
		case "poi":
			tags = convertPoiAttrToOSMTag(num, fields, shapeReader)
		case "place":
			tags = convertPointPlaceAttrToOSMTag(num, fields, shapeReader)
		}
		node := osm.Node{
			ID:        osmID,
			Lat:       util.TruncateFloat64(point.Y),
			Lon:       util.TruncateFloat64(point.X),
			Tags:      tags,
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

func convertPoiAttrToOSMTag(num int, fields []shp.Field, reader *shp.Reader) (tags osm.Tags) {
	var key, value, str_type_uz, str_type_ru, str_type_en string

	for k, f := range fields {
		field := strings.ToUpper(f.String())
		switch field {
		case "TYP_STR_UZ":
			str_type_uz = reader.ReadAttribute(num, k)
		case "TYP_STR_RU":
			str_type_ru = reader.ReadAttribute(num, k)
		case "TYP_STR_EN":
			str_type_en = reader.ReadAttribute(num, k)
		}
	}

	for k, f := range fields {
		attr := reader.ReadAttribute(num, k)

		if attr == "" {
			continue
		}

		key = ""
		field := strings.ToUpper(f.String())

		switch field {
		case "NAME_UZ":
			key = "name:uz"
			value = attr
		case "NAME", "NAME_RU":
			key = "name:ru"
			value = attr
		case "NAME_EN":
			key = "name:en"
			value = attr
		case "STREET_UZ":
			key = "addr:street:uz"
			value = attr
			if str_type_uz != "" {
				value += " " + str_type_uz
			}
		case "STREET_RU":
			key = "addr:street:ru"
			value = attr
			if str_type_ru != "" {
				value += " " + str_type_ru
			}
		case "STREET_EN":
			key = "addr:street:en"
			value = attr
			if str_type_en != "" {
				value += " " + str_type_en
			}
		case "ADDRESS_UZ":
			key = "addr:housenumber"
			value = attr
		case "CLASS":
			key = "amenity"
			value = strings.ToLower(attr)
		case "ADDRESS":
			key = "addr:full"
			value = attr
		case "PHONE":
			key = "phone"
			value = attr
		case "WORKTIME":
			key = "opening_hours"
			value = attr
		case "HTTP":
			key = "website"
			value = attr
		case "EMAIL":
			key = "email"
			value = attr
		case "ZIPCODE":
			key = "addr:postcode"
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

func convertPointPlaceAttrToOSMTag(num int, fields []shp.Field, reader *shp.Reader) (tags osm.Tags) {
	var key, value, str_type_uz, str_type_ru, str_type_en string

	for k, f := range fields {
		field := strings.ToUpper(f.String())
		switch field {
		case "TYP_STR_UZ":
			str_type_uz = reader.ReadAttribute(num, k)
		case "TYP_STR_RU":
			str_type_ru = reader.ReadAttribute(num, k)
		case "TYP_STR_EN":
			str_type_en = reader.ReadAttribute(num, k)
		}
	}

	for k, f := range fields {
		attr := reader.ReadAttribute(num, k)

		if attr == "" {
			continue
		}

		key = ""
		field := strings.ToUpper(f.String())

		switch field {
		case "TYP_COD":
			key = "place"
			switch attr {
			case "20":
				tag := osm.Tag{
					Key:   "admin_level",
					Value: "6",
				}
				tags = append(tags, tag)
				value = "city"
			case "10":
				value = "town"
			case "11":
				value = "village"
			case "134", "130", "723":
				value = "neighbourhood"
			case "12":
				value = "district"
			case "21":
				value = "city"
			case "9":
				value = "region"
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
		case "STREET_UZ":
			key = "addr:street:uz"
			value = attr
			if str_type_uz != "" {
				value += " " + str_type_uz
			}
		case "STREET_RU":
			key = "addr:street:ru"
			value = attr
			if str_type_ru != "" {
				value += " " + str_type_ru
			}
		case "STREET_EN":
			key = "addr:street:en"
			value = attr
			if str_type_en != "" {
				value += " " + str_type_en
			}
		case "ADDRESS_UZ":
			key = "addr:housenumber"
			value = attr
		case "CLASS":
			key = "amenity"
			value = strings.ToLower(attr)
		case "ADDRESS":
			key = "addr:full"
			value = attr
		case "PHONE":
			key = "phone"
			value = attr
		case "WORKTIME":
			key = "opening_hours"
			value = attr
		case "HTTP":
			key = "website"
			value = attr
		case "EMAIL":
			key = "email"
			value = attr
		case "ZIPCODE":
			key = "addr:postcode"
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
