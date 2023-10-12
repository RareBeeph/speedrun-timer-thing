package timing

/*
	TODO: make splithandler idle check not consider index 0 to always be idle

	idle: Timer.Idle(), SplitHandler.Idle()
	active: Timer.Running(), SplitHandler.Active()
	paused: Timer.Paused(), SplitHandler.Active()
	cancelled: Timer.Stopped(), SplitHandler.Active()
	finished: Timer.Stopped(), SplitHandler.Finished()

					Split		Pause		Stop
				-----------------------------------------
	idle		|	active	|	-		|	-			|
	active		|	(a/f)	|	paused	|	cancelled	|	// this row is shorthand for splithandler substates
	paused		|	-		|	active	|	cancelled	|	// this row is shorthand for splithandler substates
	cancelled	|	-		|	-		|	idle		|
	finished	|	-		|	-		|	idle		|
				-----------------------------------------

	no other states are valid
*/
