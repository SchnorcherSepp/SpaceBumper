package core

import (
	"io"
	"math"
)

// Ship represents a player ship.
// (see NewShip)
type Ship struct {
	world    *WorldMap
	playerID int
	name     string
	color    string
	remoteRW io.ReadWriter // optional

	position     *Vector
	velocity     *Vector
	acceleration *Vector
	score        int

	factorAccel   float64
	factorVeloc   float64
	factorRollRes float64
	factorBoost   float64
	factorSlow    float64

	lastCollider *Ship
	oldPosition  *Vector
}

// NewShip create a new ship without spawning.
// (used by WorldMap.AddPlayer)
func NewShip(world *WorldMap, playerID int, name, color string, remote io.ReadWriter) *Ship {

	ship := &Ship{
		world:         world,
		playerID:      playerID,
		name:          name,
		color:         color,
		remoteRW:      remote,
		position:      new(Vector),
		velocity:      new(Vector),
		acceleration:  new(Vector),
		score:         100,   // [default 100] start health points of the ship
		factorAccel:   0.15,  // [default 0.05] how fast is acceleration converted to velocity
		factorVeloc:   0.15,  // [default 0.1] how much does velocity change the position per tick
		factorRollRes: 0.985, // [default 0.97] rolling resistance limit the max. speed
		factorBoost:   1.05,  // [default 1.1] effect of the boost cell
		factorSlow:    0.95,  // [default 0.9] effect of the slow cell
		lastCollider:  nil,
		oldPosition:   new(Vector),
	}

	return ship
}

//--------  Getter  --------------------------------------------------------------------------------------------------//

// PlayerID returns the unique player ID of this ship.
func (s *Ship) PlayerID() int {
	return s.playerID
}

// Name returns the unique name of this ship.
func (s *Ship) Name() string {
	return s.name
}

// Color returns the ship color.
// (see WorldMap.AddPlayer)
func (s *Ship) Color() string {
	return s.color
}

// Remote If set, then this ship is controlled remotely
// and the world status is sent back via this channel continuously.
func (s *Ship) Remote() io.ReadWriter {
	return s.remoteRW
}

// Position returns the ship position on the grid.
// The value is immutable (vector clone).
func (s *Ship) Position() *Vector {
	return s.position.Clone()
}

// Velocity returns the speed and move direction of the ship.
// The value is immutable (vector clone).
func (s *Ship) Velocity() *Vector {
	return s.velocity.Clone()
}

// Acceleration returns the movement command (vector) set by the player.
// The vector is converted to velocity continuously.
// The value is immutable (vector clone).
func (s *Ship) Acceleration() *Vector {
	return s.acceleration.Clone()
}

// Score are the current health points of the ship.
func (s *Ship) Score() int {
	return s.score
}

// Angle return the angle in radians (rad).
func (s *Ship) Angle() float64 {
	return s.velocity.Angle()
}

// TouchingCells returns all cells touched by a ship.
// (see WorldMap.TouchingCells)
func (s *Ship) TouchingCells() []*Cell {
	return s.world.TouchingCells(s.position)
}

// IsAlive return true if the ship score is not 0.
func (s *Ship) IsAlive() bool {
	return s.score > 0
}

//--------  Setter  --------------------------------------------------------------------------------------------------//

// Move accelerate the ship in any directions.
// The vector strength (length) is limited from 0 to 1.
// The strength is calculated as sqrt(X*X + Y*Y)
func (s *Ship) Move(acceleration *Vector) {
	// The strength is calculated as sqrt(x*x+y*y)
	strength := acceleration.Length()

	if strength > 1.0 {
		acceleration.Normalize()
	}
	s.acceleration = acceleration
}

// Spawn set the ship to a random spawner (see WorldMap.FreeSpawn).
// velocity and acceleration are reset.
func (s *Ship) Spawn() {
	spawn := s.world.FreeSpawn()
	s.velocity = new(Vector)
	s.acceleration = new(Vector)
	s.position = spawn.Clone()
}

// Collide returns true if there is a collision with the given ship.
func (s *Ship) Collide(o *Ship) bool {
	tmp := s.position.Clone()
	tmp.Add(o.position, -1)
	l := tmp.Length()
	return l < 2*CellRadius
}

//--------  UPDATE  --------------------------------------------------------------------------------------------------//

// Update the ship stats and process interactions with other objects in WorldMap.
// This function is called from WorldMap.Update().
func (s *Ship) Update() {

	// DIE: Destroyed ships are no longer updated
	// and relocated outside the grit.
	// All velocity and acceleration are reset constantly.
	//-----------------------------------------------------
	if !s.IsAlive() {
		s.velocity = new(Vector)
		s.acceleration = new(Vector)
		s.position = NewVector(-1000, -1000)
		return // EXIT
	}

	// SPEED: Acceleration is converted to velocity.
	//-----------------------------------------------------
	s.velocity.Add(s.acceleration, s.factorAccel)
	s.velocity.Multi(s.factorRollRes)

	// POSITION: Velocity moves the ship (new position).
	//-----------------------------------------------------
	s.oldPosition = s.position.Clone()
	s.position.Add(s.velocity, s.factorVeloc)
	currentCell := s.world.CellByVector(s.position)

	// OTHER SHIPS: Colliding with other ships damages the slower ship
	// and heals the faster ship. The collision changes the ship's
	// course and reset the acceleration.
	//-----------------------------------------------------
	for _, o := range s.world.Players() {
		// check collisions
		if s.playerID == o.playerID || !s.Collide(o) {
			continue
		}

		// set last collider
		s.lastCollider = o
		o.lastCollider = s

		// reset acceleration
		s.acceleration = new(Vector)
		o.acceleration = new(Vector)

		// winner / loser
		winner := o
		loser := s
		if s.velocity.Length() > o.velocity.Length() {
			winner = s
			loser = o
		}

		// score
		if winner.velocity.Length() > 7 {
			winner.score += 5
			loser.score -= 5
		}

		// loser get off the road
		loser.velocity.Add(winner.velocity.Clone(), 1.5)

		// winner bounce back
		winner.position = winner.oldPosition.Clone()
		winner.velocity.Add(winner.velocity, -1.5)
	}

	// NONE cell interaction (die).
	// A ship loses points if it falls into the void and respawn.
	// If there was previously a collision with another ship,
	// the other ship gets points.
	//-----------------------------------------------------
	if currentCell.Type() == None {
		if s.lastCollider != nil {
			s.lastCollider.score += 50
			s.lastCollider = nil
		}
		s.score -= 30
		if s.IsAlive() {
			s.Spawn()
			return //EXIT
		}
	}

	// STAR cell interaction (good).
	// The star is collected when touched and gives points.
	//-----------------------------------------------------
	if currentCell.Type() == Star {
		s.score += 50
		currentCell.SetType(Tile) // remove star
	}

	// ANTI-STAR cell interaction (bad).
	// The star is collected when touched and removes points.
	//-----------------------------------------------------
	if currentCell.Type() == Anti {
		s.score -= 30
		currentCell.SetType(Tile) // remove anti star
	}

	// NEAR CELLs: interaction with BLOCK, BOOST and SLOW.
	// These cells affect the ship as long as they are touched.
	//-----------------------------------------------------
	for _, c := range s.TouchingCells() {

		// crash and rebound
		if c.Type() == Blocked {
			// calc damage
			damage := math.RoundToEven(s.velocity.Length() * 0.3)
			s.score -= int(damage)
			// bounce back
			s.position = s.oldPosition.Clone()
			s.velocity.Add(s.velocity, -1.5)
			break

			// speed up
		} else if c.Type() == Boost {
			s.velocity.Multi(s.factorBoost)

			// slow down
		} else if c.Type() == Slow {
			s.velocity.Multi(s.factorSlow)
		}
	}
}
