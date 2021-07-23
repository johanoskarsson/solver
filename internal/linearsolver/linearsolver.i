// Copyright 2010-2018 Google LLC
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This .i file exposes the linear and integer programming APIs. It was adapted
// from ortools/linear_solver/python/linear_solver.i.

%include "ortools/base/base.i"

%include "std_string.i"
%include "stdint.i"

%include "absl/base/attributes.h"

// We need to forward-declare the proto here, so that the PROTO_* macros
// involving them work correctly. The order matters very much: this declaration
// needs to be before the %{ #include ".../linear_solver.h" %}.
namespace operations_research {
class MPModelProto;
class MPModelRequest;
class MPSolutionResponse;
}  // namespace operations_research

%{
#include "absl/status/status.h"
#include "ortools/linear_solver/linear_solver.h"
#include "ortools/linear_solver/model_exporter.h"
#include "ortools/linear_solver/model_exporter_swig_helper.h"
#include "ortools/linear_solver/model_validator.h"

typedef ::absl::Status Status;
%}


%extend operations_research::MPVariable {
  std::string __str__() {
    return $self->name();
  }
  std::string __repr__() {
    return $self->name();
  }
}

%extend operations_research::MPSolver {
  // Change the API of LoadModelFromProto() to simply return the error message:
  // it will always be empty iff the model was valid.
  std::string LoadModelFromProto(const operations_research::MPModelProto& input_model) {
    std::string error_message;
    $self->LoadModelFromProto(input_model, &error_message);
    return error_message;
  }

  // Change the API of ExportModelAsLpFormat() to simply return the model.
  std::string ExportModelAsLpFormat(bool obfuscated) {
    operations_research::MPModelExportOptions options;
    options.obfuscate = obfuscated;
    operations_research::MPModelProto model;
    $self->ExportModelToProto(&model);
    return ExportModelAsLpFormat(model, options).value_or("");
  }

  // Change the API of ExportModelAsMpsFormat() to simply return the model.
  std::string ExportModelAsMpsFormat(bool fixed_format, bool obfuscated) {
    operations_research::MPModelExportOptions options;
    options.obfuscate = obfuscated;
    operations_research::MPModelProto model;
    $self->ExportModelToProto(&model);
    return ExportModelAsMpsFormat(model, options).value_or("");
  }

  /// Set a hint for solution.
  ///
  /// If a feasible or almost-feasible solution to the problem is already known,
  /// it may be helpful to pass it to the solver so that it can be used. A
  /// solver that supports this feature will try to use this information to
  /// create its initial feasible solution.
  ///
  /// Note that it may not always be faster to give a hint like this to the
  /// solver. There is also no guarantee that the solver will use this hint or
  /// try to return a solution "close" to this assignment in case of multiple
  /// optimal solutions.
  void SetHint(const std::vector<operations_research::MPVariable*>& variables,
               const std::vector<double>& values) {
    if (variables.size() != values.size()) {
      LOG(FATAL) << "Different number of variables and values when setting "
                 << "hint.";
    }
    std::vector<std::pair<const operations_research::MPVariable*, double> >
        hint(variables.size());
    for (int i = 0; i < variables.size(); ++i) {
      hint[i] = std::make_pair(variables[i], values[i]);
    }
    $self->SetHint(hint);
  }

  /// Sets the number of threads to be used by the solver.
  bool SetNumThreads(int num_theads) {
    return $self->SetNumThreads(num_theads).ok();
  }


// Catch runtime exceptions in class methods
%exception operations_research::MPSolver {
    try {
      $action
    } catch ( std::runtime_error& e ) {
      SWIG_exception(SWIG_RuntimeError, e.what());
    }
  }


  static double Infinity() { return operations_research::MPSolver::infinity(); }
  void SetTimeLimit(int64 x) { $self->set_time_limit(x); }
  int64 WallTime() const { return $self->wall_time(); }
  int64 Iterations() const { return $self->iterations(); }
}  // extend operations_research::MPSolver

%extend operations_research::MPVariable {
  double SolutionValue() const { return $self->solution_value(); }
  bool Integer() const { return $self->integer(); }
  double Lb() const { return $self->lb(); }
  double Ub() const { return $self->ub(); }
  void SetLb(double x) { $self->SetLB(x); }
  void SetUb(double x) { $self->SetUB(x); }
  double ReducedCost() const { return $self->reduced_cost(); }
}  // extend operations_research::MPVariable

%extend operations_research::MPConstraint {
  double Lb() const { return $self->lb(); }
  double Ub() const { return $self->ub(); }
  void SetLb(double x) { $self->SetLB(x); }
  void SetUb(double x) { $self->SetUB(x); }
  double DualValue() const { return $self->dual_value(); }
}  // extend operations_research::MPConstraint

%extend operations_research::MPObjective {
  double Offset() const { return $self->offset();}
}  // extend operations_research::MPObjective





%ignoreall

%unignore operations_research;

// Strip the "MP" prefix from the exposed classes.
%rename (Solver) operations_research::MPSolver;
%rename (Solver) operations_research::MPSolver::MPSolver;
%rename (Constraint) operations_research::MPConstraint;
%rename (Variable) operations_research::MPVariable;
%rename (Objective) operations_research::MPObjective;

// Expose the MPSolver::OptimizationProblemType enum.
%unignore operations_research::MPSolver::OptimizationProblemType;
%unignore operations_research::MPSolver::GLOP_LINEAR_PROGRAMMING;
%unignore operations_research::MPSolver::CLP_LINEAR_PROGRAMMING;
%unignore operations_research::MPSolver::GLPK_LINEAR_PROGRAMMING;
%unignore operations_research::MPSolver::SCIP_MIXED_INTEGER_PROGRAMMING;
%unignore operations_research::MPSolver::CBC_MIXED_INTEGER_PROGRAMMING;
%unignore operations_research::MPSolver::GLPK_MIXED_INTEGER_PROGRAMMING;
%unignore operations_research::MPSolver::BOP_INTEGER_PROGRAMMING;
%unignore operations_research::MPSolver::SAT_INTEGER_PROGRAMMING;
// These aren't unit tested, as they only run on machines with a Gurobi license.
%unignore operations_research::MPSolver::GUROBI_LINEAR_PROGRAMMING;
%unignore operations_research::MPSolver::GUROBI_MIXED_INTEGER_PROGRAMMING;
%unignore operations_research::MPSolver::CPLEX_LINEAR_PROGRAMMING;
%unignore operations_research::MPSolver::CPLEX_MIXED_INTEGER_PROGRAMMING;
%unignore operations_research::MPSolver::XPRESS_LINEAR_PROGRAMMING;
%unignore operations_research::MPSolver::XPRESS_MIXED_INTEGER_PROGRAMMING;


// Expose the MPSolver::ResultStatus enum.
%unignore operations_research::MPSolver::ResultStatus;
%unignore operations_research::MPSolver::OPTIMAL;
%unignore operations_research::MPSolver::FEASIBLE;  // No unit test
%unignore operations_research::MPSolver::INFEASIBLE;
%unignore operations_research::MPSolver::UNBOUNDED;  // No unit test
%unignore operations_research::MPSolver::ABNORMAL;
%unignore operations_research::MPSolver::NOT_SOLVED;  // No unit test
%rename (SolverStatus) operations_research::MPSolver::ResultStatus;
%rename (StatusOptimal) operations_research::MPSolver::OPTIMAL;
%rename (StatusFeasible) operations_research::MPSolver::FEASIBLE;  // No unit test
%rename (StatusInfeasible) operations_research::MPSolver::INFEASIBLE;
%rename (StatusUnbounded) operations_research::MPSolver::UNBOUNDED;  // No unit test
%rename (StatusAbnormal) operations_research::MPSolver::ABNORMAL;
%rename (StatusNotSolved) operations_research::MPSolver::NOT_SOLVED;  // No unit test

// Expose the MPSolver's basic API, with some renames.
%rename (Objective) operations_research::MPSolver::MutableObjective;
%rename (BoolVar) operations_research::MPSolver::MakeBoolVar;  // No unit test
%rename (IntVar) operations_research::MPSolver::MakeIntVar;
%rename (NumVar) operations_research::MPSolver::MakeNumVar;
%rename (Var) operations_research::MPSolver::MakeVar;
// We intentionally don't expose MakeRowConstraint(LinearExpr), because this
// "natural language" API is specific to C++: other languages may add their own
// syntactic sugar on top of MPSolver instead of this.
%rename (Constraint) operations_research::MPSolver::MakeRowConstraint(double, double);
%rename (Constraint) operations_research::MPSolver::MakeRowConstraint();
%rename (Constraint) operations_research::MPSolver::MakeRowConstraint(double, double, const std::string&);
%rename (Constraint) operations_research::MPSolver::MakeRowConstraint(const std::string&);
%unignore operations_research::MPSolver::~MPSolver;
%unignore operations_research::MPSolver::Solve;
%unignore operations_research::MPSolver::VerifySolution;
%unignore operations_research::MPSolver::infinity;
%unignore operations_research::MPSolver::set_time_limit;  // No unit test

// Proto-based API of the MPSolver. Use is encouraged.
%unignore operations_research::MPSolver::SolveWithProto;
%unignore operations_research::MPSolver::ExportModelToProto;
%unignore operations_research::MPSolver::FillSolutionResponseProto;
// LoadModelFromProto() is also visible: it's overridden by an %extend, above.
%unignore operations_research::MPSolver::LoadSolutionFromProto;  // No test

// Expose some of the more advanced MPSolver API.
%unignore operations_research::MPSolver::InterruptSolve;
%unignore operations_research::MPSolver::SupportsProblemType;  // No unit test
%unignore operations_research::MPSolver::wall_time;  // No unit test
%unignore operations_research::MPSolver::Clear;  // No unit test
%unignore operations_research::MPSolver::constraints;
%unignore operations_research::MPSolver::variables;
%unignore operations_research::MPSolver::NumConstraints;
%unignore operations_research::MPSolver::NumVariables;
%unignore operations_research::MPSolver::EnableOutput;  // No unit test
%unignore operations_research::MPSolver::SuppressOutput;  // No unit test
%rename (LookupConstraint)
    operations_research::MPSolver::LookupConstraintOrNull;
%rename (LookupVariable) operations_research::MPSolver::LookupVariableOrNull;
%unignore operations_research::MPSolver::SetSolverSpecificParametersAsString;
%unignore operations_research::MPSolver::NextSolution;
// %unignore operations_research::MPSolver::ExportModelAsLpFormat;
// %unignore operations_research::MPSolver::ExportModelAsMpsFormat;

// Expose very advanced parts of the MPSolver API. For expert users only.
%unignore operations_research::MPSolver::ComputeConstraintActivities;
%unignore operations_research::MPSolver::ComputeExactConditionNumber;
%unignore operations_research::MPSolver::nodes;
%unignore operations_research::MPSolver::iterations;  // No unit test
%unignore operations_research::MPSolver::BasisStatus;
%unignore operations_research::MPSolver::FREE;  // No unit test
%unignore operations_research::MPSolver::AT_LOWER_BOUND;
%unignore operations_research::MPSolver::AT_UPPER_BOUND;
%unignore operations_research::MPSolver::FIXED_VALUE;  // No unit test
%unignore operations_research::MPSolver::BASIC;

// MPVariable: writer API.
%unignore operations_research::MPVariable::SetLb;
%unignore operations_research::MPVariable::SetUb;
%unignore operations_research::MPVariable::SetBounds;
%unignore operations_research::MPVariable::SetInteger;

// MPVariable: reader API.
%unignore operations_research::MPVariable::solution_value;
%unignore operations_research::MPVariable::lb;
%unignore operations_research::MPVariable::ub;
%unignore operations_research::MPVariable::integer;  // No unit test
%unignore operations_research::MPVariable::name;  // No unit test
%unignore operations_research::MPVariable::index;  // No unit test
%unignore operations_research::MPVariable::basis_status;
%unignore operations_research::MPVariable::reduced_cost;  // For experts only.

// MPConstraint: writer API.
%unignore operations_research::MPConstraint::SetCoefficient;
%unignore operations_research::MPConstraint::SetLb;
%unignore operations_research::MPConstraint::SetUb;
%unignore operations_research::MPConstraint::SetBounds;
%unignore operations_research::MPConstraint::set_is_lazy;
%unignore operations_research::MPConstraint::Clear;  // No unit test

// MPConstraint: reader API.
%unignore operations_research::MPConstraint::GetCoefficient;
%unignore operations_research::MPConstraint::lb;
%unignore operations_research::MPConstraint::ub;
%unignore operations_research::MPConstraint::name;
%unignore operations_research::MPConstraint::index;
%unignore operations_research::MPConstraint::basis_status;
%unignore operations_research::MPConstraint::dual_value;  // For experts only.

// MPObjective: writer API.
%unignore operations_research::MPObjective::SetCoefficient;
%unignore operations_research::MPObjective::SetMinimization;
%unignore operations_research::MPObjective::SetMaximization;
%unignore operations_research::MPObjective::SetOptimizationDirection;
%unignore operations_research::MPObjective::Clear;  // No unit test
%unignore operations_research::MPObjective::SetOffset;
%unignore operations_research::MPObjective::AddOffset;  // No unit test

// MPObjective: reader API.
%unignore operations_research::MPObjective::Value;
%unignore operations_research::MPObjective::GetCoefficient;
%unignore operations_research::MPObjective::minimization;
%unignore operations_research::MPObjective::maximization;
%unignore operations_research::MPObjective::offset;
%unignore operations_research::MPObjective::Offset;
%unignore operations_research::MPObjective::BestBound;

// MPSolverParameters API. For expert users only.
// TODO(user): also strip "MP" from the class name.
%unignore operations_research::MPSolverParameters;
%unignore operations_research::MPSolverParameters::MPSolverParameters;

// Expose the MPSolverParameters::DoubleParam enum.
%unignore operations_research::MPSolverParameters::DoubleParam;
%unignore operations_research::MPSolverParameters::RELATIVE_MIP_GAP;
%unignore operations_research::MPSolverParameters::PRIMAL_TOLERANCE;
%unignore operations_research::MPSolverParameters::DUAL_TOLERANCE;
%unignore operations_research::MPSolverParameters::GetDoubleParam;
%unignore operations_research::MPSolverParameters::SetDoubleParam;
%unignore operations_research::MPSolverParameters::kDefaultRelativeMipGap;
%unignore operations_research::MPSolverParameters::kDefaultPrimalTolerance;
%unignore operations_research::MPSolverParameters::kDefaultDualTolerance;
// TODO(user): unit test kDefaultPrimalTolerance.

// Expose the MPSolverParameters::IntegerParam enum.
%unignore operations_research::MPSolverParameters::IntegerParam;
%unignore operations_research::MPSolverParameters::PRESOLVE;
%unignore operations_research::MPSolverParameters::LP_ALGORITHM;
%unignore operations_research::MPSolverParameters::INCREMENTALITY;
%unignore operations_research::MPSolverParameters::SCALING;
%unignore operations_research::MPSolverParameters::GetIntegerParam;
%unignore operations_research::MPSolverParameters::SetIntegerParam;
%unignore operations_research::MPSolverParameters::RELATIVE_MIP_GAP;
%unignore operations_research::MPSolverParameters::kDefaultPrimalTolerance;
// TODO(user): unit test kDefaultPrimalTolerance.

// Expose the MPSolverParameters::PresolveValues enum.
%unignore operations_research::MPSolverParameters::PresolveValues;
%unignore operations_research::MPSolverParameters::PRESOLVE_OFF;
%unignore operations_research::MPSolverParameters::PRESOLVE_ON;
%unignore operations_research::MPSolverParameters::kDefaultPresolve;

// Expose the MPSolverParameters::LpAlgorithmValues enum.
%unignore operations_research::MPSolverParameters::LpAlgorithmValues;
%unignore operations_research::MPSolverParameters::DUAL;
%unignore operations_research::MPSolverParameters::PRIMAL;
%unignore operations_research::MPSolverParameters::BARRIER;

// Expose the MPSolverParameters::IncrementalityValues enum.
%unignore operations_research::MPSolverParameters::IncrementalityValues;
%unignore operations_research::MPSolverParameters::INCREMENTALITY_OFF;
%unignore operations_research::MPSolverParameters::INCREMENTALITY_ON;
%unignore operations_research::MPSolverParameters::kDefaultIncrementality;

// Expose the MPSolverParameters::ScalingValues enum.
%unignore operations_research::MPSolverParameters::ScalingValues;
%unignore operations_research::MPSolverParameters::SCALING_OFF;
%unignore operations_research::MPSolverParameters::SCALING_ON;

// Expose the model exporters.
%rename (ModelExportOptions) operations_research::MPModelExportOptions;
%rename (ModelExportOptions) operations_research::MPModelExportOptions::MPModelExportOptions;
%rename (ExportModelAsLpFormat) operations_research::ExportModelAsLpFormatReturnString;
%rename (ExportModelAsMpsFormat) operations_research::ExportModelAsMpsFormatReturnString;

// Expose the model validator.
%rename (FindErrorInModelProto) operations_research::FindErrorInMPModelProto;

%include "ortools/linear_solver/linear_solver.h"
%include "ortools/linear_solver/model_exporter.h"
%include "ortools/linear_solver/model_exporter_swig_helper.h"

namespace operations_research {
  std::string FindErrorInMPModelProto(
      const operations_research::MPModelProto& input_model);
}  // namespace operations_research

%unignoreall