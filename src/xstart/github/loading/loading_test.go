package loading

import (
    "errors"
    "testing"
    "time"
)

func TestSpin(t *testing.T) {
    title := "Test spin"
    t.Run(title, func(t *testing.T) {
        showLoading(title)
    })
}

func TestRestart(t *testing.T) {
    title := "Restarting..."

    t.Run(title, func(t *testing.T) {
        l := StartNew(title)

        l.Restart()
    })
}

func TestStart(t *testing.T) {
    l := NewLoading("Starting...")
    got := l.Start()
    want := l

    if want != got {
        t.Errorf("Want '%v', but Got '%v'", want, got)
    }
}

func TestIsColorAllowed(t *testing.T) {
    allowedColor := "blue"

    got := IsColorAllowed(allowedColor)
    want := true

    if want != got {
        t.Errorf("Want '%t', but Got '%t'", want, got)
    }

}

func TestStartNew(t *testing.T) {
    got := StartNew("Starting new")

    if got == nil {
        t.Errorf("NIL is not allowed")
    }
}

func TestSetColorErr(t *testing.T) {
    l := NewLoading("Setting color")
    want := errors.New("color not found")
    got := l.SetColor("black")

    if want == got {
        t.Errorf("Want '%v', but Got'%v'", want, got)
    }
}

func TestSetColor(t *testing.T) {
    loader := NewLoading("Setting valid color")

    colorErr := loader.SetColor("blue")

    if colorErr != nil {
        t.Errorf("ERROR DUE TO: %v", colorErr)
    }
}

func TestSetSpeed(t *testing.T) {
    loader := NewLoading("Setting speed")

    speedErr := loader.SetSpeed(200 * time.Millisecond)

    if speedErr == nil {
        t.Errorf("Loading instance is nil, due to: %v", speedErr)
    }

}

func TestSetLoader(t *testing.T) {
    loader := NewLoading("Set Loaders")

    loaders := loader.SetLoaders([]string{"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"})

    if loaders == nil {
        t.Errorf("Loaders not set due to: %v", loaders)
    }
}

func showLoading(title string) {
    loader := StartNew(title)
    defer loader.Stop()
    time.Sleep(2 * time.Second)
}
