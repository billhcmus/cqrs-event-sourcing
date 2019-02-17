```How to design the software that is less error prone```
### CQRS
- CQRS is an implementation of Command Query Separation principle to the architecture of software.
- The journey starts with something like this, this is N-Layer architecture.
    - DTO: Data transfer object
    - ORM: Object Relational Mapping

![](https://www.future-processing.com/blog/wp-content/uploads/2015/04/1_layered.png)

- If we simply separate business logic into Commands and Querries, then we get some CQS here.

![](https://www.future-processing.com/blog/wp-content/uploads/2015/04/CQRS-Simple-Architecture_2_CQS_1.png)

- The main idea behind CQS is: "A method should either change state of an object, or return a result, but not both."
- So what Command and Query is?
    - Commands - change the state of an object or entire system. Commands will return void. (modifiers or mutators). Firing command is the only way tho change the state of the system. 
    - Queries: return results and do not change the state of an object. Queries will declare return type. It reads the state of the system, filters, aggregates and transform data to deliver it in the most useful format. It can be executed multiple times and will not affect the state of the system.
- Model to Domain Model
    - Model is a group of containers for data
    - Domain Model encapsulates essential complexity of business rules.

- Many applications that use mainstream approaches consists of models which are common for read and write side. Having the same model for read and write side leads to a more complex model that could be very difficult to be maintained and optimized. So...

![](https://www.future-processing.com/blog/wp-content/uploads/2015/04/CQRS_CQS_WRITE_READ.png)

- The real strength of these two patterns it that you can separate methods that change state from those that don't.
- You can optimize the read side of the system separately from the write side. The write side is known as the domain. The domain contains all the behaviour. The read side is specialized for reporting needs.
- Another benefit of this pattern is in the case of large applications. You can split developers into smaller teams working on different sides of the system (read or write) without knowledge of the other side.

- When ORM introduces overhead, it would be useful to simplify this architechture:
![](https://www.future-processing.com/blog/wp-content/uploads/2015/04/CQRS_4_CQS_NoREAD.png)

- Now we have ```Read``` and ```Write``` model separated at the logical level, but only this. Both of them still shares common database.

- Fully separated data models:
![](https://www.future-processing.com/blog/wp-content/uploads/2015/04/CQRS_5_CQRS.png)

- So, what is projection? Let's see it bellow...

### Event Sourcing

- Event Sourcing is an idea that was presented along with CQRS, and is often identified as apart of CQRS.
- Do not store state of an object. Instead, we store all the events impacting its state. To retrieve an object state we read the different events related to this object and applied them one by one.
- When we design any app, we immediately translate some specs of this to some storage mechanism. If we use SQL we design the tables, if it's NoSQL we design documents. ```This forces us to think of everything in terms of current state```.
- Having a question as ```How do I store this thing so I can retrieve it later?```

![](https://res.cloudinary.com/practicaldev/image/fetch/s--PodGpeUR--/c_limit%2Cf_auto%2Cfl_progressive%2Cq_auto%2Cw_880/https://thepracticaldev.s3.amazonaws.com/i/bv1l8eatoljyin82dykd.png)

- We use the same model for both reads and writes. Typically we'd design our tables from a write perspective, the we figure out how to build our query on-top of those structures. When the system grows, it won't work well. This soon becomes unmaintainable and expensive to change.
- Normally we only storing the current state of the world, we have no idea how system got into that state in the first place. We can't answer some question like ```How many times has a user changed their email address?``` or ```How many people added an item to their cart, then removed it, then bought that item a month later?``` This look like useful business information, you are losing when you starting store your data.

> Ensures that all changes to application state are stored as a sequence of events.
> Instead of focussing on current sate, you focus on the changes that have occured over time.

- Example: A cart's lifecycle could be modelled as the following sequence of events:

| Event                   |
|-------------------------|
| Create Shopping Cart    |
| Item added to Cart      |
| Item added to Cart      |
| Item removed from Cart  |
| Shopping Cart Check-out |

And then, we have fully modelled as a sequence of events and that call Event Sourcing

![](https://res.cloudinary.com/practicaldev/image/fetch/s--FL_ZjOYw--/c_limit%2Cf_auto%2Cfl_progressive%2Cq_auto%2Cw_880/https://thepracticaldev.s3.amazonaws.com/i/xz6bej22iw97t46v2g6w.png)

- Any process can be modelled as a sequence of events. Infact, every process IS modelled as a sequence of events, they won't talk about "tables" and "joins", they will describle process as the series of events that can occur and the rules applied to them.
- Most business operations have constraints, a hard rule that cannot be broken.

```How do we answer some question about state in event sourced system?```

- We replay a subset of the events to get the answer we need. For example, in above example about Shopping Cart, we have business rules, an item must be in cart before it can be removed. So we need to answer the question, "Does this item exists?", how do we do that without state? We just need to check that an "Item added to Cart" event has happened for that item, then we know the item is in the cart and that it can be removed.
- The above action we usually called ```projecting``` the events, and the end result is a ```projection```

- Isn't this expensive and time consuming?
    - First thing we need the tiniest subset of events. Fetching the usefull history of a concept, this is typically a single database call.
    - You load the events and replay them in memory and ```projecting``` them to build up your dataset.
    - This is lighting fast, we doing it on the local processor, rather than making a series of SQL calls (network call).

```If every piece of state is derived from events, how do we fetch data that needs to be presented to the user? Do we need fetch all the events and build the datset each time?```
- The answer is no, instead of building it on the fly, we build it in the background, store the intermediate results in a database. Then user can query for data, and they will get it in the exact shape they need, with minimal delay, you cache the results for later use.

![](https://res.cloudinary.com/practicaldev/image/fetch/s--O4HIIlNN--/c_limit%2Cf_auto%2Cfl_progressive%2Cq_auto%2Cw_880/https://thepracticaldev.s3.amazonaws.com/i/lhhlvel6ifqaj0y4otb7.png)

- With Event Sourcing, You are no longer bound by your current table structure.  
```If you need to present data in a new shape, simply build up a new data structure designed to answer that query. This give you complete freedom to build and implement your read models in any way you want, discarding old models when they're no longer needed.```

#### The benefits

1. Empemeral data-structures
    - If you need to view data in a new way, you simply create a new projection of data and projecting it into the shape you need.
2. Easier to communicate with domain experts
    - Domain experts express business process as a series of events.
    - Building an event sourced system, we're modelling the system exactly as they describe it.
    - This makes communication a lot smoother.
3. Reports become painless
    - You have the full history of every event thathas ever happened.
4. Composing services becomes trivial
    - Plugging two systems together is usually quite tricky.
    - ES solve this problem by letting servicescommunicate via events. Need to trigger a process in another service when something happens? Write an event listener that runs whenever that event is fired.
5. Lightning fast on standard databases
    - You can use a standard MySQL to store your events, it's optimised for append only operations, which mean that storing data is fast.
6. Easy to change database implementations.

- Since each entity is represented by a stream of events there is no way to reliably query the data. This is especially evident on data intensive applications where much of business value relies on analyzing the dta. But Event Sourcing unsuitable for most applications and only relevant for very isolated and specific use cases.
- That's wehn CQRS comes in. According to above explain about CQRS, CQRS describes the concept of having two different models, one to change information and one to read it, completely separated from each other.

#### The issues

1. Eventual Consistency
    - Having a queue separates the write and the read model. That means any change will be available for read in between the next miliseconds and foreseeable future. It mean a whole is eventually consistent. It mean the system settles on a value can return stale or inconsistent data.
    - By definition that queue between the write and read model can fill up, the system can have an unforeseen peak of usage and can take more than it is expected to process.
2. Event Upgrading
    - When you fire events, shape will change, and this can be a bit trickly to handle.
    - When event changes shape, you have to write an upgrader that takes the old event and converts it into the new one.