package com.vng.eventsourcing.Account.AccountService.Command

import com.vng.eventsourcing.Common.General.ICommand

/**
 * Created by Thuyen Phan
 * Package: com.vng.eventsourcing.Account.AccountService.Command
 * User: lap11852
 * Date: 27/01/2019
 * Time: 21:18
 */

class CreateAccount(val id : String, val username : String, val balance : Int) : ICommand()

class DepositMoney(val id : String, val amount : Int) : ICommand()

class WithdrawMoney(val id : String, val amount : Int) : ICommand()