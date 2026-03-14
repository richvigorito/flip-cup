package game

import "testing"

// newTestPlayer is a helper that creates a Player with no WebSocket connection,
// suitable for unit tests that do not exercise network I/O.
func newTestPlayer(id, name string) *Player {
	return &Player{ID: id, Name: name}
}

// --- Team.AddPlayer ---

func TestTeam_AddPlayer(t *testing.T) {
	team := &Team{Name: "Test", Players: []*Player{}}
	p := newTestPlayer("1", "Alice")
	team.AddPlayer(p)

	if len(team.Players) != 1 {
		t.Fatalf("Expected 1 player, got %d", len(team.Players))
	}
	if team.Players[0].ID != "1" {
		t.Errorf("Expected player ID '1', got '%s'", team.Players[0].ID)
	}
}

// --- Team.RemovePlayer ---

func TestTeam_RemovePlayer_Found(t *testing.T) {
	p1 := newTestPlayer("1", "Alice")
	p2 := newTestPlayer("2", "Bob")
	team := &Team{Name: "Test", Players: []*Player{p1, p2}}

	team.RemovePlayer(p1)

	if len(team.Players) != 1 {
		t.Fatalf("Expected 1 player after removal, got %d", len(team.Players))
	}
	if team.Players[0].ID != "2" {
		t.Errorf("Expected remaining player ID '2', got '%s'", team.Players[0].ID)
	}
}

func TestTeam_RemovePlayer_NotFound(t *testing.T) {
	p1 := newTestPlayer("1", "Alice")
	absent := newTestPlayer("99", "Ghost")
	team := &Team{Name: "Test", Players: []*Player{p1}}

	team.RemovePlayer(absent)

	if len(team.Players) != 1 {
		t.Errorf("Expected 1 player (no change), got %d", len(team.Players))
	}
}

// --- Team.GetCurrentPlayer ---

func TestTeam_GetCurrentPlayer(t *testing.T) {
	p0 := newTestPlayer("0", "Alice")
	p1 := newTestPlayer("1", "Bob")
	team := &Team{Name: "Test", Players: []*Player{p0, p1}, Turn: 1}

	got := team.GetCurrentPlayer()
	if got.ID != "1" {
		t.Errorf("Expected current player ID '1', got '%s'", got.ID)
	}
}

// --- Team.GetPlayerIndex ---

func TestTeam_GetPlayerIndex_Found(t *testing.T) {
	p0 := newTestPlayer("0", "Alice")
	p1 := newTestPlayer("1", "Bob")
	team := &Team{Name: "Test", Players: []*Player{p0, p1}}

	if idx := team.GetPlayerIndex(p1); idx != 1 {
		t.Errorf("Expected index 1, got %d", idx)
	}
	if idx := team.GetPlayerIndex(p0); idx != 0 {
		t.Errorf("Expected index 0, got %d", idx)
	}
}

func TestTeam_GetPlayerIndex_NotFound(t *testing.T) {
	p0 := newTestPlayer("0", "Alice")
	absent := newTestPlayer("99", "Ghost")
	team := &Team{Name: "Test", Players: []*Player{p0}}

	if idx := team.GetPlayerIndex(absent); idx != -1 {
		t.Errorf("Expected -1 for missing player, got %d", idx)
	}
}

// --- Team.IsPlayerAllowedToAnswer ---

func TestTeam_IsPlayerAllowedToAnswer_Correct(t *testing.T) {
	p0 := newTestPlayer("0", "Alice")
	p1 := newTestPlayer("1", "Bob")
	team := &Team{Name: "Test", Players: []*Player{p0, p1}, Turn: 1}

	if !team.IsPlayerAllowedToAnswer(p1) {
		t.Error("Expected p1 to be allowed (it is at Turn index)")
	}
}

func TestTeam_IsPlayerAllowedToAnswer_Wrong(t *testing.T) {
	p0 := newTestPlayer("0", "Alice")
	p1 := newTestPlayer("1", "Bob")
	team := &Team{Name: "Test", Players: []*Player{p0, p1}, Turn: 1}

	if team.IsPlayerAllowedToAnswer(p0) {
		t.Error("Expected p0 not to be allowed (not at Turn index)")
	}
}

// --- Team.ExtractPlayerNames ---

func TestTeam_ExtractPlayerNames(t *testing.T) {
	p0 := newTestPlayer("0", "Alice")
	p1 := newTestPlayer("1", "Bob")
	team := &Team{Name: "Test", Players: []*Player{p0, p1}}

	names := team.ExtractPlayerNames()
	if len(names) != 2 {
		t.Fatalf("Expected 2 names, got %d", len(names))
	}
	if names[0] != "Alice" || names[1] != "Bob" {
		t.Errorf("Unexpected names: %v", names)
	}
}

// --- Team.Shuffle ---

func TestTeam_Shuffle_PreservesAllPlayers(t *testing.T) {
	players := []*Player{
		newTestPlayer("1", "Alice"),
		newTestPlayer("2", "Bob"),
		newTestPlayer("3", "Carol"),
		newTestPlayer("4", "Dave"),
	}
	team := &Team{Name: "Test", Players: players}

	team.Shuffle()

	if len(team.Players) != 4 {
		t.Fatalf("Expected 4 players after shuffle, got %d", len(team.Players))
	}

	// Verify every original player is still present.
	ids := map[string]bool{}
	for _, p := range team.Players {
		ids[p.ID] = true
	}
	for _, p := range players {
		if !ids[p.ID] {
			t.Errorf("Player %s missing after shuffle", p.ID)
		}
	}
}

// --- Team.Snapshot ---

func TestTeam_Snapshot(t *testing.T) {
	p0 := newTestPlayer("p0", "Alice")
	p1 := newTestPlayer("p1", "Bob")
	team := &Team{Name: "A-Team", Players: []*Player{p0, p1}, Turn: 1}

	snap := team.Snapshot()

	if snap.Name != "A-Team" {
		t.Errorf("Expected name 'A-Team', got '%s'", snap.Name)
	}
	if snap.Turn != 1 {
		t.Errorf("Expected turn 1, got %d", snap.Turn)
	}
	if len(snap.Players) != 2 {
		t.Fatalf("Expected 2 players in snapshot, got %d", len(snap.Players))
	}
	if snap.Players[0].ID != "p0" || snap.Players[1].ID != "p1" {
		t.Errorf("Unexpected player snapshots: %+v", snap.Players)
	}
}
