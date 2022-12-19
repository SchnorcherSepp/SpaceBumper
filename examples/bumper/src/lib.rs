use nom::{
         IResult,
         branch::alt,
         bytes::complete::{tag, take_until1},
         combinator::opt,
         character::complete::{char, digit1, line_ending, u32},
         error::{Error as NomError},
         multi::many_till,
         sequence::tuple,
};

use std::{
         io::{BufReader, Write, BufRead, Error as IoError},
         net::TcpStream,
         fmt::{Display, Formatter},
};

#[derive(Debug)]
pub enum Error {
    ConnectionError(IoError),
    UserNameLengthMissmatch,
    LoginError(IoError),
    LoginResultError(String),
    AccelerationWriteError(IoError),
    ServerDidntSendLine(IoError),
    InvalidServerMessage(String),
}

pub struct Connection {
    writer: TcpStream,
    reader: BufReader<TcpStream>,
}

impl Connection {
    pub fn connect(host: &str, port: u16) -> Result<Connection, Error> {
        connect(host, port)
    }
}

fn connect(host: &str, port: u16) -> Result<Connection, Error> {
    let writer = TcpStream::connect(format!("{host}:{port}")).map_err(Error::ConnectionError)?;
    let reader = BufReader::new(writer.try_clone().map_err(Error::ConnectionError)?);
    Ok(Connection {
        writer,
        reader
    })
}

#[derive(Debug, PartialEq, Eq)]
pub enum Color {
    Red,
    Blue,
    Green,
    Orange,
}

impl From<&str> for Color {
    fn from(text: &str) -> Color {
        match text {
            "red" => Color::Red,
            "blue" => Color::Blue,
            "green" => Color::Green,
            "orange" => Color::Orange,
            _ => Color::Red,
        }
    }
}

impl Display for Color {
    fn fmt(&self, fmt: &mut Formatter<'_>) -> Result<(), std::fmt::Error> { 
        let color = match self {
            Color::Red => "red",
            Color::Blue => "blue",
            Color::Green => "green",
            Color::Orange => "orange",
        };
        write!(fmt, "{color}")
    }
}

pub struct PlayerName {
    name: String,
}

impl Display for PlayerName {
    fn fmt(&self, fmt: &mut Formatter<'_>) -> Result<(), std::fmt::Error> { 
        let name = &self.name;
        write!(fmt, "{name}")
    }
}


impl PlayerName {
    pub fn new(name: &str) -> Result<Self, Error> {
        if name.is_empty() || name.len() >= 20 {
            Err(Error::UserNameLengthMissmatch)
        } else {
            Ok(Self { name: name.to_string() })
        }
    }
}

pub struct Game<'a, W: Write, R: BufRead> {
    pub id: usize,
    writer: &'a mut W,
    reader: &'a mut R
}

type TcpGame<'a> = Game<'a, TcpStream, BufReader<TcpStream>>;

pub trait Login<'a, W: Write, R: BufRead> {
    fn login(&'a mut self, password: &str, name: PlayerName, color: Color) -> Result<Game<'a, W, R>, Error>;
}

impl<'a> Login<'a, TcpStream, BufReader<TcpStream>> for Connection {
    fn login(&'a mut self, password: &str, name: PlayerName, color: Color) -> Result<TcpGame<'a>, Error> {
        let id = login(&mut self.writer, &mut self.reader, password, name, color)?;
        Ok(TcpGame {
           id,
           writer: &mut self.writer,
           reader: &mut self.reader
        })
    }
}

fn login<W: Write, R: BufRead>(writer: &mut W, reader: &mut R, password: &str, name: PlayerName, color: Color) -> Result<usize, Error> {
    let login_message = format!("{password}|{name}|{color}\n");
    writer.write(login_message.as_bytes()).map_err(Error::LoginError)?;
    let mut line = String::new();
    reader.read_line(&mut line).map_err(Error::LoginError)?;
    if let Some(id) = parse_login_result(&line) {
        Ok(id)
    } else {
        Err(Error::LoginResultError(line.to_string()))
    }
}

pub trait Accelerate {
    fn accelerate(&mut self, acceleration: (f32, f32)) -> Result<(), Error>;
}

impl<'a> Accelerate for TcpGame<'a> {
    fn accelerate(&mut self, acceleration: (f32, f32)) -> Result<(), Error> {
        let (x, y) = acceleration;
        self.writer.write(format!("{x}|{y}\n").as_bytes()).map_err(Error::AccelerationWriteError)?;
        self.writer.flush().map_err(Error::AccelerationWriteError)?;
        Ok(())
    }
}

#[derive(Debug, PartialEq)]
pub enum Event {
    Status(StatusBlock),
    Player(PlayerBlock),
    Map(MapBlock),
    GameEnded,
}

pub trait BlockEvents {
    fn wait_next(&mut self) -> Result<Event, Error>;
}

impl<'a> BlockEvents for TcpGame<'a> {
    fn wait_next(&mut self) -> Result<Event, Error> {
        let block = read_next_block(&mut self.reader)?;
        let block_str = block.as_str();
        let (_, event) = alt((parse_status_block, parse_player_block, parse_map_block))(block_str).map_err(|_error| Error::InvalidServerMessage(block.clone()))?;
        Ok(event)
    }
}

fn read_next_block<R: BufRead>(reader: &mut R) -> Result<String, Error> {
    let mut text = String::new();
    let mut line = String::new();
    reader.read_line(&mut line).map_err(Error::ServerDidntSendLine)?;
    if text_after_tag_postfix(line.as_str(), "START ", "\n").is_err() {
        return Err(Error::InvalidServerMessage(format!("expected 'START ' but found '{line}'")));
    }
    text.push_str(line.as_str());
    loop {
        let mut line = String::new();
        reader.read_line(&mut line).map_err(Error::ServerDidntSendLine)?;
        text.push_str(line.as_str());
        if text_after_tag_postfix(line.as_str(), "END ", "\n").is_ok() {
            return Ok(text);
        }
    }
}

#[derive(Debug, PartialEq, Eq)]
pub struct StatusBlock {
    pub iteration: usize,
    pub end_time: usize,
    pub max_update_time: String,
    pub max_players: usize,
}

#[derive(Debug, PartialEq)]
pub struct Player {
    pub id: usize,
    pub name: String,
    pub color: Color,
    pub position: (f32, f32),
    pub velocity: (f32, f32),
    pub acceleration: (f32, f32),
    pub score: usize,
    pub angle: f32,
    pub touching_cells: Vec<(usize, usize)>,
    pub alive: bool,
}

#[derive(Debug, PartialEq)]
pub struct PlayerBlock {
    pub players: Vec<Player>,
}

#[derive(Copy, Clone, PartialEq, Eq, Debug)]
pub enum Cell {
    Ground,
    FallDown,
    Block,
    Slow,
    Boost,
    Star,
    AntiStar,
    Spawn,
}


#[derive(Debug, PartialEq, Eq)]
pub struct MapBlock {
    pub rows: Vec<Vec<Cell>>,
}

impl From<char> for Cell {
    fn from(cell: char) -> Cell {
        match cell {
            '.' => Cell::Ground,
            ' ' => Cell::FallDown,
            '#' => Cell::Block,
            's' => Cell::Slow,
            'b' => Cell::Boost,
            'x' => Cell::Star,
            'a' => Cell::AntiStar,
            'o' => Cell::Spawn,
            _ => panic!("unexcpected cell {cell}"),
        }
    }
}

fn parse_status_block(block: &str) -> IResult<&str, Event> {
    let (block, _start) = needle(block, "START STATUS\n")?;
    let (block, iteration) = number_after_tag(block, "Iteration:")?;
    let (block, end_time) = number_after_tag(block, "Endtime:")?;
    let (block, max_update_time) = text_after_tag_postfix(block, "MaxUpdateTime:", "\n")?;
    let (block, max_players) = number_after_tag(block, "MaxPlayers:")?;
    let (block, _end) = needle(block, "END STATUS\n")?;
    
    Ok((block, Event::Status(StatusBlock {
        iteration,
        end_time,
        max_update_time: max_update_time.to_string(),
        max_players,
    })))
}

fn parse_player_block(block: &str) -> IResult<&str, Event> {
    let (block, _start) = needle(block, "START PLAYER\n")?;
    let (block, (players, _)) = many_till(parse_player, tag("END PLAYER\n"))(block)?;
    Ok((block, Event::Player(PlayerBlock {
        players
    })))
}

fn parse_map_block(block: &str) -> IResult<&str, Event> {
    let (block, _start) = needle(block, "START MAP\n")?;
    let (block, (rows, _)) = many_till(parse_map_row, tag("END MAP\n"))(block)?;
    Ok((block, Event::Map(MapBlock {
        rows
    })))
}

fn parse_player(block: &str) -> IResult<&str, Player> {
    let (block, id) = number_after_tag_postfix(block, "PlayerID:", "|")?;
    let (block, name) = text_after_tag_postfix(block, "Name:", "|")?;
    let (block, color) = color_after_tag_postfix(block, "Color:", "|")?;
    let (block, position) = position_after_tag_postfix(block, "Position:", "|")?;
    let (block, velocity) = position_after_tag_postfix(block, "Velocity:", "|")?;
    let (block, acceleration) = position_after_tag_postfix(block, "Acceleration:", "|")?;
    let (block, score) = number_after_tag_postfix(block, "Score:", "|")?;
    let (block, angle) = float_after_tag_postfix(block, "Angle:", "|")?;
    let (block, cells) = touching_cells(block)?;
    let (block, alive) = bool_after_tag_postfix(block, "IsAlive:", "\n")?;

    Ok((block, Player {
        id,
        name: name.to_string(),
        color,
        position: position.into(),
        velocity: velocity.into(),
        acceleration: acceleration.into(),
        score,
        angle,
        touching_cells: cells,
        alive: alive.0
    }))
}

fn parse_map_row(block: &str) -> IResult<&str, Vec<Cell>> {
    let (block, (row_cells, _)) = many_till(alt((char('.'), char(' '), char('#'), char('s'), char('b'), char('x'), char('a'), char('o'))), line_ending)(block)?;
    Ok((block, row_cells.into_iter().map(Cell::from).collect()))
}

fn needle<'a>(haystack: &'a str, needle: &str) -> IResult<&'a str, ()> {
    let (remainder, _needle) = tag(needle)(haystack)?; 
    Ok((remainder, ()))
}

fn number_after_tag<'a>(haystack: &'a str, prefix: &str) -> IResult<&'a str, usize> {
    number_after_tag_postfix(haystack, prefix, "\n")
}

fn number_after_tag_postfix<'a>(haystack: &'a str, prefix: &str, postfix: &str) -> IResult<&'a str, usize> {
    let (remainder, (_, number, _)) = tuple((tag(prefix), digit1, tag(postfix)))(haystack)?;
    Ok((remainder, number.parse().expect("impossible")))
}

fn float_after_tag_postfix<'a>(haystack: &'a str, prefix: &str, postfix: &str) -> IResult<&'a str, f32> {
    let (remainder, (_, opt_prefix, pre, _, post, _)) = tuple((tag(prefix), opt(tag("-")), digit1, tag("."), digit1, tag(postfix)))(haystack)?;
    let opt_prefix = opt_prefix.unwrap_or("+");
    Ok((remainder, format!("{opt_prefix}{pre}.{post}").parse().expect("impossible")))
}

fn text_after_tag_postfix<'a>(haystack: &'a str, prefix: &str, postfix: &str) -> IResult<&'a str, &'a str> {
    any_after_tag_postfix(haystack, prefix, postfix)
}

struct BumperBool(bool);

impl From<&str> for BumperBool {
    fn from(text: &str) -> BumperBool {
        match text {
            "true" | "True" => BumperBool(true),
            _ => BumperBool(false)
        }
    }
}

fn bool_after_tag_postfix<'a>(haystack: &'a str, prefix: &str, postfix: &str) -> IResult<&'a str, BumperBool> {
    any_after_tag_postfix(haystack, prefix, postfix)
}

fn color_after_tag_postfix<'a>(haystack: &'a str, prefix: &str, postfix: &str) -> IResult<&'a str, Color> {
    any_after_tag_postfix(haystack, prefix, postfix)
}

struct FloatPair(f32, f32);
impl From<&str> for FloatPair {
    fn from(text: &str) -> FloatPair {
        let (text, first) = float_after_tag_postfix(text, "", ",").expect("don't use this outside the routines defined here!");
        let (_text, second) = float_after_tag_postfix(text, "", "").expect("don't use this outside the routines defined here!");
        FloatPair(first, second)
    }
}

impl From<FloatPair> for (f32, f32) {
    fn from(pair: FloatPair) -> Self {
        (pair.0, pair.1)
    }
}

fn position_after_tag_postfix<'a>(haystack: &'a str, prefix: &str, postfix: &str) -> IResult<&'a str, FloatPair> {
    any_after_tag_postfix(haystack, prefix, postfix)
}

fn touching_cells(haystack: &str) -> IResult<&str, Vec<(usize, usize)>> {
    let (remainder,(_, (cells, _))) = tuple((tag("TouchingCells:"),many_till(tuple((u32, tag(","), u32, tag(";"))), tag("|"))))(haystack)?;
    Ok((remainder, cells.into_iter().map(|(x, _, y, _)| (x as usize, y as usize)).collect()))
}


fn any_after_tag_postfix<'a, OUT: From<&'a str>>(haystack: &'a str, prefix: &str, postfix: &str) -> IResult<&'a str, OUT> {
    let (remainder, (_, text, _)) = tuple((tag(prefix), take_until1(postfix), tag(postfix)))(haystack)?;
    Ok((remainder, text.into()))
}

fn parse_login_result(login: &str) -> Option<usize> {
    tuple::<_, _, NomError<_>, _>((tag("PLAYERID:"), digit1, tag("\n")))
        (login).map(|(_next, result)| result.1).map(|num| num.parse::<usize>().expect("impossible")).ok()
}

#[cfg(test)]
mod tests {
use super::*;
    #[test]
    fn should_parse_valid_login_result() {
        assert_eq!(parse_login_result("PLAYERID:1\n"), Some(1))
    }
    
    #[test]
    fn should_not_parse_invalid_login_result() {
        assert_eq!(parse_login_result("PLAYERID:\n"), None)
    }

    #[test]
    fn should_parse_negative_floats() {
        assert_eq!(float_after_tag_postfix("NEG:-123.2345\n", "NEG:", "\n"), Ok(("", -123.2345)))
    }

    #[test]
    fn should_parse_status_block() {
        let block = "START STATUS
Iteration:0
Endtime:33572
MaxUpdateTime:0s
MaxPlayers:4
END STATUS
";
        assert_eq!(parse_status_block(block), Ok(("", Event::Status(StatusBlock {
            iteration: 0,
            end_time: 33572,
            max_update_time: "0s".to_string(),
            max_players: 4
        }))))
    }

    #[test]
    fn should_parse_player_block() {
        let block = "START PLAYER
PlayerID:0|Name:Der rote Baron|Color:red|Position:820.000000,180.000000|Velocity:0.000000,0.000000|Acceleration:0.000000,0.000000|Score:100|Angle:0.000000|TouchingCells:20,4;|IsAlive:true
PlayerID:1|Name:asdads|Color:red|Position:740.000000,580.000000|Velocity:0.000000,0.000000|Acceleration:0.000000,0.000000|Score:100|Angle:0.000000|TouchingCells:18,14;5,3;|IsAlive:true
END PLAYER
";

        assert_eq!(parse_player_block(block), Ok(("", Event::Player(PlayerBlock {
            players: vec![
                Player {
                    id: 0,
                    name: "Der rote Baron".to_string(),
                    color: Color::Red,
                    position: (820., 180.),
                    velocity: (0., 0.),
                    acceleration: (0., 0.),
                    score: 100,
                    angle: 0.,
                    touching_cells: vec![(20, 4)],
                    alive: true,
                },
                Player {
                    id: 1,
                    name: "asdads".to_string(),
                    color: Color::Red,
                    position: (740., 580.),
                    velocity: (0., 0.),
                    acceleration: (0., 0.),
                    score: 100,
                    angle: 0.,
                    touching_cells: vec![(18, 14),(5, 3)],
                    alive: true,
                },
            ]
        }))))
    }

    #[test]
    fn should_parse_map_block() {
        let map = "START MAP
#xb
. o
as.
END MAP
";
        assert_eq!(parse_map_block(map), Ok(("",
            Event::Map(MapBlock {
                rows: vec![
                    vec![Cell::Block, Cell::Star, Cell::Boost],
                    vec![Cell::Ground, Cell::FallDown, Cell::Spawn],
                    vec![Cell::AntiStar, Cell::Slow, Cell::Ground],
                ],
            })))
        );
    }
}
