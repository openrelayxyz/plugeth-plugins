## Consensus Engine Plugin

This is an mvp and proof of concept plugin. The PluGeth hooks and injections allow for the implementation of a custom consensus engine. Our vision is that this plugin can be used to create development chains as well as provide the basic infrastructure to avoid having to fork Geth in order to implement a heavily altered protocol. 

This package provides a bash script which instantiates two separate nodes: a miner and a passive node. The `BlockChain` method in `main.go` then sends two contracts to the chain triggering the mining of two blocks. 

In addition to foundation Geth we support other forks as well. As such there are multiple bash scripts contained in this project to initiate the plugin to work on the commiserate chain. A not about separate chains: Ultimately this plugin was designed to be used with foundation Geth. Our vision for the plugin is to support modifying foundation Geth to enable desperate implementation and thus new chains. As such the functionality is diminished for our plugeth-etc project. Use with caution.  

In order to use the plugin navigate to the `/chain` directory. Change the permissions on the appropriate bash file.  Then point the executable file to a Geth binary like so: `/path/to/geth/geth` `./blockchain.sh `.  After a few seconds the application should close with exit code 0. 

Note: depending on where the script is deployed Geth may complain that the path to the `.ipc` file is too long. Renaming of the directories or moving the project may be necessary.