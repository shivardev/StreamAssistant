package utils

import (
	"fmt"
	"sort"
)

type Leader struct {
	UserId   string `json:"userId"`
	UserName string `json:"userName"`
	Points   int    `json:"points"`
}

var LeaderBoard = make(map[string]Leader)

func GetTopTwoLeaderBoard() []Leader {
	var sortedLeaderBoard []Leader
	for _, user := range LeaderBoard {
		sortedLeaderBoard = append(sortedLeaderBoard, user)
	}

	sort.Slice(sortedLeaderBoard, func(i, j int) bool {
		return sortedLeaderBoard[i].Points > sortedLeaderBoard[j].Points
	})

	// Return top 2 users
	if len(sortedLeaderBoard) > 2 {
		return sortedLeaderBoard[:2]
	}
	fmt.Println(sortedLeaderBoard)
	return sortedLeaderBoard
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
	fmt.Println(LeaderBoard)
}

func ResetLeaderBoard() {
	LeaderBoard = make(map[string]Leader)
}
