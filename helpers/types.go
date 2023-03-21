package helpers

type PlayerID int64

type Button struct {
	Text string
	CallbackData string
}

type Response struct {
	Text string
	Keyboard [][]Button
	ChatsID []int64
}