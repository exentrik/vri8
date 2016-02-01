package main

import "fmt"
import "math/rand"
import (
	"time"
)

type suit int
type player []card
type card struct {
	suit  string
	value int
	name string
}

var turn = "Player1"

var deck = make([]card,52)
var table =[]card{}
var player1 player
var player2 player
var currentPlayer player
var currentPlayerHasFinishedHisTurn bool = false
var currentPlayerDrawnCards = 0
var reShuffleCount = 0
var simulateMoveCount = 0
var cardsDrawnCounter = 0
var drawCardToPlayerTotalCount = 0
var placedCardsOnTableCount = 0
var passesCount = 0

func main() {
	rand.Seed(time.Now().Unix())
	setupTable()
	shuffleDeck()
	deal()
	drawCardToTable()
	play()
}


func setupTable(){
	d := 0
	var suits = []string{"Spades","Hearts","Clubs","Diamonds"}
	for i := 0; i < 4; i++ {
		for y := 2; y < 15; y++ {
			deck[d] = card{suit: suits[i], value: y,name:getNameFromValue(y,suits[i])}
			d++
		}
	}
	player1 =[]card{} //reset player1
	player2 =[]card{} //reset player2
	table =[]card{} // reset table
}


func getNameFromValue(value int, suit string) string{
	switch value {
	case 2: return "Two of "+suit
	case 3: return "Three of "+suit
	case 4: return "Four of "+suit
	case 5: return "Five of "+suit
	case 6: return "Six of "+suit
	case 7: return "Seven of "+suit
	case 8: return "Eight of "+suit
	case 9: return "Nine of "+suit
	case 10: return "Ten of "+suit
	case 11: return "Jack of "+suit
	case 12: return "Queen of "+suit
	case 13: return "King of "+suit
	case 14: return "Ace of "+suit
	}
return "blabla"
}

func drawCardToTable(){
	card := drawCard()
	fmt.Println("Dealer drew "+card.name +" to the table")
	table = append(table,card)
}

func takeAllButOneCardFromTableAndReshuffleDeck(){
	for  1 < len(table) { //all but the last card
		deck = append(deck,table[0])//put in deck the first of table each time
		table = append(table[:0],table[1:]...) //Delete from table
	}
	shuffleDeck()
	reShuffleCount++
}

func drawCard() card{
	cardsDrawnCounter++
	if(len(deck)==0) {
		fmt.Println("\nDealer: Hold on i am reshuffling the deck\n")
		takeAllButOneCardFromTableAndReshuffleDeck()
	}
		card := deck[0]
		deck = append(deck[:0],deck[1:]...)
		return card

}

func printCards(stack []card){
	for _, element := range stack {
		fmt.Printf("%s - %s - %d\n", element.name, element.suit, element.value)
	}
}

func shuffleDeck() {
	for i := range deck {
		j := rand.Intn(i + 1)
		deck[i], deck[j] = deck[j], deck[i]
	}
}

func deal(){
	for i := 0; i < 8; i++ {
		card := drawCard()
		player1 = append(player1,card)
		card = drawCard()
		player2 = append(player2,card)
	}
	fmt.Println("Dealer dealt 8 card to each player")
}
func play(){

	play:=true
	for play==true{
		if(turn=="Player1"){
			currentPlayer = player1
		}else{
			currentPlayer = player2
		}
		for !currentPlayerHasFinishedHisTurn{
			simulateMove()
		}
		currentPlayerHasFinishedHisTurn=false
		currentPlayerDrawnCards=0

		if(turn=="Player1"){
			player1 = currentPlayer
			turn="Player2"
		}else{
			player2 = currentPlayer
			turn="Player1"
		}
		//fmt.Printf("Player1 has cards: %d\n",len(player1))
		//fmt.Printf("Player2 has cards: %d\n",len(player2))

		if(len(player1)==0){
			fmt.Println("\nPlayer1 WON\n")
			play = false

		}else if(len(player2)==0){
			fmt.Println("\nPlayer2 WON\n")
			play = false
		}

	}
	fmt.Printf("Deck was reshuffled %d times\n",reShuffleCount)
	fmt.Printf("Game was finished in %d valid moves\n",simulateMoveCount)
	fmt.Printf("which players played %d cards to the table \n",placedCardsOnTableCount)
	fmt.Printf("and had to pass %d times \n",passesCount)
	fmt.Printf("Cards drawn from deck: %d including first 17 \n",cardsDrawnCounter)
	fmt.Printf("Cards drawn to a player: %d \n",drawCardToPlayerTotalCount)



}
func passAfterThreeDrawnCards(){
	//check if he has drawn three cards
	if(currentPlayerDrawnCards==3){
		passesCount++
		currentPlayerHasFinishedHisTurn = true
		fmt.Println(turn + " passed")
	}else{
		fmt.Println("Dealer: "+ turn+ " you cant pass before you have drawn 3 cards")
	}
}


//Simulate move will be eventually replaced with a client player.
func simulateMove(){
	simulateMoveCount++
	var card card
	var validMove = true
	for _, element := range currentPlayer {
		card = element
		validMove = playCard(card)
		if(validMove){
			break
		}
	}
	// pull up to three cards  //TODO  need to test on server that it really pulls 3 cards
	i:=0;
	for !validMove && i<3{
		i++
		//fmt.Printf(turn + " drew a card")
		drawCardToPlayer()
		validMove = playCard(currentPlayer[len(currentPlayer)-1])
	}
	if(!validMove){
		fmt.Println(turn +" wants to pass")
		passAfterThreeDrawnCards()
	}
}

func drawCardToPlayer(){
	if(currentPlayerDrawnCards==3){
		fmt.Println("Dealer: "+turn + " you cant draw more than 3 cards!" )
	}else{
		currentPlayerDrawnCards++
		drawCardToPlayerTotalCount++
		card := drawCard()
		currentPlayer = append(currentPlayer,card)
		if currentPlayerDrawnCards==1 {
			fmt.Println(turn + " drew one card")
		}else if currentPlayerDrawnCards==2{
			fmt.Println(turn + " drew a second card")
		}else{
			fmt.Println(turn + " drew the third card")}
	}

}

func playCard(playedCard card) bool{

	if(doesCurrentPlayerHaveThisCard(playedCard)){
		if(placeCardOnTable(playedCard)){
			removeFromCurrentPlayersHand(playedCard)
			currentPlayerHasFinishedHisTurn = true
			return true
		}

	}
	return false
}

func placeCardOnTable(card card) bool{
	if(isThisCardValidToPutOnTable(card)){
		placedCardsOnTableCount++
		table = append(table,card)
		fmt.Printf(turn+" put %v\n", table[len(table)-1].name)
		return true
	}
	return false
}

func isThisCardValidToPutOnTable(card card) bool{
	cardOnTable := table[len(table)-1]
	if(cardOnTable.suit == card.suit || cardOnTable.value == card.value){
		return true
	}
	return false
}

func doesCurrentPlayerHaveThisCard(playedCard card) bool{
	for _, card := range currentPlayer {
		if(card==playedCard){
			return true
		}
	}
	return false
}

func removeFromCurrentPlayersHand(playedCard card){
		for i := 0; i < len(currentPlayer); i++ {
			if(currentPlayer[i] == playedCard){
				currentPlayer = append(currentPlayer[:i],currentPlayer[i+1:]...)
			}
		}
}
