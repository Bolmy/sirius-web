package _game_server

import (
	"fmt"
	"math/rand"
	"time"
)

//------------------------------------------------------------------------------------------------------
// Trading Functions

// MaritimeTrade handles Player <-> Bank transactions
func (p *Player) MaritimeTrade(give Resource, take Resource, ratio int, b *Bank) error {
	if p.Resources[give] < ratio {
		return fmt.Errorf("not enough %v to trade (need %d)", give, ratio)
	}
	if b.Resources[take] <= 0 {
		return fmt.Errorf("bank is out of %v", take)
	}

	// Execute transaction
	p.Resources[give] -= ratio
	b.Resources[give] += ratio

	b.Resources[take]--
	p.Resources[take]++

	return nil
}

// TradeOffer defines what is being swapped
type TradeOffer struct {
	SenderID   int
	ReceiverID int
	Give       ResourceMap
	Receive    ResourceMap
}

func ExecutePlayerTrade(p1 *Player, p2 *Player, offer TradeOffer) error {
	// 1. Verify Sender has what they offered
	if !p1.CanAfford(offer.Give) {
		return fmt.Errorf("sender cannot afford the trade")
	}

	// 2. Verify Receiver has what was requested
	if !p2.CanAfford(offer.Receive) {
		return fmt.Errorf("receiver cannot afford the trade")
	}

	// 3. Perform Swap
	// Move resources from p1 to p2
	for res, amount := range offer.Give {
		p1.Resources[res] -= amount
		p2.Resources[res] += amount
	}

	// Move resources from p2 to p1
	for res, amount := range offer.Receive {
		p2.Resources[res] -= amount
		p1.Resources[res] += amount
	}

	return nil
}

//------------------------------------------------------------------------------------------------------
// Player Pay Functions

func (p *Player) CanAfford(cost ResourceMap) bool {
	for res, amount := range cost {
		if p.Resources[res] < amount {
			return false
		}
	}
	return true
}

func (g *GameSavestate) BuildSettlementOrCity(playerID int, cornerID int) error {
	player := g.Players[playerID]
	corner := g.Board.Corners[cornerID]

	// 1. Basic Validations
	if corner.OwnerID != -1 {
		return fmt.Errorf("corner already occupied")
	} else if (corner.OwnerID == playerID) && (corner.IsCity != true) {
		// upgrade to city
		// 2.1 Pay the bank
		cost := Costs["city"]
		if err := player.Pay(cost, g.Bank); err != nil {
			return err
		}

		//3.1 Place the piece
		corner.OwnerID = playerID
		corner.IsCity = true
		player.Points += 1

	} else if corner.OwnerID == -1 {
		// place settlement
		// 2.2 Pay the bank
		cost := Costs["settlement"]
		if err := player.Pay(cost, g.Bank); err != nil {
			return err
		}

		// 3.2 Place the piece
		corner.OwnerID = playerID
		player.Points += 1

		// 4. CRITICAL: Recalculate Longest Road for EVERYONE
		// Because this settlement might have split an opponent's path
		for _, p := range g.Players {
			p.LongestRoad = p.GetLongestRoad(g.Board)
		}

		// 5. Update the Trophy
		g.UpdateLongestRoadTrophy()
	}

	return nil
}

func (p *Player) BuildRoad(edgeID int, b *Board, bank *Bank, isFree bool) error {
	edge, exists := b.Edges[edgeID]
	cost := Costs["road"]

	if !exists {
		return fmt.Errorf("invalid edge ID")
	}

	// 1. Check if occupied
	if edge.OwnerID != -1 {
		return fmt.Errorf("edge already occupied")
	}

	// 2. Check Costs (unless it's from a Road Building card/starting phase)
	if !isFree {
		if !p.CanAfford(cost) {
			return fmt.Errorf("cannot afford road")
		}
	}

	// 3. Check Connectivity
	// A road must touch a corner or another edge owned by the player
	connected := false
	for _, cornerID := range edge.Corners {
		corner := b.Corners[cornerID]
		// If player owns the corner (Settlement/City), it's connected
		if corner.OwnerID == p.ID {
			connected = true
			break
		}

		// If the corner is empty or owned by player, check other edges touching it
		if corner.OwnerID == -1 || corner.OwnerID == p.ID {
			// Find all edges that touch this corner
			for _, otherEdge := range b.Edges {
				if (otherEdge.Corners[0] == cornerID || otherEdge.Corners[1] == cornerID) &&
					otherEdge.OwnerID == p.ID && otherEdge.ID != edgeID {
					connected = true
					break
				}
			}
		}
	}

	if !connected {
		return fmt.Errorf("road must connect to your existing network")
	}

	// 4. Place Road
	edge.OwnerID = p.ID
	// Execute payment
	p.Pay(cost, bank)

	// 5. Set Longest Road of Player
	p.LongestRoad = p.GetLongestRoad(b)

	return nil
}

func (p *Player) Pay(cost ResourceMap, b *Bank) error {
	// 1. Validation Phase (Check everything before changing anything)
	if !p.CanAfford(cost) {
		return fmt.Errorf("transaction failed: insufficient resources, cannot afford")
	}

	// 2. Execution Phase
	for res, amount := range cost {
		p.Resources[res] -= amount
		b.Resources[res] += amount
	}

	return nil
}

//------------------------------------------------------------------------------------------------------
// Action Cards Effects

func (p *Player) DrawActionCard(b *Bank) (ActionCard, error) {
	// 1. Check cost (Sheep: 1, Wheat: 1, Rock: 1)
	cost := Costs["actionCard"]
	if !p.CanAfford(cost) {
		return -1, fmt.Errorf("insufficient resources to buy a development card")
	}

	// 2. Check if deck is empty
	if len(b.ActionCards) == 0 {
		return -1, fmt.Errorf("no development cards left in the bank")
	}

	// 3. Pay the cost
	for res, amount := range cost {
		p.Resources[res] -= amount
		b.Resources[res] += amount
	}

	// 4. Pop card from deck (Top of the "stack")
	card := b.ActionCards[len(b.ActionCards)-1]
	b.ActionCards = b.ActionCards[:len(b.ActionCards)-1]

	// 5. Add to player's hand
	p.ActionCards = append(p.ActionCards, card)

	return card, nil
}

func (p *Player) PlayKnight(b *Board, targetHexID int, allPlayers map[int]*Player) (Resource, error) {
	// 1. Move the Robber
	for id, hex := range b.Hexes {
		if id == targetHexID && hex.HasRobber {
			return None, fmt.Errorf("robber is already on this hex")
		}
		hex.HasRobber = (id == targetHexID)
	}

	// 2. Identify potential victims (Players with settlements on this hex)
	victimIDs := make(map[int]bool)
	for _, corner := range b.Corners {
		if corner.OwnerID != -1 && corner.OwnerID != p.ID {
			for _, hexID := range corner.AdjacentHexes {
				if hexID == targetHexID {
					victimIDs[corner.OwnerID] = true
				}
			}
		}
	}

	// 3. Logic for stealing a random card
	var stolenRes Resource = None
	if len(victimIDs) > 0 {
		// Pick a random victim
		ids := []int{}
		for id := range victimIDs {
			ids = append(ids, id)
		}
		rand.New(rand.NewSource(time.Now().UnixNano()))
		victimID := ids[rand.Intn(len(ids))]
		victim := allPlayers[victimID]

		// Extract all resources into a slice to pick one randomly
		hand := []Resource{}
		for res, count := range victim.Resources {
			for i := 0; i < count; i++ {
				hand = append(hand, res)
			}
		}

		if len(hand) > 0 {
			stolenRes = hand[rand.Intn(len(hand))]
			victim.Resources[stolenRes]--
			p.Resources[stolenRes]++
		}
	}

	// 4. Update Knight count for Largest Army
	p.KnightsPlayed++

	return stolenRes, nil
}

func (p *Player) PlayYearOfPlenty(res1 Resource, res2 Resource, b *Bank) error {
	if b.Resources[res1] <= 0 || b.Resources[res2] <= 0 {
		return fmt.Errorf("bank does not have requested resources")
	}

	p.Resources[res1]++
	b.Resources[res1]--

	p.Resources[res2]++
	b.Resources[res2]--

	return nil
}

func (p *Player) PlayMonopoly(target Resource, allPlayers map[int]*Player) int {
	totalStolen := 0
	for id, otherPlayer := range allPlayers {
		if id == p.ID {
			continue
		}
		count := otherPlayer.Resources[target]
		otherPlayer.Resources[target] = 0
		totalStolen += count
	}
	p.Resources[target] += totalStolen
	return totalStolen
}
