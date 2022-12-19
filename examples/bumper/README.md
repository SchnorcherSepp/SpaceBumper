# Bumper

A Rust client library for Space Bumper (https://github.com/SchnorcherSepp/SpaceBumper).

## License

MIT, see LICENSE file

## Testing

You may run some basic unit tests and integration tests by running `cargo tests` - 
the integration tests will fail unless you set the environment variable `BUMPER_EXECUTABLE` before pointing to the 
`SpaceBumper_1.1` executable provided by `SchnorcherSepp` found at https://github.com/SchnorcherSepp/SpaceBumper/releases/download/v1.1/SpaceBumper_1.1.zip

As an example on the Windows command prompt:
```cmd
set BUMPER_EXECUTABLE=C:/SuperBumper/SpaceBumper_1.1.exe
cargo test
```

## Usage

As I don't intend to release this library on crates.io, you have to add a dependency directly via the path option, like this:
 
```toml
[dependencies]
bumper = { path = "PATH_TO_THIS_FOLDER" }
```

See `examples/random_direction_bumper.rs` and `cargo doc` on how to use this library.

