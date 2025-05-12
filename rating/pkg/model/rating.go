package model

type RecordId string

type RecordType string

const (
	RecordTypeMovie  = RecordType("movie")
	RecordTypeSeries = RecordType("series")
)

type UserID string

type RatingValue int

type Rating struct {
	RecordID   RecordId    `json:"recordId"`
	RecordType RecordType  `json:"recordType"`
	UserID     UserID      `json:"userId"`
	Value      RatingValue `json:"value"`
}
