#[cfg(test)]
mod tests {
    use bumper::{PlayerName, Connection, Color, Login, Accelerate, BlockEvents};

    #[test]
    fn should_connect_to_running_server() {
        let _child = start_server(1234);
        assert!(Connection::connect("127.0.0.1", 1234).is_ok())
    }

    #[test]
    fn should_login_to_server() {
        let _child = start_server(1235);
        let mut connection = Connection::connect("127.0.0.1", 1235).expect("cannot establish connection");
        assert!(connection.login("pass", PlayerName::new("Bumper").expect("invalid username"), Color::Red).is_ok())
    }

    #[test] 
    fn should_accelerate() {
        let _child = start_server(1236);
        let mut connection = Connection::connect("127.0.0.1", 1236).expect("cannot establish connection");
        let mut game = connection.login("pass", PlayerName::new("Bumper").expect("invalid username"), Color::Red).expect("failed to login");
        assert!(game.accelerate((1., 0.)).is_ok());
    }
    
    #[test] 
    fn should_retrieve_twenty_events() {
        let _child = start_server(1237);
        let mut connection = Connection::connect("127.0.0.1", 1237).expect("cannot establish connection");
        let mut game = connection.login("pass", PlayerName::new("Bumper").expect("invalid username"), Color::Red).expect("failed to login");
        for _i in 0 .. 20 {
            assert!(game.wait_next().is_ok());
        }
    }

    fn start_server(port: u16) -> Server {
        let server_executable = path_to_server_exe();
        let child = std::process::Command::new(server_executable)
                .args(["-port", &format!("{port}"), "-headless", "-remote"])
                .spawn()
                .expect("failed to start the bumper executable");
        Server {
            child
        }
    }

    struct Server {
        child: std::process::Child,
    }

    impl Drop for Server {
        fn drop(&mut self) {
            self.child.kill().expect("failed to stop game server process");
        }
    }

    fn path_to_server_exe() -> String {
        std::env::var("BUMPER_EXECUTABLE").expect("You must set the BUMPER_EXECUTABLE environment variable before running the integration tests!")
    }
}
