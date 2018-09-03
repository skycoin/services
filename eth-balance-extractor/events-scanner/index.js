'use strict';

// import Web3 from 'web3';
// import parse from 'csv-parse';

var parse = require('csv-parse');
var Web3 = require('web3');
var fs = require('fs');

var web3 = new Web3();

/// #Parameters:
const nodeUrl = 'http://localhost:8545'; // URL to the Ethereum node API
const contractAddress = "0xf230b790e05390fc8295f4d3f60332c93bed42e2"; // Address of the smart contract

const enableTransactionCaching = false; // Saves transaction hashes if they from field is equal to event from field.
const firstBlock = 4212160; // Block that contains transaction with deploying of smart contract
const lastBlock = 5950982; // Last added block (if useLatestBlock=true this variable must be set to the nearest latest block of blockchain).
const offset = 30000; // Amount of blocks that are processed in one transaction

const useLatestBlock = true; // If true - last iteration of the scanning process will use last block (from Eth node, not from lastBlock variable above)
const logProcessedEvents = false; // If true - all processed events will be saved

// ABI of the smart contract.
// See this answer (https://ethereum.stackexchange.com/a/3150/42340) for more details on how to obtain ABI.
const abi = [
	{
		"constant": true,
		"inputs": [],
		"name": "name",
		"outputs": [
			{
				"name": "",
				"type": "string"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [],
		"name": "stop",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "_spender",
				"type": "address"
			},
			{
				"name": "_value",
				"type": "uint256"
			}
		],
		"name": "approve",
		"outputs": [
			{
				"name": "success",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "totalSupply",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "_from",
				"type": "address"
			},
			{
				"name": "_to",
				"type": "address"
			},
			{
				"name": "_value",
				"type": "uint256"
			}
		],
		"name": "transferFrom",
		"outputs": [
			{
				"name": "success",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "decimals",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "_value",
				"type": "uint256"
			}
		],
		"name": "burn",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"name": "balanceOf",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "stopped",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "symbol",
		"outputs": [
			{
				"name": "",
				"type": "string"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "_to",
				"type": "address"
			},
			{
				"name": "_value",
				"type": "uint256"
			}
		],
		"name": "transfer",
		"outputs": [
			{
				"name": "success",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [],
		"name": "start",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "_name",
				"type": "string"
			}
		],
		"name": "setName",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "",
				"type": "address"
			},
			{
				"name": "",
				"type": "address"
			}
		],
		"name": "allowance",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"name": "_addressFounder",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "constructor"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"name": "_from",
				"type": "address"
			},
			{
				"indexed": true,
				"name": "_to",
				"type": "address"
			},
			{
				"indexed": false,
				"name": "_value",
				"type": "uint256"
			}
		],
		"name": "Transfer",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"name": "_owner",
				"type": "address"
			},
			{
				"indexed": true,
				"name": "_spender",
				"type": "address"
			},
			{
				"indexed": false,
				"name": "_value",
				"type": "uint256"
			}
		],
		"name": "Approval",
		"type": "event"
	}
];

/// End of #Parameters section
web3.setProvider(new web3.providers.HttpProvider(nodeUrl));
var contract = new web3.eth.Contract(abi, contractAddress);

let eventsCount = 0;

const wallets = {};

const processEvents = events => {
	events.forEach(async e => {
		if (e.removed
			|| !e.returnValues._value) {
			console.log(e);
		}

		const from = e.returnValues._from.toLowerCase();
		const to = e.returnValues._to.toLowerCase();

		if (e.transactionHash !== '0x9299ff8878c8a22306be0b79a740b6dcc76a4dc5307717a292879923930ba9bc') {
			if (wallets[from]) {
				wallets[from].balance = wallets[from].balance.add(web3.utils.toBN(e.returnValues._value).neg());
				wallets[from].transactions++;
				if (enableTransactionCaching && !wallets[from].transactionHash) {
					try {
						const t = await web3.eth.getTransaction(e.transactionHash);
						if (t.from.toLowerCase() === from)
							wallets[from].transactionHash = e.transactionHash;
					}
					catch (e) { }
				}
			} else {
				wallets[from] = { balance: web3.utils.toBN(e.returnValues._value).neg(), transactionHash: null, transactions: 1 };

				if (enableTransactionCaching) {
					try {
						const t = await web3.eth.getTransaction(e.transactionHash);
						if (t.from.toLowerCase() === from)
							wallets[from].transactionHash = e.transactionHash;
					}
					catch (e) { }
				}
			}

			if (wallets[to]) {
				wallets[to].balance = wallets[to].balance.add(web3.utils.toBN(e.returnValues._value));
				wallets[to].transactions++;
			} else {
				wallets[to] = { balance: web3.utils.toBN(e.returnValues._value), transactionHash: null, transactions: 1, }
			}
		}
	});
}

const saveEvents = (iteration, events) => {
	const data = events.map(e => JSON.stringify(e) + "\n").reduce((acc, i) => (`${acc}${i}`));
	fs.writeFile('events/' + iteration + '.json', data, 'utf8', err => {
		console.error(err);
	});
}

const walletsToCSV = (filename, wallets) => {
	const data = wallets
		.map(w => `${w.wallet},${w.transactionHash},${w.balance},${w.transactions}\n`)
		.reduce((acc, i) => `${acc}${i}`);

	fs.writeFile(filename, data, 'utf8', err => {
		console.error(err);
	});
};

const saveWallets = wallets => {
	const keys = Object.keys(wallets);
	const array = keys.map(k => ({
		wallet: k,
		balance: web3.utils.toBN(wallets[k].balance),
		transactions: wallets[k].transactions,
		transactionHash: wallets[k].transactionHash,
	}));

	const zero = web3.utils.toBN(0);
	const zeroWallets = array.filter(w => w.balance.eq(zero));
	const gtZeroWallets = array.filter(w => w.balance.gt(zero));
	const ltZeroWallets = array.filter(w => w.balance.lt(zero));

	walletsToCSV('zero_wallets.csv', zeroWallets);
	walletsToCSV('positive_balance_wallets.csv', gtZeroWallets);
	walletsToCSV('negative_balance_wallets.csv', ltZeroWallets);
};

const run = async () => {
	for (var i = firstBlock; i < lastBlock;) {
		console.log("Iteration: " + i);
		try {
			const toBlock = i + offset > lastBlock
				? useLatestBlock ? 'latest' : lastBlock
				: i + offset;
			const events = await contract.getPastEvents("Transfer",
				{
					fromBlock: i + 1,
					toBlock,
				});

			console.log("Total: " + eventsCount + " Current iteration: " + events.length);
			console.log("Wallets count: " + Object.keys(wallets).length);
			eventsCount += events.length;
			await processEvents(events);
			if (logProcessedEvents) {
				saveEvents(i, events);
			}

			i += offset;
		} catch (e) {
			console.error(e);
		}
	}

	saveWallets(wallets);
}

run();
