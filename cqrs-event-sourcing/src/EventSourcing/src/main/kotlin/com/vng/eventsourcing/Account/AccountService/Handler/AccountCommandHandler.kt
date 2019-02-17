package com.vng.eventsourcing.Account.AccountService.Handler

import com.vng.eventsourcing.Account.AccountService.Command.CreateAccount
import com.vng.eventsourcing.Account.AccountService.Command.DepositMoney
import com.vng.eventsourcing.Account.AccountService.Command.WithdrawMoney
import com.vng.eventsourcing.Account.AccountService.Domain.Model.Account
import com.vng.eventsourcing.Common.Repository.IRepository

/**
 * Created by Thuyen Phan
 * Package: com.vng.eventsourcing.Account.AccountService.Handler
 * User: lap11852
 * Date: 27/01/2019
 * Time: 21:17
 */

class AccountCommandHandler(val repository : IRepository) {
    fun Handle(command : CreateAccount) {
        // do validation

        // Process
        val series = Account(
            command.id,
            command.username,
            command.balance
        )
        repository.Save(series)
    }

    fun Handle(command : DepositMoney) {
        // validation
        val account = repository.GetById(command.id)
        account.DepositMoney(command.amount)
        repository.Save(account)
        // store and publish event
    }

    fun Handle(command : WithdrawMoney) {
        // validation
        val account = repository.GetById(command.id)
        account.WithdrawMoney(command.amount)
        repository.Save(account)
    }
}