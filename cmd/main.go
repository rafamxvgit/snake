package main

import (
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"atomicgo.dev/cursor"
)

const right int = 100
const up int = 119
const left int = 97
const down int = 115

const screeHeight int = 13
const screenWidth int = 30

var playerIntention int = 100

type Pos struct {
	x int
	y int
}

func main() {
	clear()
	jogo()
}

func jogo() {
	cursor.Hide()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()

	var gotPoint bool
	player := []Pos{{x: 1, y: 1}, {x: 1, y: 1}, {x: 1, y: 1}, {x: 1, y: 1}}

	points := []Pos{{x: 4, y: 4}}

	go getInput()

	for {
		//move o player de acordo com a entrada de teclado
		movePlayer(&player, playerIntention, &gotPoint)
		checkPlayerColl(&player[len(player)-1], &points, &gotPoint)

		cursor.Move(0, screeHeight+3) //eu acho que esse +3 eh opcional mas eu tô com preguiça de checar
		printBoard(&player, &points)  //printa o mapa

		time.Sleep(time.Millisecond * 500) // espera um tempinho
	}
}

func movePlayer(player *[]Pos, playerIntention int, gotPoint *bool) {

	newPos := (*player)[len(*player)-1]
	switch playerIntention {
	case right:
		newPos.x++
		*player = append(*player, newPos)
	case up:
		newPos.y--
		*player = append(*player, newPos)
	case left:
		newPos.x--
		*player = append(*player, newPos)
	case down:
		newPos.y++
		*player = append(*player, newPos)
	}

	if !*gotPoint {
		*player = pop(*player, 0)
	} else {
		*gotPoint = false
	}
}

func checkPlayerColl(plrPos *Pos, points *[]Pos, gotPoint *bool) {
	//checagem da colisão com as bordas
	if plrPos.x == -1 || plrPos.y == -1 || plrPos.x == screenWidth || plrPos.y == screeHeight {
		os.Exit(0)
	}

	for i, pt := range *points {
		if plrPos.x == pt.x && plrPos.y == pt.y {
			*gotPoint = true
			*points = pop(*points, i)
			*points = append(*points, Pos{x: rand.Int() % screenWidth, y: rand.Int() % screeHeight})
			break
		}
	}

}

func printBoard(player, points *[]Pos) {
	var mapa strings.Builder
	for y := 0; y < screeHeight; y++ {
		mapa.WriteByte('#')
		for x := 0; x < screenWidth; x++ {
			amIHere := false
			for _, pos := range *player {
				if y == pos.y && x == pos.x {
					mapa.WriteString("█")
					amIHere = true
					break
				}
			}
			if !amIHere {
				for _, pos := range *points {
					if y == pos.y && x == pos.x {
						mapa.WriteString("º")
						amIHere = true
						break
					}
				}
			}
			if !amIHere {
				mapa.WriteByte(' ')
			}
		}
		mapa.WriteByte('#')
		mapa.WriteByte('\n')
	}
	print(mapa.String())
}

func getInput() {
	var char []byte = make([]byte, 1)
	for {
		os.Stdin.Read(char)
		number := int(char[0])
		if number == right || number == up || number == left || number == down {
			playerIntention = number
		}
	}
}

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func pop[T any](slice []T, element int) []T {
	return append(slice[:element], slice[element+1:]...)
}
