# depatureboard
Command line depature board using data from the National Rail site written in Golang. Depature board will poll the National rail site for 
updates on depatures for a selected station.

## How to use?
```
go run main.go -station WWA
```
Simply, run or compile the program and pass your three character station code to receive depature updates for your station in question.

## Options
```
-destination - string - will filter results down to the provided three character destination code
-refresh - integer - will change the frequency in which the national rail site is polled, default being one minute
-results - integer - will change the maxinum number of results to display, default being 30
```
