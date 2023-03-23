package models

type Secret struct {
	Name    string `json:"data1"`
	Data    string `json:"data2"`
	Meta    string `json:"data3"`
	Version int    `json:"version"`
	Type    int    `json:"type"`
	Uid     string `json:"uid"`
}
