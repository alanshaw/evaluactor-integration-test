# evaluactor-integration-test

A test for the [evaluactor](https://github.com/alanshaw/evaluactor).

## What did you do?

1. I used [test vectors](https://github.com/filecoin-project/test-vectors) to generate a state tree (`state.car`).
    * In that state tree the evaluactor is at address `f097`.
    * There's also an account at `f0100` with `1,000,000,000,000,000` attoFIL in their wallet.
1. I copied bits of the lotus conformance test driver, and loaded the `state.car` into a blockstore.
1. I instantiated a Lotus VM and registered the evaluactor actor.
1. I crafted a message to `f097` from `f0100` calling method 2 (`evaluactor.Eval()`) with some lua in the params.
1. I call `vm.ApplyMessage()` to have the VM execute the message and then print out the result.
