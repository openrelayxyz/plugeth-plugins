## Plugeth Test Plugin 

This plugin works in a very similar fashion as the Consensus Engine plugin in that there is a bash script to run which instantiates individual nodes which execute a pseudo blockchain. Throughout the execution of said chain many of the PluGeth hooks and injections are engaged which in turn trigger a logging function which tracks the invoked methods and complains if any are not called. In this way we can confirm that the PluGeth application is fully implemented on the target client. 

In order to use the plugin navigate to the `/test` directory. Change the permissions on `run-test.sh` to enable execution. Then point the executable file to a Geth binary like so: `/path/to/geth/geth` `./blockchain.sh `. The test takes roughly 4.5 mins to run. If successful it should close with an exit code of 0. 

There are four methods note covered by testing at this time all within the blockTracer project: LiveCaptureFault(), LiveCaptureEnter(), LiveCaptureExit(), LiveTracerResult(). Also, there are several injections which fall outside of the testing parameters of this application: core/ NewSideBlock(), Reorg(), BlockProcessingError(), core/rawdb/ ModifyAncients(), AppendAncient(), cmd/geth/ OnShutdown(). These are covered by stand alone standard go tests which can all be run by navigating to the respective directories 

Note: depending on where the script is deployed Geth may complain that the path to the `.ipc` file is too long. Renaming of the directories or moving the project may be necessary. 