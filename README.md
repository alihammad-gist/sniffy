# Todo

Instead of filtering and sending event to one reciever. Take {Filter, EventChan} combo as NewWatcher argument. And if event successfuly passes through filter send the Event to corresponding EventChan. 