use bumper::*;
use rand::prelude::*;
use std::{env::args, 
        time::Instant};


struct Args {
    name: PlayerName,
    password: String,
    color: Color,
    host: String,
    port: u16,
}

fn main() -> Result<(), ExampleError> {
    let arguments = parse_env_args()?;

    let mut connection = Connection::connect(arguments.host.as_str(), arguments.port)?;
    let mut game = connection.login(arguments.password.as_str(), arguments.name, arguments.color)?;
    
    let mut rng = rand::thread_rng();
    loop {
        // drive in random direction max vector 1x1
        game.accelerate((rng.gen::<f32>() * 2. - 1., rng.gen::<f32>() * 2. - 1.))?;

        // interpret 5 seconds events, then accelerate again
        let now = Instant::now();
        while now.elapsed().as_secs() < 5 {
            match game.wait_next()? {
                Event::Player(block) => println!("{block:?}"),
                Event::Map(block) => println!("{block:?}"),
                Event::Status(block) => println!("{block:?}"),
                Event::GameEnded => break,
            }
        }
    }
}

fn parse_env_args() -> Result<Args, ExampleError> {
    let args: Vec<String> = args().collect();
    if args.len() != 6 {
        print_usage();
        return Err(ExampleError::ArgumentError);
    }
    let name = PlayerName::new(args[1].as_str())?;
    let password = args[2].clone();
    let color = args[3].as_str().into();
    let host = args[4].clone();
    let port = args[5].parse().expect("[PORT] must be a positive number");

    Ok(Args {
        name,
        password,
        color,
        host,
        port
    })
}

fn print_usage() {
    println!("must run random_direction_bumper with arguments [NAME] [PASSWORD] [COLOR] [HOST] and [PORT]");
}

#[derive(Debug)]
enum ExampleError {
    BumperError(Error),
    ArgumentError,
}

impl From<Error> for ExampleError {
    fn from(err: Error) -> Self {
        ExampleError::BumperError(err)
    }
}
