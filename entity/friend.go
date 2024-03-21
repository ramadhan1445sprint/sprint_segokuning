package entity

type FriendPair struct {
	FriendPairId string `db:"id"`
	FriendIdA    string `db:"user_id1"`
	FriendIdZ    string `db:"user_id2"`
}

type AddDeleteFriendPayload struct {
	FriendId string `json:"userId"`
}
