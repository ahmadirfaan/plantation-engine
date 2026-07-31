package model

import "time"

type Estate struct {
	Id        string
	Name      *string
	Width     int
	Length    int
	ExtInfo   *string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type Block struct {
	Id        string
	EstateId  string
	BlockX    int
	BlockY    int
	ExtInfo   *string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type Tree struct {
	Id       string
	EstateId string
	BlockId  string
	XAxis    int
	YAxis    int
	Height   int
}

type EstateStats struct {
	EstateId           string
	MinHeightTree      int
	MaxHeightTree      int
	MedianHeightTree   float64
	SumTree            int
	TotalDistanceDrone int
	Width              int
	Length             int
	ExtInfo            *string
	CreatedAt          *time.Time
	UpdatedAt          *time.Time
}
