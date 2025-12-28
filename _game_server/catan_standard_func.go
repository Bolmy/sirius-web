package _game_server

import (
	"fmt"
	"math/rand"
	"time"
)

//------------------------------------------------------------------------------------------------------
// Standard functions

func RollDice() int {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	roll := r.Intn(7) + r.Intn(7)

	return roll
}

func (b *Board) DistributeResources(roll int, players map[int]*Player, bank *Bank) {
	for _, corner := range b.Corners {
		if corner.OwnerID == -1 {
			continue
		}

		player := players[corner.OwnerID]
		for _, hexID := range corner.AdjacentHexes {
			hex := b.Hexes[hexID]
			if hex.Value == roll && !hex.HasRobber {
				amount := 1
				if corner.IsCity {
					amount = 2
				}
				// Logic: Bank -> Player
				if bank.Resources[hex.ResourceType] >= amount {
					bank.Resources[hex.ResourceType] -= amount
					player.Resources[hex.ResourceType] += amount
				}
			}
		}
	}
}

//------------------------------------------------------------------------------------------------------
// Handle the Robber

func (p *Player) TotalResources() int {
	total := 0
	for _, count := range p.Resources {
		total += count
	}
	return total
}

func (p *Player) HandleSevenRoll(discard ResourceMap, b *Bank) error {
	total := p.TotalResources()
	if total <= 7 {
		return nil // No action needed
	}

	neededToDiscard := total / 2
	actualDiscardCount := 0
	for _, count := range discard {
		actualDiscardCount += count
	}

	if actualDiscardCount != neededToDiscard {
		return fmt.Errorf("must discard exactly %d cards", neededToDiscard)
	}

	// Verify they actually have the cards they are trying to discard
	if !p.CanAfford(discard) {
		return fmt.Errorf("cannot discard resources you don't have")
	}

	// Subtract from player, add to bank
	for res, amount := range discard {
		p.Resources[res] -= amount
		b.Resources[res] += amount
	}

	return nil
}
