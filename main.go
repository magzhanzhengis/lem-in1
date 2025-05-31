package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"strings"

	"lem-in/lib"
)

type RoomJSON struct {
	Name    string `json:"name"`
	X       int    `json:"x"`
	Y       int    `json:"y"`
	IsStart bool   `json:"isStart,omitempty"`
	IsEnd   bool   `json:"isEnd,omitempty"`
}

type TunnelJSON struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type MoveJSON struct {
	AntID  string `json:"antId"`
	ToRoom string `json:"toRoom"`
}

type TurnJSON struct {
	Moves []MoveJSON `json:"moves"`
}

type SimulationJSON struct {
	Rooms           []RoomJSON   `json:"rooms"`
	Tunnels         []TunnelJSON `json:"tunnels"`
	AntsCount       int          `json:"antsCount"`
	StartRoom       string       `json:"startRoom"`
	EndRoom         string       `json:"endRoom"`
	SimulationTurns []TurnJSON   `json:"simulationTurns"`
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/simulate", simulateHandler)
	fmt.Println("Server running on http://localhost:8080")
	go openBrowser("http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func openBrowser(url string) {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	case "windows":
		cmd = "rundll32"
		args = append(args, "url.dll,FileProtocolHandler")
	default:
		cmd = "xdg-open"
	}
	args = append(args, url)
	exec.Command(cmd, args...).Start()
}

func simulateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read input", http.StatusBadRequest)
		return
	}
	lines := strings.Split(string(data), "\n")
	field := lib.Field{}
	err = lib.ResolveInput(lines, &field)
	if err != nil {
		http.Error(w, "Invalid format: "+err.Error(), http.StatusBadRequest)
		return
	}

	visited := make(map[string]bool)
	var paths [][]string
	lib.FindShortestPaths(field, "", []string{}, visited, &paths)
	paths = lib.RemoveTooLongPaths(paths)

	if len(paths) == 0 {
		http.Error(w, "No valid paths", http.StatusUnprocessableEntity)
		return
	}

	turnLog := [][]MoveJSON{}
	for _, ant := range field.Ants {
		ant.CurrentRoom = field.StartRoomName
		ant.IsFinished = false
	}

	for done := false; !done; {
		done = true
		turn := []MoveJSON{}
		used := map[string]bool{}
		for _, ant := range field.Ants {
			if ant.CurrentRoom == field.EndRoomName {
				ant.IsFinished = true
				continue
			}
			done = false
			next := getNext(ant.CurrentRoom, paths, used, field)
			if next != "" {
				used[ant.CurrentRoom+"-"+next] = true
				turn = append(turn, MoveJSON{AntID: fmt.Sprintf("L%d", ant.ID+1), ToRoom: next})
				ant.CurrentRoom = next
				if next == field.EndRoomName {
					ant.IsFinished = true
				}
			}
		}
		if len(turn) > 0 {
			turnLog = append(turnLog, turn)
		}
	}

	var rooms []RoomJSON
	for _, r := range field.Rooms {
		x, y := 0, 0
		fmt.Sscanf(r.Name+" 0 0", "%s %d %d", &r.Name, &x, &y)
		rooms = append(rooms, RoomJSON{
			Name:    r.Name,
			X:       x,
			Y:       y,
			IsStart: r.Name == field.StartRoomName,
			IsEnd:   r.Name == field.EndRoomName,
		})
	}

	tunnels := []TunnelJSON{}
	seen := make(map[string]bool)
	for _, r := range field.Rooms {
		for _, conn := range r.ConnectedRooms {
			key := r.Name + "-" + conn
			rev := conn + "-" + r.Name
			if !seen[key] && !seen[rev] {
				tunnels = append(tunnels, TunnelJSON{From: r.Name, To: conn})
				seen[key] = true
			}
		}
	}

	out := SimulationJSON{
		Rooms:           rooms,
		Tunnels:         tunnels,
		AntsCount:       len(field.Ants),
		StartRoom:       field.StartRoomName,
		EndRoom:         field.EndRoomName,
		SimulationTurns: []TurnJSON{},
	}
	for _, turn := range turnLog {
		out.SimulationTurns = append(out.SimulationTurns, TurnJSON{Moves: turn})
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(out)
}

func getNext(current string, paths [][]string, used map[string]bool, field lib.Field) string {
	for _, path := range paths {
		for i, room := range path {
			if room == current && i+1 < len(path) {
				next := path[i+1]
				if used[current+"-"+next] {
					continue
				}
				occupied := false
				for _, ant := range field.Ants {
					if ant.CurrentRoom == next {
						occupied = true
						break
					}
				}
				if !occupied || next == field.EndRoomName {
					return next
				}
			}
		}
	}
	return ""
}
