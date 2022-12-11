package megagrid

type Ship struct {
	World    *WorldMap
	PlayerID int
	Name     string
	Color    string

	Position     *Vector
	Velocity     *Vector
	Acceleration *Vector
	Score        int
}

func NewShip(world *WorldMap, playerID int, name, color string) *Ship {

	ship := &Ship{
		World:        world,
		PlayerID:     playerID,
		Name:         name,
		Color:        color,
		Position:     new(Vector),
		Velocity:     new(Vector),
		Acceleration: new(Vector),
		Score:        100,
	}

	return ship
}
