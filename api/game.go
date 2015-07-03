package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/fortytw2/embercrest/api/util"
	"github.com/fortytw2/embercrest/gc"
	"github.com/julienschmidt/httprouter"
)

// Queue enters the user into the queue
func Queue(gc *gc.GC) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		user, err := util.Authenticate(r, gc.UserService)
		if err != nil {
			util.JSONError(w, err, http.StatusUnauthorized)
			return
		}

		err = gc.EnterQueue(user)
		if err != nil {
			util.JSONError(w, err, http.StatusInternalServerError)
			return
		}

		length, err := gc.Q.Length()
		if err != nil {
			util.JSONError(w, err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "queue entered", "length": length})
	}
}

// CancelQueue removes the user from the queue
func CancelQueue(gc *gc.GC) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		user, err := util.Authenticate(r, gc.UserService)
		if err != nil {
			util.JSONError(w, err, http.StatusUnauthorized)
			return
		}

		err = gc.ExitQueue(user.Username)
		if err != nil {
			util.JSONError(w, err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "queue left"})
	}
}

type draftRequest struct {
	Level int    `json:"level"`
	Class string `json:"class"`
}

// Draft handles a users draft
func Draft(gc *gc.GC) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, err := util.Authenticate(r, gc.UserService)
		if err != nil {
			util.JSONError(w, err, http.StatusUnauthorized)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			util.JSONError(w, err, http.StatusBadRequest)
			return
		}

		var dr draftRequest
		err = json.Unmarshal(body, &dr)
		if err != nil {
			util.JSONError(w, err, http.StatusBadRequest)
			return
		}

		match, err := gc.GetMatch(ps.ByName("id"))
		if err != nil {
			util.JSONError(w, err, http.StatusNotFound)
			return
		}

		if match.Players[0] != user.Username && match.Players[1] != user.Username {
			util.JSONError(w, errors.New("you can't draft for another player"), http.StatusUnauthorized)
			return
		}

		if match.PlayerTurn != user.Username {
			util.JSONError(w, errors.New("it is not your turn to draft"), http.StatusNotFound)
			return
		}

		class, err := gc.GetClass(dr.Class)
		if err != nil {
			util.JSONError(w, errors.New("class not found"), http.StatusNotFound)
			return
		}

		match.Draft(dr.Class, dr.Level, class)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "unit created", "funds remaining": match.Funds[user.Username]})
	}
}

// MatchHistory returns all of a users past matches
func MatchHistory(gc *gc.GC) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		inactiveMatches, err := gc.GetUsersMatches(ps.ByName("user"), false)
		if err != nil {
			util.JSONError(w, err, http.StatusInternalServerError)
			return
		}

		if len(inactiveMatches) == 0 {
			w.WriteHeader(http.StatusNotFound)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{"message": "no matches found"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(inactiveMatches)
	}
}

// Match returns a single match by id
func Match(gc *gc.GC) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, err := util.Authenticate(r, gc.UserService)
		if err != nil {
			util.JSONError(w, err, http.StatusUnauthorized)
			return
		}

		match, err := gc.GetMatch(ps.ByName("id"))
		if err != nil {
			util.JSONError(w, err, http.StatusInternalServerError)
			return
		}

		if match.Players[0] != user.Username && match.Players[1] != user.Username {
			util.JSONError(w, errors.New("cannot view active match not belonging to you"), http.StatusUnauthorized)
			return
		}
		if match == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{"message": "no match found"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(match)
	}
}

// Turn handles all logic wrapping IO for a user to make a turn
func Turn(gc *gc.GC) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "not implemented yet"})
	}
}

// Concede allows a user to end the game prematurely
func Concede(gc *gc.GC) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		_, err := util.Authenticate(r, gc.UserService)
		if err != nil {
			util.JSONError(w, err, http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "match conceded"})
	}
}

// Leaderboard gets the top 25 users
func Leaderboard(gc *gc.GC) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		users, err := gc.GetLeaderboard()
		if err != nil {
			util.JSONError(w, err, http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

// Classes returns all classes currently in this copy of embercrest
func Classes(gc *gc.GC) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		classes, _ := gc.GetClasses()
		json.NewEncoder(w).Encode(classes)
	}
}

// Tiles returns all tiles currently in this copy of embercrest
func Tiles(gc *gc.GC) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		tiles, _ := gc.GetTiles()
		json.NewEncoder(w).Encode(tiles)
	}
}

// Version returns the current version
func Version(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"version": os.Getenv("VERSION")})
}
