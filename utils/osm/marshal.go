package osm

import (
	"encoding/xml"

	"shape2osm/utils/osm/internal/osmpb"
)

const locMultiple = 10000000.0

var memberTypeMap = map[Type]osmpb.Relation_MemberType{
	TypeNode:     osmpb.Relation_NODE,
	TypeWay:      osmpb.Relation_WAY,
	TypeRelation: osmpb.Relation_RELATION,
}

// xmlNameJSONTypeNode is kind of a hack to encode the proper json
// object type attribute for this struct type.
type xmlNameJSONTypeNode xml.Name

// xmlNameJSONTypeWay is kind of a hack to encode the proper json
// object type attribute for this struct type.
type xmlNameJSONTypeWay xml.Name

// xmlNameJSONTypeRel is kind of a hack to encode the proper json
// object type attribute for this struct type.
type xmlNameJSONTypeRel xml.Name

func marshalNodes(nodes Nodes, ss *stringSet, includeChangeset bool) *osmpb.DenseNodes {
	dense := denseNodesValues(nodes)
	encoded := &osmpb.DenseNodes{
		Ids: encodeInt64(dense.IDs),
		/*
			DenseInfo: &osmpb.DenseInfo{
				Versions:   dense.Versions,
				Timestamps: encodeInt64(dense.Timestamps),
				Committeds: encodeInt64(dense.Committeds),
				Visibles:   dense.Visibles,
			},
		*/
		Lats: encodeInt64(dense.Lats),
		Lons: encodeInt64(dense.Lons),
	}

	if dense.TagCount > 0 {
		encoded.KeysVals = encodeNodesTags(nodes, ss, dense.TagCount)
	}
	/*
		if includeChangeset {
			csinfo := nodesChangesetInfo(nodes, ss)
			encoded.DenseInfo.ChangesetIds = encodeInt64(csinfo.Changesets)
			encoded.DenseInfo.UserIds = encodeInt32(csinfo.UserIDs)
			encoded.DenseInfo.UserSids = encodeInt32(csinfo.UserSids)
		}
	*/
	return encoded
}

func marshalWay(way *Way, ss *stringSet, includeChangeset bool) *osmpb.Way {
	keys, vals := way.Tags.keyValues(ss)
	encoded := &osmpb.Way{
		Id:   int64(way.ID),
		Keys: keys,
		Vals: vals,
		/*
			Info: &osmpb.Info{
				Version:   int32(way.Version),
				Timestamp: timeToUnix(way.Timestamp),
				Visible:   proto.Bool(way.Visible),
			},
			Updates: marshalUpdates(way.Updates),
		*/
	}
	/*
		if way.Committed != nil {
			encoded.Info.Committed = timeToUnixPointer(*way.Committed)
		}
	*/
	if len(way.Nodes) > 0 {
		encoded.Refs = encodeWayNodeIDs(way.Nodes)
		/*
			if way.Nodes[0].Version != 0 {
				encoded.DenseMembers = encodeDenseWayNodes(way.Nodes)
			}
		*/
	}
	/*
		if includeChangeset {
			encoded.Info.ChangesetId = int64(way.ChangesetID)
			encoded.Info.UserId = int32(way.UserID)
			encoded.Info.UserSid = ss.Add(way.User)
		}
	*/
	return encoded
}

func marshalRelation(relation *Relation, ss *stringSet, includeChangeset bool) *osmpb.Relation {
	l := len(relation.Members)
	roles := make([]uint32, l)
	refs := make([]int64, l)
	types := make([]osmpb.Relation_MemberType, l)

	// interestingMember := false
	for i, m := range relation.Members {
		roles[i] = ss.Add(m.Role)
		refs[i] = m.Ref
		types[i] = memberTypeMap[m.Type]
		/*
			if m.Version != 0 {
				interestingMember = true
			}
		*/
	}

	keys, vals := relation.Tags.keyValues(ss)
	encoded := &osmpb.Relation{
		Id:   int64(relation.ID),
		Keys: keys,
		Vals: vals,
		/*
			Info: &osmpb.Info{
				Version:   int32(relation.Version),
				Timestamp: timeToUnix(relation.Timestamp),
				Visible:   proto.Bool(relation.Visible),
			},
		*/
		// Roles: roles,
		// Refs:  encodeInt64(refs),
		Types: types,
		// Updates: marshalUpdates(relation.Updates),
	}
	/*
			if relation.Committed != nil {
				encoded.Info.Committed = timeToUnixPointer(*relation.Committed)
			}

		if interestingMember {
			// relations can be partial annotated, in that case we still
			// want to save the annotation data.
			encoded.DenseMembers = encodeDenseMembers(relation.Members)
		}
			if includeChangeset {
				encoded.Info.ChangesetId = int64(relation.ChangesetID)
				encoded.Info.UserId = int32(relation.UserID)
				encoded.Info.UserSid = ss.Add(relation.User)
			}
	*/
	return encoded
}

type denseNodesResult struct {
	IDs  []int64
	Lats []int64
	Lons []int64
	// Timestamps []int64
	// Committeds []int64
	// Versions   []int32
	// Visibles   []bool
	TagCount int
}

func denseNodesValues(ns Nodes) denseNodesResult {
	l := len(ns)
	ds := denseNodesResult{
		IDs:  make([]int64, l),
		Lats: make([]int64, l),
		Lons: make([]int64, l),
		// Timestamps: make([]int64, l),
		// Committeds: make([]int64, l),
		// Versions:   make([]int32, l),
		// Visibles:   make([]bool, l),
	}

	// cc := 0
	for i, n := range ns {
		ds.IDs[i] = int64(n.ID)
		ds.Lats[i] = geoToInt64(n.Lat)
		ds.Lons[i] = geoToInt64(n.Lon)
		// ds.Timestamps[i] = n.Timestamp.Unix()
		// ds.Versions[i] = int32(n.Version)
		// ds.Visibles[i] = n.Visible
		ds.TagCount += len(n.Tags)
		/*
			if n.Committed != nil {
				ds.Committeds[i] = timeToUnix(*n.Committed)
				cc++
			}
		*/
	}
	/*
		if cc == 0 {
			ds.Committeds = nil
		}
	*/
	return ds
}

func encodeNodesTags(ns Nodes, ss *stringSet, count int) []uint32 {
	r := make([]uint32, 0, 2*count+len(ns))
	for _, n := range ns {
		for _, t := range n.Tags {
			r = append(r, ss.Add(t.Key))
			r = append(r, ss.Add(t.Value))
		}
		r = append(r, 0)
	}

	return r
}

func encodeWayNodeIDs(waynodes WayNodes) []int64 {
	result := make([]int64, len(waynodes))
	var prev int64

	for i, r := range waynodes {
		result[i] = int64(r.ID) - prev
		prev = int64(r.ID)
	}

	return result
}

/*
func encodeDenseMembers(members Members) *osmpb.DenseMembers {
	l := len(members)
	versions := make([]int32, l)
	changesetIDs := make([]int64, l)
	orientations := make([]int32, l)
	lats := make([]int64, l)
	lons := make([]int64, l)

	locCount := 0
	orientCount := 0
	for i, m := range members {
		if m.Lat != 0 || m.Lon != 0 {
			locCount++
		}

		lats[i] = geoToInt64(m.Lat)
		lons[i] = geoToInt64(m.Lon)

		versions[i] = int32(m.Version)
		// changesetIDs[i] = int64(m.ChangesetID)

		if m.Orientation != 0 {
			orientations[i] = int32(m.Orientation)
			orientCount++
		}
	}

	result := &osmpb.DenseMembers{
		Versions:     versions,
		ChangesetIds: encodeInt64(changesetIDs),
	}

	if locCount > 0 {
		result.Lats = encodeInt64(lats)
		result.Lons = encodeInt64(lons)
	}

	if orientCount > 0 {
		result.Orientation = orientations
	}

	return result
}
*/
func encodeInt64(vals []int64) []int64 {
	var prev int64
	for i, v := range vals {
		vals[i] = v - prev
		prev = v
	}

	return vals
}

func geoToInt64(l float64) int64 {
	// on rounding errors
	//
	// It is the case that 32.850314 * 10e6 = 32850313.999999996
	// Simpily casting this as an int will truncate towards zero
	// and result in an off by one. The true solution is to round
	// the scaled result, like so:
	//
	// int64(math.Floor(stream.BaseData[i][0]*factor + 0.5))
	//
	// However, the code below does the same thing in this context,
	// and is twice as fast:
	sign := 0.5
	if l < 0 {
		sign = -0.5
	}

	return int64(l*locMultiple + sign)
}
