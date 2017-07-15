Feature: Compiling Escape plan
    As a Infrastructure engineer
    In order to validate my changes
    I need to be able to compile Escape plans

    Scenario: Compiling the default Escape plan
      Given a new Escape plan called "my-release"
      When I compile the plan
      Then I should have valid release metadata
       And the metadata should have its "name" set to "my-release"
       And the metadata should have its "project" set to "_"
       And the metadata should have its "version" set to "0.0.0"

    Scenario: Compiling with a project name
      Given a new Escape plan called "my-project/my-release"
      When I compile the plan
      Then I should have valid release metadata
       And the metadata should have its "name" set to "my-release"
       And the metadata should have its "project" set to "my-project"
       And the metadata should have its "version" set to "0.0.0"

    Scenario: Compiling auto versions (@)
      Given a new Escape plan called "my-release"
        And it has "version" set to "@"
      When I compile the plan
      Then I should have valid release metadata
       And the metadata should have its "version" set to "0"

    Scenario: Compiling auto versions (0.@)
      Given a new Escape plan called "my-release"
        And it has "version" set to "0.@"
      When I compile the plan
      Then I should have valid release metadata
       And the metadata should have its "version" set to "0.0"

    Scenario: Compiling auto versions (0.1.@)
      Given a new Escape plan called "my-release"
        And it has "version" set to "0.1.@"
      When I compile the plan
      Then I should have valid release metadata
       And the metadata should have its "version" set to "0.1.0"
