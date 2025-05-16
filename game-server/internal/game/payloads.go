//game/game.go package game
package game

///////////////////////////////////
//////// inbound  ///////////////// 
///////////////////////////////////

// (initial player) Step 1a, create a game
type StartGamePayload struct {
	QuizFilename string `json:"file"`
	Name string `json:"name"`
}

// (subsequent players) Step 1b, enter an existing game
type JoinExistingGamePayload struct {
	GameID string `json:"game_id"`
	Name string `json:"name"`
}


// Step 2, 'join' aka give us your name
type AddPlayerPayload struct {
	Name string `json:"name"`
}


type AnswerPayload struct {
	Answer string `json:"answer"`
}


///////////////////////////////////
//////// outbound  ///////////////// 
///////////////////////////////////

// response to either 1a or 1b
type PlayerJoinedPayload struct {
	PlayerID   string `json:"player_id"`
	Name 			 string `json:"name"`
}


