#!/usr/bin/env node
import { ethers, ContractFactory } from "ethers"
import * as fs from 'fs'

const g_server = "http://localhost:8545"


async function getWallet(index: number): Promise<ethers.Wallet> {
    let path = `m/44'/60'/0'/0/${index}`;
    const mymnemonics = process.env.MYMNEMONICS ?? ''
    const provider = new ethers.providers.JsonRpcProvider(g_server);
    let walletMnemonic = ethers.Wallet.fromMnemonic(mymnemonics, path);
    let w = walletMnemonic.connect(provider)
    return w
}
async function getAddress(index: number): Promise<string> {
    let walletMnemonic = await getWallet(index)
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

async function createContract(): Promise<string> {
    const myaddress = await getAddress(0)
    const balance = await getBalance(myaddress)
    console.log("Hello World")
    console.log(`${myaddress} balance ${balance}`)

    const contractByteCode = fs.readFileSync('hello_sol_Hello.bin', 'utf-8')
    console.log(`binary ${contractByteCode}`)

    const contractAbi = JSON.parse(fs.readFileSync('hello_sol_Hello.abi', 'utf-8'))
    console.log(`abi ${contractAbi}`)

    const provider = new ethers.providers.JsonRpcProvider(g_server);
    const signer = provider.getSigner()
    console.log(`address ${await signer.getAddress()}`)
    const factory = new ContractFactory(contractAbi, contractByteCode, signer)

    // If your contract requires constructor args, you can specify them here
    const contract = await factory.deploy()

    console.log(contract)
    // console.log(contract.deployTransaction)
    return contract.address

}


async function processContract(contractAddress: string) {
    const contractByteCode = fs.readFileSync('hello_sol_Hello.bin', 'utf-8')

    const contractAbi = JSON.parse(fs.readFileSync('hello_sol_Hello.abi', 'utf-8'))
    const provider = new ethers.providers.JsonRpcProvider(g_server);
    console.log(`contract address ${contractAddress}`)
    const contractInstance = new ethers.Contract(contractAddress, contractAbi, provider)

    let currentValue = await contractInstance.retrieve();
    console.log(currentValue);

    let wallet = await getWallet(0)
    let contractWithSigner = contractInstance.connect(wallet);
    let tx = await contractWithSigner.store(ethers.BigNumber.from("0x15"))
    console.log(tx)


    currentValue = await contractInstance.retrieve();
    console.log(currentValue);

}


async function run() {

    let contractAddress = await createContract()
    console.log(`contract address ${contractAddress}`)
    processContract(contractAddress)
}

run()