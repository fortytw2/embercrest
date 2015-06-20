package gc

import (
	"log"
	"time"

	"github.com/fortytw2/embercrest/datastore"
	"github.com/fortytw2/embercrest/game"
	"github.com/fortytw2/embercrest/user"
)

// GC is the game coordinator for embercrest
type GC struct {
	*datastore.Datastore
	Q Queue
}

// EnterQueue puts a user into the matchmaking pool
func (gc *GC) EnterQueue(u *user.User) error {
	return gc.Q.Enter(u)
}

// ExitQueue removes a user from the queue
func (gc *GC) ExitQueue(username string) error {
	return gc.Q.Exit(username)
}

// Matchmake creates matches from users in queue
func (gc *GC) Matchmake() {
	log.Println("gc: starting matchmaking")
	ticker := time.NewTicker(time.Second * 5)
	for _ = range ticker.C {
		length, err := gc.Q.Length()
		if err != nil {
			log.Fatalln(err)
			return
		}

		log.Println("gc:", gc.Q)

		if length >= 2 {
			// match up the top two users in the queue
			user1, err := gc.Q.Pop()
			if err != nil {
				log.Fatalln("gc: something went wrong:", err)
			}
			user2, err := gc.Q.Pop()
			if err != nil {
				log.Fatalln("gc: something went wrong:", err)
			}

			tiles, err := gc.GetTiles()
			if err != nil {
				log.Fatalln("gc: something went wrong:", err)
			}

			// put them in a new match
			m := game.NewMatch(user1, user2, tiles)

			log.Println("gc: match", m.UUID, "started between", user1, "and", user2)
			gc.CreateMatch(m)
		}
	}
}

// EnforceTimeLimits moves the match to the next players turn if
// one player takes more than 60s to send their turn
func (gc *GC) EnforceTimeLimits() {
	ticker := time.NewTicker(time.Second * 1)
	for _ = range ticker.C {
		matches, err := gc.ActiveMatches()
		if err != nil {
			log.Fatalln("gc: something went wrong:", err)
		}

		for _, m := range matches {
			// if no one has made a move yet, no reason to bother
			if m.LastTurn.IsZero() {
				continue
			}

			if time.Now().Unix()-m.LastTurn.Unix() > 60 {
				log.Println("gc: time limit exceeded by", m.PlayerTurn, "in match", m.UUID)
				m.NextTurn()
				gc.UpdateMatch(&m)
			}
		}

	}
}

// Noxville elo snippet
// final static double K_WEIGHTING = 64.0
// final static double MEAN_RATING = 1000.0
//
// def newRadElo = calculateElo(radElo, direElo, (winner == "Radiant"))
// def newDireElo = calculateElo(direElo, radElo, (winner == "Dire"))
//
// double calculateElo(Elo thisTeam, Elo otherTeam, boolean win) {
//
//         double diff = (win ? 1.0 : -1.0) * (thisTeam.rating - otherTeam.rating) / 400.00
//         double _d = 1.0 / (1.0 + Math.pow(10.0, diff))
//         double _e = K_WEIGHTING * _d
//
//         return win ? thisTeam.rating + _e : thisTeam.rating - _e
// }
