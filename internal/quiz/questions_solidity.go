package quiz

// SolidityQuestions returns all Solidity quiz questions organized into 6 modules.
func SolidityQuestions() []Question {
	var all []Question
	all = append(all, solModule1Foundations()...)
	all = append(all, solModule2ERC20()...)
	all = append(all, solModule3Signatures()...)
	all = append(all, solModule4Gasless()...)
	all = append(all, solModule5Advanced()...)
	all = append(all, solModule6X402()...)
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
