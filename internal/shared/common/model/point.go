package model

// type PointType int32

// const (
// 	NONE_POINT = iota
// 	LINE_POINT
// 	SMALL_TREE_POINT
// )

// var (
// 	pointMapper = make(map[PointType]*Point)
// )

// func init() {
// 	pointMapper = map[PointType]*Point{
// 		NONE_POINT: {
// 			PointType: NONE_POINT,
// 			PointName: "無",
// 		},
// 		LINE_POINT: {
// 			PointType: LINE_POINT,
// 			PointName: "LINE POINT點數",
// 		},
// 		SMALL_TREE_POINT: {
// 			PointType: SMALL_TREE_POINT,
// 			PointName: "小樹點",
// 		},
// 	}
// }

// type Point struct {
// 	PointType PointType `json:"pointType"`
// 	PointName string    `json:"pointName"`
// }

// func GetAllPointTypes() []*Point {

// 	points := []*Point{}

// 	for _, v := range pointMapper {
// 		points = append(points, v)
// 	}
// 	return points
// }

// func GetPointType(pointType int32) (*Point, error) {

// 	point, ok := pointMapper[PointType(pointType)]

// 	if !ok {
// 		return nil, errors.New("Cannot find point type")
// 	}
// 	return point, nil
// }
