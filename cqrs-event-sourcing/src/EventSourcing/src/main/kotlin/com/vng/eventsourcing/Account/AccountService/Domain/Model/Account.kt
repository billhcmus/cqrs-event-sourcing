package com.vng.eventsourcing.Account.AccountService.Domain.Model

import com.vng.eventsourcing.Account.Common.Event.AccountCreateEvent
import com.vng.eventsourcing.Account.Common.Event.AccountDepositEvent
import com.vng.eventsourcing.Account.Common.Event.AccountWithdrawEvent
import com.vng.eventsourcing.Common.General.Aggregate

/**
 * Created by Thuyen Phan
 * Package: com.vng.eventsourcing.Account.AccountService.Domain
 * User: lap11852
 * Date: 27/01/2019
 * Time: 21:17
 */


class Account(id : String, var username : String, var balance : Int) : Aggregate() {
    init {
        // do validation

        ApplyEvent(AccountCreateEvent(id, username, balance))
    }

    // Apply event effect to Accout
    private fun Apply(e : AccountCreateEvent) {
        Id = e.Id
        username = e.username
        balance = e.balance
    }

    private fun Apply(e : AccountDepositEvent) {
        balance += e.amount
    }

    private fun Apply(e : AccountWithdrawEvent) {
        balance -= e.amount
    }

    //

    fun DepositMoney(amount : Int) {
        ApplyEvent(AccountDepositEvent(Id, amount))
    }

    fun WithdrawMoney(amount : Int) {
        ApplyEvent(AccountWithdrawEvent(Id, amount))
    }
}