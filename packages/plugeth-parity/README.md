# PluGeth-Parity trace plugin suite.

This plugin provides facsimiles of four tracing methods available from the [OpenEthereum](https://openethereum.github.io/JSONRPC-trace-module) project. At this point we have implementations for the following methods:

```
trace_call

trace_rawTransaction

trace_replayTransaction

trace_replayBlockTransactions
```

The plugin can be [built](https://docs.plugeth.org/en/latest/build.html) like any other PluGeth plugin. Once built just point towards a PluGeth node and they will take the same arguments as the OpenEthereum documentation specifies.

 ## Trace Variants

 Each method can be executed with one to three of the following diagnostics:

 ```
 trace

 vmTrace

 stateDiff
 ```
 #### Known Issues

 This is a beta release and as such we encourage any users to test the plugin before deploying into production.

 Throughout our development process we came to the conclusion that OpenEthereum's *tracers* do not properly implement [EIP-2929](https://eips.ethereum.org/EIPS/eip-2929), OpenEthereum still seems to be able to process blocks post-EIP-2929. As a result the ``used``(gas used) reported on contract calls in ``vmTrace`` is not accurate. Also, ``stateDiff`` on Clique networks incorrectly reports the miner address as the zero address and the balance change as a change to the balance of the zero address. All of our development was done on the Goerli test net and we believe that this issue will only effect Clique networks. We have chosen to not recreate either of these behaviors. If any users have a need that these behaviors remain intact we invite them to fork the project and develop their own versions.  

 During development we included opcodes in the return values for ``vmTrace``. We left the implementation intact so that it could be used for future development or debugging. Opcode reporting can be turned on by eliminating the hyphen from ``json:"-"`` on line 27 of vmTrace.go.

 We encourage all users and developers to get in touch with us on [discord](https://docs.plugeth.org/en/latest/contact.html) to help us continue to refine the accuracy of the plugin and to learn about how the plugin is being used.  
