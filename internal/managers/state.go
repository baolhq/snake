package managers

type GameState int

const (
	GameRunning GameState = iota
	GamePaused
	GameOver
)

var State = newStateManager()

type stateManager struct {
	current GameState
}

func newStateManager() *stateManager {
	return &stateManager{current: GameRunning}
}

func (s *stateManager) Set(state GameState) {
	s.current = state
}

func (s *stateManager) Current() GameState {
	return s.current
}

func (s *stateManager) Is(state GameState) bool {
	return s.current == state
}
