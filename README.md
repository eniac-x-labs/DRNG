# DRNG - Dicentralized Random Number Gernerator

DRNG is the abbreviation for Distributed Random Number Generator, which provides users with an unpredictable and
uncontrollable amount of random numbers. The random numbers are generated and supervised jointly by distributed nodes,
and the process is recorded by the blockchain to support auditing.

DRNG requires one leader node and at least one generator node.

DRNG provides APIs and SDKs to facilitate user calls and integration.

DRNG stores the hash of the random numbers generated during the process and the hash of the final random number results
on the blockchain. This ensures that the random numbers are authentic, transparent, and verifiable without exposing the
numbers themselves. Storing the hash on the blockchain is not a mandatory requirement for the system to operate. When
conditions do not permit, the blockchain component can be configured as non-essential. However, it can also be set as
essential, in which case users need to deploy smart contracts on the blockchain in advance and set the contract address
in the configuration file.

- Build

```shell
make build
```

The compiled code and executable files are placed by default in the `bin` directory.

- Run Node

```shell
./bin/DRNG-service --confdir [config-dir-path] --confname [config-file-name]
```

Such as `./bin/DRNG-service --confdir ./tests/conf --confname config-1`

- Apis

| Funcs                                      | Type |
|:-------------------------------------------|:-----|
| /getSingleRandomNumber                     | GET  |
| /getMultiRandomNumbers/{num}               | GET  |
| /getPoolUpperLimit                         | GET  |
| /getPoolLatestSeq                          | GET  |
| /getLatestProcessingSeq                    | GET  |
| /getGeneratorsNum                          | GET  |
| /getLastDiscardSeq                         | GET  |
| /getSizeofPool                             | GET  |
| /isEvil/{node}                             | GET  |
| /getEvils                                  | GET  |
| /getBlackListLen                           | GET  |
| /freeFromBlackList/{node}                  | GET  |
| /getFailureTimesThreshold                  | GET  |
| /examineRandNumWithLeaderAndGeneratorsLogs | POST |

