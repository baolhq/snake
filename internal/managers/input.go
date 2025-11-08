package managers

import "github.com/hajimehoshi/ebiten/v2"

// KeyState stores the state of a single key.
type KeyState struct {
	IsDown      bool // true if the key is currently held down
	WasPressed  bool // true if the key was pressed this frame
	WasReleased bool // true if the key was released this frame
}

// Action represents a named logical action in the game.
type Action string

// Predefined actions.
const (
	ActionUp    Action = "MoveUp"
	ActionDown  Action = "MoveDown"
	ActionLeft  Action = "MoveLeft"
	ActionRight Action = "MoveRight"
	ActionPause Action = "Pause"
	ActionEnter Action = "Enter"
)

// InputManager tracks key states globally and maps actions to keys.
type InputManager struct {
	keyStates map[ebiten.Key]*KeyState
	keyMap    map[Action][]ebiten.Key
}

// Input is the global instance of InputManager.
var Input = &InputManager{
	keyStates: make(map[ebiten.Key]*KeyState),
	keyMap:    make(map[Action][]ebiten.Key),
}

// Registers default actions and their associated keys.
func init() {
	Input.RegisterAction(ActionUp, ebiten.KeyUp, ebiten.KeyW)
	Input.RegisterAction(ActionDown, ebiten.KeyDown, ebiten.KeyS)
	Input.RegisterAction(ActionLeft, ebiten.KeyLeft, ebiten.KeyA)
	Input.RegisterAction(ActionRight, ebiten.KeyRight, ebiten.KeyD)
	Input.RegisterAction(ActionPause, ebiten.KeyEscape)
	Input.RegisterAction(ActionEnter, ebiten.KeyEnter)
}

// RegisterAction maps a logical action to one or more physical keys.
func (i *InputManager) RegisterAction(action Action, keys ...ebiten.Key) {
	i.keyMap[action] = keys
	for _, key := range keys {
		if _, ok := i.keyStates[key]; !ok {
			i.keyStates[key] = &KeyState{}
		}
	}
}

// Update updates the state of all tracked keys. Call this once per frame.
func (i *InputManager) Update() {
	for key, state := range i.keyStates {
		pressed := ebiten.IsKeyPressed(key)
		state.WasPressed = !state.IsDown && pressed
		state.WasReleased = state.IsDown && !pressed
		state.IsDown = pressed
	}
}

// WasPressed returns true if any key mapped to the action was pressed this frame.
func (i *InputManager) WasPressed(action Action) bool {
	return i.checkKeys(i.keyMap[action], func(s *KeyState) bool { return s.WasPressed })
}

// WasReleased returns true if any key mapped to the action was released this frame.
func (i *InputManager) WasReleased(action Action) bool {
	return i.checkKeys(i.keyMap[action], func(s *KeyState) bool { return s.WasReleased })
}

// IsDown returns true if any key mapped to the action is currently held down.
func (i *InputManager) IsDown(action Action) bool {
	return i.checkKeys(i.keyMap[action], func(s *KeyState) bool { return s.IsDown })
}

// Iterates over a list of keys and returns true if any key satisfies the checker function.
func (i *InputManager) checkKeys(keys []ebiten.Key, checker func(*KeyState) bool) bool {
	for _, key := range keys {
		if state, ok := i.keyStates[key]; ok && checker(state) {
			return true
		}
	}
	return false
}
