SolVM is a Lua Runtime written in Go.
<div align="center">
  <img src="https://github.com/user-attachments/assets/420980fd-902e-4167-b8ba-b9ad8cb2b50c" alt="solvm-icon" width="150" style="border-radius: 15px; box-shadow: 0 4px 10px rgba(0,0,0,0.2);" />
  <h1 style="margin-top: 10px;">Welcome to SolVM</h1>
  <p><em>Simple. Elegant. Powerful.</em></p>
</div>
<p align="center">
  <strong>Love KleeStore? Give it a star on GitHub! ⭐ Your support helps!</strong><br/>
  <a href="https://github.com/kleeedolinux/solvm">
    <img src="https://img.shields.io/github/stars/kleeedolinux/solvm?style=social" alt="GitHub stars">
  </a>
</p>

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


<p align="center">
  <strong>Found SolVM useful? Don't forget to <a href="https://github.com/kleeedolinux/solvm">star the repository</a>!</strong>
</p>


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

## Benchmark
```bash
go run benchmark/benchmark.go
```
```bash
Running benchmark: Simple Arithmetic       
  Iteration 1: 31.1518ms
  Iteration 2: 23.5224ms
  Iteration 3: 24.1932ms
  Iteration 4: 26.413ms
  Iteration 5: 25.0678ms

Running benchmark: Table Operations        
  Iteration 1: 20.4479ms
  Iteration 2: 18.9732ms
  Iteration 3: 17.5482ms
  Iteration 4: 16.2188ms
  Iteration 5: 17.5624ms

  Iteration 5: 1.11ms

Running benchmark: String Operations       
  Iteration 1: 111.8972ms
  Iteration 2: 100.1595ms
  Iteration 3: 96.002ms
  Iteration 4: 101.6598ms
  Iteration 5: 99.4323ms

Benchmark Results Summary:
=========================
Table Operations:
  Average: 18.1501ms
  Min: 16.2188ms
  Max: 20.4479ms
Function Calls:
  Average: 1.53918ms
  Min: 1.0022ms
  Max: 2.035ms
String Operations:
  Average: 101.83016ms
  Min: 96.002ms
  Max: 111.8972ms
Simple Arithmetic:
  Average: 26.06964ms
  Min: 23.5224ms
  Max: 31.1518ms
```	

## License
MIT

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
