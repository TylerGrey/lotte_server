package model

// ReservationList 예약 리스트 조회 API 응답
type ReservationList struct {
	ID            int64  `json:"id"`
	RoomID        int64  `json:"room_id"`
	BookerID      int64  `json:"booker_id"`
	StartDatetime string `json:"start_datetime"`
	EndDatetime   string `json:"end_datetime"`
	Title         string `json:"title"`
	UserName      string `json:"user_name"`
}
