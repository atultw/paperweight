# paperweight

## Extensible static site generator as a package

* Choose your level of abstraction
* Define your logic in code - no config files or binaries
  * A config format + cli tool might come later 
## Extensible architecture
### Input 



### Output

Implement a function with the `Renderer` signature for RSS, HTML, JSON etc
* Any `io.Writer` can be used - write to a remote api, local disk, etc
* Built-in: HTML
