package com.vng.eventsourcing.Common.General

/**
 * Created by Thuyen Phan
 * Package: com.vng.eventsourcing.Common.General
 * User: lap11852
 * Date: 27/01/2019
 * Time: 21:55
 */

abstract class Event {
    abstract var Id : String
    abstract var Version : Int
}
