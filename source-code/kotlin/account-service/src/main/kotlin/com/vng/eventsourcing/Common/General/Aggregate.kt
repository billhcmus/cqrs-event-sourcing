package com.vng.eventsourcing.Common.General

/**
 * Created by Thuyen Phan
 * Package: com.vng.eventsourcing.Common.Aggregate
 * User: lap11852
 * Date: 27/01/2019
 * Time: 21:35
 */
abstract class Aggregate(val events : MutableList<Event>, private var version : Int = -1) {
    lateinit var Id : String
    protected set

    constructor() : this(mutableListOf(), -1)

    fun ApplyEvent(event : Event, isNew : Boolean) {
        if (isNew) {
            events.add(event)
            version++
        }
    }

    fun ApplyEvent(event : Event) {
        ApplyEvent(event, true)
    }
}