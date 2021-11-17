# Wallet : (This is in progress)

Service has to provide good management of transactions in any currency between wallet with 1.5% of fee.

### Instances
  * Transaction
  * Wallet
  * User

### Todo
###### level 1
1. Implement authorization for user
2. Create minimum 2 users 
3. Create at least 1 wallet for user
4. Make a transaction between wallets through API 
###### level 2
2. Implement fee collection for each transaction (amount of fee 1.5% or from config file). The amount of fee must be stored in the wallet UUID '85aa7525-4fdb-4436-a600-66ffc55e0f65'
3. Implementation of transactions has not to have any critical sections (https://en.wikipedia.org/wiki/Critical_section)
4. The Transaction must be idempotent
5. Priority queue (depends on type of transaction, Priority: 0 - low , 1 - high)

### Hints
* Use transaction in the context of databases (ACID) for making the transactions between wallets
* Use abstract data structure for transactions, such as Queue, Priority queue, Stacks 
* Use JWT for authorization at other endpoints
### For development

```bash
$ git clone git@github.com:GolangUA/workshops.git
$ cd wallet
$ docker-compose up -d #for creating DB and API doc server
```
**Note:** API documentation [here](http://localhost:8080).