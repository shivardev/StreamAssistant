package utils

import "github.com/emirpasic/gods/sets/hashset"

var WelcomeMsgs = []string{
	"Hey, this dude is new to the hood 👉 %s! Welcome! 😎",
	"Look who just showed up! It's %s! Ready to roll? 🎉",
	"Hey %s! You seem new here... Don’t worry, we don’t bite! 😜",
	"Newbie alert! 🚨 Welcome to the madness, %s! 🥳",
	"Yo %s, fresh to the squad! Let the fun begin! 🎊",
	"Well, look who decided to join us! Welcome, %s! 🎉",
	"Guess who’s the new kid on the block? 😏 It’s %s!",
	"New face in the house! Welcome, %s! 😎",
	"Hey %s, looks like you're new here... Better buckle up! 🏎️",
	"Yo %s, you're officially part of the crew now! 🔥",
}

// var CommandHandler = map[string]string{
// 	"!cmds":    "list is available here 🔗 https://blazingbane.com/chatcmds",
// 	"!cmd":     "list is available here 🔗 https://blazingbane.com/chatcmds",
// 	"!discord": "join the discord 🔗 https://discord.com/invite/U5jDNUPceM",
// 	"!setup":   "setup is available here 🔗 https://blazingbane.com/setup",
// 	"!blog":    "Check out all blogs! - 🔗 https://blazingbane.com/blog",
// }

var SubscriberMsgs = []string{
	"Look who decided to show up! %s, back with %d points... Don’t get too cocky! 😜",
	"Well, well, well, %s! You’ve got %d points and somehow still haven’t improved! 😆",
	"Oh, it’s %s! Back with %d points... Are you sure you didn’t just get lucky? 😏",
	"Yo %s, back again with %d points? Guess the stream wasn’t boring enough for you! 😜",
	"Oh great, %s is back with %d points... Let’s see if you can actually do something useful this time! 😂",
	"Look who’s here, %s! With %d points, you still somehow manage to make it look easy... 🙄",
	"Hey %s, back again with your %d points... Don’t get too comfortable, you’re still not a pro. 😜",
	"Guess who’s back? %s with %d points! Let’s hope you actually do something impressive today! 😆",
	"Oh, it’s %s with %d points... You’ve been here before, but let’s see if you can do better this time! 😜",
	"Look who came crawling back! %s with %d points... Don’t let your ego get too big! 😂",
}
var IgnoredUsers = hashset.New("Nightbot", "YouTube", "Blazing Bane", "Relangi mama", "Streamlabs")
var CongratsMessages = []string{
	"Congrats, %s! You actually got %d points? I'm impressed... kinda. 😏",
	"Well, look at you, %s! %d points, huh? You’re almost good at this! 😜",
	"Wow, %s, %d points already? Did you accidentally cheat or something? 😂",
	"Nice job, %s! %d points, but don’t let it go to your head... Oh wait, too late! 😆",
	"Congrats, %s! You got %d points. Don’t get too excited—your high score is still pathetic! 😜",
	"Whoa, %s, %d points? Are you sure you're not just lucky? 😏",
	"Look at you go, %s! %d points, but don’t worry, it’s just beginner’s luck! 🙃",
	"Hey %s, you got %d points... Now let’s see if you can do it again without tripping over your own feet! 😅",
	"Not bad, %s! %d points. Keep it up and you might actually become a pro... maybe. 😜",
	"Look who’s climbing the leaderboard! %s with %d points. Keep it up, but don’t get cocky. 😎",
}
var PointsMsgs = []string{
	"%s -> %d Points",
}
