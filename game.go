package poker

import "time"

func NewTexasHoldem(alerter BlindAlerter, store PlayerStore) *TexasHoldem {
	return &TexasHoldem{
		Alerter: alerter,
		Store:   store,
	}
}

type TexasHoldem struct {
	Alerter BlindAlerter
	Store   PlayerStore
}

type Game interface {
	Start(numberOfPlayers int)
	Finish(winner string)
}

func (p *TexasHoldem) Start(numberOfPlayers int) {
	blindIncrement := time.Duration(5+numberOfPlayers) * time.Minute

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Minute
	for _, blind := range blinds {
		p.Alerter.ScheduleAlertAt(blindTime, blind)
		blindTime = blindTime + blindIncrement
	}
}

func (p *TexasHoldem) Finish(winner string) {
	p.Store.RecordWin(winner)
}
