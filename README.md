SolVM is a Lua Runtime written in Go.

> Perhaps of the name SolVM, SolVM is not a VM, it's a runtime. The real VM is the gopher-lua, and SolVM is just a runtime that uses gopher-lua. I put the name of SolVM because the first idea of this project is to be a VM, but then i changed my mind.

## Usage

```bash
solvm main.lua
```

## Features
- Concurrency with goroutines  
- Channel-based communication  
- Timers and intervals  
- Cron job scheduling  
- Lua VM isolation per routine  
- Error handling and monitoring  
- File system operations  
- HTTP client module  
- Environment variable loading (.env)  
- CSV read/write  
- JSONC support  
- TOML parsing  
- YAML parsing  
- UUID generation  
- Random utilities  
- Cryptography functions  
- DateTime manipulation  
- File transfer utilities  
- INI file handling  
- TAR archive management  
- Text processing  
- Scheduler system  
- Debugging support  
- Network utilities  
- Import system  
- Monitoring and resource checking  

## Installation
[Release](https://github.com/kleeedolinux/SolVM/releases)
```bash
git clone https://github.com/kleeedolinux/SolVM
cd SolVM
```

If Linux
```bash
go build -o solvm main.go
```

If Linux
```bash
chmod +x build.sh
./build.sh
```

If Windows
```bash
./build.ps1
```

## Usage
```bash
./solvm main.lua
./solvm.exe main.lua
```

## Documentation
[DOC.md](DOC.md)

## License
MIT

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
