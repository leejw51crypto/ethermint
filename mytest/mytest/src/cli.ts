#!/usr/bin/env node
import { ethers } from "ethers"


async function test() {
    const provider = new ethers.providers.JsonRpcProvider("http://localhost:8545");
    const b = await provider.getBlockNumber();
    console.log(`${b}`)
    const balance = await provider.getBalance("0x59E9A2D1F17E1970BF3843F11CE0BB7419E48D43")
    console.log("OK")
    console.log(`balance ${balance}`)
}

test();