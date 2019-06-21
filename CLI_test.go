package poker_test

import (
	"github.com/djangulo/go-fast"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	t.Run("record Denis win from user input", func(t *testing.T) {
		in := strings.NewReader("Denis wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in, dummySpyAlerter)
		cli.PlayPoker()
		poker.AssertPlayerWin(t, playerStore, "Denis")
	})

	t.Run("record letty win from user input", func(t *testing.T) {
		in := strings.NewReader("Letty wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in, dummySpyAlerter)
		cli.PlayPoker()
		poker.AssertPlayerWin(t, playerStore, "Letty")
	})
}

var dummySpyAlerter = &poker.SpyBlindAlerter{}
