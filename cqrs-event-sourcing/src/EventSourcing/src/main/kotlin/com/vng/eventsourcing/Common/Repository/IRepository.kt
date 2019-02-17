package com.vng.eventsourcing.Common.Repository

import com.vng.eventsourcing.Account.AccountService.Domain.Model.Account
import com.vng.eventsourcing.Common.General.Aggregate

/**
 * Created by Thuyen Phan
 * Package: com.vng.eventsourcing.Common.Repository
 * User: lap11852
 * Date: 27/01/2019
 * Time: 21:33
 */
interface IRepository {
    fun Save(aggregate : Aggregate)
    fun GetById(Id : String) : Account
}

/*
Save method:
access list event in Aggregate then persistent to event store
then publish to message bus
 */