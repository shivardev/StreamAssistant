package utils

import (
	"sort"
)

type Leader struct {
	UserId   string `json:"userId"`
	UserName string `json:"userName"`
	Points   int    `json:"points"`
}

var LeaderBoard = make(map[string]Leader)

func GetTopLeaderBoard(isVerified bool) []Leader {
	var sortedLeaderBoard []Leader
	for _, user := range LeaderBoard {
		sortedLeaderBoard = append(sortedLeaderBoard, user)
	}

	sort.Slice(sortedLeaderBoard, func(i, j int) bool {
		return sortedLeaderBoard[i].Points > sortedLeaderBoard[j].Points
	})

	// Return top 2 users
	if isVerified {
		return sortedLeaderBoard[:2]
	} else {
		return sortedLeaderBoard[:4]
	}
}

func SetLeaderBoard(user User) {
	if existingUser, exists := LeaderBoard[user.UserId]; exists {
		if user.Points > existingUser.Points {
			newLeader := Leader{
				UserId:   existingUser.UserId,
				UserName: existingUser.UserName,
				Points:   existingUser.Points + 1,
			}
			LeaderBoard[user.UserId] = newLeader
		}
	} else {
		newLeader := Leader{
			UserId:   user.UserId,
			UserName: user.UserName,
			Points:   1,
		}
		LeaderBoard[user.UserId] = newLeader
	}
}

func ResetLeaderBoard() {
	LeaderBoard = make(map[string]Leader)
}
