package com.vng.eventsourcing.Account.Common.Event

import com.vng.eventsourcing.Common.General.Event

/**
 * Created by Thuyen Phan
 * Package: com.vng.eventsourcing.Account.Common.Event
 * User: lap11852
 * Date: 27/01/2019
 * Time: 21:05
 */

class AccountCreateEvent(id: String, var username: String, var balance: Int) : Event() {
    override var Id: String
        get() = TODO("not implemented") //To change initializer of created properties use File | Settings | File Templates.
        set(value) {}
    override var Version: Int
        get() = TODO("not implemented") //To change initializer of created properties use File | Settings | File Templates.
        set(value) {}
    init {
        Id = id
    }

    constructor(id : String, username: String, balance: Int, version : Int) : this(id, username, balance) {
        Version = version
    }
}