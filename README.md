# SpaceBumper

CloudWars is a GO implementation of the
Gathering [Hardcore Programming Compo case 2012](https://archive.gathering.org/tg12/en/creative/competitions/hardcore-programming-competition/hardcore-programming-compo-case/)
. This implementation offers more configuration options, a larger game board and the server-client protocol has been
changed.

![img](https://archive.gathering.org/tg12/files/content/images/thumbs/spacebumper_screenshot2.jpg)

## Summary

**The assignment is to create an Artificial Intelligence (AI) that plays a game called SpaceBumper. In the game every
player controls a bumpership (almost like a bumpercar but in space), and the goal is to get the highest score.**

There are multiple things that has impact on your score:

- Bump into other ships. The ship with the highest speed (velocity.length) at the impact gets 5 points and the other
  loses 5 points. However, this only happens if a minimum speed (7) is exceeded.
- Bump other ships out of the map. If your ship is the last one to bump into another ship that falls off the map, you
  get 50 points.
- If you go outside the map and re-spawns, you lose 30 points.
- Collect the stars located along the map and get 50 points per star.
- Collect the anti-stars located along the map and lose 30 points per anti-star.

The map consists of both barriers and advantages. First of all, if you drive outside the map, you will lose power and
just float out in space. Luckily we have some highly skilled towing ships that will catch you and drive you back to a
random spawn position.

The next obstacle (or you may use them to your advantage) are pillars. If you bump into them your direction will be
reflected. This deals damage based on your speed (damage = math.RoundToEven(velocity.Length * 0.3).

One set of tiles are some sort of space dirt and will slow your vehicle down until it has passed. Other tiles are full
of energy and speed you up like a rocket.

The game ends after a specified number of iterations (~ 3 min). This will be adjusted according to the number of entries
in the competition.

## How to compete

Download the GO source from [Github](https://github.com/SchnorcherSepp/SpaceBumper/) or the fully compiled binaries from
the [release page](https://github.com/SchnorcherSepp/SpaceBumper/releases). The simulator will act as a TCP/IP game
server and simulate and visualize the game according to the formal game rules, found below.

Each player AI is a separate application written in
the [language of your choice](https://github.com/SchnorcherSepp/SpaceBumper/tree/master/examples) that connects to the
simulator via TCP/IP. The clients (player AIs) and server communicate via a simple ASCII protocol over this connection.
This protocol is described in the formal game rules.

The simulator supports several game modes (AI vs AI, AI vs Human). Feel free to try or train your AI against human
players or AI's made by others entering the competition ahead of the compo tournament.

The source code for the simulator is also provided. Feel free to modify it to accommodate any type of testing process
you prefer. You are also free to create your own simulator from scratch, if you wish to do so.

## Formal game rules

Game is played in a grid of W x H cells in a coordinate system spanning X=0..W-1 and Y=0..H-1.
This grid is called the map. The map is static throughout a game session. The players and the stars are dynamic.

Each cell has dimensions 40 units x 40 units.

Cell [X, Y] has its upper-left corner in (X, Y) and its lower-right corner in (X+40, Y+40)

Each cell can be one of the following cell types:

```
   ground     (represented by '.' in the map)
   fall_down  (represented by ' ' in the map)
   block      (represented by '#' in the map)
   slow       (represented by 's' in the map)
   boost      (represented by 'b' in the map)
   star       (represented by 'x' in the map)
   anti-star  (represented by 'a' in the map)
   spawn      (represented by 'o' in the map)
```

Each player controls a bumpership with a radius of 20 units.

A bumpership collides with a cell if the circle defined by the bumperships position and radius overlaps the interior of
the cell.

A bumpership collides with another bumpership if the distance between the centers of the two bumperships is lower than
the sum of the radiuses of the bumperships.

The map contains an arbitrary number of stars. Stars don't respawn during a game session. The common case is that they
are stationary but disappear when picked up by players. A star is modeled by a circle with a radius of 20 units. A
bumpership overlaps a star if their circles overlap, same way as bumperships overlap other bumperships.

The player AIs communicate with the game simulator over a TCP/IP socket connection. The communication is done through
commands in ASCII where a newline '\n' character marks the end of a command. The player AIs are clients while the game
simulator is the server. After connecting to the server, the game status is continuously transmitted to the player.

Numbers are sent as floating point ASCII with a period ('.') as the decimal separator. Exponential form numbers are not
allowed. Command arguments are separated by '|'.

The ACCELERATION socket command (move) is the only way to interact with the game. A set command usually does not need
to be updated and is valid until it is changed or a collision with another player occurs. You are only allowed to send
a maximum of 4 commands per iteration.

A position is spawnable if no other bumperships have their position within a distance of 20 units from the position.

The game first runs the init() procedure, then executes update() a given number of times. The bumpership with the
highest score after this wins.

## Network protocol specification

### General conventions

1) The client sends a command to the server as a single line of text.
2) There are only two commands. Initial the login command. Only movement commands are sent during the game.
3) The server responds by sending a single line of text.
4) After that, the server continuously sends the world status during the game.
5) line of text must always be a string of ASCII characters terminated by a single, unix-style new line character:
   `'\n'`
6) All floating point numbers are represented as ASCII text on the form 13.37

### Initialization

When a client connects to the server, the server reads exactly one line.

The initial login command must be structured as follows:

```
{pass}|{name}|{color}\n
```

The first argument `{pass}` can be anything and is ignored.
The second argument `{name}` is the unique player name and must be between 1 and 20 characters long.
The third argument `{color}` is the player color (red, blue, green or orange).
Command arguments are separated by '|' and end with new line.

Only if the command is successful the server respond with a single line with your player ID. Otherwise the error is
returned and the client has to reconnect.

```
PLAYERID:{id}
```

The server waits for other players until the configured number is reached.
When the server enters the in-game phase, it continuously sends the world status to the clients.

### In-game phase

The server sends three status blocks type. Each block sent replaces the previous one. 

#### Status

Status returns global statistics: 
- Iteration is the current round. There are about 60 rounds per second.
- Endtime is the last Iteration. After that the game ends.
- MaxUpdateTime returns the maximum runtime of a round. This value should not exceed 16ms.
- MaxPlayers provides the spawn points count of this map and the max. supported number of players.

```
START STATUS
Iteration:0
Endtime:33572
MaxUpdateTime:0s
MaxPlayers:4
END STATUS
```

#### Player

The Player block is the most important one and describes the position and status of all players.
You can check your own ship with the PLAYERID from the init-phase. Each line is a player.
The attributes are separated by '|'. The float numbers of the vectors are separated with by ','.
The elements of lists are separated by ';'.

Player has the following attributes:
- PlayerID is the unique player ID of this ship. (int)
- Name is the unique name of this ship. (String)
- Color is the ship color. (String)
- Position is the ship position on the grid. (float Vector)
- Velocity is the speed and move direction of the ship. (float Vector)
- Acceleration is the movement command set by the player. It is converted to velocity continuously. (float Vector)
- Score are the current health points of the ship. (int)
- Angle is the angle in radians unit (rad). (float)
- TouchingCells all cells touched by a ship. (list of int [x,y] coordinates)
- IsAlive is true if the ship score is not 0. (bool))

```
START PLAYER
PlayerID:0|Name:Der rote Baron|Color:red|Position:820.000000,180.000000|Velocity:0.000000,0.000000|Acceleration:0.000000,0.000000|Score:100|Angle:0.000000|TouchingCells:20,4;|IsAlive:true
PlayerID:1|Name:asdads|Color:red|Position:740.000000,580.000000|Velocity:0.000000,0.000000|Acceleration:0.000000,0.000000|Score:100|Angle:0.000000|TouchingCells:18,14;|IsAlive:true
END PLAYER
```

#### Map

The map is static throughout a game session.
But the stars and the anti-stars are dynamic.

```
START MAP
######################
.                    ..#
x.....................x....#
########.......................#
#...s..a............o..........#
#...s..a.......b...............#
#...b..a.......b...............#
#...b..........b...............#
#...b..........##..........o...#
#...b..........##..............#
#..............................#
#....o.........................#
#..............................#
#..............................#
#.................o............#
#.....x..................x.....#
#..........                ....#
################################
END MAP
```


#### Command: move

Only one command can be sent in the game phase:

```
{float x}|{float y}\n
```

Acceleration is a vector with a maximum length of 1.
It does not have to be sent continuously and is set permanently.
After a collision the value can be reset. Check your player status and renew the command.
There is no server response in case of an error.
