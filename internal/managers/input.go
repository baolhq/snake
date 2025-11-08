package managers

import "github.com/hajimehoshi/ebiten/v2"

type KeyState struct {
	IsDown      bool
	WasPressed  bool
	WasReleased bool
}

type Action string

const (
	ActionUp    Action = "MoveUp"
	ActionDown  Action = "MoveDown"
	ActionLeft  Action = "MoveLeft"
	ActionRight Action = "MoveRight"
	ActionPause Action = "Pause"
	ActionEnter Action = "Enter"
)

type InputManager struct {
	keyStates map[ebiten.Key]*KeyState
	keyMap    map[Action][]ebiten.Key
}

var Input = &InputManager{
	keyStates: make(map[ebiten.Key]*KeyState),
	keyMap:    make(map[Action][]ebiten.Key),
}

func init() {
	Input.RegisterAction(ActionUp, ebiten.KeyUp, ebiten.KeyW)
	Input.RegisterAction(ActionDown, ebiten.KeyDown, ebiten.KeyS)
	Input.RegisterAction(ActionLeft, ebiten.KeyLeft, ebiten.KeyA)
	Input.RegisterAction(ActionRight, ebiten.KeyRight, ebiten.KeyD)
	Input.RegisterAction(ActionPause, ebiten.KeyEscape)
	Input.RegisterAction(ActionEnter, ebiten.KeyEnter)
}

func (i *InputManager) RegisterAction(action Action, keys ...ebiten.Key) {
	i.keyMap[action] = keys
	for _, key := range keys {
		if _, exists := i.keyStates[key]; !exists {
			i.keyStates[key] = &KeyState{}
		}
	}
}

func (i *InputManager) Update() {
	for key, state := range i.keyStates {
		pressed := ebiten.IsKeyPressed(key)
		state.WasPressed = !state.IsDown && pressed
		state.WasReleased = state.IsDown && !pressed
		state.IsDown = pressed
	}
}

func (i *InputManager) WasPressed(action Action) bool {
	return i.checkKeys(i.keyMap[action], func(s *KeyState) bool { return s.WasPressed })
}

func (i *InputManager) WasReleased(action Action) bool {
	return i.checkKeys(i.keyMap[action], func(s *KeyState) bool { return s.WasReleased })
}

func (i *InputManager) IsDown(action Action) bool {
	return i.checkKeys(i.keyMap[action], func(s *KeyState) bool { return s.IsDown })
}

func (i *InputManager) checkKeys(keys []ebiten.Key, checker func(*KeyState) bool) bool {
	for _, key := range keys {
		if state, ok := i.keyStates[key]; ok && checker(state) {
			return true
		}
	}
	return false
}
