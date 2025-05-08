## [inbound] join, which puts them in a lobby
{"type":"join", "name":"SOME-NAME-1"}
- only available while the game is inactive 
- returns[outbound] to all players game state snapshot
{"type":"join", "name":"SOME-NAME-2"}
- only available while the game is inactive 
- returns[outbound] to all players game state snapshot
{"type":"join", "name":"SOME-NAME-3"}
- only available while the game is inactive 
- returns[outbound] to all players game state snapshot
{"type":"join", "name":"SOME-NAME-4"}
- only available while the game is inactive 
- returns[outbound] to all players game state snapshot

## [inbound] show players:
{"type":"show_players"}
- any player can send this message at any time, though idk if its necessary for the UI
- anyone can send
- returns[outbound] to all players game state snapshot
## [inbound] assigns players:
{"type":"assign_teams"}
- only available while the game is inactive 
- anyone can send
- upon receiving: game server randomly selects/assigns team teams.
- returns[outbound] to all players game state snapshot

## [inbound]start game:
- {"type":"start"}
- only available while the game is inactive 
- anyone can send
- only available if no on is the lobby and teams have members and team numbers are even
- upon receiving: game server starts the inactive game
-   returns[outbound] to all players game state snapshot
-   emits question to play from each team whose turn it is

## [outbound]Emit Question:
{"type":"question","name":"What is the capital of France?"}
- game emits this message to specific player

## [inbound]Answer:
- {"type":"answer", "answer": "Paris"}                                    
- game validates correct person is answer (is_valid)
- IF answer is valid, game moves question pointer.
- IF no more questions for the team then gam over and :
-   returns[outbound] winner message to all players whatever team one
- ELSE
-   asks next person on that team their question


## Winner                                                             
{"type":"winner","name":"A-Team"}                                                     
- game emits this message when a team wins    ... people should join in a lobby , then assign teams, then start ... then the table w/ the cups and an answer input
-
## restart game                                                                             
- {"type":"restart"}                                                                        
- restarts game, keeping teams as-is (they can reassign tho)                                
- anyone can do this and any time                                                           
- returns to all players game state snapshot

-
