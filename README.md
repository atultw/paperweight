# paperweight

## Extensible static site generator as a package

* Choose your level of abstraction
* Define your logic in code - no config files or binaries
  * A config format + cli tool might come later 
## Extensible architecture
### Input 

* From a file glob, go array, single file, or web source
  * Add your own easily by implementing `MultiSource` or `Source`
  * Works like `sql.Rows`
### Output

Implement a function with the `Renderer` signature for RSS, HTML, JSON etc
* Any `io.Writer` can be used - write to a remote api, local disk, etc
* Built-in: HTML

The build process is centered around pipelines. Pass any `MultiSource` input and `io.Writer` output to the `MultiPipeline` struct and call `Run()`
