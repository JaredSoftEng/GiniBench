# GiniPre
GiniPre is a SAT preprocessor written entirely in GO. It leverages the efficient data structures present in the gini SAT solver (designed by irifrance). GiniPre has a dependency on our fork of gini in order to access some methods which were previously hidden internally.

It enables efficient preprocessing techniques such as:
- Subsumption
- Self-Subsuming Resolution

The command-line interface can be directly used to solve files using the gini CDCL Sat solver backend, or with the -cnf option to just perform pre-processing on an input file without applying the solver. A -ui option enables a file-picker and the ability to run the command on a directory of files while pushing the results to a .csv report. 

GiniPre accepts DIMACS and AIGER files (.cnf, .bz2, .gz, .aag, .aig).

A sample command:
  ginipre -fullsub -model "%PWD%\testfile.cnf"
  
  Would output statistics showing the number of literals reduced through self-subsumption, the clauses removed by full subsumption and if the resulting solve is satisfiable by the gini SAT solver, would output the satisfying literal assignment model.
