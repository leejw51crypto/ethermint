#!/usr/bin/env node
import { ethers } from "ethers"


async function test() {
    const provider = new ethers.providers.JsonRpcProvider("http://localhost:8545");
    const b = await provider.getBlockNumber();
    console.log(`${b}`)
    console.log("OK")
}

test();