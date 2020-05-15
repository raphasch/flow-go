import FungibleToken from 0x%s
import FlowToken from 0x%s
import FeeContract from 0x%s

pub contract ServiceAccount {

    pub let transactionFee: UFix64
    pub let accountCreationFee: UFix64

    pub fun deductTransactionFee(_ acct: AuthAccount) {
        let tokenVault = self.defaultTokenVault(acct)
        let feeVault <- tokenVault.withdraw(amount: self.transactionFee)

        FeeContract.deposit(from: <- feeVault)
    }

    pub fun deductAccountCreationFee(_ acct: AuthAccount) {
        let tokenVault = self.defaultTokenVault(acct)
        let feeVault <- tokenVault.withdraw(amount: self.accountCreationFee)

        FeeContract.deposit(from: <- feeVault)
    }

    pub fun initDefaultToken(_ acct: AuthAccount) {
        // Create a new FlowToken Vault and put it in storage
        acct.save(<-FlowToken.createEmptyVault(), to: /storage/flowTokenVault)

        // Create a public capability to the Vault that only exposes
        // the deposit function through the Receiver interface
        acct.link<&{FungibleToken.Receiver}>(
            /public/flowTokenReceiver,
            target: /storage/flowTokenVault
        )

        // Create a public capability to the Vault that only exposes
        // the balance field through the Balance interface
        acct.link<&{FungibleToken.Balance}>(
            /public/flowTokenBalance,
            target: /storage/flowTokenVault
        )
    }

    pub fun defaultTokenBalance(_ acct: PublicAccount): UFix64 {
        let balanceRef = acct
            .getCapability(/public/flowTokenBalance)!
            .borrow<&{FungibleToken.Balance}>()!

        return balanceRef.balance
    }

    pub fun defaultTokenVault(_ acct: AuthAccount): &FlowToken.Vault {
        return acct.borrow<&FlowToken.Vault>(from: /storage/flowTokenVault) ?? panic("Unable to borrow reference to the default token vault")
    }

    init() {
        self.transactionFee = 0.001
        self.accountCreationFee = 0.001
    }
}
