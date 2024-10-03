// Branch matrix refactoring
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

	var gotPoint bool

	cursor.Hide()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()

	player := []Pos{{x: 1, y: 1}, {x: 1, y: 1}, {x: 1, y: 1}, {x: 1, y: 1}}
	points := []Pos{{x: 4, y: 4}, {x: 9, y: 3}, {x: 2, y: 2}}

	go getInput()

	for {
		//move o player de acordo com a entrada de teclado
		movePlayer(&player, playerIntention, &gotPoint)
		checkPlayerColl(&player, &points, &gotPoint)

		cursor.Move(0, screeHeight+3) //eu acho que esse +3 eh opcional mas eu tô com preguiça de checar
		printBoard(&player, &points)  //printa o mapa

		time.Sleep(time.Millisecond * 500) // espera um tempinho
	}
}

func movePlayer(player *[]Pos, playerIntention int, gotPoint *bool) {

	//essa é a variável se tornará a nova posição da cabeça da cobrinha
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

	//se não pegar ponto, apaga a "cauda" da cobrinha
	if !*gotPoint {
		*player = pop(*player, 0)
	} else {
		*gotPoint = false
	}
}

func checkPlayerColl(player *[]Pos, points *[]Pos, gotPoint *bool) {
	head := (*player)[len(*player)-1]

	//checagem da colisão com as bordas
	if head.x == -1 || head.y == -1 || head.x == screenWidth || head.y == screeHeight {
		os.Exit(0)
	}

	for i, pt := range *points {
		if head.x == pt.x && head.y == pt.y {
			*gotPoint = true
			*points = pop(*points, i)
			*points = append(*points, Pos{x: rand.Int() % screenWidth, y: rand.Int() % screeHeight})
			break
		}
	}

	for _, segment := range (*player)[:len(*player)-1] {
		if segment.x == head.x && segment.y == head.y {
			os.Exit(0)
		}
	}

}

func printBoard(player, points *[]Pos) {

	//O tabuleiro que será printado
	var BoardStr strings.Builder

	//a matrix que representa o tabuleiro do jogo
	board := make([][]int, screeHeight)
	for i := range board {
		board[i] = make([]int, screenWidth)
	}

	for _, point := range *points {
		board[point.y][point.x] = 1
	}

	for _, segment := range *player {
		board[segment.y][segment.x] = 2
	}

	for _, line := range board {

		BoardStr.WriteByte('#')
		for _, element := range line {
			switch element {
			case 1:
				BoardStr.WriteString("º")
			case 2:
				BoardStr.WriteString("█")
			default:
				BoardStr.WriteString(" ")
			}
		}
		BoardStr.WriteString("#\n")
	}

	print(BoardStr.String())
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
