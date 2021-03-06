# Ontology Layer2

中文|[English](README.md)

## 名称解释

### Layer2交易

用户在Layer2进行转账或者执行合约的请求，用户已经对其签名。这个交易可以和Ontology主链的交易格式一样，也可以不一样。

### Node

Node是Layer2交易收集器，它负责收集用户的Layer2交易，验证并执行交易，每生成一个Layer2区块，Node负责执行区块中的交易，更新状态，并生成Layer2合约可以解释的、用于安全性保证的状态证明。

### Layer2区块

Node周期性的打包收集到的Layer2交易，生成一个包含这期间所有Layer2交易的区块，产生一个新的Layer2区块。

### Layer2 State

Node执行Layer2区块中打包的交易，更新状态，将所有更新的状态数据排序生成一个Merkle树，计算Merkle树的根hash，该根hash即为该区块的Layer2 State。

### Operator

Operator是Layer2的安全守护程序，负责监听Ontology主链是否有到Layer2的代币转移或者Layer2都Ontology主链的代币转移交易，同时Operator还负责周期性的将Layer2的状态证明提交到Ontology主网作为证明。

### Challenger

负责验证Operator提交到Ontology主链的状态证明。这要求Challenger从Operator或者链上同步Layer2交易，维护完整的全局状态。在Challenger同步执行交易并更新状态后，可以验证Operator提交在Ontology主网的状态证明正确性，如果不正确，Challenger可以生成Layer2合约可以解释的欺诈证明挑战Operator。

### 账户状态证明

包括账户状态信息以及其Merkle证明，可以从Operator和Challenger查询来获取。只有他们维护有完整的全局状态。

### 欺诈证明

欺诈证明包含当前Layer2区块更新之前账户状态证明。因为有以前Layer2区块的状态证明以及提交的账户状态证明，可以证明更新之前旧状态的合法性，在旧状态合法的情况下，运行当前区块即可证明新的状态证明合法性。

## 工作流程

###	Deposit到Layer2

1.	用户在Ontology主链进行deposit操作，主链合约锁定用户Deposit的资金，记录这笔资金在Layer2的状态，此时状态为“未释放”。

2.	Operator查询到主链Ontology上有deposit操作，Operator会提交deposit到Layer2 Node。

3.	Node将该交易和收集的其他用户交易一起打包到一个新的Layer2区块，执行Layer2区块，生成新的Layer2 State。

4.	Operator会监听Layer2 Node新的区块，提交新的Layer2 State到Ontology主链时，会附带提交deposit已经释放的请求。

5.	主链合约执行deposit释放操作，修改deposit资金状态为“已释放”。

<div align=center><img width="360" height="450" src="doc/pic/user_deposit.png"/></div>

### Withdraw到Ontology

1.	用户构造Withdraw的Layer2交易并提交给Node。

2.	Node根据Withdraw修改其State，同时打包该Withdraw交易以及其他用户交易一起到一个Layer2区块。

3.	Operator提交这个Layer2区块State到Ontology主链时，会附带提交withdraw请求。

4.	主链合约执行Withdraw请求，记录一笔withdraw资金记录，并设置状态为“未释放”。

5.	在State确认后，主链合约执行withdraw释放请求，给目标账户转账，同时设置withdraw记录为“已释放“

<div align=center><img width="499" height="450" src="doc/pic/user_withdraw.png"/></div>

###	用户Layer交易

1.	用户构造Transfer的Layer2交易并提交给Node。

2.	Node打包该transfer交易以及其他交易到一个Layer2区块，执行区块中的交易，提交这个Layer2区块State到Ontology主链。

3.	等待State确认。

### 用户Layer交易以及安全保证流程

<div align=center><img width="650" height="450" src="doc/pic/system.png"/></div>

## 账户模型

账户使用Merkle树的方式来组织，但这是一个可以跟踪更新的Merkle树。Merkle树包含了更新之前的Root Hash和更新的账户.

<div align=center><img width="450" height="660" src="doc/pic/account.png"/></div>

每个State Root都固定对应一个高度，从0开始从下往上一次递增。我们需要在链上记录每个高度上的State Root。

如何证明一个账户的状态？

如果账户在高度为H时账户状态有更新，那么在高度为H的state树中包含有其账户状态，这是一个Merkle树，从全局看又是一个子树。但我们在链上有这个子树的root hash，这样可以在这个子树上生成Merkle proof，结合这个子树的root hash可以证明这个账户的状态。

但这个账户可能不是最新的，因为后续的更新产生的新的Merkle树中包含了更新的账户状态。所以有挑战机制，挑战者只需要提交这个账户的Merkle proof，其root hash所在的高度更高，那么其挑战成功。

为什么要这种Merkle树？

在链上只有每个高度的state root，有时候需要链上验证state变更的有效性。状态转换可以写成如下形式：

S = F（S‘，Txs）

其中S‘是上一个状态，Txs是交易，S为执行Txs后生成的新状态。

最简单的链上验证状态转换方式是提交全局状态S‘和Txs，根据链上的S’的state root来验证全局状态S‘，在这个全局状态下执行Txs，生成一个新的全局状态S，根据链上S的state root来验证状态转换有效性。

但这个方式有许多问题：

1.	全局状态往往很大，无论数据量大小还是在这个全局状态下执行交易，都有很大的限制要求。

Txs往往不会影响到全局所有状态，只是其中很小部分，只有很小部分状态有更新，以上账户模型只跟踪状态的更新，这里有一个新的实现方式。

不提交全局状态的S‘，仅仅提交S’中会被Txs更新的局部状态以及其Merkle proof，在上面已经介绍过如何证明这个状态。在局部状态下执行Txs，生成更新的局部状态，再加上S‘的state root可以计算新的S的state root，从而验证状态转换的有效性。

有哪些好处：

1.	无论状态有多少，但只要每次更新的状态不大，那么其子树很小，其Merkle proof也比较小

2.	验证状态转换代价小，效率高，只需要提交较少的有更新需求的状态以及其很小的Merkle proof，就可以验证状态转换的有效性。

3.	可以证明一个状态更新过程