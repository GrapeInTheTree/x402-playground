package quiz

// SolidityQuestions returns all Solidity quiz questions organized into 7 modules.
func SolidityQuestions() []Question {
	var all []Question
	all = append(all, solModule1Foundations()...)
	all = append(all, solModule2ERC20()...)
	all = append(all, solModule3Signatures()...)
	all = append(all, solModule4Gasless()...)
	all = append(all, solModule5Advanced()...)
	all = append(all, solModule6X402()...)
	all = append(all, solModule7ERC8004()...)
	return all
}

// ============================================================
// MODULE 1: Foundations
// ============================================================

func solModule1Foundations() []Question {
	return []Question{
		solTypesVars(),
		solFunctions(),
		solControlFlow(),
		solMsgBlock(),
	}
}

func solTypesVars() Question {
	return Question{
		ID: "sol-types-vars", Title: "Types & Variables",
		Difficulty: "easy", Category: "M1: Foundations", Language: LangSolidity,
		Description: `Solidity is a statically-typed language for the EVM. The core value types are:

- uint256: unsigned 256-bit integer, the default for most on-chain math
- address: 20-byte Ethereum address, with .balance and .transfer() methods
- bool: true/false
- bytes32: fixed-size 32-byte value, used for hashes and nonces
- mapping: hash table from key type to value type (storage only)

Implement a simple storage contract that demonstrates these types.
Store an owner address, a counter, and a mapping of addresses to values.`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract TypesDemo {
    // TODO 1: Declare a public address variable called 'owner'
    // TODO 2: Declare a public uint256 variable called 'counter'
    // TODO 3: Declare a public bool variable called 'active'
    // TODO 4: Declare a public bytes32 variable called 'dataHash'
    // TODO 5: Declare a mapping from address to uint256 called 'balances' (public)

    constructor() {
        // TODO 6: Set owner to msg.sender
        // TODO 7: Set active to true
    }

    function increment() external {
        // TODO 8: Increase counter by 1
    }

    function setBalance(address account, uint256 amount) external {
        // TODO 9: Set balances[account] = amount
    }

    function setDataHash(bytes32 hash) external {
        // TODO 10: Set dataHash = hash
    }
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

contract TypesDemoTest is Test {
    TypesDemo demo;
    address deployer = address(1);

    function setUp() public {
        vm.prank(deployer);
        demo = new TypesDemo();
    }

    function test_Owner() public view {
        assertEq(demo.owner(), deployer);
    }

    function test_Active() public view {
        assertTrue(demo.active());
    }

    function test_Increment() public {
        assertEq(demo.counter(), 0);
        demo.increment();
        assertEq(demo.counter(), 1);
        demo.increment();
        assertEq(demo.counter(), 2);
    }

    function test_SetBalance() public {
        demo.setBalance(address(2), 1000);
        assertEq(demo.balances(address(2)), 1000);
    }

    function test_SetDataHash() public {
        bytes32 h = keccak256("hello");
        demo.setDataHash(h);
        assertEq(demo.dataHash(), h);
    }
}
`,
		Hints: []string{
			"address public owner; uint256 public counter; bool public active;",
			"bytes32 public dataHash; mapping(address => uint256) public balances;",
			"constructor: owner = msg.sender; active = true;",
		},
	}
}

func solFunctions() Question {
	return Question{
		ID: "sol-functions", Title: "Functions & Visibility",
		Difficulty: "easy", Category: "M1: Foundations", Language: LangSolidity,
		Description: `Solidity functions have four visibility levels:

- public: callable from anywhere (generates automatic getter for state vars)
- external: only callable from outside the contract (more gas efficient for large args)
- internal: only callable from this contract and derived contracts
- private: only callable from this contract

Function modifiers: view (reads state, no writes), pure (no state access), payable (accepts ETH).

Implement a calculator contract with different visibility functions.`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract Calculator {
    uint256 public result;

    // TODO 1: Write an external function 'add(uint256 a, uint256 b)'
    //         that stores a+b in result and returns the sum

    // TODO 2: Write an external function 'multiply(uint256 a, uint256 b)'
    //         that stores a*b in result and returns the product

    // TODO 3: Write a pure public function 'pureAdd(uint256 a, uint256 b)'
    //         that returns a+b without touching storage

    // TODO 4: Write a view public function 'getResult()'
    //         that returns the current result
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

contract CalculatorTest is Test {
    Calculator calc;

    function setUp() public {
        calc = new Calculator();
    }

    function test_Add() public {
        uint256 sum = calc.add(3, 7);
        assertEq(sum, 10);
        assertEq(calc.result(), 10);
    }

    function test_Multiply() public {
        uint256 prod = calc.multiply(4, 5);
        assertEq(prod, 20);
        assertEq(calc.result(), 20);
    }

    function test_PureAdd() public pure {
        assertEq(Calculator(address(0)).pureAdd(100, 200), 300);
    }

    function test_GetResult() public {
        calc.add(10, 20);
        assertEq(calc.getResult(), 30);
    }
}
`,
		Hints: []string{
			"function add(uint256 a, uint256 b) external returns (uint256) { result = a + b; return result; }",
			"pure functions cannot read or write state: function pureAdd(...) public pure returns (uint256)",
			"view functions can read but not write: function getResult() public view returns (uint256)",
		},
	}
}

func solControlFlow() Question {
	return Question{
		ID: "sol-control-flow", Title: "Control Flow & Events",
		Difficulty: "easy", Category: "M1: Foundations", Language: LangSolidity,
		Description: `Solidity has three ways to handle errors:

- require(condition, "message"): reverts with message if condition is false
- revert("message"): always reverts
- Custom errors: error InsufficientBalance(uint256 available, uint256 required)

Events are the primary way contracts communicate with off-chain applications:
- event Transfer(address indexed from, address indexed to, uint256 value)
- emit Transfer(from, to, value)
- 'indexed' parameters are searchable in logs (max 3 per event)

Implement a vault with deposits, withdrawals, and proper error handling.`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract Vault {
    mapping(address => uint256) public deposits;

    // TODO 1: Define a custom error 'InsufficientBalance(uint256 available, uint256 required)'
    // TODO 2: Define an event 'Deposited(address indexed user, uint256 amount)'
    // TODO 3: Define an event 'Withdrawn(address indexed user, uint256 amount)'

    function deposit() external payable {
        // TODO 4: require msg.value > 0
        // TODO 5: Add msg.value to deposits[msg.sender]
        // TODO 6: Emit Deposited event
    }

    function withdraw(uint256 amount) external {
        // TODO 7: If deposits[msg.sender] < amount, revert with InsufficientBalance
        // TODO 8: Subtract amount from deposits[msg.sender]
        // TODO 9: Transfer ETH back to msg.sender
        // TODO 10: Emit Withdrawn event
    }
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

contract VaultTest is Test {
    Vault vault;
    address alice = address(1);

    function setUp() public {
        vault = new Vault();
        vm.deal(alice, 10 ether);
    }

    function test_Deposit() public {
        vm.prank(alice);
        vault.deposit{value: 1 ether}();
        assertEq(vault.deposits(alice), 1 ether);
    }

    function test_DepositZero() public {
        vm.prank(alice);
        vm.expectRevert();
        vault.deposit{value: 0}();
    }

    function test_Withdraw() public {
        vm.prank(alice);
        vault.deposit{value: 2 ether}();

        uint256 before = alice.balance;
        vm.prank(alice);
        vault.withdraw(1 ether);

        assertEq(vault.deposits(alice), 1 ether);
        assertEq(alice.balance, before + 1 ether);
    }

    function test_WithdrawInsufficient() public {
        vm.prank(alice);
        vm.expectRevert();
        vault.withdraw(1 ether);
    }
}
`,
		Hints: []string{
			"error InsufficientBalance(uint256 available, uint256 required);",
			"event Deposited(address indexed user, uint256 amount); emit Deposited(msg.sender, msg.value);",
			"revert InsufficientBalance(deposits[msg.sender], amount); payable(msg.sender).transfer(amount);",
		},
	}
}

func solMsgBlock() Question {
	return Question{
		ID: "sol-msg-block", Title: "msg & block Globals",
		Difficulty: "easy", Category: "M1: Foundations", Language: LangSolidity,
		Description: `Solidity provides global variables for transaction and block context:

- msg.sender: the address calling this function (can be an EOA or contract)
- msg.value: amount of wei sent with the call (only in payable functions)
- block.timestamp: current block's Unix timestamp in seconds
- block.number: current block number

These are essential for access control, time-based logic, and payment handling.
In x402, block.timestamp is used for validAfter/validBefore in EIP-3009.

Build a time-lock contract that releases funds after a deadline.`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract TimeLock {
    address public beneficiary;
    uint256 public releaseTime;
    uint256 public depositBlock;

    // TODO 1: Write a constructor that takes (address _beneficiary, uint256 _releaseTime)
    //         Store both, and also store block.number as depositBlock

    function deposit() external payable {
        // TODO 2: require msg.value > 0
        // TODO 3: require msg.sender is not the beneficiary
    }

    function release() external {
        // TODO 4: require block.timestamp >= releaseTime
        // TODO 5: require address(this).balance > 0
        // TODO 6: Send all ETH to beneficiary using .transfer()
    }

    function getBalance() external view returns (uint256) {
        // TODO 7: Return this contract's ETH balance
        return 0;
    }

    function timeRemaining() external view returns (uint256) {
        // TODO 8: If releaseTime > block.timestamp, return the difference; else return 0
        return 0;
    }
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

contract TimeLockTest is Test {
    TimeLock lock;
    address alice = address(1);
    address bob = address(2);

    function setUp() public {
        vm.warp(1000);
        lock = new TimeLock(alice, 2000);
        vm.deal(bob, 10 ether);
    }

    function test_Constructor() public view {
        assertEq(lock.beneficiary(), alice);
        assertEq(lock.releaseTime(), 2000);
        assertGt(lock.depositBlock(), 0);
    }

    function test_Deposit() public {
        vm.prank(bob);
        lock.deposit{value: 1 ether}();
        assertEq(lock.getBalance(), 1 ether);
    }

    function test_ReleaseBeforeTime() public {
        vm.prank(bob);
        lock.deposit{value: 1 ether}();
        vm.warp(1500);
        vm.expectRevert();
        lock.release();
    }

    function test_ReleaseAfterTime() public {
        vm.prank(bob);
        lock.deposit{value: 1 ether}();
        vm.warp(2000);
        uint256 before = alice.balance;
        lock.release();
        assertEq(alice.balance, before + 1 ether);
    }

    function test_TimeRemaining() public {
        vm.warp(1500);
        assertEq(lock.timeRemaining(), 500);
        vm.warp(2500);
        assertEq(lock.timeRemaining(), 0);
    }
}
`,
		Hints: []string{
			"constructor(address _beneficiary, uint256 _releaseTime) { beneficiary = _beneficiary; releaseTime = _releaseTime; depositBlock = block.number; }",
			"require(block.timestamp >= releaseTime); payable(beneficiary).transfer(address(this).balance);",
			"return address(this).balance; / if (releaseTime > block.timestamp) return releaseTime - block.timestamp; else return 0;",
		},
	}
}

// ============================================================
// MODULE 2: ERC-20 Token
// ============================================================

func solModule2ERC20() []Question {
	return []Question{
		solERC20Basic(),
		solERC20Allowance(),
		solERC20Events(),
		solERC20Metadata(),
	}
}

func solERC20Basic() Question {
	return Question{
		ID: "sol-erc20-basic", Title: "Basic Token",
		Difficulty: "easy", Category: "M2: ERC-20", Language: LangSolidity,
		Description: `ERC-20 is the standard interface for fungible tokens on Ethereum.
USDC, the token used in x402 payments, is an ERC-20 token.

The core ERC-20 functions are:
- totalSupply(): total number of tokens in existence
- balanceOf(address): token balance of an account
- transfer(address to, uint256 amount): move tokens to another address

Implement a minimal ERC-20 with a constructor that mints the
initial supply to the deployer.`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract SimpleToken {
    mapping(address => uint256) private _balances;
    uint256 public totalSupply;

    constructor(uint256 initialSupply) {
        // TODO 1: Mint initialSupply to msg.sender
        // Set _balances[msg.sender] = initialSupply
        // Set totalSupply = initialSupply
    }

    function balanceOf(address account) public view returns (uint256) {
        // TODO 2: Return the balance of account
        return 0;
    }

    function transfer(address to, uint256 amount) public returns (bool) {
        // TODO 3: Transfer amount from msg.sender to 'to'
        // 1. Check msg.sender has enough balance (require)
        // 2. Subtract from sender, add to receiver
        // 3. Return true
        return false;
    }
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

contract SimpleTokenTest is Test {
    SimpleToken token;
    address alice = address(1);
    address bob = address(2);

    function setUp() public {
        vm.prank(alice);
        token = new SimpleToken(1000000);
    }

    function test_InitialBalance() public view {
        assertEq(token.balanceOf(alice), 1000000);
        assertEq(token.balanceOf(bob), 0);
    }

    function test_TotalSupply() public view {
        assertEq(token.totalSupply(), 1000000);
    }

    function test_Transfer() public {
        vm.prank(alice);
        bool ok = token.transfer(bob, 100000);
        assertTrue(ok);
        assertEq(token.balanceOf(alice), 900000);
        assertEq(token.balanceOf(bob), 100000);
    }

    function test_TransferInsufficientBalance() public {
        vm.prank(bob);
        vm.expectRevert();
        token.transfer(alice, 1);
    }
}
`,
		Hints: []string{
			"constructor: _balances[msg.sender] = initialSupply; totalSupply = initialSupply;",
			"balanceOf: return _balances[account];",
			"transfer: require(_balances[msg.sender] >= amount); then subtract and add",
		},
	}
}

func solERC20Allowance() Question {
	return Question{
		ID: "sol-erc20-allowance", Title: "Approval System",
		Difficulty: "medium", Category: "M2: ERC-20", Language: LangSolidity,
		Description: `The ERC-20 approval system lets a token owner authorize a third party
(the "spender") to transfer tokens on their behalf. This is the foundation
for all DeFi protocols — Uniswap, Aave, and Permit2 all rely on it.

Three functions work together:
- approve(spender, amount): owner authorizes spender for up to amount tokens
- allowance(owner, spender): query the current allowance
- transferFrom(from, to, amount): spender moves tokens from owner to recipient

In x402, the Facilitator uses transferFrom to settle payments after
verifying the client's EIP-712 signature.`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract TokenWithApproval {
    mapping(address => uint256) private _balances;
    mapping(address => mapping(address => uint256)) private _allowances;
    uint256 public totalSupply;

    constructor(uint256 initialSupply) {
        _balances[msg.sender] = initialSupply;
        totalSupply = initialSupply;
    }

    function balanceOf(address account) public view returns (uint256) {
        return _balances[account];
    }

    function transfer(address to, uint256 amount) public returns (bool) {
        require(_balances[msg.sender] >= amount, "insufficient balance");
        _balances[msg.sender] -= amount;
        _balances[to] += amount;
        return true;
    }

    function approve(address spender, uint256 amount) public returns (bool) {
        // TODO 1: Set allowance for spender to spend msg.sender's tokens
        // Store in _allowances[msg.sender][spender]
        return false;
    }

    function allowance(address owner, address spender) public view returns (uint256) {
        // TODO 2: Return the allowance
        return 0;
    }

    function transferFrom(address from, address to, uint256 amount) public returns (bool) {
        // TODO 3: Transfer tokens from 'from' to 'to' on behalf of msg.sender
        // 1. Check allowance: _allowances[from][msg.sender] >= amount
        // 2. Check balance: _balances[from] >= amount
        // 3. Subtract allowance, subtract balance, add to receiver
        return false;
    }
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

contract TokenApprovalTest is Test {
    TokenWithApproval token;
    address alice = address(1);
    address bob = address(2);
    address charlie = address(3);

    function setUp() public {
        vm.prank(alice);
        token = new TokenWithApproval(1000000);
    }

    function test_Approve() public {
        vm.prank(alice);
        assertTrue(token.approve(bob, 500000));
        assertEq(token.allowance(alice, bob), 500000);
    }

    function test_TransferFrom() public {
        vm.prank(alice);
        token.approve(bob, 200000);

        vm.prank(bob);
        assertTrue(token.transferFrom(alice, charlie, 100000));

        assertEq(token.balanceOf(alice), 900000);
        assertEq(token.balanceOf(charlie), 100000);
        assertEq(token.allowance(alice, bob), 100000);
    }

    function test_TransferFrom_ExceedAllowance() public {
        vm.prank(alice);
        token.approve(bob, 100);

        vm.prank(bob);
        vm.expectRevert();
        token.transferFrom(alice, charlie, 200);
    }

    function test_TransferFrom_InsufficientBalance() public {
        vm.prank(alice);
        token.approve(bob, type(uint256).max);

        vm.prank(bob);
        vm.expectRevert();
        token.transferFrom(alice, charlie, 2000000);
    }
}
`,
		Hints: []string{
			"approve: _allowances[msg.sender][spender] = amount; return true;",
			"allowance: return _allowances[owner][spender];",
			"transferFrom: check allowance, check balance, then _allowances[from][msg.sender] -= amount",
		},
	}
}

func solERC20Events() Question {
	return Question{
		ID: "sol-erc20-events", Title: "Full ERC-20 with Events",
		Difficulty: "medium", Category: "M2: ERC-20", Language: LangSolidity,
		Description: `A complete ERC-20 token must emit events for all state changes.
The ERC-20 standard defines two required events:

- event Transfer(address indexed from, address indexed to, uint256 value)
  Emitted on transfer() and transferFrom(). For minting, 'from' is address(0).

- event Approval(address indexed owner, address indexed spender, uint256 value)
  Emitted on approve().

Events allow off-chain applications (wallets, explorers, indexers) to track
token movements without querying every block's state. The 'indexed' keyword
makes parameters searchable in Ethereum logs.

Add proper events to a complete ERC-20 implementation.`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract EventToken {
    mapping(address => uint256) private _balances;
    mapping(address => mapping(address => uint256)) private _allowances;
    uint256 public totalSupply;

    // TODO 1: Define event Transfer(address indexed from, address indexed to, uint256 value)
    // TODO 2: Define event Approval(address indexed owner, address indexed spender, uint256 value)

    constructor(uint256 initialSupply) {
        _balances[msg.sender] = initialSupply;
        totalSupply = initialSupply;
        // TODO 3: Emit Transfer from address(0) to msg.sender for minting
    }

    function balanceOf(address account) public view returns (uint256) {
        return _balances[account];
    }

    function transfer(address to, uint256 amount) public returns (bool) {
        require(_balances[msg.sender] >= amount, "insufficient balance");
        _balances[msg.sender] -= amount;
        _balances[to] += amount;
        // TODO 4: Emit Transfer event
        return true;
    }

    function approve(address spender, uint256 amount) public returns (bool) {
        _allowances[msg.sender][spender] = amount;
        // TODO 5: Emit Approval event
        return true;
    }

    function allowance(address owner, address spender) public view returns (uint256) {
        return _allowances[owner][spender];
    }

    function transferFrom(address from, address to, uint256 amount) public returns (bool) {
        require(_allowances[from][msg.sender] >= amount, "insufficient allowance");
        require(_balances[from] >= amount, "insufficient balance");
        _allowances[from][msg.sender] -= amount;
        _balances[from] -= amount;
        _balances[to] += amount;
        // TODO 6: Emit Transfer event for the transfer
        return true;
    }
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

contract EventTokenTest is Test {
    EventToken token;
    address alice = address(1);
    address bob = address(2);

    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);

    function setUp() public {
        vm.prank(alice);
        token = new EventToken(1000000);
    }

    function test_MintEvent() public {
        vm.expectEmit(true, true, false, true);
        emit Transfer(address(0), address(this), 500);
        vm.prank(address(this));
        new EventToken(500);
    }

    function test_TransferEvent() public {
        vm.expectEmit(true, true, false, true);
        emit Transfer(alice, bob, 100);
        vm.prank(alice);
        token.transfer(bob, 100);
    }

    function test_ApprovalEvent() public {
        vm.expectEmit(true, true, false, true);
        emit Approval(alice, bob, 500);
        vm.prank(alice);
        token.approve(bob, 500);
    }

    function test_TransferFromEvent() public {
        vm.prank(alice);
        token.approve(bob, 1000);

        vm.expectEmit(true, true, false, true);
        emit Transfer(alice, address(3), 200);
        vm.prank(bob);
        token.transferFrom(alice, address(3), 200);
    }
}
`,
		Hints: []string{
			"event Transfer(address indexed from, address indexed to, uint256 value);",
			"In constructor: emit Transfer(address(0), msg.sender, initialSupply);",
			"In transfer: emit Transfer(msg.sender, to, amount);",
		},
	}
}

func solERC20Metadata() Question {
	return Question{
		ID: "sol-erc20-metadata", Title: "Token Metadata",
		Difficulty: "easy", Category: "M2: ERC-20", Language: LangSolidity,
		Description: `ERC-20 tokens have optional metadata functions that are universally used:

- name(): human-readable name (e.g., "USD Coin")
- symbol(): trading symbol (e.g., "USDC")
- decimals(): number of decimal places (USDC uses 6, most tokens use 18)

These are crucial for x402: the EIP-712 domain separator's 'name' field
must match the token contract's name() return value EXACTLY.
For Base Sepolia USDC, name() returns "USDC" (not "USD Coin").

Implement a token with configurable metadata.`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract MetadataToken {
    mapping(address => uint256) private _balances;
    uint256 public totalSupply;

    // TODO 1: Declare private string variables _name and _symbol
    // TODO 2: Declare a private uint8 variable _decimals

    constructor(string memory tokenName, string memory tokenSymbol, uint8 tokenDecimals, uint256 initialSupply) {
        // TODO 3: Store name, symbol, decimals
        // TODO 4: Mint initialSupply to msg.sender (adjust by decimals: initialSupply * 10**tokenDecimals)
    }

    function name() public view returns (string memory) {
        // TODO 5: Return _name
        return "";
    }

    function symbol() public view returns (string memory) {
        // TODO 6: Return _symbol
        return "";
    }

    function decimals() public view returns (uint8) {
        // TODO 7: Return _decimals
        return 0;
    }

    function balanceOf(address account) public view returns (uint256) {
        return _balances[account];
    }
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

contract MetadataTokenTest is Test {
    MetadataToken usdc;
    MetadataToken dai;
    address deployer = address(1);

    function setUp() public {
        vm.prank(deployer);
        usdc = new MetadataToken("USDC", "USDC", 6, 1000);
        vm.prank(deployer);
        dai = new MetadataToken("Dai Stablecoin", "DAI", 18, 100);
    }

    function test_USDCMetadata() public view {
        assertEq(usdc.name(), "USDC");
        assertEq(usdc.symbol(), "USDC");
        assertEq(usdc.decimals(), 6);
    }

    function test_USDCSupply() public view {
        assertEq(usdc.totalSupply(), 1000 * 10**6);
        assertEq(usdc.balanceOf(deployer), 1000 * 10**6);
    }

    function test_DAIMetadata() public view {
        assertEq(dai.name(), "Dai Stablecoin");
        assertEq(dai.symbol(), "DAI");
        assertEq(dai.decimals(), 18);
    }

    function test_DAISupply() public view {
        assertEq(dai.totalSupply(), 100 * 10**18);
    }
}
`,
		Hints: []string{
			"string private _name; string private _symbol; uint8 private _decimals;",
			"constructor: _name = tokenName; _symbol = tokenSymbol; _decimals = tokenDecimals;",
			"uint256 supply = initialSupply * 10**uint256(tokenDecimals); _balances[msg.sender] = supply; totalSupply = supply;",
		},
	}
}

// ============================================================
// MODULE 3: Signatures & Hashing
// ============================================================

func solModule3Signatures() []Question {
	return []Question{
		solKeccak256(),
		solEcrecover(),
		solEIP712Domain(),
		solEIP712Struct(),
	}
}

func solKeccak256() Question {
	return Question{
		ID: "sol-keccak256", Title: "Keccak256 Hashing",
		Difficulty: "medium", Category: "M3: Signatures", Language: LangSolidity,
		Description: `Keccak256 is the hash function used throughout Ethereum:
- Address derivation from public keys
- EIP-712 type hashes and struct hashes
- Function selectors (first 4 bytes of keccak256 of signature)

Two encoding methods exist:
- abi.encode(...): ABI-encodes values with padding (deterministic, no collisions)
- abi.encodePacked(...): tightly packed, no padding (shorter, but risk of hash collisions)

Rule of thumb: use abi.encode for hashing structs (EIP-712), abi.encodePacked
only when you need compact encoding and inputs have fixed length.

Implement hashing utilities using both methods.`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract HashUtils {
    // TODO 1: Write a pure function 'hashPacked(string memory a, string memory b)'
    //         that returns keccak256(abi.encodePacked(a, b))

    // TODO 2: Write a pure function 'hashEncoded(address addr, uint256 value)'
    //         that returns keccak256(abi.encode(addr, value))

    // TODO 3: Write a pure function 'functionSelector(string memory sig)'
    //         that returns the first 4 bytes of keccak256(bytes(sig))
    //         Return type: bytes4

    // TODO 4: Write a view function 'hashWithSender(uint256 nonce)'
    //         that returns keccak256(abi.encode(msg.sender, nonce))
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

contract HashUtilsTest is Test {
    HashUtils utils;

    function setUp() public {
        utils = new HashUtils();
    }

    function test_HashPacked() public view {
        bytes32 h = utils.hashPacked("hello", "world");
        assertEq(h, keccak256(abi.encodePacked("hello", "world")));
    }

    function test_HashEncoded() public view {
        bytes32 h = utils.hashEncoded(address(1), 100);
        assertEq(h, keccak256(abi.encode(address(1), 100)));
    }

    function test_FunctionSelector() public view {
        bytes4 sel = utils.functionSelector("transfer(address,uint256)");
        assertEq(sel, bytes4(0xa9059cbb));
    }

    function test_HashWithSender() public {
        vm.prank(address(42));
        bytes32 h = utils.hashWithSender(1);
        assertEq(h, keccak256(abi.encode(address(42), uint256(1))));
    }
}
`,
		Hints: []string{
			"function hashPacked(string memory a, string memory b) public pure returns (bytes32) { return keccak256(abi.encodePacked(a, b)); }",
			"function functionSelector(string memory sig) public pure returns (bytes4) { return bytes4(keccak256(bytes(sig))); }",
			"hashWithSender uses msg.sender: keccak256(abi.encode(msg.sender, nonce))",
		},
	}
}

func solEcrecover() Question {
	return Question{
		ID: "sol-ecrecover", Title: "ECDSA Recovery",
		Difficulty: "hard", Category: "M3: Signatures", Language: LangSolidity,
		Description: `ECDSA (Elliptic Curve Digital Signature Algorithm) is how Ethereum verifies
that a transaction or message was signed by a specific private key.

ecrecover(hash, v, r, s) is a built-in Solidity function that recovers the
signer's address from a message hash and signature components:
- hash: the 32-byte message that was signed
- v: recovery identifier (27 or 28)
- r, s: the signature components (32 bytes each)

For EIP-712 signed data, the hash is: keccak256("\x19\x01" || domainSeparator || structHash)

Security: always check that ecrecover doesn't return address(0), which
indicates an invalid signature.

Build a signature verification contract.`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract SigVerifier {
    // TODO 1: Write a pure function 'recoverSigner(bytes32 hash, uint8 v, bytes32 r, bytes32 s)'
    //         that returns the recovered address
    //         Must revert if recovered address is address(0)

    // TODO 2: Write a pure function 'toEthSignedHash(bytes32 messageHash)'
    //         that returns keccak256("\x19Ethereum Signed Message:\n32" || messageHash)
    //         This is what wallets sign for personal_sign

    // TODO 3: Write a pure function 'verifySignature(bytes32 hash, uint8 v, bytes32 r, bytes32 s, address expected)'
    //         that returns true if the recovered signer matches expected
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

contract SigVerifierTest is Test {
    SigVerifier verifier;
    uint256 constant PRIV_KEY = 0xA11CE;
    address signer;

    function setUp() public {
        verifier = new SigVerifier();
        signer = vm.addr(PRIV_KEY);
    }

    function test_RecoverSigner() public view {
        bytes32 hash = keccak256("test message");
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(PRIV_KEY, hash);
        address recovered = verifier.recoverSigner(hash, v, r, s);
        assertEq(recovered, signer);
    }

    function test_RecoverSigner_Invalid() public {
        bytes32 hash = keccak256("test");
        vm.expectRevert();
        verifier.recoverSigner(hash, 27, bytes32(0), bytes32(0));
    }

    function test_VerifySignature() public view {
        bytes32 hash = keccak256("verify me");
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(PRIV_KEY, hash);
        assertTrue(verifier.verifySignature(hash, v, r, s, signer));
        assertFalse(verifier.verifySignature(hash, v, r, s, address(99)));
    }

    function test_EthSignedHash() public view {
        bytes32 msgHash = keccak256("hello");
        bytes32 ethHash = verifier.toEthSignedHash(msgHash);
        assertEq(ethHash, keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", msgHash)));
    }
}
`,
		Hints: []string{
			"address recovered = ecrecover(hash, v, r, s); require(recovered != address(0));",
			`return keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", messageHash));`,
			"verifySignature: return recoverSigner(hash, v, r, s) == expected;",
		},
	}
}

func solEIP712Domain() Question {
	return Question{
		ID: "sol-eip712-domain", Title: "EIP-712 Domain Separator",
		Difficulty: "hard", Category: "M3: Signatures", Language: LangSolidity,
		Description: `EIP-712 defines typed structured data signing. The domain separator prevents
signature replay across contracts and chains. It's computed as:

  domainSeparator = keccak256(abi.encode(
      DOMAIN_TYPEHASH,
      keccak256(bytes(name)),
      keccak256(bytes(version)),
      chainId,
      verifyingContract
  ))

Where DOMAIN_TYPEHASH = keccak256(
  "EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"
)

For x402: USDC on Base Sepolia uses name="USDC", version="2".
The domain separator is computed once at deployment and cached.

Implement an EIP-712 domain separator contract.`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract EIP712DomainDemo {
    // TODO 1: Define a constant bytes32 DOMAIN_TYPEHASH = keccak256(
    //   "EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"
    // )

    bytes32 public immutable DOMAIN_SEPARATOR;
    string public name;
    string public version;

    constructor(string memory _name, string memory _version) {
        name = _name;
        version = _version;
        // TODO 2: Compute DOMAIN_SEPARATOR using abi.encode(
        //   DOMAIN_TYPEHASH, keccak256(bytes(_name)), keccak256(bytes(_version)),
        //   block.chainid, address(this)
        // )
    }

    // TODO 3: Write a pure function 'computeDomainSeparator' that takes
    //   (string memory _name, string memory _version, uint256 chainId, address verifyingContract)
    //   and returns the domain separator bytes32
    //   This is useful for off-chain computation verification

    // TODO 4: Write a view function 'hashTypedData(bytes32 structHash)'
    //   that returns keccak256(abi.encodePacked("\x19\x01", DOMAIN_SEPARATOR, structHash))
    //   This is the final hash that gets signed in EIP-712
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

contract EIP712DomainTest is Test {
    EIP712DomainDemo domain;

    function setUp() public {
        domain = new EIP712DomainDemo("USDC", "2");
    }

    function test_DomainTypehash() public view {
        bytes32 expected = keccak256("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)");
        assertEq(domain.DOMAIN_TYPEHASH(), expected);
    }

    function test_DomainSeparator() public view {
        bytes32 expected = keccak256(abi.encode(
            keccak256("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"),
            keccak256(bytes("USDC")),
            keccak256(bytes("2")),
            block.chainid,
            address(domain)
        ));
        assertEq(domain.DOMAIN_SEPARATOR(), expected);
    }

    function test_ComputeDomainSeparator() public view {
        bytes32 result = domain.computeDomainSeparator("USDC", "2", block.chainid, address(domain));
        assertEq(result, domain.DOMAIN_SEPARATOR());
    }

    function test_HashTypedData() public view {
        bytes32 structHash = keccak256("test struct");
        bytes32 expected = keccak256(abi.encodePacked("\x19\x01", domain.DOMAIN_SEPARATOR(), structHash));
        assertEq(domain.hashTypedData(structHash), expected);
    }
}
`,
		Hints: []string{
			`bytes32 public constant DOMAIN_TYPEHASH = keccak256("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)");`,
			"DOMAIN_SEPARATOR = keccak256(abi.encode(DOMAIN_TYPEHASH, keccak256(bytes(_name)), keccak256(bytes(_version)), block.chainid, address(this)));",
			`return keccak256(abi.encodePacked("\x19\x01", DOMAIN_SEPARATOR, structHash));`,
		},
	}
}

func solEIP712Struct() Question {
	return Question{
		ID: "sol-eip712-struct", Title: "EIP-712 Struct Hash",
		Difficulty: "hard", Category: "M3: Signatures", Language: LangSolidity,
		Description: `In EIP-712, every struct type has a type hash and each instance has a struct hash:

  TYPEHASH = keccak256("TypeName(type1 name1,type2 name2,...)")
  structHash = keccak256(abi.encode(TYPEHASH, field1, field2, ...))

For EIP-3009's TransferWithAuthorization:
  TYPEHASH = keccak256("TransferWithAuthorization(address from,address to,uint256 value,uint256 validAfter,uint256 validBefore,bytes32 nonce)")

The struct hash is then combined with the domain separator:
  digest = keccak256("\x19\x01" || domainSeparator || structHash)

This digest is what gets signed by the client's private key.

Implement struct hash computation for a transfer authorization.`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract StructHashDemo {
    // TODO 1: Define constant TRANSFER_TYPEHASH = keccak256(
    //   "TransferWithAuthorization(address from,address to,uint256 value,uint256 validAfter,uint256 validBefore,bytes32 nonce)"
    // )

    struct TransferAuth {
        address from;
        address to;
        uint256 value;
        uint256 validAfter;
        uint256 validBefore;
        bytes32 nonce;
    }

    // TODO 2: Write a pure function 'hashTransferAuth(TransferAuth memory auth)'
    //   that returns keccak256(abi.encode(TRANSFER_TYPEHASH, auth.from, auth.to,
    //   auth.value, auth.validAfter, auth.validBefore, auth.nonce))

    // TODO 3: Write a pure function 'hashPermit(address owner, address spender, uint256 value, uint256 nonce, uint256 deadline)'
    //   Define PERMIT_TYPEHASH inline as keccak256("Permit(address owner,address spender,uint256 value,uint256 nonce,uint256 deadline)")
    //   Return keccak256(abi.encode(PERMIT_TYPEHASH, owner, spender, value, nonce, deadline))
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

contract StructHashTest is Test {
    StructHashDemo demo;

    function setUp() public {
        demo = new StructHashDemo();
    }

    function test_TransferTypehash() public view {
        bytes32 expected = keccak256("TransferWithAuthorization(address from,address to,uint256 value,uint256 validAfter,uint256 validBefore,bytes32 nonce)");
        assertEq(demo.TRANSFER_TYPEHASH(), expected);
    }

    function test_HashTransferAuth() public view {
        StructHashDemo.TransferAuth memory auth = StructHashDemo.TransferAuth({
            from: address(1),
            to: address(2),
            value: 100000,
            validAfter: 0,
            validBefore: 999999,
            nonce: bytes32(uint256(42))
        });

        bytes32 expected = keccak256(abi.encode(
            demo.TRANSFER_TYPEHASH(),
            address(1), address(2), uint256(100000),
            uint256(0), uint256(999999), bytes32(uint256(42))
        ));
        assertEq(demo.hashTransferAuth(auth), expected);
    }

    function test_HashPermit() public view {
        bytes32 permitHash = keccak256("Permit(address owner,address spender,uint256 value,uint256 nonce,uint256 deadline)");
        bytes32 expected = keccak256(abi.encode(permitHash, address(1), address(2), uint256(500), uint256(0), uint256(9999)));
        assertEq(demo.hashPermit(address(1), address(2), 500, 0, 9999), expected);
    }
}
`,
		Hints: []string{
			`bytes32 public constant TRANSFER_TYPEHASH = keccak256("TransferWithAuthorization(address from,address to,uint256 value,uint256 validAfter,uint256 validBefore,bytes32 nonce)");`,
			"return keccak256(abi.encode(TRANSFER_TYPEHASH, auth.from, auth.to, auth.value, auth.validAfter, auth.validBefore, auth.nonce));",
			"For Permit: define the typehash inline, then abi.encode all fields",
		},
	}
}

// ============================================================
// MODULE 4: Gasless Transactions
// ============================================================

func solModule4Gasless() []Question {
	return []Question{
		solEIP2612Permit(),
		solEIP3009(),
		solNonceMgmt(),
	}
}

func solEIP2612Permit() Question {
	return Question{
		ID: "sol-eip2612-permit", Title: "EIP-2612 Permit",
		Difficulty: "hard", Category: "M4: Gasless", Language: LangSolidity,
		Description: `EIP-2612 adds a permit() function to ERC-20 tokens, allowing gasless approvals:

Instead of: owner calls approve(spender, amount) on-chain (costs gas)
With permit: owner signs a message off-chain, anyone submits permit() on-chain

The permit function verifies an EIP-712 signature and sets the allowance:
  permit(owner, spender, value, deadline, v, r, s)

It checks:
1. deadline hasn't passed
2. Signature is valid (signed by owner)
3. Nonce matches (prevents replay)
4. Sets allowance just like approve()

This is the precursor to EIP-3009 and is used by many DeFi protocols.

Implement a simplified permit function (signature verification via ecrecover).`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract PermitToken {
    mapping(address => uint256) private _balances;
    mapping(address => mapping(address => uint256)) private _allowances;
    mapping(address => uint256) public nonces;
    uint256 public totalSupply;

    bytes32 public immutable DOMAIN_SEPARATOR;
    bytes32 public constant PERMIT_TYPEHASH = keccak256(
        "Permit(address owner,address spender,uint256 value,uint256 nonce,uint256 deadline)"
    );

    constructor(uint256 initialSupply) {
        _balances[msg.sender] = initialSupply;
        totalSupply = initialSupply;
        DOMAIN_SEPARATOR = keccak256(abi.encode(
            keccak256("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"),
            keccak256(bytes("PermitToken")),
            keccak256(bytes("1")),
            block.chainid,
            address(this)
        ));
    }

    function balanceOf(address account) public view returns (uint256) { return _balances[account]; }
    function allowance(address owner, address spender) public view returns (uint256) { return _allowances[owner][spender]; }

    function permit(
        address owner, address spender, uint256 value,
        uint256 deadline, uint8 v, bytes32 r, bytes32 s
    ) external {
        // TODO 1: require deadline >= block.timestamp
        // TODO 2: Compute structHash = keccak256(abi.encode(PERMIT_TYPEHASH, owner, spender, value, nonces[owner], deadline))
        // TODO 3: Compute digest = keccak256(abi.encodePacked("\x19\x01", DOMAIN_SEPARATOR, structHash))
        // TODO 4: Recover signer via ecrecover(digest, v, r, s)
        // TODO 5: require signer == owner && signer != address(0)
        // TODO 6: Increment nonces[owner]
        // TODO 7: Set _allowances[owner][spender] = value
    }
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

contract PermitTokenTest is Test {
    PermitToken token;
    uint256 constant OWNER_KEY = 0xA11CE;
    address owner;
    address spender = address(2);

    function setUp() public {
        owner = vm.addr(OWNER_KEY);
        vm.prank(owner);
        token = new PermitToken(1000000);
    }

    function _signPermit(uint256 value, uint256 deadline) internal view returns (uint8 v, bytes32 r, bytes32 s) {
        bytes32 structHash = keccak256(abi.encode(
            token.PERMIT_TYPEHASH(), owner, spender, value, token.nonces(owner), deadline
        ));
        bytes32 digest = keccak256(abi.encodePacked("\x19\x01", token.DOMAIN_SEPARATOR(), structHash));
        (v, r, s) = vm.sign(OWNER_KEY, digest);
    }

    function test_Permit() public {
        uint256 deadline = block.timestamp + 1 hours;
        (uint8 v, bytes32 r, bytes32 s) = _signPermit(500000, deadline);

        token.permit(owner, spender, 500000, deadline, v, r, s);
        assertEq(token.allowance(owner, spender), 500000);
        assertEq(token.nonces(owner), 1);
    }

    function test_PermitExpired() public {
        vm.warp(1000);
        (uint8 v, bytes32 r, bytes32 s) = _signPermit(500000, 500);

        vm.expectRevert();
        token.permit(owner, spender, 500000, 500, v, r, s);
    }

    function test_PermitWrongSigner() public {
        uint256 deadline = block.timestamp + 1 hours;
        (uint8 v, bytes32 r, bytes32 s) = _signPermit(500000, deadline);

        vm.expectRevert();
        token.permit(address(99), spender, 500000, deadline, v, r, s);
    }
}
`,
		Hints: []string{
			`require(deadline >= block.timestamp, "permit expired");`,
			`bytes32 structHash = keccak256(abi.encode(PERMIT_TYPEHASH, owner, spender, value, nonces[owner], deadline));`,
			`bytes32 digest = keccak256(abi.encodePacked("\x19\x01", DOMAIN_SEPARATOR, structHash)); address signer = ecrecover(digest, v, r, s);`,
		},
	}
}

func solEIP3009() Question {
	return Question{
		ID: "sol-eip3009", Title: "EIP-3009 transferWithAuthorization",
		Difficulty: "hard", Category: "M4: Gasless", Language: LangSolidity,
		Description: `EIP-3009 is the standard that powers USDC gasless transfers in x402.
Unlike EIP-2612 (which only does gasless approvals), EIP-3009 does
gasless TRANSFERS directly via transferWithAuthorization.

Key differences from EIP-2612:
- Random nonce (bytes32) instead of sequential — no front-running risk
- validAfter/validBefore time window instead of a single deadline
- Transfers tokens directly, not just sets allowances
- Anyone can submit the transaction (the facilitator in x402)

The facilitator never touches the tokens — it just relays the signed
authorization to the USDC contract, which moves funds from Client to PayTo.

Implement the authorization storage, time validation, and nonce tracking.
(Signature verification is simplified for this exercise.)`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract EIP3009Token {
    mapping(address => uint256) private _balances;
    mapping(address => mapping(bytes32 => bool)) private _usedNonces;
    uint256 public totalSupply;

    constructor(uint256 initialSupply) {
        _balances[msg.sender] = initialSupply;
        totalSupply = initialSupply;
    }

    function balanceOf(address account) public view returns (uint256) {
        return _balances[account];
    }

    // Simplified transferWithAuthorization (no actual signature verification)
    // In real USDC, this verifies an EIP-712 signature.
    function transferWithAuthorization(
        address from,
        address to,
        uint256 value,
        uint256 validAfter,
        uint256 validBefore,
        bytes32 nonce
    ) external returns (bool) {
        // TODO 1: Check block.timestamp > validAfter
        // TODO 2: Check block.timestamp < validBefore
        // TODO 3: Check nonce not already used: !_usedNonces[from][nonce]
        // TODO 4: Mark nonce as used
        // TODO 5: Check from has sufficient balance
        // TODO 6: Transfer: subtract from, add to
        // TODO 7: Return true
        return false;
    }

    // Check if a nonce has been used
    function isNonceUsed(address account, bytes32 nonce) public view returns (bool) {
        // TODO 8: Return whether this nonce was already used
        return false;
    }
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

contract EIP3009Test is Test {
    EIP3009Token token;
    address alice = address(1);
    address bob = address(2);
    address facilitator = address(3);

    function setUp() public {
        vm.prank(alice);
        token = new EIP3009Token(1000000);
    }

    function test_TransferWithAuth() public {
        vm.warp(100);
        vm.prank(facilitator);
        bool ok = token.transferWithAuthorization(
            alice, bob, 100000,
            0,         // validAfter
            200,       // validBefore
            bytes32(uint256(1))  // nonce
        );
        assertTrue(ok);
        assertEq(token.balanceOf(alice), 900000);
        assertEq(token.balanceOf(bob), 100000);
    }

    function test_NonceReuse() public {
        bytes32 nonce = bytes32(uint256(42));
        vm.warp(100);

        vm.prank(facilitator);
        token.transferWithAuthorization(alice, bob, 1000, 0, 200, nonce);

        vm.prank(facilitator);
        vm.expectRevert();
        token.transferWithAuthorization(alice, bob, 1000, 0, 200, nonce);
    }

    function test_Expired() public {
        vm.warp(300);
        vm.prank(facilitator);
        vm.expectRevert();
        token.transferWithAuthorization(alice, bob, 1000, 0, 200, bytes32(uint256(2)));
    }

    function test_NotYetValid() public {
        vm.warp(50);
        vm.prank(facilitator);
        vm.expectRevert();
        token.transferWithAuthorization(alice, bob, 1000, 100, 200, bytes32(uint256(3)));
    }

    function test_IsNonceUsed() public {
        bytes32 nonce = bytes32(uint256(99));
        assertFalse(token.isNonceUsed(alice, nonce));

        vm.warp(100);
        vm.prank(facilitator);
        token.transferWithAuthorization(alice, bob, 100, 0, 200, nonce);

        assertTrue(token.isNonceUsed(alice, nonce));
    }
}
`,
		Hints: []string{
			`require(block.timestamp > validAfter, "not yet valid");`,
			`require(!_usedNonces[from][nonce], "nonce used"); then _usedNonces[from][nonce] = true;`,
			`require(_balances[from] >= value); _balances[from] -= value; _balances[to] += value;`,
		},
	}
}

func solNonceMgmt() Question {
	return Question{
		ID: "sol-nonce-mgmt", Title: "Nonce Management",
		Difficulty: "medium", Category: "M4: Gasless", Language: LangSolidity,
		Description: `Nonce management prevents replay attacks in signed messages. Two approaches:

Sequential nonces (EIP-2612):
- Simple counter per address: 0, 1, 2, 3...
- Must be used in order — can't skip ahead
- Risk: front-running can invalidate pending transactions

Random nonces (EIP-3009):
- Random bytes32 value per authorization
- Can be used in any order — no front-running risk
- Stored in mapping(address => mapping(bytes32 => bool))
- Can be cancelled by pre-marking as used

Both approaches also support cancellation:
- Sequential: increment nonce to invalidate pending authorizations
- Random: mark nonce as used without executing the transfer

Implement both nonce strategies and cancellation logic.`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract NonceManager {
    // Sequential nonces (like EIP-2612)
    mapping(address => uint256) public sequentialNonces;

    // Random nonces (like EIP-3009)
    mapping(address => mapping(bytes32 => bool)) public randomNonceUsed;

    // TODO 1: Write 'useSequentialNonce(address account)' that:
    //   - Returns the current nonce
    //   - Increments it by 1
    //   (external function, returns uint256)

    // TODO 2: Write 'useRandomNonce(address account, bytes32 nonce)' that:
    //   - Requires the nonce hasn't been used
    //   - Marks it as used
    //   (external function)

    // TODO 3: Write 'cancelRandomNonce(bytes32 nonce)' that:
    //   - Marks a nonce as used for msg.sender without executing any transfer
    //   - Requires the nonce hasn't been used yet
    //   (external function)

    // TODO 4: Write 'bumpSequentialNonce()' that:
    //   - Increments msg.sender's sequential nonce (invalidates pending ops)
    //   (external function)

    function isRandomNonceUsed(address account, bytes32 nonce) external view returns (bool) {
        return randomNonceUsed[account][nonce];
    }
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

contract NonceManagerTest is Test {
    NonceManager mgr;
    address alice = address(1);

    function setUp() public {
        mgr = new NonceManager();
    }

    function test_SequentialNonce() public {
        uint256 n0 = mgr.useSequentialNonce(alice);
        assertEq(n0, 0);
        uint256 n1 = mgr.useSequentialNonce(alice);
        assertEq(n1, 1);
        assertEq(mgr.sequentialNonces(alice), 2);
    }

    function test_RandomNonce() public {
        bytes32 nonce = bytes32(uint256(42));
        mgr.useRandomNonce(alice, nonce);
        assertTrue(mgr.isRandomNonceUsed(alice, nonce));
    }

    function test_RandomNonce_Replay() public {
        bytes32 nonce = bytes32(uint256(42));
        mgr.useRandomNonce(alice, nonce);
        vm.expectRevert();
        mgr.useRandomNonce(alice, nonce);
    }

    function test_CancelRandomNonce() public {
        bytes32 nonce = bytes32(uint256(99));
        vm.prank(alice);
        mgr.cancelRandomNonce(nonce);
        assertTrue(mgr.isRandomNonceUsed(alice, nonce));
    }

    function test_BumpSequentialNonce() public {
        vm.prank(alice);
        mgr.bumpSequentialNonce();
        assertEq(mgr.sequentialNonces(alice), 1);
    }
}
`,
		Hints: []string{
			"useSequentialNonce: uint256 current = sequentialNonces[account]; sequentialNonces[account]++; return current;",
			"useRandomNonce: require(!randomNonceUsed[account][nonce]); randomNonceUsed[account][nonce] = true;",
			"cancelRandomNonce: require(!randomNonceUsed[msg.sender][nonce]); randomNonceUsed[msg.sender][nonce] = true;",
		},
	}
}

// ============================================================
// MODULE 5: Advanced Patterns
// ============================================================

func solModule5Advanced() []Question {
	return []Question{
		solCREATE2(),
		solPermit2(),
		solProxy(),
		solAccessControl(),
		solReentrancy(),
	}
}

func solCREATE2() Question {
	return Question{
		ID: "sol-create2", Title: "CREATE2 Deployment",
		Difficulty: "hard", Category: "M5: Advanced", Language: LangSolidity,
		Description: `CREATE2 enables deterministic contract deployment — the address is known
before deployment. The address is computed as:

  address = keccak256(0xff || deployer || salt || keccak256(bytecode))[12:]

This is how Permit2 has the same address (0x000000000022D473030F116dDEE9F6B43aC78BA3)
on every EVM chain. The salt and bytecode are fixed, so the address is the same
regardless of which chain it's deployed on.

Uses in x402:
- Permit2 contract: same address on all chains (Base, Polygon, etc.)
- x402Permit2Proxy: deterministic address via CREATE2

Implement a factory that uses CREATE2 for deterministic deployment.`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract SimpleChild {
    address public owner;
    uint256 public value;

    constructor(address _owner, uint256 _value) {
        owner = _owner;
        value = _value;
    }
}

contract CREATE2Factory {
    event Deployed(address indexed addr, bytes32 salt);

    // TODO 1: Write 'deploy(bytes32 salt, address owner, uint256 value)'
    //   that deploys SimpleChild using CREATE2 (new SimpleChild{salt: salt}(owner, value))
    //   Emit Deployed event and return the address
    //   (external function, returns address)

    // TODO 2: Write a view function 'computeAddress(bytes32 salt, address owner, uint256 value)'
    //   that computes the CREATE2 address without deploying
    //   Use: address(uint160(uint256(keccak256(abi.encodePacked(
    //       bytes1(0xff), address(this), salt,
    //       keccak256(abi.encodePacked(type(SimpleChild).creationCode, abi.encode(owner, value)))
    //   )))))
    //   (public view function, returns address)
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

contract CREATE2FactoryTest is Test {
    CREATE2Factory factory;

    function setUp() public {
        factory = new CREATE2Factory();
    }

    function test_Deploy() public {
        bytes32 salt = bytes32(uint256(1));
        address deployed = factory.deploy(salt, address(this), 42);
        assertTrue(deployed != address(0));
        assertEq(SimpleChild(deployed).owner(), address(this));
        assertEq(SimpleChild(deployed).value(), 42);
    }

    function test_ComputeAddress() public {
        bytes32 salt = bytes32(uint256(2));
        address predicted = factory.computeAddress(salt, address(this), 100);
        address actual = factory.deploy(salt, address(this), 100);
        assertEq(predicted, actual);
    }

    function test_DeterministicAddress() public view {
        bytes32 salt = bytes32(uint256(3));
        address a1 = factory.computeAddress(salt, address(1), 50);
        address a2 = factory.computeAddress(salt, address(1), 50);
        assertEq(a1, a2);
    }

    function test_DifferentSalt() public view {
        address a1 = factory.computeAddress(bytes32(uint256(1)), address(1), 50);
        address a2 = factory.computeAddress(bytes32(uint256(2)), address(1), 50);
        assertTrue(a1 != a2);
    }
}
`,
		Hints: []string{
			"address child = address(new SimpleChild{salt: salt}(owner, value));",
			"keccak256(abi.encodePacked(type(SimpleChild).creationCode, abi.encode(owner, value))) gives the bytecode hash",
			"address(uint160(uint256(keccak256(abi.encodePacked(bytes1(0xff), address(this), salt, bytecodeHash)))));",
		},
	}
}

func solPermit2() Question {
	return Question{
		ID: "sol-permit2", Title: "Permit2 SignatureTransfer",
		Difficulty: "hard", Category: "M5: Advanced", Language: LangSolidity,
		Description: `Permit2 is Uniswap's universal token approval protocol. Instead of each
dApp needing its own approve() call, users approve Permit2 once, then
sign EIP-712 messages to authorize individual transfers.

The key function is permitWitnessTransferFrom:
  permitWitnessTransferFrom(permit, transferDetails, owner, witness, witnessTypeString, signature)

In x402:
1. Client approves Permit2 for USDC (one-time on-chain tx)
2. Client signs a Permit2 message authorizing transfer to PayTo
3. Facilitator submits the signed permit to the x402Permit2Proxy
4. Proxy calls Permit2's permitWitnessTransferFrom
5. Permit2 transfers USDC from Client to PayTo

Implement a mock Permit2 with the core transfer-from logic.`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

interface IERC20 {
    function transferFrom(address from, address to, uint256 amount) external returns (bool);
    function allowance(address owner, address spender) external view returns (uint256);
}

contract MockPermit2 {
    // Transfer tokens from owner to recipient via Permit2.
    // This contract must have been approved by the owner on the token contract.
    function transferFrom(
        address token,
        address owner,
        address recipient,
        uint256 amount
    ) external returns (bool) {
        // TODO 1: Check that owner has approved this contract (Permit2) for >= amount
        //          Use IERC20(token).allowance(owner, address(this))
        // TODO 2: Call IERC20(token).transferFrom(owner, recipient, amount)
        // TODO 3: Return true
        return false;
    }

    // Check if owner has sufficient allowance for Permit2
    function hasAllowance(
        address token,
        address owner,
        uint256 amount
    ) external view returns (bool) {
        // TODO 4: Return whether allowance >= amount
        return false;
    }
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

// Simple ERC20 for testing
contract TestToken {
    mapping(address => uint256) public balanceOf;
    mapping(address => mapping(address => uint256)) public allowance;

    constructor(address to, uint256 amount) {
        balanceOf[to] = amount;
    }

    function approve(address spender, uint256 amount) external returns (bool) {
        allowance[msg.sender][spender] = amount;
        return true;
    }

    function transferFrom(address from, address to, uint256 amount) external returns (bool) {
        require(allowance[from][msg.sender] >= amount, "allowance");
        require(balanceOf[from] >= amount, "balance");
        allowance[from][msg.sender] -= amount;
        balanceOf[from] -= amount;
        balanceOf[to] += amount;
        return true;
    }
}

contract MockPermit2Test is Test {
    MockPermit2 permit2;
    TestToken token;
    address alice = address(1);
    address bob = address(2);

    function setUp() public {
        permit2 = new MockPermit2();
        token = new TestToken(alice, 1000000);

        // Alice approves Permit2 contract
        vm.prank(alice);
        token.approve(address(permit2), 500000);
    }

    function test_HasAllowance() public view {
        assertTrue(permit2.hasAllowance(address(token), alice, 100000));
        assertFalse(permit2.hasAllowance(address(token), alice, 600000));
    }

    function test_TransferFrom() public {
        bool ok = permit2.transferFrom(address(token), alice, bob, 100000);
        assertTrue(ok);
        assertEq(token.balanceOf(bob), 100000);
        assertEq(token.balanceOf(alice), 900000);
    }

    function test_TransferFrom_NoAllowance() public {
        vm.expectRevert();
        permit2.transferFrom(address(token), bob, alice, 100);
    }
}
`,
		Hints: []string{
			"IERC20(token).allowance(owner, address(this)) gives the Permit2 allowance",
			"require(IERC20(token).allowance(owner, address(this)) >= amount);",
			"return IERC20(token).transferFrom(owner, recipient, amount);",
		},
	}
}

func solProxy() Question {
	return Question{
		ID: "sol-proxy", Title: "Minimal Proxy (EIP-1167)",
		Difficulty: "hard", Category: "M5: Advanced", Language: LangSolidity,
		Description: `EIP-1167 Minimal Proxy (Clone) is a lightweight pattern for deploying
many contracts with the same logic. Instead of deploying full bytecode
each time, a tiny proxy delegates all calls to an implementation contract.

The proxy is only ~45 bytes of bytecode, making deployment extremely cheap.
It uses delegatecall internally — the proxy's storage is used, but the
implementation's code is executed.

This pattern is used in:
- Token factories (deploy new tokens cheaply)
- Account abstraction wallets
- NFT collection factories

For this exercise, implement a clone factory that creates minimal
proxies and initializes them.`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract Implementation {
    address public owner;
    uint256 public value;
    bool private _initialized;

    function initialize(address _owner, uint256 _value) external {
        require(!_initialized, "already initialized");
        owner = _owner;
        value = _value;
        _initialized = true;
    }

    function setValue(uint256 _value) external {
        require(msg.sender == owner, "not owner");
        value = _value;
    }
}

contract CloneFactory {
    event Cloned(address indexed clone);

    // TODO 1: Write 'clone(address implementation)' that creates a minimal proxy
    //   Use assembly to deploy the EIP-1167 bytecode:
    //   The minimal proxy bytecode pattern (where <impl> is the 20-byte address):
    //     3d602d80600a3d3981f3363d3d373d3d3d363d73<impl>5af43d82803e903d91602b57fd5bf3
    //
    //   Steps:
    //   a) In assembly, use mstore to write the bytecode:
    //      mstore(0x00, 0x3d602d80600a3d3981f3363d3d373d3d3d363d73000000000000000000000000)
    //      mstore(0x14, shl(96, implementation))
    //      mstore(0x28, 0x5af43d82803e903d91602b57fd5bf30000000000000000000000000000000000)
    //   b) instance := create(0, 0x00, 0x37)
    //   c) require instance != address(0)
    //   (external function, returns address)

    // TODO 2: Write 'cloneAndInitialize(address implementation, address owner, uint256 value)'
    //   that clones the implementation, then calls initialize on the clone
    //   (external function, returns address)
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

contract CloneFactoryTest is Test {
    CloneFactory factory;
    Implementation impl;

    function setUp() public {
        factory = new CloneFactory();
        impl = new Implementation();
    }

    function test_Clone() public {
        address c = factory.clone(address(impl));
        assertTrue(c != address(0));
    }

    function test_CloneAndInitialize() public {
        address c = factory.cloneAndInitialize(address(impl), address(this), 42);
        assertEq(Implementation(c).owner(), address(this));
        assertEq(Implementation(c).value(), 42);
    }

    function test_ClonesAreIndependent() public {
        address c1 = factory.cloneAndInitialize(address(impl), address(1), 100);
        address c2 = factory.cloneAndInitialize(address(impl), address(2), 200);
        assertEq(Implementation(c1).value(), 100);
        assertEq(Implementation(c2).value(), 200);
        assertTrue(c1 != c2);
    }
}
`,
		Hints: []string{
			"The EIP-1167 bytecode is 55 bytes (0x37). Use assembly with mstore and create.",
			"mstore(0x00, 0x3d602d80600a3d3981f3...); mstore(0x14, shl(96, implementation)); mstore(0x28, 0x5af43d82803e903d91602b57fd5bf3...);",
			"For cloneAndInitialize: address c = clone(implementation); Implementation(c).initialize(owner, value);",
		},
	}
}

func solAccessControl() Question {
	return Question{
		ID: "sol-access-control", Title: "Access Control",
		Difficulty: "medium", Category: "M5: Advanced", Language: LangSolidity,
		Description: `Access control restricts who can call sensitive functions. Two patterns:

1. Ownable: single owner with special privileges
   - constructor sets owner = msg.sender
   - onlyOwner modifier restricts functions
   - transferOwnership to hand over control

2. Role-based: multiple roles for different permissions
   - ADMIN_ROLE, MINTER_ROLE, PAUSER_ROLE, etc.
   - grantRole/revokeRole for managing permissions
   - hasRole for checking

In x402: the Facilitator contract uses access control to ensure
only authorized accounts can submit payment settlements.

Implement both Ownable and role-based access control.`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract AccessControlled {
    address public owner;
    mapping(bytes32 => mapping(address => bool)) private _roles;

    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");

    // TODO 1: Define modifier 'onlyOwner' that requires msg.sender == owner

    // TODO 2: Define modifier 'onlyRole(bytes32 role)' that requires _roles[role][msg.sender]

    constructor() {
        // TODO 3: Set owner to msg.sender
        // TODO 4: Grant ADMIN_ROLE to msg.sender
    }

    function transferOwnership(address newOwner) external onlyOwner {
        // TODO 5: require newOwner != address(0)
        // TODO 6: Set owner = newOwner
    }

    function grantRole(bytes32 role, address account) external onlyRole(ADMIN_ROLE) {
        // TODO 7: Set _roles[role][account] = true
    }

    function revokeRole(bytes32 role, address account) external onlyRole(ADMIN_ROLE) {
        // TODO 8: Set _roles[role][account] = false
    }

    function hasRole(bytes32 role, address account) public view returns (bool) {
        // TODO 9: Return _roles[role][account]
        return false;
    }

    // Protected function: only MINTER_ROLE can call
    function mint() external onlyRole(MINTER_ROLE) pure returns (string memory) {
        return "minted";
    }
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

contract AccessControlTest is Test {
    AccessControlled ac;
    address admin = address(1);
    address minter = address(2);
    address nobody = address(3);

    function setUp() public {
        vm.prank(admin);
        ac = new AccessControlled();
        vm.prank(admin);
        ac.grantRole(ac.MINTER_ROLE(), minter);
    }

    function test_Owner() public view {
        assertEq(ac.owner(), admin);
    }

    function test_AdminRole() public view {
        assertTrue(ac.hasRole(ac.ADMIN_ROLE(), admin));
    }

    function test_MintWithRole() public {
        vm.prank(minter);
        assertEq(ac.mint(), "minted");
    }

    function test_MintWithoutRole() public {
        vm.prank(nobody);
        vm.expectRevert();
        ac.mint();
    }

    function test_TransferOwnership() public {
        vm.prank(admin);
        ac.transferOwnership(address(4));
        assertEq(ac.owner(), address(4));
    }

    function test_TransferOwnership_NotOwner() public {
        vm.prank(nobody);
        vm.expectRevert();
        ac.transferOwnership(nobody);
    }
}
`,
		Hints: []string{
			"modifier onlyOwner() { require(msg.sender == owner); _; }",
			"modifier onlyRole(bytes32 role) { require(_roles[role][msg.sender]); _; }",
			"constructor: owner = msg.sender; _roles[ADMIN_ROLE][msg.sender] = true;",
		},
	}
}

func solReentrancy() Question {
	return Question{
		ID: "sol-reentrancy", Title: "Reentrancy Guard",
		Difficulty: "medium", Category: "M5: Advanced", Language: LangSolidity,
		Description: `Reentrancy is the most infamous smart contract vulnerability. It occurs when
a contract makes an external call before updating its state, allowing the
called contract to re-enter and exploit the stale state.

The famous DAO hack (2016) exploited this: withdraw() sent ETH before
updating the balance, so the attacker could withdraw repeatedly.

Two defenses:
1. Checks-Effects-Interactions (CEI) pattern:
   - Check conditions, update state, THEN make external calls
2. ReentrancyGuard: a lock that prevents re-entry
   - Uses a state variable (_locked) to block nested calls

Implement a vault with proper reentrancy protection.`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract SecureVault {
    mapping(address => uint256) public balances;
    bool private _locked;

    // TODO 1: Define modifier 'nonReentrant' that:
    //   - requires !_locked
    //   - sets _locked = true
    //   - executes the function body (_)
    //   - sets _locked = false

    function deposit() external payable {
        // TODO 2: Add msg.value to balances[msg.sender]
    }

    // Use BOTH reentrancy guard AND checks-effects-interactions
    function withdraw(uint256 amount) external nonReentrant {
        // TODO 3: Check: require balances[msg.sender] >= amount
        // TODO 4: Effect: subtract amount from balances[msg.sender] BEFORE transfer
        // TODO 5: Interaction: send ETH via call (not transfer)
        //   (bool success, ) = payable(msg.sender).call{value: amount}("");
        //   require(success);
    }

    function getBalance() external view returns (uint256) {
        return address(this).balance;
    }
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

contract Attacker {
    SecureVault public vault;
    uint256 public count;

    constructor(address _vault) {
        vault = SecureVault(_vault);
    }

    function attack() external payable {
        vault.deposit{value: msg.value}();
        vault.withdraw(msg.value);
    }

    receive() external payable {
        count++;
        if (count < 3 && address(vault).balance >= 1 ether) {
            vault.withdraw(1 ether);
        }
    }
}

contract SecureVaultTest is Test {
    SecureVault vault;
    address alice = address(1);

    function setUp() public {
        vault = new SecureVault();
        vm.deal(alice, 5 ether);
    }

    function test_DepositAndWithdraw() public {
        vm.prank(alice);
        vault.deposit{value: 2 ether}();
        assertEq(vault.balances(alice), 2 ether);

        vm.prank(alice);
        vault.withdraw(1 ether);
        assertEq(vault.balances(alice), 1 ether);
    }

    function test_ReentrancyBlocked() public {
        // Fund vault with some ETH first
        vm.prank(alice);
        vault.deposit{value: 3 ether}();

        Attacker attacker = new Attacker(address(vault));
        vm.deal(address(attacker), 1 ether);

        vm.expectRevert();
        attacker.attack{value: 1 ether}();
    }

    function test_WithdrawInsufficient() public {
        vm.prank(alice);
        vm.expectRevert();
        vault.withdraw(1 ether);
    }
}
`,
		Hints: []string{
			"modifier nonReentrant() { require(!_locked); _locked = true; _; _locked = false; }",
			"Checks-Effects-Interactions: update balances BEFORE sending ETH",
			"(bool success, ) = payable(msg.sender).call{value: amount}(\"\"); require(success);",
		},
	}
}

// ============================================================
// MODULE 6: x402 Protocol
// ============================================================

func solModule6X402() []Question {
	return []Question{
		solX402Settle(),
		solX402Permit2Proxy(),
		solX402VerifySettle(),
	}
}

func solX402Settle() Question {
	return Question{
		ID: "sol-x402-settle", Title: "Settlement Contract",
		Difficulty: "hard", Category: "M6: x402", Language: LangSolidity,
		Description: `In x402, the Facilitator settles payments by calling transferWithAuthorization
on the USDC contract. The settlement process:

1. Client signs an EIP-3009 authorization (off-chain)
2. Facilitator verifies the payment requirements
3. Facilitator calls USDC.transferWithAuthorization with the signed data
4. USDC verifies the signature and transfers tokens from Client to PayTo

The settlement contract wraps this process, adding:
- Event logging for tracking settlements
- Fee collection for the facilitator
- Payment validation (amount, token, recipient)

Implement a settlement contract that wraps transferWithAuthorization.`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

interface IEIP3009 {
    function transferWithAuthorization(
        address from, address to, uint256 value,
        uint256 validAfter, uint256 validBefore, bytes32 nonce,
        uint8 v, bytes32 r, bytes32 s
    ) external;
}

contract Settlement {
    event PaymentSettled(
        address indexed from,
        address indexed to,
        uint256 value,
        bytes32 nonce
    );

    address public token;
    address public facilitator;

    constructor(address _token, address _facilitator) {
        token = _token;
        facilitator = _facilitator;
    }

    // TODO 1: Define modifier 'onlyFacilitator' requiring msg.sender == facilitator

    // TODO 2: Write 'settle' function (onlyFacilitator, external) that:
    //   - Takes: from, to, value, validAfter, validBefore, nonce, v, r, s
    //   - Calls IEIP3009(token).transferWithAuthorization(from, to, value, validAfter, validBefore, nonce, v, r, s)
    //   - Emits PaymentSettled(from, to, value, nonce)

    // TODO 3: Write a view function 'isValidPayment(address from, address to, uint256 value)'
    //   that returns true if from != address(0) && to != address(0) && value > 0
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

contract MockEIP3009 {
    mapping(address => uint256) public balanceOf;
    bool public called;
    address public lastFrom;
    address public lastTo;
    uint256 public lastValue;

    constructor(address holder, uint256 amount) {
        balanceOf[holder] = amount;
    }

    function transferWithAuthorization(
        address from, address to, uint256 value,
        uint256, uint256, bytes32,
        uint8, bytes32, bytes32
    ) external {
        require(balanceOf[from] >= value, "insufficient");
        balanceOf[from] -= value;
        balanceOf[to] += value;
        called = true;
        lastFrom = from;
        lastTo = to;
        lastValue = value;
    }
}

contract SettlementTest is Test {
    Settlement settlement;
    MockEIP3009 mockToken;
    address facilAddr = address(1);
    address alice = address(2);
    address bob = address(3);

    event PaymentSettled(address indexed from, address indexed to, uint256 value, bytes32 nonce);

    function setUp() public {
        mockToken = new MockEIP3009(alice, 1000000);
        settlement = new Settlement(address(mockToken), facilAddr);
    }

    function test_Settle() public {
        bytes32 nonce = bytes32(uint256(42));
        vm.prank(facilAddr);
        settlement.settle(alice, bob, 100000, 0, 999999, nonce, 27, bytes32(0), bytes32(0));

        assertTrue(mockToken.called());
        assertEq(mockToken.balanceOf(bob), 100000);
        assertEq(mockToken.balanceOf(alice), 900000);
    }

    function test_SettleEvent() public {
        bytes32 nonce = bytes32(uint256(1));
        vm.expectEmit(true, true, false, true);
        emit PaymentSettled(alice, bob, 50000, nonce);

        vm.prank(facilAddr);
        settlement.settle(alice, bob, 50000, 0, 999999, nonce, 27, bytes32(0), bytes32(0));
    }

    function test_OnlyFacilitator() public {
        vm.prank(alice);
        vm.expectRevert();
        settlement.settle(alice, bob, 100, 0, 999, bytes32(0), 27, bytes32(0), bytes32(0));
    }

    function test_IsValidPayment() public view {
        assertTrue(settlement.isValidPayment(alice, bob, 100));
        assertFalse(settlement.isValidPayment(address(0), bob, 100));
        assertFalse(settlement.isValidPayment(alice, bob, 0));
    }
}
`,
		Hints: []string{
			"modifier onlyFacilitator() { require(msg.sender == facilitator); _; }",
			"IEIP3009(token).transferWithAuthorization(from, to, value, validAfter, validBefore, nonce, v, r, s);",
			"isValidPayment: return from != address(0) && to != address(0) && value > 0;",
		},
	}
}

func solX402Permit2Proxy() Question {
	return Question{
		ID: "sol-x402-permit2-proxy", Title: "x402 Permit2 Proxy",
		Difficulty: "hard", Category: "M6: x402", Language: LangSolidity,
		Description: `The x402Permit2Proxy contract enables x402 payments for ANY ERC-20 token,
not just those implementing EIP-3009. It acts as an intermediary:

1. Client approves Permit2 for the token (one-time)
2. Client signs a Permit2 message with x402 witness data
3. Facilitator calls x402Permit2Proxy.settle()
4. Proxy calls Permit2.permitWitnessTransferFrom()
5. Permit2 transfers tokens from Client to PayTo

The proxy adds a "witness" — extra data included in the signed message
that proves the transfer is specifically for an x402 payment.

The real x402Permit2Proxy address is 0x402085c248EeA27D92E8b30b2C58ed07f9E20001,
deployed via CREATE2 on all chains.

Implement a simplified version of the proxy's settlement logic.`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

interface IPermit2 {
    function transferFrom(address from, address to, uint256 amount, address token) external;
}

contract X402Permit2Proxy {
    address public immutable permit2;
    address public facilitator;

    event PaymentSettled(
        address indexed token,
        address indexed from,
        address indexed to,
        uint256 amount,
        bytes32 paymentId
    );

    constructor(address _permit2, address _facilitator) {
        permit2 = _permit2;
        facilitator = _facilitator;
    }

    // TODO 1: Write 'settle' function that takes:
    //   (address token, address from, address to, uint256 amount, bytes32 paymentId)
    //   - Require msg.sender == facilitator
    //   - Require from, to, token are not address(0)
    //   - Require amount > 0
    //   - Call IPermit2(permit2).transferFrom(from, to, amount, token)
    //   - Emit PaymentSettled event
    //   (external function)

    // TODO 2: Write a view function 'validatePayment'
    //   (address token, address from, address to, uint256 amount)
    //   Returns (bool valid, string memory reason)
    //   Check all parameters are valid and return descriptive reason if not
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

contract MockPermit2ForProxy {
    mapping(address => mapping(address => uint256)) public balances;
    bool public transferCalled;

    function setBalance(address token, address user, uint256 amount) external {
        balances[token][user] = amount;
    }

    function transferFrom(address from, address to, uint256 amount, address token) external {
        require(balances[token][from] >= amount, "insufficient");
        balances[token][from] -= amount;
        balances[token][to] += amount;
        transferCalled = true;
    }
}

contract X402Permit2ProxyTest is Test {
    X402Permit2Proxy proxy;
    MockPermit2ForProxy mockPermit2;
    address facilAddr = address(1);
    address alice = address(2);
    address bob = address(3);
    address tokenAddr = address(4);

    event PaymentSettled(address indexed token, address indexed from, address indexed to, uint256 amount, bytes32 paymentId);

    function setUp() public {
        mockPermit2 = new MockPermit2ForProxy();
        proxy = new X402Permit2Proxy(address(mockPermit2), facilAddr);
        mockPermit2.setBalance(tokenAddr, alice, 1000000);
    }

    function test_Settle() public {
        bytes32 pid = bytes32(uint256(1));
        vm.prank(facilAddr);
        proxy.settle(tokenAddr, alice, bob, 100000, pid);

        assertTrue(mockPermit2.transferCalled());
        assertEq(mockPermit2.balances(tokenAddr, bob), 100000);
    }

    function test_SettleEvent() public {
        bytes32 pid = bytes32(uint256(2));
        vm.expectEmit(true, true, true, true);
        emit PaymentSettled(tokenAddr, alice, bob, 50000, pid);

        vm.prank(facilAddr);
        proxy.settle(tokenAddr, alice, bob, 50000, pid);
    }

    function test_OnlyFacilitator() public {
        vm.prank(alice);
        vm.expectRevert();
        proxy.settle(tokenAddr, alice, bob, 100, bytes32(0));
    }

    function test_ValidatePayment() public view {
        (bool valid, ) = proxy.validatePayment(tokenAddr, alice, bob, 100);
        assertTrue(valid);

        (bool invalid, string memory reason) = proxy.validatePayment(address(0), alice, bob, 100);
        assertFalse(invalid);
        assertTrue(bytes(reason).length > 0);
    }
}
`,
		Hints: []string{
			"require(msg.sender == facilitator); require(from != address(0) && to != address(0) && token != address(0) && amount > 0);",
			"IPermit2(permit2).transferFrom(from, to, amount, token); emit PaymentSettled(token, from, to, amount, paymentId);",
			`validatePayment: if (token == address(0)) return (false, "invalid token"); etc.`,
		},
	}
}

func solX402VerifySettle() Question {
	return Question{
		ID: "sol-x402-verify-settle", Title: "Full Payment Verification",
		Difficulty: "hard", Category: "M6: x402", Language: LangSolidity,
		Description: `The complete x402 payment lifecycle combines verification and settlement:

1. VERIFY: Check that the payment meets requirements
   - Amount >= required amount
   - Token matches
   - Recipient (payTo) matches
   - Not expired (validBefore > now)
   - Not replayed (nonce not used)

2. SETTLE: Execute the on-chain transfer
   - Call transferWithAuthorization or Permit2
   - Record the settlement
   - Emit events

The Facilitator performs both steps atomically:
- First verify all conditions
- Then settle the payment
- If verification fails, the payment is rejected (HTTP 400)
- If settlement fails, the payment is reversed (HTTP 500)

Implement a complete verify+settle contract.`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract PaymentProcessor {
    struct PaymentRequest {
        address from;
        address to;
        address token;
        uint256 amount;
        uint256 validAfter;
        uint256 validBefore;
        bytes32 nonce;
    }

    mapping(bytes32 => bool) public settledPayments;
    mapping(bytes32 => bool) public usedNonces;

    event PaymentVerified(bytes32 indexed paymentId, address from, address to, uint256 amount);
    event PaymentSettled(bytes32 indexed paymentId, address from, address to, uint256 amount);

    // TODO 1: Write a view function 'verify(PaymentRequest memory req, uint256 requiredAmount)'
    //   that returns (bool valid, string memory reason)
    //   Checks:
    //   - from != address(0) && to != address(0) && token != address(0)
    //   - amount >= requiredAmount
    //   - block.timestamp > validAfter
    //   - block.timestamp < validBefore
    //   - nonce not in usedNonces
    //   - paymentId not in settledPayments
    //   paymentId = keccak256(abi.encode(req.from, req.to, req.nonce))

    // TODO 2: Write 'settle(PaymentRequest memory req)' external function
    //   - Compute paymentId
    //   - require !usedNonces[req.nonce]
    //   - require !settledPayments[paymentId]
    //   - Mark usedNonces[req.nonce] = true
    //   - Mark settledPayments[paymentId] = true
    //   - Emit PaymentSettled
    //   (In production, this would also call the token contract)

    // TODO 3: Write a view function 'getPaymentId(PaymentRequest memory req)'
    //   that returns keccak256(abi.encode(req.from, req.to, req.nonce))
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

contract PaymentProcessorTest is Test {
    PaymentProcessor processor;

    function setUp() public {
        processor = new PaymentProcessor();
    }

    function _req() internal pure returns (PaymentProcessor.PaymentRequest memory) {
        return PaymentProcessor.PaymentRequest({
            from: address(1),
            to: address(2),
            token: address(3),
            amount: 100000,
            validAfter: 0,
            validBefore: 999999,
            nonce: bytes32(uint256(42))
        });
    }

    function test_VerifyValid() public {
        vm.warp(100);
        (bool valid, ) = processor.verify(_req(), 100000);
        assertTrue(valid);
    }

    function test_VerifyInsufficientAmount() public {
        vm.warp(100);
        (bool valid, string memory reason) = processor.verify(_req(), 200000);
        assertFalse(valid);
        assertTrue(bytes(reason).length > 0);
    }

    function test_VerifyExpired() public {
        vm.warp(1000000);
        (bool valid, ) = processor.verify(_req(), 100000);
        assertFalse(valid);
    }

    function test_Settle() public {
        vm.warp(100);
        PaymentProcessor.PaymentRequest memory req = _req();
        processor.settle(req);

        bytes32 pid = processor.getPaymentId(req);
        assertTrue(processor.settledPayments(pid));
        assertTrue(processor.usedNonces(req.nonce));
    }

    function test_SettleReplay() public {
        vm.warp(100);
        processor.settle(_req());

        vm.expectRevert();
        processor.settle(_req());
    }
}
`,
		Hints: []string{
			"bytes32 paymentId = keccak256(abi.encode(req.from, req.to, req.nonce));",
			`if (req.amount < requiredAmount) return (false, "insufficient amount"); etc.`,
			"settle: require(!usedNonces[req.nonce]); usedNonces[req.nonce] = true; settledPayments[paymentId] = true;",
		},
	}
}

// ============================================================
// MODULE 7: ERC-8004 — Autonomous Agent Identity & Reputation
// ============================================================

func solModule7ERC8004() []Question {
	return []Question{
		solERC8004Identity(),
		solERC8004Metadata(),
		solERC8004Feedback(),
		solERC8004SelfPrevention(),
		solERC8004Validation(),
		solERC8004Wallet(),
		solERC8004ReputationSummary(),
		solERC8004X402Integration(),
	}
}

func solERC8004Identity() Question {
	return Question{
		ID: "sol-erc8004-identity", Title: "Agent Identity Registry",
		Difficulty: "easy", Category: "M7: ERC-8004", Language: LangSolidity,
		Description: `ERC-8004 defines an on-chain identity and reputation system for autonomous agents.
The first building block is an identity registry where each agent gets a unique ID,
an owner (the address that registered it), and a URI pointing to off-chain metadata
(e.g., an IPFS link describing the agent's capabilities).

Think of it like ERC-721 for agents: each agent is a non-fungible on-chain entity
with a numeric ID. However, unlike NFTs meant for collectibles, agent identities
are functional — they anchor reputation, wallets, and service discovery.

The registry uses a simple auto-incrementing counter starting at 1. When an address
calls register(), it mints a new agent identity, stores the URI, records the owner,
emits a Registered event, and returns the new agent ID.

Implement the core identity registry with register, ownerOf, and getAgentURI.`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract AgentIdentityRegistry {
    uint256 private _nextAgentId;
    mapping(uint256 => string) private _agentURIs;
    mapping(uint256 => address) private _owners;

    // TODO 1: Define event Registered(uint256 indexed agentId, address indexed owner, string agentURI)

    constructor() {
        // TODO 2: Initialize _nextAgentId to 1 (agent IDs start at 1, not 0)
    }

    function register(string calldata agentURI) external returns (uint256) {
        // TODO 3: Get the current _nextAgentId as the new agentId
        // TODO 4: Store the owner as msg.sender in _owners[agentId]
        // TODO 5: Store the URI in _agentURIs[agentId]
        // TODO 6: Increment _nextAgentId
        // TODO 7: Emit the Registered event
        // TODO 8: Return the agentId
        return 0;
    }

    function ownerOf(uint256 agentId) external view returns (address) {
        // TODO 9: Return the owner of the given agentId
        // Require that the agent exists (owner != address(0))
        return address(0);
    }

    function getAgentURI(uint256 agentId) external view returns (string memory) {
        // TODO 10: Require agent exists, then return the URI
        return "";
    }
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

contract AgentIdentityRegistryTest is Test {
    AgentIdentityRegistry registry;
    address alice = address(1);
    address bob = address(2);

    function setUp() public {
        registry = new AgentIdentityRegistry();
    }

    function test_RegisterAgent() public {
        vm.prank(alice);
        uint256 id = registry.register("ipfs://agent1");
        assertEq(id, 1);
        assertEq(registry.ownerOf(1), alice);
        assertEq(registry.getAgentURI(1), "ipfs://agent1");
    }

    function test_SecondAgentGetsId2() public {
        vm.prank(alice);
        registry.register("ipfs://agent1");

        vm.prank(bob);
        uint256 id2 = registry.register("ipfs://agent2");
        assertEq(id2, 2);
        assertEq(registry.ownerOf(2), bob);
    }

    function test_OwnerOfNonexistent() public {
        vm.expectRevert();
        registry.ownerOf(999);
    }

    function test_GetURINonexistent() public {
        vm.expectRevert();
        registry.getAgentURI(999);
    }

    function test_RegisterEmitsEvent() public {
        vm.prank(alice);
        vm.expectEmit(true, true, false, true);
        emit AgentIdentityRegistry.Registered(1, alice, "ipfs://agent1");
        registry.register("ipfs://agent1");
    }
}
`,
		Hints: []string{
			"constructor: _nextAgentId = 1;",
			"register: uint256 agentId = _nextAgentId; _owners[agentId] = msg.sender; _agentURIs[agentId] = agentURI; _nextAgentId++;",
			"ownerOf: require(_owners[agentId] != address(0), \"agent not found\"); return _owners[agentId];",
		},
	}
}

func solERC8004Metadata() Question {
	return Question{
		ID: "sol-erc8004-metadata", Title: "Agent Metadata Storage",
		Difficulty: "easy", Category: "M7: ERC-8004", Language: LangSolidity,
		Description: `Beyond the base URI, agents need structured key-value metadata stored on-chain.
This allows discoverability — other contracts or off-chain indexers can query an
agent's capabilities, version, supported protocols, etc.

The metadata system uses a nested mapping: agentId => key => value, where both
key and value are strings. Only the agent's owner can set metadata for that agent.

One important constraint: the key "agentWallet" is reserved for the wallet
verification system (covered in a later question). The contract must reject
attempts to set this key via the metadata function. String comparison in
Solidity is done by comparing keccak256 hashes of the strings.

Implement the metadata storage with owner-only writes and reserved key protection.`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract AgentMetadata {
    uint256 private _nextAgentId;
    mapping(uint256 => address) private _owners;
    mapping(uint256 => mapping(string => string)) private _metadata;

    constructor() {
        _nextAgentId = 1;
    }

    function register() external returns (uint256) {
        uint256 agentId = _nextAgentId;
        _owners[agentId] = msg.sender;
        _nextAgentId++;
        return agentId;
    }

    function ownerOf(uint256 agentId) public view returns (address) {
        require(_owners[agentId] != address(0), "agent not found");
        return _owners[agentId];
    }

    function setMetadata(uint256 agentId, string calldata key, string calldata value) external {
        // TODO 1: Require that msg.sender is the owner of agentId
        // TODO 2: Require that the key is NOT "agentWallet" (reserved)
        //         Use keccak256(abi.encodePacked(key)) != keccak256(abi.encodePacked("agentWallet"))
        // TODO 3: Store the value in _metadata[agentId][key]
    }

    function getMetadata(uint256 agentId, string calldata key) external view returns (string memory) {
        // TODO 4: Return the metadata value for the given agent and key
        return "";
    }
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

contract AgentMetadataTest is Test {
    AgentMetadata meta;
    address alice = address(1);
    address bob = address(2);

    function setUp() public {
        meta = new AgentMetadata();
        vm.prank(alice);
        meta.register(); // agentId = 1
    }

    function test_SetAndGetMetadata() public {
        vm.prank(alice);
        meta.setMetadata(1, "version", "1.0.0");
        assertEq(meta.getMetadata(1, "version"), "1.0.0");
    }

    function test_OverwriteMetadata() public {
        vm.prank(alice);
        meta.setMetadata(1, "protocol", "x402");
        vm.prank(alice);
        meta.setMetadata(1, "protocol", "x402-v2");
        assertEq(meta.getMetadata(1, "protocol"), "x402-v2");
    }

    function test_ReservedKeyReverts() public {
        vm.prank(alice);
        vm.expectRevert();
        meta.setMetadata(1, "agentWallet", "0xabc");
    }

    function test_OnlyOwnerCanSet() public {
        vm.prank(bob);
        vm.expectRevert();
        meta.setMetadata(1, "version", "hacked");
    }

    function test_GetUnsetMetadata() public view {
        string memory val = meta.getMetadata(1, "nonexistent");
        assertEq(bytes(val).length, 0);
    }
}
`,
		Hints: []string{
			"require(msg.sender == ownerOf(agentId), \"not owner\");",
			"require(keccak256(abi.encodePacked(key)) != keccak256(abi.encodePacked(\"agentWallet\")), \"reserved key\");",
			"_metadata[agentId][key] = value; / return _metadata[agentId][key];",
		},
	}
}

func solERC8004Feedback() Question {
	return Question{
		ID: "sol-erc8004-feedback", Title: "Reputation Feedback",
		Difficulty: "medium", Category: "M7: ERC-8004", Language: LangSolidity,
		Description: `The reputation system in ERC-8004 allows anyone to leave feedback on an agent.
Each feedback entry is a struct containing the provider (who left it), the target
agent ID, a signed integer value (positive or negative), a free-text tag describing
the context (e.g., "payment", "response-quality"), and a timestamp.

Feedback is stored per-agent in an array, allowing enumeration. The feedbackCount
function returns the total number of feedbacks for an agent, and getFeedback returns
a specific entry by index.

Crucially, feedback can be revoked — but only by the original provider. This prevents
spam while allowing providers to correct mistakes. Revocation sets the feedback value
to zero and clears the tag, but preserves the slot to maintain index stability.

Implement the reputation feedback system with give, get, count, and revoke operations.`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract ReputationRegistry {
    struct Feedback {
        address provider;
        uint256 agentId;
        int256 value;
        string tag;
        uint64 timestamp;
    }

    // TODO 1: Declare a mapping from uint256 (agentId) to an array of Feedback structs
    //         e.g., mapping(uint256 => Feedback[]) private _feedbacks;

    function giveFeedback(uint256 agentId, int256 value, string calldata tag) external {
        // TODO 2: Create a new Feedback struct with msg.sender as provider,
        //         the given agentId, value, tag, and block.timestamp as uint64
        // TODO 3: Push the struct into the feedbacks array for this agentId
    }

    function getFeedback(uint256 agentId, uint256 index)
        external view
        returns (address provider, int256 value, string memory tag, uint64 timestamp)
    {
        // TODO 4: Require index < length of the feedbacks array for this agentId
        // TODO 5: Return the provider, value, tag, and timestamp fields
        return (address(0), 0, "", 0);
    }

    function getFeedbackCount(uint256 agentId) external view returns (uint256) {
        // TODO 6: Return the length of the feedbacks array for this agentId
        return 0;
    }

    function revokeFeedback(uint256 agentId, uint256 index) external {
        // TODO 7: Require index is valid
        // TODO 8: Require msg.sender == feedback.provider (only the original provider can revoke)
        // TODO 9: Set the feedback value to 0 and tag to "" (revoked)
    }
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

contract ReputationRegistryTest is Test {
    ReputationRegistry rep;
    address alice = address(1);
    address bob = address(2);

    function setUp() public {
        rep = new ReputationRegistry();
    }

    function test_GiveFeedback() public {
        vm.warp(1000);
        vm.prank(alice);
        rep.giveFeedback(1, 80, "good-service");
        assertEq(rep.getFeedbackCount(1), 1);
    }

    function test_GetFeedbackFields() public {
        vm.warp(2000);
        vm.prank(alice);
        rep.giveFeedback(1, -20, "slow-response");

        (address provider, int256 value, string memory tag, uint64 ts) = rep.getFeedback(1, 0);
        assertEq(provider, alice);
        assertEq(value, -20);
        assertEq(tag, "slow-response");
        assertEq(ts, 2000);
    }

    function test_MultipleFeedbacks() public {
        vm.prank(alice);
        rep.giveFeedback(1, 50, "ok");
        vm.prank(bob);
        rep.giveFeedback(1, 90, "excellent");
        assertEq(rep.getFeedbackCount(1), 2);
    }

    function test_RevokeFeedback() public {
        vm.prank(alice);
        rep.giveFeedback(1, 80, "good");

        vm.prank(alice);
        rep.revokeFeedback(1, 0);

        (, int256 value, string memory tag,) = rep.getFeedback(1, 0);
        assertEq(value, 0);
        assertEq(bytes(tag).length, 0);
    }

    function test_RevokeNotProviderReverts() public {
        vm.prank(alice);
        rep.giveFeedback(1, 80, "good");

        vm.prank(bob);
        vm.expectRevert();
        rep.revokeFeedback(1, 0);
    }
}
`,
		Hints: []string{
			"mapping(uint256 => Feedback[]) private _feedbacks;",
			"giveFeedback: _feedbacks[agentId].push(Feedback(msg.sender, agentId, value, tag, uint64(block.timestamp)));",
			"revokeFeedback: require(msg.sender == _feedbacks[agentId][index].provider); then set value=0 and tag=\"\"",
		},
	}
}

func solERC8004SelfPrevention() Question {
	return Question{
		ID: "sol-erc8004-self-prevention", Title: "Self-Feedback Prevention",
		Difficulty: "medium", Category: "M7: ERC-8004", Language: LangSolidity,
		Description: `A critical integrity rule in ERC-8004's reputation system is that an agent's
owner must not be able to give feedback to their own agent. Without this check,
owners could inflate their agent's reputation by leaving positive reviews.

This guard is implemented as a separate contract that cross-references the
identity registry. The ReputationGuard takes the identity registry address
in its constructor and calls ownerOf() on it to check whether the caller
owns the target agent.

This pattern demonstrates contract composability: the reputation contract
depends on an external identity registry via an interface. In production, these
would be separate deployed contracts, with the reputation system querying
identity ownership on every feedback submission.

Implement the self-feedback prevention guard with an interface to the identity registry.`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

// TODO 1: Define interface IIdentityRegistry with a single function:
//         ownerOf(uint256 agentId) external view returns (address)

contract ReputationGuard {
    struct Feedback {
        address provider;
        uint256 agentId;
        int256 value;
        string tag;
    }

    // TODO 2: Declare a state variable for the identity registry (IIdentityRegistry)
    // TODO 3: Declare a mapping from uint256 (agentId) to Feedback array

    constructor(address _identityRegistry) {
        // TODO 4: Store the identity registry reference
        //         Cast address to IIdentityRegistry
    }

    function giveFeedback(uint256 agentId, int256 value, string calldata tag) external {
        // TODO 5: Get the owner of the agentId from the identity registry
        // TODO 6: Require that msg.sender is NOT the owner (prevent self-feedback)
        // TODO 7: Store the feedback
    }

    function getFeedbackCount(uint256 agentId) external view returns (uint256) {
        // TODO 8: Return the count of feedbacks for this agentId
        return 0;
    }
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

// Mock identity registry for testing
contract MockIdentityRegistry {
    mapping(uint256 => address) private _owners;

    function setOwner(uint256 agentId, address owner) external {
        _owners[agentId] = owner;
    }

    function ownerOf(uint256 agentId) external view returns (address) {
        require(_owners[agentId] != address(0), "not found");
        return _owners[agentId];
    }
}

contract ReputationGuardTest is Test {
    MockIdentityRegistry mockRegistry;
    ReputationGuard guard;
    address alice = address(1);
    address bob = address(2);

    function setUp() public {
        mockRegistry = new MockIdentityRegistry();
        mockRegistry.setOwner(1, alice);
        guard = new ReputationGuard(address(mockRegistry));
    }

    function test_NonOwnerCanGiveFeedback() public {
        vm.prank(bob);
        guard.giveFeedback(1, 80, "great-agent");
        assertEq(guard.getFeedbackCount(1), 1);
    }

    function test_OwnerCannotSelfFeedback() public {
        vm.prank(alice);
        vm.expectRevert();
        guard.giveFeedback(1, 100, "im-the-best");
    }

    function test_MultipleFeedbacksFromDifferentUsers() public {
        vm.prank(bob);
        guard.giveFeedback(1, 50, "ok");

        address charlie = address(3);
        vm.prank(charlie);
        guard.giveFeedback(1, 90, "excellent");

        assertEq(guard.getFeedbackCount(1), 2);
    }

    function test_NegativeFeedbackAllowed() public {
        vm.prank(bob);
        guard.giveFeedback(1, -30, "poor-service");
        assertEq(guard.getFeedbackCount(1), 1);
    }
}
`,
		Hints: []string{
			"interface IIdentityRegistry { function ownerOf(uint256 agentId) external view returns (address); }",
			"IIdentityRegistry public identityRegistry; constructor: identityRegistry = IIdentityRegistry(_identityRegistry);",
			"require(msg.sender != identityRegistry.ownerOf(agentId), \"self-feedback not allowed\");",
		},
	}
}

func solERC8004Validation() Question {
	return Question{
		ID: "sol-erc8004-validation", Title: "Validation Request & Response",
		Difficulty: "medium", Category: "M7: ERC-8004", Language: LangSolidity,
		Description: `ERC-8004 includes a validation mechanism where agent owners can request
third-party validation of their agent's capabilities. A designated validator
reviews the agent and responds with a score (0-100) and a textual reason.

The flow works like this:
1. The agent owner calls requestValidation(), specifying a validator address
   and arbitrary data describing what should be validated.
2. The designated validator calls respond() with a score and reason.
3. Anyone can query the validation result by request ID.

This on-chain validation creates a verifiable audit trail. Validators could be
trusted third parties, DAOs, or even other autonomous agents with established
reputation.

Only the agent's original registrant can request validation, and only the
designated validator can respond. Scores must be between 0 and 100 inclusive.

Implement the validation request and response system.`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract ValidationRegistry {
    struct ValidationRequest {
        uint256 agentId;
        address requester;
        address validator;
        bytes data;
    }

    struct ValidationResponse {
        uint8 score;
        string reason;
        bool responded;
    }

    mapping(uint256 => address) private _agentOwners;
    uint256 private _nextAgentId;
    uint256 private _nextRequestId;

    // TODO 1: Declare a mapping from uint256 (requestId) to ValidationRequest
    // TODO 2: Declare a mapping from uint256 (requestId) to ValidationResponse

    constructor() {
        _nextAgentId = 1;
        _nextRequestId = 1;
    }

    function registerAgent() external returns (uint256) {
        uint256 agentId = _nextAgentId;
        _agentOwners[agentId] = msg.sender;
        _nextAgentId++;
        return agentId;
    }

    function requestValidation(uint256 agentId, address validator, bytes calldata data)
        external returns (uint256)
    {
        // TODO 3: Require msg.sender is the owner of the agent
        // TODO 4: Create the request with the current _nextRequestId
        // TODO 5: Store the ValidationRequest (agentId, msg.sender, validator, data)
        // TODO 6: Increment _nextRequestId and return the request ID
        return 0;
    }

    function respond(uint256 requestId, uint8 score, string calldata reason) external {
        // TODO 7: Require the request exists (requester != address(0))
        // TODO 8: Require msg.sender is the designated validator
        // TODO 9: Require score <= 100
        // TODO 10: Require not already responded
        // TODO 11: Store the ValidationResponse (score, reason, responded=true)
    }

    function getResponse(uint256 requestId)
        external view returns (uint8 score, string memory reason, bool responded)
    {
        // TODO 12: Return the response fields
        return (0, "", false);
    }
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

contract ValidationRegistryTest is Test {
    ValidationRegistry registry;
    address alice = address(1);
    address bob = address(2);
    address validator = address(3);

    function setUp() public {
        registry = new ValidationRegistry();
        vm.prank(alice);
        registry.registerAgent(); // agentId = 1
    }

    function test_RequestAndRespond() public {
        vm.prank(alice);
        uint256 reqId = registry.requestValidation(1, validator, "check-api");
        assertEq(reqId, 1);

        vm.prank(validator);
        registry.respond(1, 85, "API works well");

        (uint8 score, string memory reason, bool responded) = registry.getResponse(1);
        assertEq(score, 85);
        assertEq(reason, "API works well");
        assertTrue(responded);
    }

    function test_WrongValidatorReverts() public {
        vm.prank(alice);
        registry.requestValidation(1, validator, "check");

        vm.prank(bob);
        vm.expectRevert();
        registry.respond(1, 80, "not my job");
    }

    function test_ScoreAbove100Reverts() public {
        vm.prank(alice);
        registry.requestValidation(1, validator, "check");

        vm.prank(validator);
        vm.expectRevert();
        registry.respond(1, 101, "too high");
    }

    function test_NonOwnerCannotRequest() public {
        vm.prank(bob);
        vm.expectRevert();
        registry.requestValidation(1, validator, "hack");
    }

    function test_DoubleRespondReverts() public {
        vm.prank(alice);
        registry.requestValidation(1, validator, "check");

        vm.prank(validator);
        registry.respond(1, 90, "good");

        vm.prank(validator);
        vm.expectRevert();
        registry.respond(1, 95, "even better");
    }
}
`,
		Hints: []string{
			"mapping(uint256 => ValidationRequest) private _requests; mapping(uint256 => ValidationResponse) private _responses;",
			"requestValidation: require(_agentOwners[agentId] == msg.sender); _requests[reqId] = ValidationRequest(agentId, msg.sender, validator, data);",
			"respond: require(_requests[requestId].validator == msg.sender); require(score <= 100); require(!_responses[requestId].responded);",
		},
	}
}

func solERC8004Wallet() Question {
	return Question{
		ID: "sol-erc8004-wallet", Title: "Agent Wallet Verification via EIP-712",
		Difficulty: "hard", Category: "M7: ERC-8004", Language: LangSolidity,
		Description: `Each agent in ERC-8004 can have an associated wallet address. This wallet is
where the agent receives and sends payments (e.g., via x402). Setting the wallet
requires cryptographic proof that the wallet holder consents — you cannot assign
an arbitrary address as your agent's wallet.

The verification uses EIP-712 typed data signatures. The wallet holder signs a
message containing the agentId, their wallet address, and a deadline. The contract
verifies this signature on-chain using ecrecover and the EIP-712 digest.

The EIP-712 domain includes the contract name ("AgentWalletRegistry"), version "1",
the chain ID, and the verifying contract address. The type hash covers
"SetAgentWallet(uint256 agentId,address wallet,uint256 deadline)".

This pattern ensures only willing wallet holders can be linked to agents,
preventing impersonation. The deadline prevents replay of old signatures.

Implement the wallet verification contract with EIP-712 signature validation.`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract AgentWalletRegistry {
    mapping(uint256 => address) private _owners;
    mapping(uint256 => address) public agentWallets;
    uint256 private _nextAgentId;

    // TODO 1: Declare a public immutable bytes32 DOMAIN_SEPARATOR
    // TODO 2: Declare a public constant bytes32 SET_WALLET_TYPEHASH =
    //         keccak256("SetAgentWallet(uint256 agentId,address wallet,uint256 deadline)")

    constructor() {
        _nextAgentId = 1;
        // TODO 3: Compute DOMAIN_SEPARATOR using keccak256(abi.encode(
        //     keccak256("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"),
        //     keccak256(bytes("AgentWalletRegistry")),
        //     keccak256(bytes("1")),
        //     block.chainid,
        //     address(this)
        // ))
    }

    function registerAgent() external returns (uint256) {
        uint256 agentId = _nextAgentId;
        _owners[agentId] = msg.sender;
        _nextAgentId++;
        return agentId;
    }

    function setAgentWallet(
        uint256 agentId,
        address wallet,
        uint256 deadline,
        uint8 v, bytes32 r, bytes32 s
    ) external {
        // TODO 4: Require msg.sender is the owner of the agent
        // TODO 5: Require block.timestamp <= deadline (signature not expired)
        // TODO 6: Compute the struct hash: keccak256(abi.encode(SET_WALLET_TYPEHASH, agentId, wallet, deadline))
        // TODO 7: Compute the EIP-712 digest: keccak256(abi.encodePacked("\x19\x01", DOMAIN_SEPARATOR, structHash))
        // TODO 8: Recover the signer using ecrecover(digest, v, r, s)
        // TODO 9: Require recovered signer == wallet (wallet holder consents)
        // TODO 10: Store agentWallets[agentId] = wallet
    }

    function clearWallet(uint256 agentId) external {
        // TODO 11: Require msg.sender is the owner of the agent
        // TODO 12: Delete agentWallets[agentId]
    }
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

contract AgentWalletRegistryTest is Test {
    AgentWalletRegistry registry;
    address alice = address(1);
    uint256 walletPk = 0xBEEF;
    address walletAddr;

    function setUp() public {
        registry = new AgentWalletRegistry();
        walletAddr = vm.addr(walletPk);
        vm.prank(alice);
        registry.registerAgent(); // agentId = 1
    }

    function _sign(uint256 agentId, address wallet, uint256 deadline)
        internal view returns (uint8 v, bytes32 r, bytes32 s)
    {
        bytes32 structHash = keccak256(abi.encode(
            registry.SET_WALLET_TYPEHASH(),
            agentId, wallet, deadline
        ));
        bytes32 digest = keccak256(abi.encodePacked(
            "\x19\x01",
            registry.DOMAIN_SEPARATOR(),
            structHash
        ));
        (v, r, s) = vm.sign(walletPk, digest);
    }

    function test_SetAgentWallet() public {
        vm.warp(100);
        (uint8 v, bytes32 r, bytes32 s) = _sign(1, walletAddr, 200);

        vm.prank(alice);
        registry.setAgentWallet(1, walletAddr, 200, v, r, s);

        assertEq(registry.agentWallets(1), walletAddr);
    }

    function test_ExpiredDeadlineReverts() public {
        vm.warp(100);
        (uint8 v, bytes32 r, bytes32 s) = _sign(1, walletAddr, 50);

        vm.prank(alice);
        vm.expectRevert();
        registry.setAgentWallet(1, walletAddr, 50, v, r, s);
    }

    function test_WrongSignerReverts() public {
        vm.warp(100);
        // Sign with a different key
        uint256 otherPk = 0xDEAD;
        bytes32 structHash = keccak256(abi.encode(
            registry.SET_WALLET_TYPEHASH(),
            1, walletAddr, 200
        ));
        bytes32 digest = keccak256(abi.encodePacked(
            "\x19\x01",
            registry.DOMAIN_SEPARATOR(),
            structHash
        ));
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(otherPk, digest);

        vm.prank(alice);
        vm.expectRevert();
        registry.setAgentWallet(1, walletAddr, 200, v, r, s);
    }

    function test_ClearWallet() public {
        vm.warp(100);
        (uint8 v, bytes32 r, bytes32 s) = _sign(1, walletAddr, 200);
        vm.prank(alice);
        registry.setAgentWallet(1, walletAddr, 200, v, r, s);

        vm.prank(alice);
        registry.clearWallet(1);
        assertEq(registry.agentWallets(1), address(0));
    }
}
`,
		Hints: []string{
			"bytes32 public immutable DOMAIN_SEPARATOR; bytes32 public constant SET_WALLET_TYPEHASH = keccak256(\"SetAgentWallet(uint256 agentId,address wallet,uint256 deadline)\");",
			"bytes32 structHash = keccak256(abi.encode(SET_WALLET_TYPEHASH, agentId, wallet, deadline)); bytes32 digest = keccak256(abi.encodePacked(\"\\x19\\x01\", DOMAIN_SEPARATOR, structHash));",
			"address signer = ecrecover(digest, v, r, s); require(signer == wallet, \"invalid signature\");",
		},
	}
}

func solERC8004ReputationSummary() Question {
	return Question{
		ID: "sol-erc8004-reputation-summary", Title: "Reputation Summary with WAD",
		Difficulty: "hard", Category: "M7: ERC-8004", Language: LangSolidity,
		Description: `Aggregating reputation into a single summary is essential for making trust
decisions on-chain. The summary includes the average feedback score, total count,
and breakdowns of positive vs negative feedback.

Since Solidity has no floating-point numbers, we use WAD (Word as Decimal) math:
multiply values by 1e18 before division to preserve precision. For example,
an average of 40 is represented as 40 * 10^18 = 40000000000000000000.

The formula: averageWAD = (sum * WAD) / count, where WAD = 1e18.

This is the same precision pattern used throughout DeFi — Uniswap, Aave, and
MakerDAO all use WAD or similar fixed-point representations. In x402, USDC
amounts use 6 decimals, but reputation scores use 18-decimal WAD for maximum
precision in on-chain calculations.

The summary function should handle the zero-feedback case by returning all zeros.
Positive feedback has value > 0, negative has value < 0, and zero is neutral.

Implement the reputation summary with WAD-precision averaging.`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract ReputationSummary {
    struct Feedback {
        address provider;
        int256 value;
    }

    // TODO 1: Declare a constant uint256 WAD = 1e18
    // TODO 2: Declare a mapping from uint256 (agentId) to Feedback array

    function giveFeedback(uint256 agentId, int256 value) external {
        // TODO 3: Push a Feedback struct with msg.sender and value
    }

    function getSummary(uint256 agentId)
        external view
        returns (int256 averageWAD, uint256 count, uint256 positiveCount, uint256 negativeCount)
    {
        // TODO 4: Get the feedbacks array for this agentId
        // TODO 5: If count is 0, return (0, 0, 0, 0)
        // TODO 6: Loop through all feedbacks:
        //         - Accumulate sum (int256) of all values
        //         - Count positive (value > 0) entries
        //         - Count negative (value < 0) entries
        // TODO 7: Compute averageWAD = (sum * int256(WAD)) / int256(count)
        // TODO 8: Return (averageWAD, count, positiveCount, negativeCount)
        return (0, 0, 0, 0);
    }
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

contract ReputationSummaryTest is Test {
    ReputationSummary rep;
    address alice = address(1);
    address bob = address(2);
    address charlie = address(3);

    uint256 constant WAD = 1e18;

    function setUp() public {
        rep = new ReputationSummary();
    }

    function test_EmptySummary() public view {
        (int256 avg, uint256 count, uint256 pos, uint256 neg) = rep.getSummary(1);
        assertEq(avg, 0);
        assertEq(count, 0);
        assertEq(pos, 0);
        assertEq(neg, 0);
    }

    function test_ThreeFeedbacks() public {
        vm.prank(alice);
        rep.giveFeedback(1, 80);
        vm.prank(bob);
        rep.giveFeedback(1, -20);
        vm.prank(charlie);
        rep.giveFeedback(1, 60);

        (int256 avg, uint256 count, uint256 pos, uint256 neg) = rep.getSummary(1);
        // average = (80 + (-20) + 60) / 3 = 120 / 3 = 40
        assertEq(avg, 40 * int256(WAD));
        assertEq(count, 3);
        assertEq(pos, 2);
        assertEq(neg, 1);
    }

    function test_AllNegative() public {
        vm.prank(alice);
        rep.giveFeedback(1, -10);
        vm.prank(bob);
        rep.giveFeedback(1, -30);

        (int256 avg, uint256 count, uint256 pos, uint256 neg) = rep.getSummary(1);
        // average = (-10 + -30) / 2 = -20
        assertEq(avg, -20 * int256(WAD));
        assertEq(count, 2);
        assertEq(pos, 0);
        assertEq(neg, 2);
    }

    function test_SingleFeedback() public {
        vm.prank(alice);
        rep.giveFeedback(1, 100);

        (int256 avg, uint256 count, uint256 pos, uint256 neg) = rep.getSummary(1);
        assertEq(avg, 100 * int256(WAD));
        assertEq(count, 1);
        assertEq(pos, 1);
        assertEq(neg, 0);
    }
}
`,
		Hints: []string{
			"uint256 public constant WAD = 1e18;",
			"int256 sum = 0; for (uint256 i = 0; i < feedbacks.length; i++) { sum += feedbacks[i].value; if (feedbacks[i].value > 0) positiveCount++; ... }",
			"averageWAD = (sum * int256(WAD)) / int256(count); — cast WAD and count to int256 for signed division",
		},
	}
}

func solERC8004X402Integration() Question {
	return Question{
		ID: "sol-erc8004-x402-integration", Title: "x402 + Reputation Combined",
		Difficulty: "hard", Category: "M7: ERC-8004", Language: LangSolidity,
		Description: `The ultimate vision of ERC-8004 is agents that transact autonomously and build
reputation through those transactions. This question combines x402 payment
settlement with the reputation system into a single atomic operation.

The PayAndRate contract orchestrates two external calls:
1. settlement.settle(from, to, value, nonce) — executes the payment
2. reputation.giveFeedback(agentId, value, tag) — records the feedback

Both calls happen in a single transaction, ensuring atomicity: if the payment
fails, no feedback is recorded. This prevents feedback spam on failed payments.

The contract emits a PaymentRated event linking the payment nonce to the agent
and feedback value, creating an on-chain audit trail that connects payments
to reputation scores.

This pattern represents how autonomous agents will interact in production:
pay for a service via x402, then rate the service provider, all in one tx.

Implement the combined settle-and-rate contract with proper error propagation.`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

// TODO 1: Define interface ISettlement with function:
//         settle(address from, address to, uint256 value, bytes32 nonce) external

// TODO 2: Define interface IReputation with function:
//         giveFeedback(uint256 agentId, int256 value, string calldata tag) external

contract PayAndRate {
    // TODO 3: Declare public state variables for settlement (ISettlement) and reputation (IReputation)

    // TODO 4: Define event PaymentRated(bytes32 indexed nonce, uint256 indexed agentId, int256 feedbackValue)

    constructor(address _settlement, address _reputation) {
        // TODO 5: Store the settlement and reputation contract references
    }

    function settleAndRate(
        address from,
        address to,
        uint256 value,
        bytes32 nonce,
        uint256 agentId,
        int256 feedbackValue,
        string calldata tag
    ) external {
        // TODO 6: Call settlement.settle(from, to, value, nonce)
        //         If this reverts, the whole transaction reverts (no feedback recorded)
        // TODO 7: Call reputation.giveFeedback(agentId, feedbackValue, tag)
        // TODO 8: Emit the PaymentRated event
    }
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

contract MockSettlement {
    bool public shouldFail;
    bool public settleCalled;
    address public lastFrom;
    address public lastTo;
    uint256 public lastValue;
    bytes32 public lastNonce;

    function setFail(bool _fail) external {
        shouldFail = _fail;
    }

    function settle(address from, address to, uint256 value, bytes32 nonce) external {
        require(!shouldFail, "settlement failed");
        settleCalled = true;
        lastFrom = from;
        lastTo = to;
        lastValue = value;
        lastNonce = nonce;
    }
}

contract MockReputation {
    bool public feedbackCalled;
    uint256 public lastAgentId;
    int256 public lastValue;
    string public lastTag;

    function giveFeedback(uint256 agentId, int256 value, string calldata tag) external {
        feedbackCalled = true;
        lastAgentId = agentId;
        lastValue = value;
        lastTag = tag;
    }
}

contract PayAndRateTest is Test {
    MockSettlement settlement;
    MockReputation reputation;
    PayAndRate payAndRate;
    address alice = address(1);
    address bob = address(2);

    function setUp() public {
        settlement = new MockSettlement();
        reputation = new MockReputation();
        payAndRate = new PayAndRate(address(settlement), address(reputation));
    }

    function test_SettleAndRate() public {
        bytes32 nonce = keccak256("nonce1");
        payAndRate.settleAndRate(alice, bob, 1000, nonce, 1, 80, "good-service");

        assertTrue(settlement.settleCalled());
        assertEq(settlement.lastFrom(), alice);
        assertEq(settlement.lastTo(), bob);
        assertEq(settlement.lastValue(), 1000);
        assertEq(settlement.lastNonce(), nonce);

        assertTrue(reputation.feedbackCalled());
        assertEq(reputation.lastAgentId(), 1);
        assertEq(reputation.lastValue(), 80);
        assertEq(reputation.lastTag(), "good-service");
    }

    function test_SettlementFailRevertsAll() public {
        settlement.setFail(true);
        bytes32 nonce = keccak256("nonce2");

        vm.expectRevert();
        payAndRate.settleAndRate(alice, bob, 1000, nonce, 1, 80, "good");

        assertFalse(reputation.feedbackCalled());
    }

    function test_EmitsPaymentRatedEvent() public {
        bytes32 nonce = keccak256("nonce3");

        vm.expectEmit(true, true, false, true);
        emit PayAndRate.PaymentRated(nonce, 1, 80);
        payAndRate.settleAndRate(alice, bob, 1000, nonce, 1, 80, "great");
    }

    function test_NegativeFeedback() public {
        bytes32 nonce = keccak256("nonce4");
        payAndRate.settleAndRate(alice, bob, 500, nonce, 2, -30, "bad-service");

        assertEq(reputation.lastValue(), -30);
        assertEq(reputation.lastAgentId(), 2);
    }
}
`,
		Hints: []string{
			"interface ISettlement { function settle(address from, address to, uint256 value, bytes32 nonce) external; }",
			"ISettlement public settlement; IReputation public reputation; constructor: settlement = ISettlement(_settlement); reputation = IReputation(_reputation);",
			"settleAndRate: settlement.settle(from, to, value, nonce); reputation.giveFeedback(agentId, feedbackValue, tag); emit PaymentRated(nonce, agentId, feedbackValue);",
		},
	}
}
