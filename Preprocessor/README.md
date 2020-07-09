# ECE_653_Project
ECE653 Project Involving implementing a SAT pre-processor in GO for the GINI SAT solver.

## Important Notes:

Found some old pre-processing code in GOPHERSAT's original upload version. Was able to combine the pre-processing code I had already written and used the formatting, structure, and functions methods that are contained in some of GOPHERSAT's source files. Extracting these and giving it standalone functionality took quite a bit of refactoring/reformatting. 

Currently subsumption and self-subsuming resolution is implemented. This should be re-confirmed.

## TO DO:

1. Test current pre-processing to make sure it's working properly. Need better test examples.
2. Currently for some of the examples the pre-processing increases the number of clauses.. This has to be looked at.
3. Need to implement other pre-processing techniques.
4. Needs some work to obviously integrate with Gini specifically as the Gini package has its clauses, variables and literals stored differently.
