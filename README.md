# gnark-keep

   Using gnark v0.7 to complete:
1. cubic equations
2. score up
3. rollup of transfer

#### motivation
gnark has been applied in zkevm linea of ConsenSys and zkbnb on binance's layer2 network, and has a high proof efficiency;
Here we will try to build a simple zk-rollup app, The purpose is to support client-side(web、mobile devices) zk proof generation.

#### architecture
We will use golang to implement the account asset operation circuit. In order to deal with zk-friendly hash, we will use Wasm implements browser wallet;