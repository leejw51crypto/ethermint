#!/usr/bin/env node
import { ethers } from "ethers"
import * as fs from 'fs'

const g_server = "http://localhost:8545"

async function getAddress(index: number): Promise<string> {
    let path = `m/44'/60'/0'/0/${index}`;
    const mymnemonics = process.env.MYMNEMONICS ?? ''

    let walletMnemonic = ethers.Wallet.fromMnemonic(mymnemonics, path);

    const myaddress = await walletMnemonic.getAddress()
    return myaddress;
}

async function getBalance(myaddress: string): Promise<string> {
    const provider = new ethers.providers.JsonRpcProvider(g_server);
    const mybalance = await provider.getBalance(myaddress)
    return mybalance.toString()
}
async function getBlockNumber(): Promise<string> {
    const provider = new ethers.providers.JsonRpcProvider(g_server);
    const b = await provider.getBlockNumber();
    return b.toString()
}

async function sendTx(fromaddr: string, toaddr: string): Promise<string> {
    const provider = new ethers.providers.JsonRpcProvider(g_server);
    const signer = provider.getSigner()
    const tx = await signer.sendTransaction({
        from: fromaddr,
        to: toaddr,
        value: ethers.utils.parseEther("1.0")
    });
    console.log(`tx ${JSON.stringify(tx)}`)


    return ""
}
async function test2() {

    const myaddress = await getAddress(0)
    const myaddress2 = await getAddress(1)

    const b = await getBlockNumber();
    console.log(`${b}`)
    const balance = await getBalance(myaddress)
    const balance2 = await getBalance(myaddress2)
    console.log(`${myaddress} balance ${balance}`)
    console.log(`${myaddress2} balance ${balance2}`)

    await sendTx(myaddress, myaddress2)
    //const v1 = ethers.utils.parseEther("1.0")
    //console.log(`v =${v1}`)

}

async function testHelloWord() {
    const myaddress = await getAddress(0)
    const balance = await getBalance(myaddress)
    console.log("Hello World")
    console.log(`${myaddress} balance ${balance}`)

    const contractBinary = fs.readFileSync('hello_sol_Hello.bin')
    console.log(`binary ${contractBinary}`)

    const contractAbi = fs.readFileSync('hello_sol_Hello.abi')
    console.log(`abi ${contractAbi}`)

}


testHelloWord();
