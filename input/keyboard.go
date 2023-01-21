package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"log"
)

type KeyEvent struct {
	Key     ebiten.Key
	Pressed bool
}

//commenttime
type Manager struct {
	KeyStates map[ebiten.Key]bool
	Stream    chan *KeyEvent
}

func New() *Manager {
	return &Manager{
		KeyStates: make(map[ebiten.Key]bool),
		Stream:    make(chan *KeyEvent, 60),
	}
}

func (m *Manager) RegisterKey(k ebiten.Key) {
	log.Printf("registered key %v", k)
	m.KeyStates[k] = ebiten.IsKeyPressed(k)
}

func (m *Manager) Update() {
	for k := range m.KeyStates {
		m.KeyStates[k] = ebiten.IsKeyPressed(k)
		if inpututil.IsKeyJustPressed(k) {
			m.Stream <- &KeyEvent{
				Key:     k,
				Pressed: true,
			}
		}
		if inpututil.IsKeyJustReleased(k) {
			m.Stream <- &KeyEvent{
				Key:     k,
				Pressed: false,
			}
		}
	}
}
